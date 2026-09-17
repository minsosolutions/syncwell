// Command admin is the human's half of Syncwell: the approval queue, the findings browser,
// and the button that starts a Run.
//
// net/http, html/template and a vendored htmx — one binary, no npm. A Run never waits for
// this UI; it leaves Proposals behind and this is where they are approved, applied through
// the same writer path the Run uses. See docs/adr/0003-approval-gates-on-blast-radius.md.
//
//	go run ./cmd/admin
package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/minsosolutions/syncwell/internal/publish"
	"github.com/minsosolutions/syncwell/internal/routing"
	"github.com/minsosolutions/syncwell/internal/run"
	"github.com/minsosolutions/syncwell/internal/state"
)

//go:embed templates/*.html assets/*
var assets embed.FS

func main() {
	addr := flag.String("addr", ":8100", "listen address")
	dataDir := flag.String("data", "data", "data directory")
	configPath := flag.String("config", "config/syncwell.json", "config file")
	flag.Parse()

	cfg, err := routing.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	s := &server{cfg: cfg, dataDir: *dataDir, runs: map[string]*runState{}}
	s.tmpl = template.Must(template.New("").Funcs(funcs).ParseFS(assets, "templates/*.html"))

	log.Printf("syncwell admin on http://localhost%s", *addr)
	if err := http.ListenAndServe(*addr, s.routes()); err != nil {
		log.Fatal(err)
	}
}

type server struct {
	cfg     *routing.Config
	dataDir string
	tmpl    *template.Template

	mu   sync.Mutex
	runs map[string]*runState
}

type runState struct {
	Started time.Time
	Done    bool
	Err     string
	Summary *run.Summary
}

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /assets/", http.FileServerFS(assets))
	mux.HandleFunc("GET /{$}", s.index)
	mux.HandleFunc("GET /c/{slug}", s.customer)
	mux.HandleFunc("GET /c/{slug}/record", s.record)
	mux.HandleFunc("POST /c/{slug}/run", s.startRun)
	mux.HandleFunc("GET /c/{slug}/run", s.runStatus)
	mux.HandleFunc("POST /c/{slug}/proposal/{hash}/{decision}", s.decide)
	return mux
}

func (s *server) customerDir(slug string) string {
	return filepath.Join(s.dataDir, "customers", slug)
}

// The index is the one view that names more than one Customer, and it names nothing about
// them beyond how to open one. Everything else in this UI is inside a single Customer.
func (s *server) index(w http.ResponseWriter, r *http.Request) {
	type row struct {
		Slug, Name, LastRun string
		Pending, Findings   int
	}
	var rows []row
	for _, c := range s.cfg.Customers {
		item := row{Slug: c.Slug, Name: c.Name}
		for _, p := range s.proposals(c.Slug) {
			if p.State == "pending" {
				item.Pending++
			}
		}
		if sums, _ := run.Summaries(s.dataDir, c.Slug); len(sums) > 0 {
			item.LastRun, item.Findings = sums[0].RunAt, sums[0].Findings
		}
		rows = append(rows, item)
	}
	s.render(w, "index", map[string]any{"Customers": rows})
}

func (s *server) customer(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	cust := s.cfg.Customer(slug)
	if cust == nil {
		http.NotFound(w, r)
		return
	}
	sums, err := run.Summaries(s.dataDir, slug)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var latest *run.Summary
	if len(sums) > 0 {
		latest = &sums[0]
	}
	m, err := state.LoadManifest(s.customerDir(slug), slug)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	s.render(w, "customer", map[string]any{
		"Customer":  cust,
		"Summary":   latest,
		"History":   sums,
		"Proposals": s.proposals(slug),
		"Manifest":  m,
		"Run":       s.runStateOf(slug),
	})
}

func (s *server) proposals(slug string) []state.Proposal {
	ps, err := state.LoadProposals(s.customerDir(slug))
	if err != nil {
		log.Printf("loading proposals for %s: %v", slug, err)
	}
	return ps
}

// record serves one Source Record so the evidence behind a Finding is one click away. The
// path is joined to the Customer Directory and refused if it leaves — the same boundary the
// Routing Rule draws, enforced again at the only place this UI reads files.
func (s *server) record(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if s.cfg.Customer(slug) == nil {
		http.NotFound(w, r)
		return
	}
	base, err := filepath.Abs(s.customerDir(slug))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	full := filepath.Join(base, filepath.Clean("/"+r.URL.Query().Get("path")))
	if !strings.HasPrefix(full, base+string(os.PathSeparator)) {
		http.Error(w, "outside the Customer Directory", http.StatusForbidden)
		return
	}
	b, err := os.ReadFile(full)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write(b)
}

func (s *server) runStateOf(slug string) *runState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.runs[slug]
}

// startRun launches a Run in the background and answers immediately. The Run's lifetime is
// not tied to anyone watching this page.
func (s *server) startRun(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if s.cfg.Customer(slug) == nil {
		http.NotFound(w, r)
		return
	}
	s.mu.Lock()
	if st := s.runs[slug]; st != nil && !st.Done {
		s.mu.Unlock()
		s.render(w, "runstatus", map[string]any{"Slug": slug, "Run": st})
		return
	}
	st := &runState{Started: time.Now()}
	s.runs[slug] = st
	s.mu.Unlock()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
		defer cancel()
		sum, err := run.Once(ctx, s.cfg, slug, run.Options{DataDir: s.dataDir})
		s.mu.Lock()
		defer s.mu.Unlock()
		st.Summary, st.Done = sum, true
		if err != nil {
			st.Err = err.Error()
		}
	}()
	s.render(w, "runstatus", map[string]any{"Slug": slug, "Run": st})
}

func (s *server) runStatus(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	s.render(w, "runstatus", map[string]any{"Slug": slug, "Run": s.runStateOf(slug)})
}

// decide is the gate itself. Approval applies the artifact there and then, through the same
// writer a Run uses; a rejection becomes a Suppression, revocable in one click.
func (s *server) decide(w http.ResponseWriter, r *http.Request) {
	slug, hash, decision := r.PathValue("slug"), r.PathValue("hash"), r.PathValue("decision")
	cust := s.cfg.Customer(slug)
	if cust == nil {
		http.NotFound(w, r)
		return
	}
	dir := s.customerDir(slug)
	m, err := state.LoadManifest(dir, slug)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var p state.Proposal
	switch decision {
	case "approve":
		p, err = publish.Approve(dir, s.cfg.SMTP.Addr, s.cfg.SMTP.From, hash, m)
	case "reject":
		p, err = publish.Reject(dir, hash, r.FormValue("note"), m)
	case "revoke":
		if err = publish.Revoke(dir, hash, m); err == nil {
			p, err = state.LoadProposal(dir, hash)
		}
	default:
		http.Error(w, "unknown decision", http.StatusBadRequest)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusOK) // htmx swaps the row either way; the error belongs in it
		s.render(w, "proposal", map[string]any{"Slug": slug, "P": p, "Error": err.Error()})
		return
	}
	s.render(w, "proposal", map[string]any{"Slug": slug, "P": p})
}

func (s *server) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tmpl.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("rendering %s: %v", name, err)
	}
}

var funcs = template.FuncMap{
	// dict lets a page hand a partial the two things it needs without a struct per partial.
	"dict": func(kv ...any) map[string]any {
		m := map[string]any{}
		for i := 0; i+1 < len(kv); i += 2 {
			m[fmt.Sprint(kv[i])] = kv[i+1]
		}
		return m
	},
	"short": func(s string) string {
		if len(s) > 8 {
			return s[:8]
		}
		return s
	},
	"clock": func(s string) string {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return s
		}
		return t.Local().Format("2 Jan 15:04")
	},
	"plural": func(n int, word string) string {
		if n == 1 {
			return fmt.Sprintf("%d %s", n, word)
		}
		return fmt.Sprintf("%d %ss", n, word)
	},
}
