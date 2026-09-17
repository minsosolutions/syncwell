package agent

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/minsosolutions/syncwell/internal/output"
)

// writePreviousFindings leaves the next Run the keys and titles this one surfaced, and
// nothing else: the Manifest stays out of the Agent's view. See issue #10.
func WritePreviousFindings(customerDir string, r output.Report) error {
	type previous struct {
		Key     string `json:"key"`
		Title   string `json:"title"`
		Summary string `json:"summary"`
	}
	items := make([]previous, 0, len(r.Items))
	for _, f := range r.Items {
		items = append(items, previous{Key: f.Key, Title: f.Title, Summary: f.Summary})
	}
	b, err := json.MarshalIndent(map[string]any{"run_at": r.RunAt, "findings": items}, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Join(customerDir, "state")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "previous-findings.json"), b, 0o644)
}
