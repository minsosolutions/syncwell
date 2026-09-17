package publish

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/minsosolutions/syncwell/internal/linear"
	"github.com/minsosolutions/syncwell/internal/output"
	"github.com/minsosolutions/syncwell/internal/state"
)

// Action is one write a Run made, for the run summary.
type Action struct {
	Verb       string `json:"verb"` // created | updated | commented | skipped
	Identifier string `json:"identifier"`
	Key        string `json:"key"`
	Detail     string `json:"detail"`
}

func (a Action) String() string {
	return strings.TrimSpace(fmt.Sprintf("%s %s %s %s", a.Verb, a.Identifier, a.Key, a.Detail))
}

type LinearOptions struct {
	Customer  string
	ProjectID string
	RunAt     string
}

// Linear writes the Findings into the Customer's project and reconciles what a human did
// since the last Run. Creating and updating issues inside our own project is working
// material and proceeds unattended; nothing here crosses to the Customer (ADR-0003).
func Linear(ctx context.Context, c *linear.Client, m *state.Manifest, r output.Report, o LinearOptions) ([]Action, []Drift, error) {
	all, err := c.Issues(ctx)
	if err != nil {
		return nil, nil, err
	}
	live := linear.InProject(all, o.ProjectID)
	drifts := Reconcile(m, live)

	byID := map[string]linear.Issue{}
	for _, i := range live {
		byID[i.ID] = i
	}

	var actions []Action
	seen := map[string]bool{}
	for _, f := range r.Items {
		key, entry := identify(m, f)
		seen[key] = true
		act, err := applyFinding(ctx, c, m, f, key, entry, byID, o)
		if err != nil {
			return actions, drifts, err
		}
		actions = append(actions, act)
	}

	agedOut, err := commentOnAgedOut(ctx, c, m, byID, seen, o)
	return append(actions, agedOut...), drifts, err
}

// identify resolves the Finding to the Output it continues, if any. The Agent proposes the
// key; it is honoured only when the Finding still stands on a record the stored one cited.
// A key that no longer shares evidence is a different concern wearing an old name (#10).
func identify(m *state.Manifest, f output.Finding) (string, *state.Output) {
	records := recordKeys(f)
	entry := m.ByKey("linear_issue", f.Key)
	if entry == nil || entry.SharesRecord(records) {
		return f.Key, entry
	}
	for n := 2; ; n++ {
		key := fmt.Sprintf("%s-%d", f.Key, n)
		e := m.ByKey("linear_issue", key)
		if e == nil {
			return key, nil
		}
		if e.SharesRecord(records) {
			return key, e
		}
	}
}

func applyFinding(ctx context.Context, c *linear.Client, m *state.Manifest, f output.Finding, key string, entry *state.Output, byID map[string]linear.Issue, o LinearOptions) (Action, error) {
	// Prose is what Stamp will leave behind, so the Snapshot records what we actually wrote.
	// Storing the untrimmed body instead makes every later Run read its own write as a
	// human's edit.
	title, body, priority := f.Title, state.Prose(issueBody(f)), priorityOf(f.Severity)
	marker := state.Marker{Run: o.RunAt, Customer: o.Customer, Item: key}

	if entry == nil {
		issue, err := c.CreateIssue(ctx, linear.IssueInput{
			Title: title, Description: state.Stamp(body, marker), Priority: priority, ProjectID: o.ProjectID,
		})
		if err != nil {
			return Action{}, err
		}
		m.Put(state.Output{
			Kind: "linear_issue", Key: key, ID: issue.ID, Identifier: issue.Identifier,
			AuthoredRun: o.RunAt, LastRun: o.RunAt, Records: recordKeys(f),
			Snapshot: state.Snapshot{Title: title, Description: body, Priority: priority},
		})
		return Action{"created", issue.Identifier, key, title}, nil
	}

	if entry.Deleted {
		return Action{"skipped", entry.Identifier, key, "suppressed: a human deleted it"}, nil
	}
	issue, ok := byID[entry.ID]
	if !ok {
		return Action{"skipped", entry.Identifier, key, "not in the project any more"}, nil
	}
	if _, stewarded := state.FindMarker(issue.Description); !stewarded {
		return Action{"skipped", entry.Identifier, key, "stewardship withdrawn: the marker was stripped"}, nil
	}

	// Fields a human has taken are never overwritten. New information about them arrives as
	// a run-stamped comment beside the human's words (#9).
	fields := map[string]any{}
	var toldInComment, toldNames []string
	for _, ch := range []struct {
		name          string
		changed       bool
		value         any
		humanReadable string
	}{
		{"title", title != entry.Snapshot.Title, title, "Title we would write now: " + title},
		{"description", body != entry.Snapshot.Description, state.Stamp(body, marker), body},
		{"priority", priority != entry.Snapshot.Priority, priority, fmt.Sprintf("Severity is now %s (priority %d).", f.Severity, priority)},
	} {
		if !ch.changed {
			continue
		}
		if slices.Contains(entry.HumanFields, ch.name) {
			toldInComment = append(toldInComment, ch.humanReadable)
			toldNames = append(toldNames, ch.name)
			continue
		}
		fields[ch.name] = ch.value
	}

	entry.LastRun = o.RunAt
	entry.Records = recordKeys(f)
	entry.AgedRun = ""
	entry.Snapshot = state.Snapshot{Title: title, Description: body, Priority: priority}

	if len(toldInComment) > 0 {
		note := "Run " + o.RunAt + " has new information for fields a human now owns.\n\n" +
			strings.Join(toldInComment, "\n\n")
		if _, err := c.CreateComment(ctx, entry.ID, state.Stamp(note, marker)); err != nil {
			return Action{}, err
		}
	}
	if len(fields) == 0 {
		verb := "skipped"
		detail := "unchanged"
		if len(toldInComment) > 0 {
			verb, detail = "commented", "human-owned fields: "+strings.Join(toldNames, ", ")
		}
		return Action{verb, entry.Identifier, key, detail}, nil
	}
	if _, err := c.UpdateIssue(ctx, entry.ID, fields); err != nil {
		return Action{}, err
	}
	return Action{"updated", entry.Identifier, key, strings.Join(sortedKeys(fields), ", ")}, nil
}

// commentOnAgedOut says so on an Output whose Finding this Run did not surface, and does not
// close it. A Run that stops seeing something is not the same as the thing being resolved,
// and only a human knows which happened (#10).
func commentOnAgedOut(ctx context.Context, c *linear.Client, m *state.Manifest, byID map[string]linear.Issue, seen map[string]bool, o LinearOptions) ([]Action, error) {
	var actions []Action
	for i := range m.Outputs {
		entry := &m.Outputs[i]
		if entry.Kind != "linear_issue" || entry.Deleted || seen[entry.Key] || entry.AgedRun != "" {
			continue
		}
		issue, ok := byID[entry.ID]
		if !ok {
			continue
		}
		if _, stewarded := state.FindMarker(issue.Description); !stewarded {
			continue
		}
		note := "Run " + o.RunAt + " did not surface this finding again. Left open deliberately: a run " +
			"that stops seeing something has not resolved it."
		marker := state.Marker{Run: o.RunAt, Customer: o.Customer, Item: entry.Key}
		if _, err := c.CreateComment(ctx, entry.ID, state.Stamp(note, marker)); err != nil {
			return actions, err
		}
		entry.AgedRun = o.RunAt
		actions = append(actions, Action{"commented", entry.Identifier, entry.Key, "not surfaced by this run"})
	}
	return actions, nil
}

// recordKeys is the Finding's identity material: the records it stands on, as path#locator.
func recordKeys(f output.Finding) []string {
	var out []string
	for _, e := range f.Evidence {
		if e.Kind == "record" && e.Path != "" {
			out = append(out, e.Path+"#"+e.Locator)
		}
	}
	slices.Sort(out)
	return slices.Compact(out)
}

// priorityOf maps severity to Linear's 1-4 scale. Nothing a Run writes is Urgent(1) by
// itself; that is a human's call.
func priorityOf(severity string) int {
	switch strings.ToLower(severity) {
	case "high":
		return 2
	case "medium":
		return 3
	default:
		return 4
	}
}

func issueBody(f output.Finding) string {
	b := &strings.Builder{}
	b.WriteString(f.Summary)
	if f.Action != "" {
		fmt.Fprintf(b, "\n\n**Suggested action:** %s", f.Action)
	}
	if len(f.Evidence) > 0 {
		b.WriteString("\n\n**Evidence**\n")
		for _, e := range f.Evidence {
			if e.Kind == "absence" {
				fmt.Fprintf(b, "\n- Nothing found: %s (searched %s since %s)", e.Claim, strings.Join(e.Searched, ", "), e.Since)
				continue
			}
			fmt.Fprintf(b, "\n- `%s` %s %s#%s\n  > %s", e.Source, e.At, e.Path, e.Locator, e.Quote)
		}
	}
	return b.String()
}

func sortedKeys(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	slices.Sort(out)
	return out
}
