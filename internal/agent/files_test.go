package agent

import (
	"path/filepath"
	"slices"
	"testing"
)

// The Agent has Read and nothing else — Claude Code 2.1.274 ships no Grep or Glob tool, and
// Bash is what we removed to make isolation structural. So Go does the discovery.
func TestFilesListsTheWorkingSetRelatively(t *testing.T) {
	files, err := Files(filepath.Join("..", "..", "data", "customers", "ebbahus"))
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(files, "email/2026-09-09-ebbahus-saknade-utbetalningar.md") {
		t.Errorf("the working set does not list its own email: %v", files)
	}
	if !slices.IsSorted(files) {
		t.Error("the listing is not stable between Runs")
	}
	if slices.Contains(files, "state/manifest.json") {
		t.Error("the Manifest is in the listing; reconciliation is Go's business (#3, #10)")
	}
	for _, f := range files {
		if filepath.IsAbs(f) {
			t.Errorf("%q is absolute; paths are cited relative to the Customer Directory", f)
		}
	}
}
