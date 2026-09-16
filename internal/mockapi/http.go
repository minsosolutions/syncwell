package mockapi

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// requireBearer accepts any non-empty bearer token. Real credentials are out of scope for
// the kit, but collectors still have to send an Authorization header, which is the habit
// worth building.
func requireBearer(w http.ResponseWriter, r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") || strings.TrimSpace(auth[7:]) == "" {
		writeJSON(w, http.StatusUnauthorized, obj{"error": "missing or malformed Authorization: Bearer <token> header"})
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

// cursorPage slices items Slack-style: an opaque cursor, a page size, and a next_cursor
// that is empty on the last page. Defaults are deliberately small so a collector that
// ignores pagination silently collects only the first page.
func cursorPage(items []obj, cursor string, limit, defLimit, maxLimit int) (page []obj, next string) {
	if limit <= 0 {
		limit = defLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	start := 0
	if cursor != "" {
		if b, err := base64.StdEncoding.DecodeString(cursor); err == nil {
			if n, err := strconv.Atoi(string(b)); err == nil {
				start = n
			}
		}
	}
	if start > len(items) {
		start = len(items)
	}
	end := start + limit
	if end >= len(items) {
		return items[start:], ""
	}
	return items[start:end], base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(end)))
}

// offsetPage slices items Zammad-style: ?page=1&per_page=10, one-based.
func offsetPage(items []obj, page, perPage, defPerPage, maxPerPage int) []obj {
	if perPage <= 0 {
		perPage = defPerPage
	}
	if perPage > maxPerPage {
		perPage = maxPerPage
	}
	if page <= 0 {
		page = 1
	}
	start := (page - 1) * perPage
	if start >= len(items) {
		return []obj{}
	}
	end := start + perPage
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}

func intParam(r *http.Request, name string) int {
	n, _ := strconv.Atoi(r.URL.Query().Get(name))
	return n
}
