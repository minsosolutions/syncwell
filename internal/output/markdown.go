// Package output holds the reference Output writer.
//
// Only markdown is implemented. Spreadsheet, Linear and email Outputs are documented in
// the README and left for you to write.
package output

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

type Item struct {
	Title   string   `json:"title"`
	Summary string   `json:"summary"`
	Action  string   `json:"action"`
	Sources []string `json:"sources"`
	// Severity is one of high, medium, low. Anything else sorts last.
	Severity string `json:"severity"`
}

type Report struct {
	Customer string `json:"customer"`
	RunAt    string `json:"run_at"`
	Items    []Item `json:"items"`
}

var severityRank = map[string]int{"high": 0, "medium": 1, "low": 2}

func rank(s string) int {
	if r, ok := severityRank[strings.ToLower(s)]; ok {
		return r
	}
	return len(severityRank)
}

// Write renders the Report and stamps it with a Provenance Marker, so a later Run can tell
// its own output from a human's edits. Whether the Marker or the Manifest wins when they
// disagree is deliberately not decided here.
func Write(w io.Writer, r Report) error {
	items := make([]Item, len(r.Items))
	copy(items, r.Items)
	sort.SliceStable(items, func(i, j int) bool { return rank(items[i].Severity) < rank(items[j].Severity) })

	b := &strings.Builder{}
	fmt.Fprintf(b, "# Needs attention — %s\n\n", r.Customer)
	fmt.Fprintf(b, "_Generated %s. %d item(s)._\n", r.RunAt, len(items))

	if len(items) == 0 {
		b.WriteString("\nNothing needs attention.\n")
	}
	for _, it := range items {
		fmt.Fprintf(b, "\n## %s\n\n", it.Title)
		if s := strings.ToLower(it.Severity); s != "" {
			fmt.Fprintf(b, "**Severity:** %s\n\n", s)
		}
		if it.Summary != "" {
			fmt.Fprintf(b, "%s\n", it.Summary)
		}
		if it.Action != "" {
			fmt.Fprintf(b, "\n**Suggested action:** %s\n", it.Action)
		}
		if len(it.Sources) > 0 {
			b.WriteString("\n**Evidence:**\n\n")
			for _, s := range it.Sources {
				fmt.Fprintf(b, "- %s\n", s)
			}
		}
	}
	fmt.Fprintf(b, "\n<!-- syncwell:run=%s customer=%s item=report -->\n", r.RunAt, r.Customer)

	_, err := io.WriteString(w, b.String())
	return err
}
