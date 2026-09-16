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
