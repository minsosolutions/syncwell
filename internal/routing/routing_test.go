package routing

import (
	"strings"
	"testing"
)

func testConfig() *Config {
	return &Config{Customers: []Customer{
		{Slug: "nordstad", SlackChannelPrefix: "cust-nordstad-"},
		{Slug: "ebbahus", SlackChannelPrefix: "cust-ebbahus-"},
	}}
}

// The property that matters: exactly one Customer, or none. Never a guess.
func TestBySlackChannel(t *testing.T) {
	c := testConfig()
	for _, tc := range []struct{ channel, want string }{
		{"cust-nordstad-support", "nordstad"},
		{"cust-ebbahus-support", "ebbahus"},
	} {
		got, err := c.BySlackChannel(tc.channel)
		if err != nil {
			t.Fatalf("%s: %v", tc.channel, err)
		}
		if got.Slug != tc.want {
			t.Fatalf("%s routed to %s, want %s", tc.channel, got.Slug, tc.want)
		}
	}

	// The internal channel must quarantine: it discusses several Customers at once.
	if _, err := c.BySlackChannel("syncwell-internt"); err == nil {
		t.Fatal("internal channel routed to a customer; that is a leak")
	}

	// A prefix collision must quarantine rather than pick the first match.
	c.Customers = append(c.Customers, Customer{Slug: "nordstad-vast", SlackChannelPrefix: "cust-nordstad-"})
	if _, err := c.BySlackChannel("cust-nordstad-support"); err == nil {
		t.Fatal("ambiguous channel routed to a single customer")
	}
}

func transcriptConfig() *Config {
	c := &Config{Customers: []Customer{
		{Slug: "nordstad", Domains: []string{"nordstad.se"}, ZammadOrganization: "Nordstad kommun"},
		{Slug: "ebbahus", Domains: []string{"ebbahus.se"}, ZammadOrganization: "Ebbahus Stiftelse"},
		{Slug: "utvecklingsverket", Domains: []string{"utvecklingsverket.se", "ruv.se"}, ZammadOrganization: "Regionalt Utvecklingsverket"},
	}}
	c.Vendor.Domain = "minso.se"
	return c
}

// Strip the Vendor's own domain; exactly one Customer must remain.
func TestByAttendeeDomains(t *testing.T) {
	c := transcriptConfig()

	got, err := c.ByAttendeeDomains([]string{"anna.lind@nordstad.se", "per.ek@nordstad.se", "robin@minso.se"})
	if err != nil {
		t.Fatalf("nordstad meeting: %v", err)
	}
	if got.Slug != "nordstad" {
		t.Fatalf("routed to %s, want nordstad", got.Slug)
	}

	// One Customer, two domains: an alias is not a second Customer.
	got, err = c.ByAttendeeDomains([]string{"hakan.nyberg@utvecklingsverket.se", "tobias.ahl@ruv.se", "robin@minso.se"})
	if err != nil {
		t.Fatalf("alias domain: %v", err)
	}
	if got.Slug != "utvecklingsverket" {
		t.Fatalf("routed to %s, want utvecklingsverket", got.Slug)
	}

	// An internal meeting about a Customer has no Customer in the room. Quarantine.
	if _, err := c.ByAttendeeDomains([]string{"robin@minso.se", "mia.berg@minso.se"}); err == nil {
		t.Fatal("vendor-only meeting routed to a customer")
	}

	// Two Customers in one room: attributing it to either is a leak.
	if _, err := c.ByAttendeeDomains([]string{"anna.lind@nordstad.se", "nina.dahl@ebbahus.se", "robin@minso.se"}); err == nil {
		t.Fatal("two customers routed to a single customer")
	}

	// An attendee we cannot place is not evidence the rest of the room is safe.
	if _, err := c.ByAttendeeDomains([]string{"anna.lind@nordstad.se", "consultant@okand.se"}); err == nil {
		t.Fatal("unknown attendee domain routed anyway")
	}

	// Malformed frontmatter must not be read as "no customer" and quietly quarantined as internal.
	if _, err := c.ByAttendeeDomains([]string{"anna.lind"}); err == nil {
		t.Fatal("address without a domain routed anyway")
	}
}

// Email routes on the same property as a meeting — which Customers are in the room — but the
// room is from/to/cc and the failure that matters is different: a personal address.
func TestByPeerDomains(t *testing.T) {
	c := transcriptConfig()

	got, err := c.ByPeerDomains([]string{"anna.lind@nordstad.se", "robin@minso.se", "per.ek@nordstad.se"})
	if err != nil {
		t.Fatalf("nordstad thread: %v", err)
	}
	if got.Slug != "nordstad" {
		t.Fatalf("routed to %s, want nordstad", got.Slug)
	}

	// Q1: a Customer contact writing from a personal address. We know who Gustav is; the rule
	// does not, and must not guess him into ebbahus.
	_, err = c.ByPeerDomains([]string{"g.ebbe.privat@gmail.com", "robin@minso.se"})
	if err == nil {
		t.Fatal("personal address routed to a customer")
	}
	if !strings.Contains(err.Error(), "gmail.com") {
		t.Fatalf("reason does not name the domain a human must act on: %v", err)
	}

	// cc is what pulls two Customers into one thread.
	if _, err := c.ByPeerDomains([]string{"anna.lind@nordstad.se", "robin@minso.se", "nina.dahl@ebbahus.se"}); err == nil {
		t.Fatal("cc across two customers routed to one; that is a leak")
	}

	// Vendor talking to itself about a Customer.
	if _, err := c.ByPeerDomains([]string{"robin@minso.se", "sara.holm@minso.se"}); err == nil {
		t.Fatal("internal mail routed to a customer")
	}
}

// Zammad's own organization is the right thing to route on when the desk has set it. The
// fallback exists because Zammad holds one domain per organization and a Customer may use two.
func TestByOrganizationThenDomain(t *testing.T) {
	c := transcriptConfig()

	got, err := c.ByOrganizationThenDomain("Nordstad kommun", "per.ek@nordstad.se")
	if err != nil {
		t.Fatalf("organization set: %v", err)
	}
	if got.Slug != "nordstad" {
		t.Fatalf("routed to %s, want nordstad", got.Slug)
	}

	// The organization wins even when it disagrees with the requester's domain: the desk
	// assigned it deliberately, a domain is a guess about a person.
	got, err = c.ByOrganizationThenDomain("Regionalt Utvecklingsverket", "tobias.ahl@ruv.se")
	if err != nil {
		t.Fatalf("second domain under one organization: %v", err)
	}
	if got.Slug != "utvecklingsverket" {
		t.Fatalf("routed to %s, want utvecklingsverket", got.Slug)
	}

	// No organization: fall back to the requester's domain.
	got, err = c.ByOrganizationThenDomain("", "elin.sund@utvecklingsverket.se")
	if err != nil {
		t.Fatalf("fallback: %v", err)
	}
	if got.Slug != "utvecklingsverket" {
		t.Fatalf("fallback routed to %s, want utvecklingsverket", got.Slug)
	}

	// Q2: no organization and a personal address. The text of that ticket names the customer;
	// the rule cannot see it and must not pretend otherwise.
	if _, err := c.ByOrganizationThenDomain("", "karin.sjo@hotmail.com"); err == nil {
		t.Fatal("personal address with no organization routed to a customer")
	}

	// An organization Zammad knows and the config does not is a configuration gap, not a guess
	// to be papered over by the domain fallback.
	if _, err := c.ByOrganizationThenDomain("Okänd Kommun", "someone@nordstad.se"); err == nil {
		t.Fatal("unknown organization fell through to the domain fallback")
	}

	if _, err := c.ByOrganizationThenDomain("", ""); err == nil {
		t.Fatal("ticket with neither organization nor requester routed anyway")
	}
}
