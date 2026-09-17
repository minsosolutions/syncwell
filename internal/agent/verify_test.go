package agent

import (
	"path/filepath"
	"testing"

	"github.com/minsosolutions/syncwell/internal/output"
)

func customerDir(t *testing.T) string {
	t.Helper()
	return filepath.Join("..", "..", "data", "customers", "ebbahus")
}

func TestVerifyKeepsAFindingWhoseQuoteIsInTheRecord(t *testing.T) {
	r := output.Report{Customer: "ebbahus", Items: []output.Finding{{
		Key:   "bookkeeping-export-disabled",
		Title: "Export off since 5 August",
		Evidence: []output.Evidence{{
			Kind: "record", Source: "email",
			Path:  "email/2026-09-09-ebbahus-saknade-utbetalningar.md",
			Quote: "There is nothing from September in the bookkeeping.",
		}},
	}}}

	kept, dropped := Verify(r, customerDir(t))
	if len(kept.Items) != 1 {
		t.Fatalf("a Finding whose quote is really in the record was dropped: %v", dropped)
	}
	if len(dropped) != 0 {
		t.Fatalf("unexpected drops: %v", dropped)
	}
}

func TestVerifyDropsAFindingWithAQuoteNobodyWrote(t *testing.T) {
	r := output.Report{Customer: "ebbahus", Items: []output.Finding{{
		Key: "invented",
		Evidence: []output.Evidence{{
			Kind: "record", Source: "email",
			Path:  "email/2026-09-09-ebbahus-saknade-utbetalningar.md",
			Quote: "We have decided to terminate the contract immediately.",
		}},
	}}}

	kept, dropped := Verify(r, customerDir(t))
	if len(kept.Items) != 0 {
		t.Fatal("a hallucinated quote reached the report")
	}
	if len(dropped) != 1 {
		t.Fatalf("the drop was not reported: %v", dropped)
	}
}

func TestVerifyDropsEvidenceEscapingTheCustomerDirectory(t *testing.T) {
	for _, bad := range []string{"../nordstad/email/2026-09-11-nordstad-gallringstid-avslag.md", "/etc/passwd", "email/../../nordstad/state/manifest.json"} {
		r := output.Report{Customer: "ebbahus", Items: []output.Finding{{
			Key:      "leak",
			Evidence: []output.Evidence{{Kind: "record", Path: bad, Quote: "anything"}},
		}}}
		kept, dropped := Verify(r, customerDir(t))
		if len(kept.Items) != 0 {
			t.Fatalf("evidence at %q was accepted; it points outside the Customer Directory", bad)
		}
		if len(dropped) != 1 {
			t.Fatalf("the escape at %q was not reported", bad)
		}
	}
}
