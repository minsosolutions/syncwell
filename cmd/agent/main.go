// Command agent runs one Customer's Agent and prints the findings it can stand behind.
//
// The Agent judges; this validates. Nothing here writes an Output: quotes are checked
// against the records they cite, absences are computed rather than believed, and what
// survives goes to stdout as JSON for an Output writer to render. Use cmd/run for a whole
// Run, Outputs included.
//
//	go run ./cmd/agent -customer ebbahus | go run ./cmd/report
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/minsosolutions/syncwell/internal/agent"
	"github.com/minsosolutions/syncwell/internal/run"
)

func main() {
	dataDir := flag.String("data", "data", "data directory")
	customer := flag.String("customer", "", "customer slug")
	timeout := flag.Duration("timeout", 15*time.Minute, "how long one Run may take")
	maxTurns := flag.Int("max-turns", 0, "cap on the Agent's turns; 0 uses the default")
	budget := flag.String("budget", "", "cap on spend in USD; empty uses the default")
	flag.Parse()
	if *customer == "" {
		log.Fatal("-customer is required")
	}

	customerDir, err := filepath.Abs(filepath.Join(*dataDir, "customers", *customer))
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(customerDir); err != nil {
		log.Fatalf("no Customer Directory for %q: %v", *customer, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	runAt := time.Now().UTC().Format(time.RFC3339)
	report, cost, denials, err := run.Findings(ctx,
		run.Options{DataDir: *dataDir, MaxTurns: *maxTurns, Budget: *budget},
		customerDir, filepath.Join(*dataDir, "references"), *customer, runAt)
	if err != nil {
		log.Fatal(err)
	}
	report.Customer, report.RunAt = *customer, runAt

	report, dropped := agent.Verify(report, customerDir)
	report, silent := agent.ResolveAbsences(report, customerDir)
	dropped = append(dropped, silent...)

	for _, d := range dropped {
		fmt.Fprintln(os.Stderr, "dropped", d)
	}
	for _, d := range denials {
		fmt.Fprintln(os.Stderr, "denied", d)
	}
	fmt.Fprintf(os.Stderr, "%d finding(s), %d dropped, $%.3f\n", len(report.Items), len(dropped), cost)

	if err := agent.WritePreviousFindings(customerDir, report); err != nil {
		log.Fatal(err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		log.Fatal(err)
	}
}
