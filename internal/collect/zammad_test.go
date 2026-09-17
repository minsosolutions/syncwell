package collect

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/minsosolutions/syncwell/internal/mockapi"
	"github.com/minsosolutions/syncwell/internal/routing"
)

// Run against the real mock over the real seeded data: refetching for articles, the null
// organization and ticket state are all properties of that data, not of a stub.
func TestZammadCollectsEveryTicket(t *testing.T) {
	srcDir := filepath.Join("..", "..", "data")
	srv := httptest.NewServer(mockapi.New(srcDir))
	defer srv.Close()

	cfg, err := routing.Load(filepath.Join("..", "..", "config", "syncwell.json"))
	if err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	client := &Client{BaseURL: srv.URL, Token: "test", HTTP: &http.Client{Timeout: 10 * time.Second}}

	rep, err := Zammad(client, cfg, out)
	if err != nil {
		t.Fatal(err)
	}

	// Independent source of truth: the tickets actually on disk in the backing store.
	onDisk, err := filepath.Glob(filepath.Join(srcDir, "sources", "zammad", "tickets", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(onDisk) < 20 {
		t.Fatalf("only %d tickets seeded; this test is not exercising pagination", len(onDisk))
	}
	if got := len(rep.Routed) + len(rep.Quarantined); got != len(onDisk) {
		t.Fatalf("accounted for %d tickets, backing store has %d", got, len(onDisk))
	}

	// Every ticket landed exactly once, somewhere.
	collected, err := filepath.Glob(filepath.Join(out, "customers", "*", "zammad", "tickets", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	quarantined, err := filepath.Glob(filepath.Join(out, "unrouted", "zammad", "tickets", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(collected)+len(quarantined) != len(onDisk) {
		t.Fatalf("%d collected + %d quarantined != %d on disk", len(collected), len(quarantined), len(onDisk))
	}

	// Q2: no organization, personal address. Named in the ticket text, invisible to the rule.
	q2 := filepath.Join(out, "unrouted", "zammad", "tickets", "10030.json")
	if _, err := os.Stat(q2); err != nil {
		t.Fatalf("Q2 was not quarantined: %v", err)
	}
	reason, err := os.ReadFile(filepath.Join(out, "unrouted", "zammad", "tickets", "10030.reason.txt"))
	if err != nil {
		t.Fatalf("Q2 quarantined without a reason: %v", err)
	}
	if !strings.Contains(string(reason), "hotmail.com") {
		t.Fatalf("reason %q does not name what a human must act on", reason)
	}

	// N2 turns on this: a ticket marked solved while another source says it still happens.
	// If state does not survive collection, the finding cannot exist.
	var ticket struct {
		State    string            `json:"state"`
		Articles []json.RawMessage `json:"articles"`
	}
	b, err := os.ReadFile(filepath.Join(out, "customers", "nordstad", "zammad", "tickets", "10042.json"))
	if err != nil {
		t.Fatalf("ticket 10042: %v", err)
	}
	if err := json.Unmarshal(b, &ticket); err != nil {
		t.Fatal(err)
	}
	if ticket.State != "solved" {
		t.Fatalf("state is %q, want solved — N2 needs it", ticket.State)
	}
	// The list endpoint omits articles; a collector that only reads the list writes empty
	// tickets and every finding loses its evidence.
	if len(ticket.Articles) < 2 {
		t.Fatalf("ticket 10042 collected with %d articles; the conversation is the content", len(ticket.Articles))
	}
}

// The seeded desk fits in one page at Zammad's ceiling, so this is what proves the walk
// continues past page one — the failure mode is silent, and a collector that stops early looks
// exactly like a quiet customer.
func TestZammadPagesWalkToExhaustion(t *testing.T) {
	const total = 250
	var requested []int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
		if perPage <= 0 || perPage > 100 {
			perPage = 100 // Zammad's real ceiling: asking for more does not get you more.
		}
		requested = append(requested, page)
		out := []map[string]int{}
		for i := (page - 1) * perPage; i < page*perPage && i < total; i++ {
			out = append(out, map[string]int{"id": i})
		}
		json.NewEncoder(w).Encode(out)
	}))
	defer srv.Close()

	c := &Client{BaseURL: srv.URL, Token: "test", HTTP: &http.Client{Timeout: 10 * time.Second}}
	got, err := c.pages("/zammad/api/v1/tickets")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != total {
		t.Fatalf("walked %d items, want %d — the walk stopped early", len(got), total)
	}
	if len(requested) < 3 {
		t.Fatalf("made %d requests for %d items; it is not paging", len(requested), total)
	}
	// Every item exactly once: an off-by-one in the page offset duplicates or skips a ticket.
	seen := map[int]bool{}
	for _, raw := range got {
		var item struct{ ID int }
		if err := json.Unmarshal(raw, &item); err != nil {
			t.Fatal(err)
		}
		if seen[item.ID] {
			t.Fatalf("item %d collected twice", item.ID)
		}
		seen[item.ID] = true
	}
}
