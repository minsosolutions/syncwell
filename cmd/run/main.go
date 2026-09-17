// Command run is one Run end to end for one Customer: the Agent judges, Go writes.
//
// It writes report.md and the Customer's Linear issues unattended, and leaves anything that
// crosses to the Customer as a Proposal for cmd/admin. It never waits for a human.
//
//	go run ./cmd/run -customer ebbahus
//	go run ./cmd/run -customer ebbahus -findings findings.json   # replay, no model
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/minsosolutions/syncwell/internal/routing"
	"github.com/minsosolutions/syncwell/internal/run"
)

func main() {
	dataDir := flag.String("data", "data", "data directory")
	configPath := flag.String("config", "config/syncwell.json", "config file")
	customer := flag.String("customer", "", "customer slug")
	findings := flag.String("findings", "", "replay this findings JSON instead of launching the Agent")
	timeout := flag.Duration("timeout", 15*time.Minute, "how long one Run may take")
	maxTurns := flag.Int("max-turns", 0, "cap on the Agent's turns; 0 uses the default")
	budget := flag.String("budget", "", "cap on spend in USD; empty uses the default")
	flag.Parse()
	if *customer == "" {
		log.Fatal("-customer is required")
	}

	cfg, err := routing.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	opts := run.Options{DataDir: *dataDir, MaxTurns: *maxTurns, Budget: *budget}
	if *findings != "" {
		f, err := os.Open(*findings)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		opts.Findings = f
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	sum, err := run.Once(ctx, cfg, *customer, opts)
	if sum != nil {
		print(sum)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func print(s *run.Summary) {
	fmt.Printf("run %s — %s, %d finding(s), $%.3f, %s\n", s.Customer, s.RunAt, s.Findings, s.CostUSD, s.Duration)
	for _, d := range s.Dropped {
		fmt.Println("  dropped  ", d)
	}
	for _, d := range s.Drift {
		fmt.Println("  drift    ", d)
	}
	for _, a := range s.Actions {
		fmt.Println("  ", a)
	}
	fmt.Println("  report   ", s.ReportPath)
	if len(s.Proposals) > 0 {
		fmt.Printf("  %d proposal(s) awaiting approval — go run ./cmd/admin\n", len(s.Proposals))
	}
}
