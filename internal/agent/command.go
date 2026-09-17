package agent

import (
	"context"
	"encoding/json"
	"os/exec"
)

type Options struct {
	CustomerDir   string // the Agent's whole world, and its working directory
	ReferencesDir string // shared, read-only
	Prompt        string
	MaxTurns      int
	MaxBudgetUSD  string
}

// Command builds the one sub-process a Run is. Every flag here is load-bearing; the
// evidence for each is in docs/research/claude-headless-contract.md (issue #2).
//
// The Agent gets three read-only tools and no shell, which is what makes ADR-0001's
// isolation a property of the tool list rather than of the model's behaviour.
func Command(ctx context.Context, o Options) *exec.Cmd {
	if o.MaxTurns == 0 {
		o.MaxTurns = 20
	}
	if o.MaxBudgetUSD == "" {
		o.MaxBudgetUSD = "1.00"
	}
	cmd := exec.CommandContext(ctx, "claude", "-p", o.Prompt,
		"--output-format", "json",
		"--json-schema", FindingsSchema,
		"--tools", "Read,Grep,Glob",
		"--permission-mode", "dontAsk",
		"--permission-prompts", "none",
		"--strict-mcp-config",
		"--bare",
		"--setting-sources", "",
		"--settings", settings(o),
		"--max-turns", itoa(o.MaxTurns),
		"--max-budget-usd", o.MaxBudgetUSD,
	)
	cmd.Dir = o.CustomerDir
	return cmd
}

// settings is the only configuration the Run gets: everything else is switched off, so a
// Customer Directory cannot configure the Agent that reads it.
func settings(o Options) string {
	b, _ := json.Marshal(map[string]any{
		"permissions": map[string]any{
			// Enforced in every permission mode, unlike the working directory on its own.
			"blockReadsOutsideWorkingDirectories": true,
			// Grants file access without loading skills or commands out of the tree, which
			// --add-dir would. References are participant-editable; don't make them executable.
			"additionalDirectories": []string{o.ReferencesDir},
			"deny": []string{
				// A leading // is absolute; a single / anchors at the settings source.
				"Edit(//" + trimLeadingSlash(o.ReferencesDir) + "/**)",
				// Reconciliation is Go's business (#3). The Agent reads state/previous-findings.json
				// and nothing else in state/.
				"Read(state/manifest.json)",
			},
		},
	})
	return string(b)
}

func trimLeadingSlash(s string) string {
	for len(s) > 0 && s[0] == '/' {
		s = s[1:]
	}
	return s
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for ; n > 0; n /= 10 {
		b = append([]byte{byte('0' + n%10)}, b...)
	}
	return string(b)
}
