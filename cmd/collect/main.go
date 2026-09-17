package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/minsosolutions/syncwell/internal/collect"
	"github.com/minsosolutions/syncwell/internal/routing"
)

func main() {
	source := flag.String("source", "slack", "source to collect: slack, zammad, transcripts, email")
	api := flag.String("api", "http://localhost:8099", "mock API base URL")
	token := flag.String("token", "dev-token", "bearer token; any non-empty value works")
	configPath := flag.String("config", "config/syncwell.json", "routing config")
	dataDir := flag.String("data", "data", "data directory")
	flag.Parse()

	cfg, err := routing.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	var rep collect.Report
	switch *source {
	case "slack":
		client := &collect.Client{BaseURL: *api, Token: *token, HTTP: &http.Client{Timeout: 10 * time.Second}}
		rep, err = collect.Slack(client, cfg, *dataDir)
	case "zammad":
		client := &collect.Client{BaseURL: *api, Token: *token, HTTP: &http.Client{Timeout: 10 * time.Second}}
		rep, err = collect.Zammad(client, cfg, *dataDir)
	case "transcripts":
		// A file drop, not an API: no client, no token.
		rep, err = collect.Transcripts(cfg, *dataDir)
	case "email":
		rep, err = collect.Email(cfg, *dataDir)
	default:
		log.Fatalf("source %q is not implemented — see data/sources/%s/README.md and write it", *source, *source)
	}
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range rep.Routed {
		fmt.Println("routed     ", line)
	}
	for _, line := range rep.Quarantined {
		fmt.Println("quarantined", line)
	}
}
