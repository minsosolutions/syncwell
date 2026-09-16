package mockapi

import (
	"fmt"
	"net/http"
	"sort"
)

// Zammad's own maximum is 100 per page; the small default is ours, for the same reason as
// Slack's. A collector that never sends ?page= sees only the first ten tickets.
const (
	zammadDefaultPerPage = 10
	zammadMaxPerPage     = 100
)

func (s *Server) zammadTickets() ([]obj, error) {
	tickets, err := s.store.readObjectDir("sources/zammad/tickets")
	if err != nil {
		return nil, err
	}
	sort.SliceStable(tickets, func(i, j int) bool {
		return fmt.Sprint(tickets[i]["id"]) < fmt.Sprint(tickets[j]["id"])
	})
	return tickets, nil
}

func (s *Server) zammadTicketList(w http.ResponseWriter, r *http.Request) {
	tickets, err := s.zammadTickets()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, obj{"error": err.Error()})
		return
	}
	page := offsetPage(tickets, intParam(r, "page"), intParam(r, "per_page"), zammadDefaultPerPage, zammadMaxPerPage)
	// Zammad's index action does not embed articles; they are a separate call. Keeping
	// that split means collectors have to make it too.
	out := make([]obj, 0, len(page))
	for _, t := range page {
		shallow := obj{}
		for k, v := range t {
			if k != "articles" {
				shallow[k] = v
			}
		}
		out = append(out, shallow)
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) zammadTicketByID(w http.ResponseWriter, id string) {
	tickets, err := s.zammadTickets()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, obj{"error": err.Error()})
		return
	}
	for _, t := range tickets {
		if fmt.Sprint(t["id"]) == id {
			writeJSON(w, http.StatusOK, t)
			return
		}
	}
	writeJSON(w, http.StatusNotFound, obj{"error": "No such Ticket"})
}

func (s *Server) zammadArticlesByTicket(w http.ResponseWriter, id string) {
	tickets, err := s.zammadTickets()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, obj{"error": err.Error()})
		return
	}
	for _, t := range tickets {
		if fmt.Sprint(t["id"]) != id {
			continue
		}
		raw, _ := t["articles"].([]any)
		out := make([]obj, 0, len(raw))
		for _, a := range raw {
			if o, ok := a.(map[string]any); ok {
				out = append(out, o)
			}
		}
		writeJSON(w, http.StatusOK, out)
		return
	}
	writeJSON(w, http.StatusNotFound, obj{"error": "No such Ticket"})
}

func (s *Server) zammadCollection(rel string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := s.store.readList(rel)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, obj{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, offsetPage(items, intParam(r, "page"), intParam(r, "per_page"), zammadMaxPerPage, zammadMaxPerPage))
	}
}
