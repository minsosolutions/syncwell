package output

import (
	"strings"
	"testing"
)

func TestWriteRendersBothKindsOfEvidence(t *testing.T) {
	b := &strings.Builder{}
	err := Write(b, Report{
		Customer: "nordstad",
		RunAt:    "2026-09-16T09:00:00Z",
		Items: []Finding{{
			Key:      "gdpr-retention-answer-outstanding",
			Title:    "Retention answer never sent",
			Severity: "high",
			Evidence: []Evidence{
				{Kind: "record", Source: "email", Path: "email/2026-09-11-x.md", Locator: "m1@nordstad.se", Quote: "I need a written statement"},
				{Kind: "absence", Claim: "No reply from anyone at Minso", Searched: []string{"email", "slack"}, Since: "2026-09-11"},
			},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	got := b.String()
	for _, want := range []string{
		"[email/2026-09-11-x.md](email/2026-09-11-x.md#m1@nordstad.se)", // a record resolves to one locator, not a whole file
		"> I need a written statement",
		"**Nothing found:** No reply from anyone at Minso", // an absence is stated, not faked as a citation
		"searched email, slack since 2026-09-11",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("rendered report is missing %q\n---\n%s", want, got)
		}
	}
}

func TestWriteSortsBySeverity(t *testing.T) {
	b := &strings.Builder{}
	if err := Write(b, Report{Customer: "c", Items: []Finding{
		{Title: "third", Severity: "weird"},
		{Title: "second", Severity: "low"},
		{Title: "first", Severity: "high"},
	}}); err != nil {
		t.Fatal(err)
	}
	got := b.String()
	if strings.Index(got, "first") > strings.Index(got, "second") || strings.Index(got, "second") > strings.Index(got, "third") {
		t.Errorf("findings are not ordered high, low, unknown:\n%s", got)
	}
}
