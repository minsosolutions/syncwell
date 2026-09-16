// Command report renders a markdown Output from the JSON an Agent produces.
//
// The split is the point: the Agent decides what needs attention, this writes it down
// deterministically. Swap this binary for a spreadsheet or Linear writer and the Agent
// does not change.
//
//	go run ./cmd/report < items.json
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/minsosolutions/syncwell/internal/output"
)

func main() {
	dataDir := flag.String("data", "data", "data directory")
	out := flag.String("out", "", "output path; defaults to data/customers/<customer>/outputs/report.md")
	flag.Parse()

	var r output.Report
	if err := json.NewDecoder(os.Stdin).Decode(&r); err != nil {
		log.Fatalf("reading report JSON from stdin: %v", err)
	}
	if r.Customer == "" {
		log.Fatal("report JSON needs a \"customer\"")
	}

	path := *out
	if path == "" {
		path = filepath.Join(*dataDir, "customers", r.Customer, "outputs", "report.md")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := output.Write(f, r); err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote", path)
}
