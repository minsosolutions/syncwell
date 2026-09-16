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
