package publish

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/minsosolutions/syncwell/internal/linear"
	"github.com/minsosolutions/syncwell/internal/state"
)

// The four seeded Drift cases are the acceptance test for reconciliation: one of each
// disagreement a human can create between the Manifest and Linear. See facilitator/SOLUTIONS.md.
func TestReconcileFindsTheFourSeededDriftCases(t *testing.T) {
	want := map[string]struct{ customer, kind string }{
		"SYN-14": {"nordstad", "deleted"},
		"SYN-21": {"ebbahus", "edited"},
		"SYN-31": {"utvecklingsverket", "unstewarded"},
		"SYN-33": {"utvecklingsverket", "unmanaged"},
	}
	projects := map[string]string{"nordstad": "prj_0001", "ebbahus": "prj_0002", "utvecklingsverket": "prj_0003"}

	got := map[string]string{}
	for customer, project := range projects {
		m, err := state.LoadManifest(filepath.Join("..", "..", "data", "customers", customer), customer)
		if err != nil {
			t.Fatal(err)
		}
		for _, d := range Reconcile(m, linear.InProject(seededIssues(t), project)) {
			got[d.Identifier] = d.Kind
		}
	}

	for ident, w := range want {
		if got[ident] != w.kind {
			t.Errorf("%s (%s): drift %q, want %q", ident, w.customer, got[ident], w.kind)
		}
	}
}

func TestDeletionBecomesASuppression(t *testing.T) {
	m, err := state.LoadManifest(filepath.Join("..", "..", "data", "customers", "nordstad"), "nordstad")
	if err != nil {
		t.Fatal(err)
	}
	Reconcile(m, linear.InProject(seededIssues(t), "prj_0001"))

	entry := m.ByKey("linear_issue", "launch-readiness")
	if entry == nil || !entry.Deleted {
		t.Fatalf("SYN-14 was deleted in Linear; the manifest entry should carry the suppression: %+v", entry)
	}
	// A second Run must not re-report it, and must never re-create it.
	for _, d := range Reconcile(m, linear.InProject(seededIssues(t), "prj_0001")) {
		if d.Identifier == "SYN-14" {
			t.Errorf("a suppression was reported twice: %s", d)
		}
	}
}

func TestEditedFieldPassesToTheHuman(t *testing.T) {
	m, err := state.LoadManifest(filepath.Join("..", "..", "data", "customers", "ebbahus"), "ebbahus")
	if err != nil {
		t.Fatal(err)
	}
	Reconcile(m, linear.InProject(seededIssues(t), "prj_0002"))

	// The human retitled it and raised the priority; both fields are theirs from now on.
	entry := m.ByKey("linear_issue", "E1")
	if entry == nil || !slices.Equal(entry.HumanFields, []string{"title", "priority"}) {
		t.Fatalf("SYN-21's edits should hand title and priority to the human: %+v", entry)
	}
}

func seededIssues(t *testing.T) []linear.Issue {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "data", "linear", "issues.json"))
	if err != nil {
		t.Fatal(err)
	}
	var issues []linear.Issue
	if err := json.Unmarshal(b, &issues); err != nil {
		t.Fatal(err)
	}
	return issues
}
