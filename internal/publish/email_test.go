package publish

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/minsosolutions/syncwell/internal/output"
	"github.com/minsosolutions/syncwell/internal/routing"
	"github.com/minsosolutions/syncwell/internal/state"
)

func ebbahus(t *testing.T) (string, *routing.Customer) {
	t.Helper()
	cfg, err := routing.Load(filepath.Join("..", "..", "config", "syncwell.json"))
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join("..", "..", "data", "customers", "ebbahus"), cfg.Customer("ebbahus")
}

// Recipients come out of the Customer's own records and are filtered to that Customer's
// configured domains: the Vendor's own addresses are in every record and are not recipients.
func TestRecipientsComeFromTheCitedRecordsAndTheCustomersDomains(t *testing.T) {
	dir, cust := ebbahus(t)
	to, err := Recipients(dir, cust, output.Finding{Evidence: []output.Evidence{
		{Kind: "record", Path: "email/2026-09-09-ebbahus-saknade-utbetalningar.md"},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(to, "nina.dahl@ebbahus.se") {
		t.Errorf("the cited mail's customer contact should be a recipient, got %v", to)
	}
	for _, a := range to {
		if !slices.Contains([]string{"nina.dahl@ebbahus.se", "gustav.ebbe@ebbahus.se"}, a) {
			t.Errorf("%q is not an ebbahus address", a)
		}
	}
}

// Evidence that cites no path at all must not produce a send to nobody, and must never reach
// outside the Customer Directory to find one.
func TestRecipientsForAFindingWithNoCustomerAddressFallsBackToTheirMail(t *testing.T) {
	dir, cust := ebbahus(t)
	to, err := Recipients(dir, cust, output.Finding{Evidence: []output.Evidence{
		{Kind: "record", Path: "slack/cust-ebbahus-support.json"},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if len(to) == 0 {
		t.Fatal("expected the customer's mail correspondents as a fallback")
	}
	for _, a := range to {
		if !slices.Contains([]string{"nina.dahl@ebbahus.se", "gustav.ebbe@ebbahus.se"}, a) {
			t.Errorf("%q is not an ebbahus address", a)
		}
	}
}

// Approval binds to the artifact. Re-wording a draft must produce a different Proposal, so
// that an approval given to the old words cannot send the new ones.
func TestAProposalIsIdentifiedByItsContent(t *testing.T) {
	a := state.NewProposal("ebbahus", "export-off", "email", "Subject", "Body", []string{"nina.dahl@ebbahus.se"}, "now")
	same := state.NewProposal("ebbahus", "export-off", "email", "Subject", "Body", []string{"nina.dahl@ebbahus.se"}, "later")
	reworded := state.NewProposal("ebbahus", "export-off", "email", "Subject", "Body, revised", []string{"nina.dahl@ebbahus.se"}, "now")

	if a.Hash != same.Hash {
		t.Error("the same artifact from a later run should be the same proposal")
	}
	if a.Hash == reworded.Hash {
		t.Error("a reworded draft must not inherit an approval given to the old wording")
	}
}

// Zammad stores display-name addresses with JSON escapes. Reading them without decoding
// derives a recipient that does not exist on a domain that does — and the domain filter
// waves it through, because the domain is right.
func TestRecipientsAreNotDerivedFromJSONEscapes(t *testing.T) {
	dir, cust := ebbahus(t)
	to, err := Recipients(dir, cust, output.Finding{Evidence: []output.Evidence{
		{Kind: "record", Path: "zammad/tickets/10010.json"},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(to, "gustav.ebbe@ebbahus.se") {
		t.Errorf("the ticket's requester should be a recipient, got %v", to)
	}
	for _, a := range to {
		if strings.HasPrefix(a, "u003") {
			t.Errorf("%q is a JSON escape read as an address", a)
		}
	}
}

// The subject is the Agent's own words in a header position. A newline in it would forge a
// Bcc past a human who approved a subject line and a body.
func TestSendRefusesAHeaderForgedInTheSubject(t *testing.T) {
	p := state.NewProposal("ebbahus", "export-off", "email",
		"September disbursements\r\nBcc: someone@elsewhere.example", "Body",
		[]string{"nina.dahl@ebbahus.se"}, "now")

	// The address is unreachable on purpose: the check must fail before any connection.
	err := Send("127.0.0.1:0", "robin@minso.se", p)
	if err == nil || !strings.Contains(err.Error(), "line break") {
		t.Fatalf("a subject with CRLF must be refused, got %v", err)
	}
}

func TestSendRefusesARecipientThatIsNotAnAddress(t *testing.T) {
	p := state.NewProposal("ebbahus", "export-off", "email", "Subject", "Body",
		[]string{"nina.dahl@ebbahus.se, root@localhost\r\nBcc: x@y.example"}, "now")

	if err := Send("127.0.0.1:0", "robin@minso.se", p); err == nil {
		t.Fatal("a recipient that is not a single address must be refused")
	}
}
