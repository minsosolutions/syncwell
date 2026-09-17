// Package run is one Run end to end: the Agent judges, and everything after it is
// deterministic Go. It is the same code path whether a Run is started from cmd/run or from
// the admin UI's trigger, so the demo and the command line cannot drift apart.
package run

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/minsosolutions/syncwell/internal/agent"
	"github.com/minsosolutions/syncwell/internal/linear"
	"github.com/minsosolutions/syncwell/internal/output"
	"github.com/minsosolutions/syncwell/internal/publish"
	"github.com/minsosolutions/syncwell/internal/routing"
	"github.com/minsosolutions/syncwell/internal/state"
)

type Options struct {
	DataDir  string
	MaxTurns int
	Budget   string
	// Findings replays a findings JSON instead of launching the Agent. The rest of the Run
	// is identical, which is what makes a deterministic rehearsal of the Output side possible.
	Findings io.Reader
}

// Summary is what a Run leaves behind for a human: what it wrote, what it proposed, and what
// a human had changed underneath it. Stored per Run so the admin UI can show a finished one.
type Summary struct {
	Customer   string           `json:"customer"`
	RunAt      string           `json:"run_at"`
	Findings   int              `json:"findings"`
	Dropped    []string         `json:"dropped,omitempty"`
	Denials    []string         `json:"denials,omitempty"`
	CostUSD    float64          `json:"cost_usd"`
	Actions    []publish.Action `json:"actions,omitempty"`
	Drift      []publish.Drift  `json:"drift,omitempty"`
	Proposals  []state.Proposal `json:"proposals,omitempty"`
	ReportPath string           `json:"report_path"`
	Report     output.Report    `json:"report"`
	Duration   string           `json:"duration"`
}

// Once runs one Customer's Agent and writes every Output the decisions allow it to write
// unattended. Gated Actions are left as Proposals; the Run does not wait for a human.
func Once(ctx context.Context, cfg *routing.Config, customer string, o Options) (*Summary, error) {
	cust := cfg.Customer(customer)
	if cust == nil {
		return nil, fmt.Errorf("%q is not a configured customer", customer)
	}
	customerDir, err := filepath.Abs(filepath.Join(o.DataDir, "customers", customer))
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(customerDir); err != nil {
		return nil, fmt.Errorf("no Customer Directory for %q: %w", customer, err)
	}

	started := time.Now()
	runAt := started.UTC().Format(time.RFC3339)

	report, cost, denials, err := Findings(ctx, o, customerDir, filepath.Join(o.DataDir, "references"), customer, runAt)
	if err != nil {
		return nil, err
	}
	report.Customer, report.RunAt = customer, runAt

	report, dropped := agent.Verify(report, customerDir)
	report, silent := agent.ResolveAbsences(report, customerDir)
	dropped = append(dropped, silent...)

	sum := &Summary{Customer: customer, RunAt: runAt, Findings: len(report.Items), CostUSD: cost, Report: report}
	for _, d := range dropped {
		sum.Dropped = append(sum.Dropped, d.String())
	}
	sum.Denials = denials

	if err := agent.WritePreviousFindings(customerDir, report); err != nil {
		return nil, err
	}

	m, err := state.LoadManifest(customerDir, customer)
	if err != nil {
		return nil, err
	}

	sum.ReportPath, err = writeReport(customerDir, report)
	if err != nil {
		return nil, err
	}
	m.Put(state.Output{Kind: "file", Key: "report", Path: "outputs/report.md", AuthoredRun: runAt, LastRun: runAt})

	client := linear.New(cfg.Linear.URL, cfg.Linear.Token)
	actions, drift, err := publish.Linear(ctx, client, m, report, publish.LinearOptions{
		Customer: customer, ProjectID: cust.LinearProjectID, RunAt: runAt,
	})
	sum.Actions, sum.Drift = actions, drift
	if err != nil {
		// Whatever was written before the failure is already in the Manifest; losing it
		// would leave Outputs we authored with no record that we did.
		_ = m.Save(customerDir)
		return sum, fmt.Errorf("writing to linear: %w", err)
	}

	proposals, proposeActions, err := publish.Propose(customerDir, cust, m, report, runAt)
	if err != nil {
		return sum, err
	}
	sum.Proposals = proposals
	sum.Actions = append(sum.Actions, proposeActions...)

	if err := m.Save(customerDir); err != nil {
		return sum, err
	}
	sum.Duration = time.Since(started).Round(time.Second).String()
	return sum, saveSummary(customerDir, sum)
}

// Findings is the Agent half of a Run: launch it under the confinement of internal/agent,
// parse the envelope, and hand back what it said. Nothing is believed yet — the caller
// verifies quotes and resolves absences before any of it becomes an Output.
func Findings(ctx context.Context, o Options, customerDir, referencesDir, customer, runAt string) (output.Report, float64, []string, error) {
	if o.Findings != nil {
		var r output.Report
		err := json.NewDecoder(o.Findings).Decode(&r)
		return r, 0, nil, err
	}

	referencesDir, err := filepath.Abs(referencesDir)
	if err != nil {
		return output.Report{}, 0, nil, err
	}
	working, err := agent.Files(customerDir)
	if err != nil {
		return output.Report{}, 0, nil, err
	}
	references, err := agent.Files(referencesDir)
	if err != nil {
		return output.Report{}, 0, nil, err
	}
	for i, r := range references {
		references[i] = filepath.Join(referencesDir, r)
	}

	cmd := agent.Command(ctx, agent.Options{
		CustomerDir:   customerDir,
		ReferencesDir: referencesDir,
		Prompt:        agent.Prompt(customer, runAt, working, references),
		MaxTurns:      o.MaxTurns,
		MaxBudgetUSD:  o.Budget,
	})
	cmd.Stderr = os.Stderr
	stdout, cmdErr := cmd.Output()
	res, err := agent.Parse(stdout)
	if err != nil {
		if cmdErr != nil {
			return output.Report{}, 0, nil, fmt.Errorf("%w (exit: %v)", err, cmdErr)
		}
		return output.Report{}, 0, nil, err
	}
	var denials []string
	for _, d := range res.Denials {
		denials = append(denials, d.ToolName+" "+d.ToolUseID)
	}
	return res.Report, res.CostUSD, denials, nil
}

func writeReport(customerDir string, r output.Report) (string, error) {
	path := filepath.Join(customerDir, "outputs", "report.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return path, output.Write(f, r)
}

func saveSummary(customerDir string, s *Summary) error {
	dir := filepath.Join(customerDir, "state", "runs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	// Two Runs inside one second only happen in a rehearsal, and losing one of their
	// summaries there is exactly when you wanted both.
	name := strings.NewReplacer(":", "-", "/", "-").Replace(s.RunAt)
	path := filepath.Join(dir, name+".json")
	for n := 2; ; n++ {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			break
		}
		path = filepath.Join(dir, fmt.Sprintf("%s-%d.json", name, n))
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

// Summaries lists a Customer's finished Runs, newest first.
func Summaries(dataDir, customer string) ([]Summary, error) {
	paths, err := filepath.Glob(filepath.Join(dataDir, "customers", customer, "state", "runs", "*.json"))
	if err != nil {
		return nil, err
	}
	sort.Sort(sort.Reverse(sort.StringSlice(paths)))
	var out []Summary
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		var s Summary
		if err := json.Unmarshal(b, &s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}
