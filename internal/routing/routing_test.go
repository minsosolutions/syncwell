package routing

import "testing"

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
		{Slug: "nordstad", Domains: []string{"nordstad.se"}},
		{Slug: "ebbahus", Domains: []string{"ebbahus.se"}},
		{Slug: "utvecklingsverket", Domains: []string{"utvecklingsverket.se", "ruv.se"}},
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
