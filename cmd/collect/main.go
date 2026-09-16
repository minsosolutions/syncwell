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
	source := flag.String("source", "slack", "source to collect (only slack is implemented)")
	api := flag.String("api", "http://localhost:8099", "mock API base URL")
	token := flag.String("token", "dev-token", "bearer token; any non-empty value works")
	configPath := flag.String("config", "config/syncwell.json", "routing config")
	dataDir := flag.String("data", "data", "data directory")
	flag.Parse()

	if *source != "slack" {
		log.Fatalf("source %q is not implemented — see data/sources/%s/README.md and write it", *source, *source)
	}

	cfg, err := routing.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	client := &collect.Client{BaseURL: *api, Token: *token, HTTP: &http.Client{Timeout: 10 * time.Second}}

	rep, err := collect.Slack(client, cfg, *dataDir)
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
