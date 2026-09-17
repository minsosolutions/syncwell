package collect

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/minsosolutions/syncwell/internal/routing"
)

// Zammad's ceiling is 100 per page. Asking for it is the difference between 4 requests and 31.
const zammadPerPage = 100

// pages walks an offset-paginated Zammad collection to exhaustion. Zammad returns bare JSON
// arrays with no envelope and no total, so the end of the data is a short page — there is
// nothing else to detect it by. Stopping after the first page is how you silently collect a
// third of the desk.
func (c *Client) pages(path string) ([]json.RawMessage, error) {
	var all []json.RawMessage
	for page := 1; ; page++ {
		var batch []json.RawMessage
		q := url.Values{"page": {strconv.Itoa(page)}, "per_page": {strconv.Itoa(zammadPerPage)}}
		if err := c.get(path, q, &batch); err != nil {
			return nil, err
		}
		all = append(all, batch...)
		if len(batch) < zammadPerPage {
			return all, nil
		}
	}
}

type zammadOrg struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type zammadUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// zammadTicket is the routing key and nothing more. The ticket is written to disk as the bytes
// Zammad gave us, so every other field — state, articles, timestamps — survives untouched.
type zammadTicket struct {
	ID             int  `json:"id"`
	OrganizationID *int `json:"organization_id"`
	CustomerID     *int `json:"customer_id"`
}

// Zammad collects the support desk. The list endpoint omits articles, exactly as the real one
// does, so each ticket is fetched again in full — the title says almost nothing and the
// conversation is where the content is. Routing is by the desk's own organization, falling
// back to the requester's email domain when the desk never assigned one.
func Zammad(c *Client, cfg *routing.Config, dataDir string) (Report, error) {
	var rep Report

	orgs, err := lookup[zammadOrg](c, "/zammad/api/v1/organizations", func(o zammadOrg) (int, string) { return o.ID, o.Name })
	if err != nil {
		return rep, err
	}
	users, err := lookup[zammadUser](c, "/zammad/api/v1/users", func(u zammadUser) (int, string) { return u.ID, u.Email })
	if err != nil {
		return rep, err
	}

	index, err := c.pages("/zammad/api/v1/tickets")
	if err != nil {
		return rep, err
	}

	for _, raw := range index {
		var shallow zammadTicket
		if err := json.Unmarshal(raw, &shallow); err != nil {
			return rep, err
		}

		// Refetch: the index deliberately has no articles.
		var full json.RawMessage
		if err := c.get("/zammad/api/v1/tickets/"+strconv.Itoa(shallow.ID), nil, &full); err != nil {
			return rep, err
		}
		var ticket zammadTicket
		if err := json.Unmarshal(full, &ticket); err != nil {
			return rep, err
		}

		name, err := segment("ticket id", strconv.Itoa(ticket.ID))
		if err != nil {
			return rep, err
		}

		var org, requester string
		if ticket.OrganizationID != nil {
			org = orgs[*ticket.OrganizationID]
		}
		if ticket.CustomerID != nil {
			requester = users[*ticket.CustomerID]
		}

		cust, routeErr := cfg.ByOrganizationThenDomain(org, requester)
		if routeErr != nil {
			var u *routing.Unroutable
			if !errors.As(routeErr, &u) {
				return rep, routeErr
			}
			if err := quarantineTicket(dataDir, name, full, u.Reason); err != nil {
				return rep, err
			}
			rep.Quarantined = append(rep.Quarantined, fmt.Sprintf("ticket %s: %s", name, u.Reason))
			continue
		}

		slug, err := segment("customer slug", cust.Slug)
		if err != nil {
			return rep, err
		}
		if err := writeFile(filepath.Join(dataDir, "customers", slug, "zammad", "tickets", name+".json"), full); err != nil {
			return rep, err
		}
		rep.Routed = append(rep.Routed, fmt.Sprintf("ticket %s -> %s", name, cust.Slug))
	}
	return rep, nil
}

// lookup walks a paginated collection into an id-keyed map.
func lookup[T any](c *Client, path string, key func(T) (int, string)) (map[int]string, error) {
	raw, err := c.pages(path)
	if err != nil {
		return nil, err
	}
	out := make(map[int]string, len(raw))
	for _, r := range raw {
		var v T
		if err := json.Unmarshal(r, &v); err != nil {
			return nil, err
		}
		id, name := key(v)
		out[id] = name
	}
	return out, nil
}

// quarantineTicket keeps the ticket as Zammad returned it and states why it is here beside it.
func quarantineTicket(dataDir, name string, body []byte, reason string) error {
	dir := filepath.Join(dataDir, "unrouted", "zammad", "tickets")
	if err := writeFile(filepath.Join(dir, name+".json"), body); err != nil {
		return err
	}
	return writeFile(filepath.Join(dir, name+".reason.txt"), []byte(reason+"\n"))
}
