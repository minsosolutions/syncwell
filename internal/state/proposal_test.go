package state

import "testing"

// A reworded draft replaces the pending one it supersedes, so a human is never looking at an
// approval button for wording that no Run would send any more. A decided Proposal survives:
// it is the audit trail.
func TestSupersedeClearsPendingButKeepsDecided(t *testing.T) {
	dir := t.TempDir()
	old := NewProposal("ebbahus", "export-off", "email", "Subject", "First wording", []string{"nina.dahl@ebbahus.se"}, "run-1")
	sent := NewProposal("ebbahus", "export-off", "email", "Subject", "Sent wording", []string{"nina.dahl@ebbahus.se"}, "run-1")
	sent.Decide("applied", "")
	fresh := NewProposal("ebbahus", "export-off", "email", "Subject", "Second wording", []string{"nina.dahl@ebbahus.se"}, "run-2")
	other := NewProposal("ebbahus", "other-finding", "email", "Subject", "Untouched", []string{"nina.dahl@ebbahus.se"}, "run-2")
	for _, p := range []Proposal{old, sent, fresh, other} {
		if err := p.Save(dir); err != nil {
			t.Fatal(err)
		}
	}

	if err := Supersede(dir, "export-off", fresh.Hash); err != nil {
		t.Fatal(err)
	}
	got, err := LoadProposals(dir)
	if err != nil {
		t.Fatal(err)
	}
	kept := map[string]bool{}
	for _, p := range got {
		kept[p.Hash] = true
	}
	if kept[old.Hash] {
		t.Error("the superseded pending proposal should be gone")
	}
	for _, p := range []Proposal{sent, fresh, other} {
		if !kept[p.Hash] {
			t.Errorf("proposal %s should have survived", p.Hash)
		}
	}
}

func TestASuppressionIsRememberedAgainstTheFindingAndAction(t *testing.T) {
	m := &Manifest{Customer: "ebbahus"}
	m.Record(Action{Key: "export-off", Kind: "email", State: "suppressed", Hash: "abc"})

	if !m.Suppressed("export-off", "email") {
		t.Error("a rejection should suppress that action for that finding")
	}
	if m.Suppressed("export-off", "linear_issue") {
		t.Error("a rejected email says nothing about writing an issue")
	}
	if m.Suppressed("other-finding", "email") {
		t.Error("a suppression belongs to one finding, not to the customer")
	}
}
