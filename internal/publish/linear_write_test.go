package publish

import (
	"context"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/minsosolutions/syncwell/internal/linear"
	"github.com/minsosolutions/syncwell/internal/mockapi"
	"github.com/minsosolutions/syncwell/internal/output"
	"github.com/minsosolutions/syncwell/internal/state"
)

// The second Run is the whole question, so this exercises three in a row against a real
// Linear API: create, no-op, and — once a human has retitled the issue — comment instead of
// overwrite. The Finding is identical every time; only the world around it changes.
func TestSecondRunDoesNotDuplicateAndThirdRespectsAHumansEdit(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "linear"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "linear", "issues.json"), []byte("[]"), 0o644); err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(mockapi.New(dir))
	defer srv.Close()
	c := linear.New(srv.URL+"/linear/graphql", "dev")

	m := &state.Manifest{Customer: "ebbahus"}
	report := output.Report{Customer: "ebbahus", Items: []output.Finding{{
		Key: "export-off", Title: "Nightly export has been off since 5 August", Severity: "high",
		Summary:  "Disabled while debugging and never re-enabled.",
		Evidence: []output.Evidence{{Kind: "record", Source: "slack", Path: "slack/cust-ebbahus-support.json", Locator: "1785919920.000274"}},
	}}}
	opts := LinearOptions{Customer: "ebbahus", ProjectID: "prj_0002", RunAt: "2026-09-17T06:00:00Z"}

	actions, _, err := Linear(context.Background(), c, m, report, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(actions) != 1 || actions[0].Verb != "created" {
		t.Fatalf("first run should create one issue: %v", actions)
	}
	created := issues(t, c)[0]
	if marker, ok := state.FindMarker(created.Description); !ok || marker.Item != "export-off" {
		t.Fatalf("a created issue must carry its provenance marker: %q", created.Description)
	}

	opts.RunAt = "2026-09-18T06:00:00Z"
	actions, _, err = Linear(context.Background(), c, m, report, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(actions) != 1 || actions[0].Verb != "skipped" {
		t.Fatalf("the same finding must not be filed twice: %v", actions)
	}
	if got := len(issues(t, c)); got != 1 {
		t.Fatalf("issue count after two runs = %d, want 1", got)
	}

	// A human retitles it. From here the title is theirs; new information goes beside it.
	if _, err := c.UpdateIssue(context.Background(), created.ID, map[string]any{"title": "URGENT: ledger gap"}); err != nil {
		t.Fatal(err)
	}
	report.Items[0].Summary = "Ulrika needs September in the ledger before month end."
	report.Items[0].Title = "Nightly export off since 5 August — ledger gap"
	opts.RunAt = "2026-09-19T06:00:00Z"
	actions, drift, err := Linear(context.Background(), c, m, report, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(drift) != 1 || drift[0].Kind != "edited" {
		t.Fatalf("the retitle should surface as drift: %v", drift)
	}
	if got := issues(t, c)[0].Title; got != "URGENT: ledger gap" {
		t.Errorf("a human's title was overwritten: %q", got)
	}
	if actions[0].Verb != "updated" || !strings.Contains(actions[0].Detail, "description") {
		t.Errorf("the description is still ours to write: %v", actions[0])
	}
	comments, err := c.Comments(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(comments) != 1 || !strings.Contains(comments[0].Body, "Nightly export off since 5 August") {
		t.Fatalf("the new title should arrive as a comment, not an overwrite: %+v", comments)
	}
	if _, ok := state.FindMarker(comments[0].Body); !ok {
		t.Error("a comment is an output too; it carries its own marker")
	}
}

func issues(t *testing.T, c *linear.Client) []linear.Issue {
	t.Helper()
	got, err := c.Issues(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) == 0 {
		t.Fatal("no issues")
	}
	return got
}
