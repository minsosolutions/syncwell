package agent

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/minsosolutions/syncwell/internal/output"
)

func awaiting(path string) output.Report {
	return output.Report{Items: []output.Finding{{
		Key: "awaiting", Title: "Someone asked and nobody answered",
		Evidence: []output.Evidence{{
			Kind: "record", Source: "email", Path: path,
			Quote: "", AwaitingReply: true,
		}},
	}}}
}

// N3: Anna asks on 11 September and nothing in the Customer Directory replies to it.
func TestResolveAbsencesStatesTheAbsenceWhenNobodyReplied(t *testing.T) {
	dir := filepath.Join("..", "..", "data", "customers", "nordstad")

	got, drops := ResolveAbsences(awaiting("email/2026-09-11-nordstad-gallringstid-avslag.md"), dir)
	if len(got.Items) != 1 {
		t.Fatalf("the Finding was dropped although nothing replied: %v", drops)
	}
	var absences []output.Evidence
	for _, e := range got.Items[0].Evidence {
		if e.Kind == "absence" {
			absences = append(absences, e)
		}
	}
	if len(absences) != 1 {
		t.Fatalf("no absence was stated: %+v", got.Items[0].Evidence)
	}
	a := absences[0]
	if !strings.Contains(a.Claim, "20260911102200.5518@nordstad.se") {
		t.Errorf("the absence does not say what went unanswered: %q", a.Claim)
	}
	if len(a.Searched) != 1 || a.Searched[0] != "email" {
		t.Errorf("absence claims sources it cannot check: %v", a.Searched)
	}
	if a.Since != "2026-09-11T10:22:00+02:00" {
		t.Errorf("absence since = %q, want the record's own date", a.Since)
	}
}

// E1's email was answered the next day, so the silence it rested on is not there.
func TestResolveAbsencesDropsTheFindingWhenAReplyExists(t *testing.T) {
	dir := filepath.Join("..", "..", "data", "customers", "ebbahus")

	got, drops := ResolveAbsences(awaiting("email/2026-09-09-ebbahus-saknade-utbetalningar.md"), dir)
	if len(got.Items) != 0 {
		t.Fatal("a Finding resting on silence survived although a reply exists")
	}
	if len(drops) != 1 || !strings.Contains(drops[0].Reason, "20260910083300.2214@minso.se") {
		t.Fatalf("the drop does not name the reply that killed it: %v", drops)
	}
}
