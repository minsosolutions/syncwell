// Package routing decides which Customer a Source Record belongs to.
//
// Every rule here is deterministic code over stable Source metadata. Nothing in this
// package may consult a model: a Routing Rule that guesses is a leak. See
// docs/adr/0001-isolation-lives-in-the-routing-rule.md.
package routing

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Customer struct {
	Slug               string   `json:"slug"`
	Name               string   `json:"name"`
	Kind               string   `json:"kind"`
	Domains            []string `json:"domains"`
	SlackChannelPrefix string   `json:"slack_channel_prefix"`
	ZammadOrganization string   `json:"zammad_organization"`
}

type Config struct {
	Vendor struct {
		Name    string `json:"name"`
		Domain  string `json:"domain"`
		Product string `json:"product"`
	} `json:"vendor"`
	Customers []Customer `json:"customers"`
}

func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// Unroutable explains why a Source Record could not be attributed to exactly one Customer.
// It is a value, not a failure: quarantining is a correct outcome.
type Unroutable struct{ Reason string }

func (u *Unroutable) Error() string { return u.Reason }

// BySlackChannel routes on the channel naming convention. A channel belonging to no
// Customer — or, impossibly but checked anyway, to two — is quarantined rather than guessed.
func (c *Config) BySlackChannel(channel string) (*Customer, error) {
	var found *Customer
	for i := range c.Customers {
		if strings.HasPrefix(channel, c.Customers[i].SlackChannelPrefix) {
			if found != nil {
				return nil, &Unroutable{fmt.Sprintf("channel %q matches both %s and %s", channel, found.Slug, c.Customers[i].Slug)}
			}
			found = &c.Customers[i]
		}
	}
	if found == nil {
		return nil, &Unroutable{fmt.Sprintf("channel %q has no customer prefix", channel)}
	}
	return found, nil
}

// ByAttendeeDomains routes a meeting on who was in the room. The Vendor's own domain is
// stripped — it is in every meeting and says nothing — and exactly one Customer must own
// what remains. Nobody left (an internal meeting about a Customer), two Customers, or an
// attendee whose domain belongs to no Customer all quarantine: there is real signal in those
// meetings and no deterministic rule that safely claims it.
func (c *Config) ByAttendeeDomains(attendees []string) (*Customer, error) {
	var found *Customer
	for _, a := range attendees {
		at := strings.LastIndex(a, "@")
		if at < 0 {
			return nil, &Unroutable{fmt.Sprintf("attendee %q has no domain", a)}
		}
		domain := strings.ToLower(strings.TrimSpace(a[at+1:]))
		if domain == strings.ToLower(c.Vendor.Domain) {
			continue
		}
		cust := c.byDomain(domain)
		if cust == nil {
			return nil, &Unroutable{fmt.Sprintf("attendee domain %q belongs to no customer", domain)}
		}
		if found != nil && found.Slug != cust.Slug {
			return nil, &Unroutable{fmt.Sprintf("attendees span %s and %s", found.Slug, cust.Slug)}
		}
		found = cust
	}
	if found == nil {
		return nil, &Unroutable{"no customer attendee; vendor-internal meeting"}
	}
	return found, nil
}

func (c *Config) byDomain(domain string) *Customer {
	for i := range c.Customers {
		for _, d := range c.Customers[i].Domains {
			if strings.EqualFold(d, domain) {
				return &c.Customers[i]
			}
		}
	}
	return nil
}
