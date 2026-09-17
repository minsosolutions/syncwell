// Package agent runs one Customer's Agent and validates what it says before any of it is
// believed. An Agent's output is a proposal, not a fact: every quote must be in the record
// it cites, and every path must be inside the Customer Directory.
package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/minsosolutions/syncwell/internal/output"
)

// Drop records a Finding that did not survive validation, for the run summary.
type Drop struct {
	Key    string
	Reason string
}

func (d Drop) String() string { return d.Key + ": " + d.Reason }

// Verify keeps the Findings whose record Evidence checks out and reports the rest. A
// Finding is dropped whole: evidence that cannot be trusted taints the judgement built on it.
func Verify(r output.Report, customerDir string) (output.Report, []Drop) {
	var drops []Drop
	kept := r
	kept.Items = nil
	for _, f := range r.Items {
		if reason := verifyFinding(f, customerDir); reason != "" {
			drops = append(drops, Drop{Key: f.Key, Reason: reason})
			continue
		}
		kept.Items = append(kept.Items, f)
	}
	return kept, drops
}

func verifyFinding(f output.Finding, customerDir string) string {
	for _, e := range f.Evidence {
		if e.Kind != "record" {
			continue
		}
		path, err := resolve(customerDir, e.Path)
		if err != nil {
			return err.Error()
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Sprintf("evidence cites %s, which cannot be read: %v", e.Path, err)
		}
		if e.Quote != "" && !quoted(string(b), e.Quote) {
			return fmt.Sprintf("quote is not in %s: %q", e.Path, e.Quote)
		}
	}
	return ""
}

// resolve joins a cited path to the Customer Directory and refuses anything that leaves it.
// This is the same boundary the Routing Rule draws one step earlier; see ADR-0001.
func resolve(customerDir, rel string) (string, error) {
	if rel == "" {
		return "", fmt.Errorf("evidence cites no path")
	}
	if filepath.IsAbs(rel) {
		return "", fmt.Errorf("evidence cites an absolute path: %s", rel)
	}
	base, err := filepath.Abs(customerDir)
	if err != nil {
		return "", err
	}
	full := filepath.Join(base, filepath.Clean(rel))
	if full != base && !strings.HasPrefix(full, base+string(os.PathSeparator)) {
		return "", fmt.Errorf("evidence cites %s, which is outside the Customer Directory", rel)
	}
	return full, nil
}

// quoted reports whether the record says these words. A JSON record — Slack, Zammad — holds
// its text escaped, so the words an Agent read are not the bytes on disk; the decoded strings
// are searched too rather than punishing an Agent for the file's encoding.
func quoted(record, quote string) bool {
	want := normalise(quote)
	if strings.Contains(normalise(record), want) {
		return true
	}
	var v any
	if json.Unmarshal([]byte(record), &v) != nil {
		return false
	}
	for _, s := range strings.Split(strings.Join(texts(v, nil), "\n"), "\n") {
		if strings.Contains(normalise(s), want) {
			return true
		}
	}
	return false
}

// texts collects every string in a decoded JSON document, at any depth.
func texts(v any, acc []string) []string {
	switch t := v.(type) {
	case string:
		acc = append(acc, t)
	case []any:
		for _, x := range t {
			acc = texts(x, acc)
		}
	case map[string]any:
		for _, x := range t {
			acc = texts(x, acc)
		}
	}
	return acc
}

// normalise collapses whitespace so a quote survives being re-wrapped, without letting any
// other difference through.
func normalise(s string) string { return strings.Join(strings.Fields(s), " ") }
