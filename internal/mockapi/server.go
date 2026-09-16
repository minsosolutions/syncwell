package mockapi

import (
	"log"
	"net/http"
	"time"
)

type Server struct {
	store *Store
	mux   *http.ServeMux
}

func New(dataDir string) *Server {
	s := &Server{store: NewStore(dataDir), mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) routes() {
	open := s.mux.HandleFunc
	auth := func(pattern string, h http.HandlerFunc) {
		s.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			if !requireBearer(w, r) {
				return
			}
			h(w, r)
		})
	}

	open("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, obj{"ok": true})
	})
	open("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, obj{"service": "syncwell mock api", "routes": routeIndex})
	})

	auth("GET /slack/api/conversations.list", s.slackConversationsList)
	auth("GET /slack/api/conversations.history", s.slackConversationsHistory)
	auth("GET /slack/api/users.list", s.slackUsersList)
	auth("POST /slack/api/chat.postMessage", s.slackPostMessage)

	auth("GET /zammad/api/v1/tickets", s.zammadTicketList)
	auth("GET /zammad/api/v1/tickets/{id}", func(w http.ResponseWriter, r *http.Request) {
		s.zammadTicketByID(w, r.PathValue("id"))
	})
	auth("GET /zammad/api/v1/ticket_articles/by_ticket/{id}", func(w http.ResponseWriter, r *http.Request) {
		s.zammadArticlesByTicket(w, r.PathValue("id"))
	})
	auth("GET /zammad/api/v1/organizations", s.zammadCollection("sources/zammad/organizations.json"))
	auth("GET /zammad/api/v1/users", s.zammadCollection("sources/zammad/users.json"))
	auth("GET /zammad/api/v1/groups", s.zammadCollection("sources/zammad/groups.json"))

	auth("POST /linear/graphql", s.linearGraphQL)
}

var routeIndex = []string{
	"GET  /slack/api/conversations.list?cursor=&limit=",
	"GET  /slack/api/conversations.history?channel=&cursor=&limit=",
	"GET  /slack/api/users.list",
	"POST /slack/api/chat.postMessage",
	"GET  /zammad/api/v1/tickets?page=&per_page=",
	"GET  /zammad/api/v1/tickets/{id}",
	"GET  /zammad/api/v1/ticket_articles/by_ticket/{id}",
	"GET  /zammad/api/v1/organizations",
	"GET  /zammad/api/v1/users",
	"GET  /zammad/api/v1/groups",
	"POST /linear/graphql",
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	s.mux.ServeHTTP(w, r)
	log.Printf("%s %s%s (%s)", r.Method, r.URL.Path, querySuffix(r), time.Since(start).Round(time.Millisecond))
}

func querySuffix(r *http.Request) string {
	if r.URL.RawQuery == "" {
		return ""
	}
	return "?" + r.URL.RawQuery
}
