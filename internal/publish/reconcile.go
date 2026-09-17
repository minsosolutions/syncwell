// Package publish turns validated Findings into Outputs. Every decision in here is
// deterministic Go over records we keep ourselves: no model reads any of it and none of it
// reaches a model's context. See docs/adr/0002-authorship-and-stewardship-are-separate-signals.md.
package publish

import (
	"fmt"
	"slices"

	"github.com/minsosolutions/syncwell/internal/linear"
	"github.com/minsosolutions/syncwell/internal/state"
)

// Drift is one divergence between what the Manifest claims and what the Output system holds.
// It is surfaced in the run summary and the admin UI, never in report.md: the Customer's
// report is about the Customer, not about our bookkeeping.
type Drift struct {
	Kind       string `json:"kind"` // deleted | edited | unstewarded | unmanaged | adopted
	Identifier string `json:"identifier"`
	Key        string `json:"key,omitempty"`
	Detail     string `json:"detail"`
}

func (d Drift) String() string { return fmt.Sprintf("%s %s: %s", d.Kind, d.Identifier, d.Detail) }

// Reconcile compares the Manifest against the Customer's live project and records what a
// human did between Runs. It mutates the Manifest — a deletion becomes a Suppression, a
// diverged field passes to the human — and returns what it found.
//
// Authorship is the union of the two traces: either the Manifest or a Marker saying "a Run
// made this" is enough. Stewardship is the Marker alone, because stripping it inside the
// Output is the one hands-off gesture a human can make without being taught a mechanism.
func Reconcile(m *state.Manifest, live []linear.Issue) []Drift {
	var drifts []Drift
	byID := map[string]linear.Issue{}
	for _, i := range live {
		byID[i.ID] = i
	}

	for idx := range m.Outputs {
		o := &m.Outputs[idx]
		if o.Kind != "linear_issue" || o.Deleted {
			continue
		}
		issue, ok := byID[o.ID]
		if !ok {
			// A human deleted work we authored. That is a Suppression, not a gap to refill.
			o.Deleted = true
			drifts = append(drifts, Drift{"deleted", o.Identifier, o.Key, "deleted in Linear; not re-created"})
			continue
		}
		marker, stewarded := state.FindMarker(issue.Description)
		if !stewarded || marker.Customer != m.Customer {
			drifts = append(drifts, Drift{"unstewarded", o.Identifier, o.Key, "provenance marker removed; Syncwell no longer maintains it"})
			continue
		}
		for _, f := range divergedFields(*o, issue) {
			if slices.Contains(o.HumanFields, f) {
				continue
			}
			o.HumanFields = append(o.HumanFields, f)
			drifts = append(drifts, Drift{"edited", o.Identifier, o.Key, f + " edited by a human; it is theirs from now on"})
		}
	}

	for _, issue := range live {
		if m.ByID(issue.ID) != nil {
			continue
		}
		marker, ok := state.FindMarker(issue.Description)
		if ok && marker.Customer == m.Customer {
			// The Marker alone establishes Authorship, so an Output missing from the
			// Manifest is adopted rather than filed beside as a duplicate.
			m.Put(state.Output{
				Kind: "linear_issue", Key: marker.Item, ID: issue.ID, Identifier: issue.Identifier,
				AuthoredRun: marker.Run,
				Snapshot:    state.Snapshot{Title: issue.Title, Description: state.Prose(issue.Description), Priority: issue.Priority},
			})
			drifts = append(drifts, Drift{"adopted", issue.Identifier, marker.Item, "carries our marker but was missing from the manifest"})
			continue
		}
		drifts = append(drifts, Drift{"unmanaged", issue.Identifier, "", "created by a human; listed, never written to"})
	}
	return drifts
}

// divergedFields names what a human changed since our last write. The comparison is against
// the Written Snapshot, not against what we would write today: otherwise every new Finding
// would read as a human edit.
func divergedFields(o state.Output, issue linear.Issue) []string {
	var out []string
	if o.Snapshot.Title != "" && issue.Title != o.Snapshot.Title {
		out = append(out, "title")
	}
	if o.Snapshot.Description != "" && state.Prose(issue.Description) != o.Snapshot.Description {
		out = append(out, "description")
	}
	if o.Snapshot.Priority != 0 && issue.Priority != o.Snapshot.Priority {
		out = append(out, "priority")
	}
	return out
}

// Unmanaged lists the issues in the Customer's project that no Run authored, so a project
// manager sees them beside ours rather than discovering the duplicate later.
func Unmanaged(drifts []Drift) []Drift {
	var out []Drift
	for _, d := range drifts {
		if d.Kind == "unmanaged" {
			out = append(out, d)
		}
	}
	return out
}
