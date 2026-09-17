// Package linear is the client for the Linear workspace a Run maintains. Against the kit it
// talks to cmd/mockapi; against the real thing it is the same GraphQL.
//
// Nothing here decides anything. What to write is settled in internal/publish, which is
// deterministic Go: no model has an opinion about an issue's identity.
package linear

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Issue struct {
	ID          string `json:"id"`
	Identifier  string `json:"identifier"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	State       string `json:"state"`
	ProjectID   string `json:"projectId"`
	UpdatedAt   string `json:"updatedAt"`
}

type Comment struct {
	ID        string `json:"id"`
	IssueID   string `json:"issueId"`
	Body      string `json:"body"`
	CreatedAt string `json:"createdAt"`
}

type Client struct {
	URL   string
	Token string
	HTTP  *http.Client
}

func New(url, token string) *Client {
	return &Client{URL: url, Token: token, HTTP: &http.Client{Timeout: 15 * time.Second}}
}

func (c *Client) do(ctx context.Context, query string, vars map[string]any, out any) error {
	body, err := json.Marshal(map[string]any{"query": query, "variables": vars})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.URL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("linear: HTTP %s", resp.Status)
	}

	var envelope struct {
		Data   json.RawMessage `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return err
	}
	// GraphQL reports failures inside a 200, so the status code is never the whole answer.
	if len(envelope.Errors) > 0 {
		return fmt.Errorf("linear: %s", envelope.Errors[0].Message)
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(envelope.Data, out)
}

// Issues returns every issue in the workspace, following the cursor to the end. A writer
// that reads only the first page believes a human deleted everything after it.
func (c *Client) Issues(ctx context.Context) ([]Issue, error) {
	var all []Issue
	cursor := ""
	for {
		var page struct {
			Issues struct {
				Nodes    []Issue `json:"nodes"`
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
			} `json:"issues"`
		}
		vars := map[string]any{"first": 100}
		if cursor != "" {
			vars["after"] = cursor
		}
		if err := c.do(ctx, `query issues($first: Int, $after: String) { issues(first: $first, after: $after) { nodes { id identifier title description priority state projectId updatedAt } pageInfo { hasNextPage endCursor } } }`, vars, &page); err != nil {
			return nil, err
		}
		all = append(all, page.Issues.Nodes...)
		if !page.Issues.PageInfo.HasNextPage {
			return all, nil
		}
		cursor = page.Issues.PageInfo.EndCursor
	}
}

// InProject filters to one Customer's project. Reconciliation is project-scoped so that an
// issue a human created beside ours is seen and named, not invisible until we duplicate it.
func InProject(issues []Issue, projectID string) []Issue {
	var out []Issue
	for _, i := range issues {
		if i.ProjectID == projectID {
			out = append(out, i)
		}
	}
	return out
}

type IssueInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	ProjectID   string `json:"projectId"`
}

func (c *Client) CreateIssue(ctx context.Context, in IssueInput) (Issue, error) {
	var out struct {
		IssueCreate struct {
			Issue Issue `json:"issue"`
		} `json:"issueCreate"`
	}
	err := c.do(ctx, `mutation issueCreate($input: IssueCreateInput!) { issueCreate(input: $input) { success issue { id identifier title description priority state projectId updatedAt } } }`,
		map[string]any{"input": in}, &out)
	return out.IssueCreate.Issue, err
}

// UpdateIssue patches named fields only. A whole-object write would clobber whatever a human
// changed in a field this Run has nothing to say about.
func (c *Client) UpdateIssue(ctx context.Context, id string, fields map[string]any) (Issue, error) {
	var out struct {
		IssueUpdate struct {
			Issue Issue `json:"issue"`
		} `json:"issueUpdate"`
	}
	err := c.do(ctx, `mutation issueUpdate($id: String!, $input: IssueUpdateInput!) { issueUpdate(id: $id, input: $input) { success issue { id identifier title description priority state projectId updatedAt } } }`,
		map[string]any{"id": id, "input": fields}, &out)
	return out.IssueUpdate.Issue, err
}

// CreateComment is how new information reaches an Output we no longer own the fields of. We
// never overwrite a human's edit; we say the new thing beside it.
func (c *Client) CreateComment(ctx context.Context, issueID, body string) (Comment, error) {
	var out struct {
		CommentCreate struct {
			Comment Comment `json:"comment"`
		} `json:"commentCreate"`
	}
	err := c.do(ctx, `mutation commentCreate($input: CommentCreateInput!) { commentCreate(input: $input) { success comment { id issueId body createdAt } } }`,
		map[string]any{"input": map[string]any{"issueId": issueID, "body": body}}, &out)
	return out.CommentCreate.Comment, err
}

func (c *Client) Comments(ctx context.Context) ([]Comment, error) {
	var out struct {
		Comments struct {
			Nodes []Comment `json:"nodes"`
		} `json:"comments"`
	}
	err := c.do(ctx, `query comments($first: Int) { comments(first: $first) { nodes { id issueId body createdAt } pageInfo { hasNextPage endCursor } } }`,
		map[string]any{"first": 100}, &out)
	return out.Comments.Nodes, err
}
