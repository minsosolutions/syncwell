package mockapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
)

// Slack's real defaults are generous (limit=100). These are deliberately small so that a
// collector which ignores response_metadata.next_cursor visibly loses messages — the
// failure mode worth meeting in a workshop rather than in production.
const (
	slackDefaultLimit = 20
	slackMaxLimit     = 200
)

func slackErr(w http.ResponseWriter, code string) {
	writeJSON(w, http.StatusOK, obj{"ok": false, "error": code})
}

func (s *Server) slackConversationsList(w http.ResponseWriter, r *http.Request) {
	channels, err := s.store.readList("sources/slack/channels.json")
	if err != nil {
		slackErr(w, "internal_error")
		return
	}
	page, next := cursorPage(channels, r.URL.Query().Get("cursor"), intParam(r, "limit"), slackDefaultLimit, slackMaxLimit)
	writeJSON(w, http.StatusOK, obj{
		"ok":                true,
		"channels":          page,
		"response_metadata": obj{"next_cursor": next},
	})
}

func (s *Server) slackUsersList(w http.ResponseWriter, r *http.Request) {
	users, err := s.store.readList("sources/slack/users.json")
	if err != nil {
		slackErr(w, "internal_error")
		return
	}
	page, next := cursorPage(users, r.URL.Query().Get("cursor"), intParam(r, "limit"), slackDefaultLimit, slackMaxLimit)
	writeJSON(w, http.StatusOK, obj{
		"ok":                true,
		"members":           page,
		"response_metadata": obj{"next_cursor": next},
	})
}

// channelName resolves a Slack channel id (C...) to the directory name holding its messages.
func (s *Server) channelName(id string) (string, bool) {
	channels, err := s.store.readList("sources/slack/channels.json")
	if err != nil {
		return "", false
	}
	for _, c := range channels {
		if c["id"] == id {
			name, _ := c["name"].(string)
			return name, name != ""
		}
		// Accept a bare channel name too; it makes hand-written curl calls bearable.
		if c["name"] == id {
			return id, true
		}
	}
	return "", false
}

func (s *Server) slackConversationsHistory(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("channel")
	if id == "" {
		slackErr(w, "channel_not_found")
		return
	}
	name, ok := s.channelName(id)
	if !ok {
		slackErr(w, "channel_not_found")
		return
	}
	msgs, err := s.store.readGlobConcat("sources/slack/" + name)
	if err != nil {
		slackErr(w, "internal_error")
		return
	}
	// Slack returns newest first.
	sort.SliceStable(msgs, func(i, j int) bool {
		return fmt.Sprint(msgs[i]["ts"]) > fmt.Sprint(msgs[j]["ts"])
	})
	page, next := cursorPage(msgs, r.URL.Query().Get("cursor"), intParam(r, "limit"), slackDefaultLimit, slackMaxLimit)
	writeJSON(w, http.StatusOK, obj{
		"ok":                true,
		"messages":          page,
		"has_more":          next != "",
		"response_metadata": obj{"next_cursor": next},
	})
}

func (s *Server) slackPostMessage(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Channel string `json:"channel"`
		Text    string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Channel == "" {
		slackErr(w, "invalid_arguments")
		return
	}
	name, ok := s.channelName(body.Channel)
	if !ok {
		slackErr(w, "channel_not_found")
		return
	}

	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	now := time.Now()
	ts := fmt.Sprintf("%d.%06d", now.Unix(), now.Nanosecond()/1000)
	msg := obj{"type": "message", "user": "B07SYNCWELL", "text": body.Text, "ts": ts}

	rel := fmt.Sprintf("sources/slack/%s/%s.json", name, now.Format("2006-01-02"))
	day, err := s.store.readList(rel)
	if err != nil {
		slackErr(w, "internal_error")
		return
	}
	if err := s.store.writeList(rel, append(day, msg)); err != nil {
		slackErr(w, "internal_error")
		return
	}
	writeJSON(w, http.StatusOK, obj{"ok": true, "channel": body.Channel, "ts": ts, "message": msg})
}
