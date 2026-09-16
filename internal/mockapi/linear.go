package mockapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// The Linear endpoint speaks GraphQL's request and response envelope but does not parse
// the query: it dispatches on the operation keyword and reads arguments from `variables`.
// ponytail: substring dispatch, not a parser. Enough for issueCreate/issueUpdate/
// issueDelete/issues/projects; a real GraphQL client library would need the real thing.
type gqlRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

func gqlError(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusOK, obj{"errors": []obj{{"message": msg}}})
}

const (
	linearIssues   = "linear/issues.json"
	linearProjects = "linear/projects.json"
)

func (s *Server) linearGraphQL(w http.ResponseWriter, r *http.Request) {
	var req gqlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		gqlError(w, "could not parse request body")
		return
	}
	q := req.Query
	switch {
	case strings.Contains(q, "issueCreate"):
		s.linearIssueCreate(w, req)
	case strings.Contains(q, "issueUpdate"):
		s.linearIssueUpdate(w, req)
	case strings.Contains(q, "issueDelete"), strings.Contains(q, "issueArchive"):
		s.linearIssueDelete(w, req)
	case strings.Contains(q, "projectCreate"):
		s.linearProjectCreate(w, req)
	case strings.Contains(q, "projects"):
		s.linearList(w, req, linearProjects, "projects")
	case strings.Contains(q, "issues"), strings.Contains(q, "issue"):
		s.linearList(w, req, linearIssues, "issues")
	default:
		gqlError(w, "unsupported operation; this mock understands issues, projects, issueCreate, issueUpdate, issueDelete, projectCreate")
	}
}

func (s *Server) linearList(w http.ResponseWriter, req gqlRequest, rel, field string) {
	items, err := s.store.readList(rel)
	if err != nil {
		gqlError(w, err.Error())
		return
	}
	cursor, _ := req.Variables["after"].(string)
	limit := 0
	if f, ok := req.Variables["first"].(float64); ok {
		limit = int(f)
	}
	page, next := cursorPage(items, cursor, limit, 25, 100)
	writeJSON(w, http.StatusOK, obj{"data": obj{field: obj{
		"nodes":    page,
		"pageInfo": obj{"hasNextPage": next != "", "endCursor": next},
	}}})
}

// nextIdentifier returns the next SYN-n, scanning existing issues.
func nextIdentifier(issues []obj) (string, int) {
	max := 0
	for _, i := range issues {
		id, _ := i["identifier"].(string)
		if n, err := strconv.Atoi(strings.TrimPrefix(id, "SYN-")); err == nil && n > max {
			max = n
		}
	}
	return fmt.Sprintf("SYN-%d", max+1), max + 1
}

func (s *Server) linearIssueCreate(w http.ResponseWriter, req gqlRequest) {
	input, _ := req.Variables["input"].(map[string]any)
	if input == nil {
		gqlError(w, "variables.input is required")
		return
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	issues, err := s.store.readList(linearIssues)
	if err != nil {
		gqlError(w, err.Error())
		return
	}
	ident, n := nextIdentifier(issues)
	now := time.Now().UTC().Format(time.RFC3339)
	issue := obj{
		"id":          fmt.Sprintf("iss_%04d", n),
		"identifier":  ident,
		"title":       input["title"],
		"description": input["description"],
		"priority":    input["priority"],
		"state":       "Todo",
		"projectId":   input["projectId"],
		"createdAt":   now,
		"updatedAt":   now,
	}
	if err := s.store.writeList(linearIssues, append(issues, issue)); err != nil {
		gqlError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, obj{"data": obj{"issueCreate": obj{"success": true, "issue": issue}}})
}

func (s *Server) linearIssueUpdate(w http.ResponseWriter, req gqlRequest) {
	id, _ := req.Variables["id"].(string)
	input, _ := req.Variables["input"].(map[string]any)
	if id == "" || input == nil {
		gqlError(w, "variables.id and variables.input are required")
		return
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	issues, err := s.store.readList(linearIssues)
	if err != nil {
		gqlError(w, err.Error())
		return
	}
	for _, i := range issues {
		if i["id"] != id && i["identifier"] != id {
			continue
		}
		for k, v := range input {
			i[k] = v
		}
		i["updatedAt"] = time.Now().UTC().Format(time.RFC3339)
		if err := s.store.writeList(linearIssues, issues); err != nil {
			gqlError(w, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, obj{"data": obj{"issueUpdate": obj{"success": true, "issue": i}}})
		return
	}
	gqlError(w, "Entity not found: Issue "+id)
}

func (s *Server) linearIssueDelete(w http.ResponseWriter, req gqlRequest) {
	id, _ := req.Variables["id"].(string)
	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	issues, err := s.store.readList(linearIssues)
	if err != nil {
		gqlError(w, err.Error())
		return
	}
	kept := make([]obj, 0, len(issues))
	found := false
	for _, i := range issues {
		if i["id"] == id || i["identifier"] == id {
			found = true
			continue
		}
		kept = append(kept, i)
	}
	if !found {
		gqlError(w, "Entity not found: Issue "+id)
		return
	}
	if err := s.store.writeList(linearIssues, kept); err != nil {
		gqlError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, obj{"data": obj{"issueDelete": obj{"success": true}}})
}

func (s *Server) linearProjectCreate(w http.ResponseWriter, req gqlRequest) {
	input, _ := req.Variables["input"].(map[string]any)
	if input == nil {
		gqlError(w, "variables.input is required")
		return
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	projects, err := s.store.readList(linearProjects)
	if err != nil {
		gqlError(w, err.Error())
		return
	}
	now := time.Now().UTC().Format(time.RFC3339)
	project := obj{
		"id":          fmt.Sprintf("prj_%04d", len(projects)+1),
		"name":        input["name"],
		"description": input["description"],
		"state":       "started",
		"targetDate":  input["targetDate"],
		"createdAt":   now,
		"updatedAt":   now,
	}
	if err := s.store.writeList(linearProjects, append(projects, project)); err != nil {
		gqlError(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, obj{"data": obj{"projectCreate": obj{"success": true, "project": project}}})
}
