// Command agent runs one Customer's Agent and prints the findings it can stand behind.
//
// The Agent judges; this validates. Nothing here writes an Output: quotes are checked
// against the records they cite, absences are computed rather than believed, and what
// survives goes to stdout as JSON for an Output writer to render.
//
//	go run ./cmd/agent -customer ebbahus
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/minsosolutions/syncwell/internal/agent"
	"github.com/minsosolutions/syncwell/internal/output"
)

func main() {
	dataDir := flag.String("data", "data", "data directory")
	customer := flag.String("customer", "", "customer slug")
	timeout := flag.Duration("timeout", 10*time.Minute, "how long one Run may take")
	flag.Parse()
	if *customer == "" {
		log.Fatal("-customer is required")
	}

	customerDir, err := filepath.Abs(filepath.Join(*dataDir, "customers", *customer))
	if err != nil {
		log.Fatal(err)
	}
	referencesDir, err := filepath.Abs(filepath.Join(*dataDir, "references"))
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(customerDir); err != nil {
		log.Fatalf("no Customer Directory for %q: %v", *customer, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	cmd := agent.Command(ctx, agent.Options{
		CustomerDir:   customerDir,
		ReferencesDir: referencesDir,
		Prompt:        agent.Prompt,
	})
	cmd.Stderr = os.Stderr
	stdout, err := cmd.Output()
	// A failure inside the run is printed as the result on stdout, so decode before
	// deciding what a non-zero exit means.
	res, perr := agent.Parse(stdout)
	if perr != nil {
		if err != nil {
			log.Fatalf("%v (exit: %v)", perr, err)
		}
		log.Fatal(perr)
	}

	report, dropped := agent.Verify(res.Report, customerDir)
	report, silent := agent.ResolveAbsences(report, customerDir)
	dropped = append(dropped, silent...)

	for _, d := range dropped {
		fmt.Fprintln(os.Stderr, "dropped", d)
	}
	for _, d := range res.Denials {
		fmt.Fprintln(os.Stderr, "denied", d.ToolName, d.ToolUseID)
	}
	fmt.Fprintf(os.Stderr, "%d finding(s), %d dropped, $%.3f\n", len(report.Items), len(dropped), res.CostUSD)

	if err := writePreviousFindings(customerDir, report); err != nil {
		log.Fatal(err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		log.Fatal(err)
	}
}

// writePreviousFindings leaves the next Run the keys and titles this one surfaced, and
// nothing else: the Manifest stays out of the Agent's view. See issue #10.
func writePreviousFindings(customerDir string, r output.Report) error {
	type previous struct {
		Key     string `json:"key"`
		Title   string `json:"title"`
		Summary string `json:"summary"`
	}
	items := make([]previous, 0, len(r.Items))
	for _, f := range r.Items {
		items = append(items, previous{Key: f.Key, Title: f.Title, Summary: f.Summary})
	}
	b, err := json.MarshalIndent(map[string]any{"run_at": r.RunAt, "findings": items}, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Join(customerDir, "state")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "previous-findings.json"), b, 0o644)
}
