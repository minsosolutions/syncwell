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

// Evidence is one leg a Finding stands on. Kind "record" cites a Source Record in the
// Customer Directory; kind "absence" states what is *not* there.
//
// An Agent writes only "record" evidence. Absence is computed in Go from the thread
// relation of a record the Agent flagged AwaitingReply, so a claim about everything is
// never a claim a model can make.
type Evidence struct {
	Kind string `json:"kind"` // record | absence

	// record, written by an Agent
	Source  string `json:"source,omitempty"`  // slack | email | transcript | zammad
	Path    string `json:"path,omitempty"`    // relative to the Customer Directory
	Locator string `json:"locator,omitempty"` // slack ts, message_id, line number, article id
	At      string `json:"at,omitempty"`
	Quote   string `json:"quote,omitempty"`
	// AwaitingReply marks a record the Agent believes went unanswered. It is a request for
	// Go to check, not a claim: a reply found here drops the whole Finding.
	AwaitingReply bool `json:"awaiting_reply,omitempty"`

	// absence, written by Go
	Claim    string   `json:"claim,omitempty"`
	Searched []string `json:"searched,omitempty"`
	Since    string   `json:"since,omitempty"`
}

// Proposal is a Gated Action the Run prepared and did not take. Recipients are absent by
// design: they are derived in Go from the Customer's own records, never proposed by an
// Agent. See docs/adr/0003-approval-gates-on-blast-radius.md.
type Proposal struct {
	Kind    string `json:"kind"` // email
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// Finding is one thing needing attention. Key is its identity across Runs and is what a
// Provenance Marker and a Manifest entry carry.
type Finding struct {
	Key      string     `json:"key"`
	Title    string     `json:"title"`
	Summary  string     `json:"summary"`
	Action   string     `json:"action"`
	Evidence []Evidence `json:"evidence"`
	Proposal *Proposal  `json:"proposal,omitempty"`
	// Severity is one of high, medium, low. Anything else sorts last.
	Severity string `json:"severity"`
}

type Report struct {
	Customer string    `json:"customer"`
	RunAt    string    `json:"run_at"`
	Items    []Finding `json:"items"`
}

var severityRank = map[string]int{"high": 0, "medium": 1, "low": 2}

func rank(s string) int {
	if r, ok := severityRank[strings.ToLower(s)]; ok {
		return r
	}
	return len(severityRank)
}

// Write renders the Report and stamps it with a Provenance Marker. The report is a
// Snapshot Output: the next Run regenerates it whole rather than reconciling it.
func Write(w io.Writer, r Report) error {
	items := make([]Finding, len(r.Items))
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
		if len(it.Evidence) > 0 {
			b.WriteString("\n**Evidence:**\n\n")
			for _, e := range it.Evidence {
				writeEvidence(b, e)
			}
		}
		if it.Proposal != nil {
			fmt.Fprintf(b, "\n**Proposed %s, awaiting approval:** %s\n", it.Proposal.Kind, it.Proposal.Subject)
		}
	}
	fmt.Fprintf(b, "\n<!-- syncwell:run=%s customer=%s item=report -->\n", r.RunAt, r.Customer)

	_, err := io.WriteString(w, b.String())
	return err
}

func writeEvidence(b *strings.Builder, e Evidence) {
	if e.Kind == "absence" {
		fmt.Fprintf(b, "- **Nothing found:** %s\n", e.Claim)
		if len(e.Searched) > 0 {
			fmt.Fprintf(b, "  - searched %s since %s\n", strings.Join(e.Searched, ", "), e.Since)
		}
		return
	}
	fmt.Fprintf(b, "- `%s` %s — [%s](%s)\n", e.Source, e.At, e.Path, evidenceLink(e))
	if e.Quote != "" {
		fmt.Fprintf(b, "  > %s\n", e.Quote)
	}
}

func evidenceLink(e Evidence) string {
	if e.Locator == "" {
		return e.Path
	}
	return e.Path + "#" + e.Locator
}
