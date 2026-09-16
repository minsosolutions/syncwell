// Package collect holds the reference Collector.
//
// Only Slack is implemented. The other three Sources are documented in
// data/sources/<source>/README.md and left for you to write — that is the exercise.
package collect

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/minsosolutions/syncwell/internal/routing"
)

type Client struct {
	BaseURL string
	Token   string
	HTTP    *http.Client
}

func (c *Client) get(path string, q url.Values, out any) error {
	req, err := http.NewRequest("GET", c.BaseURL+path+"?"+q.Encode(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: %s", path, resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

type slackPage struct {
	OK               bool              `json:"ok"`
	Error            string            `json:"error"`
	Channels         []json.RawMessage `json:"channels"`
	Messages         []json.RawMessage `json:"messages"`
	ResponseMetadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"response_metadata"`
}

// paginate walks next_cursor to exhaustion. Stopping after the first page is the classic
// way to silently collect a third of the data.
func (c *Client) paginate(path string, q url.Values, pick func(slackPage) []json.RawMessage) ([]json.RawMessage, error) {
	var all []json.RawMessage
	cursor := ""
	for {
		page := slackPage{}
		qq := url.Values{}
		for k, v := range q {
			qq[k] = v
		}
		if cursor != "" {
			qq.Set("cursor", cursor)
		}
		if err := c.get(path, qq, &page); err != nil {
			return nil, err
		}
		if !page.OK {
			return nil, errors.New("slack: " + page.Error)
		}
		all = append(all, pick(page)...)
		if page.ResponseMetadata.NextCursor == "" {
			return all, nil
		}
		cursor = page.ResponseMetadata.NextCursor
	}
}

// safeSegment guards every externally-supplied string that becomes a path component.
// The Source is a trust boundary: a channel named "../../etc" must not escape dataDir.
var safeSegment = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{0,79}$`)

func segment(kind, name string) (string, error) {
	if !safeSegment.MatchString(name) || strings.Contains(name, "..") {
		return "", fmt.Errorf("unsafe %s name %q", kind, name)
	}
	return name, nil
}

type channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Slack collects every channel, routes each one to a Customer by its name, and writes the
// messages into that Customer Directory. A channel that does not route is quarantined
// whole, with the reason recorded next to it.
func Slack(c *Client, cfg *routing.Config, dataDir string) (Report, error) {
	var rep Report
	raw, err := c.paginate("/slack/api/conversations.list", nil, func(p slackPage) []json.RawMessage { return p.Channels })
	if err != nil {
		return rep, err
	}

	for _, r := range raw {
		var ch channel
		if err := json.Unmarshal(r, &ch); err != nil {
			return rep, err
		}
		msgs, err := c.paginate("/slack/api/conversations.history", url.Values{"channel": {ch.ID}},
			func(p slackPage) []json.RawMessage { return p.Messages })
		if err != nil {
			return rep, err
		}

		name, err := segment("channel", ch.Name)
		if err != nil {
			return rep, err
		}

		cust, routeErr := cfg.BySlackChannel(ch.Name)
		if routeErr != nil {
			var u *routing.Unroutable
			if !errors.As(routeErr, &u) {
				return rep, routeErr
			}
			path := filepath.Join(dataDir, "unrouted", "slack", name+".json")
			body := map[string]any{"source": "slack", "channel": ch.Name, "reason": u.Reason, "messages": msgs}
			if err := writeJSON(path, body); err != nil {
				return rep, err
			}
			rep.Quarantined = append(rep.Quarantined, fmt.Sprintf("%s (%d messages): %s", ch.Name, len(msgs), u.Reason))
			continue
		}

		slug, err := segment("customer slug", cust.Slug)
		if err != nil {
			return rep, err
		}
		path := filepath.Join(dataDir, "customers", slug, "slack", name+".json")
		if err := writeJSON(path, msgs); err != nil {
			return rep, err
		}
		rep.Routed = append(rep.Routed, fmt.Sprintf("%s -> %s (%d messages)", ch.Name, cust.Slug, len(msgs)))
	}
	return rep, nil
}

type Report struct {
	Routed      []string
	Quarantined []string
}

func writeJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}
