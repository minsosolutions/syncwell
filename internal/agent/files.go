package agent

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Files lists the working set an Agent may read, relative to dir and in a stable order.
//
// Claude Code 2.1.274 gives a session Bash, Edit and Read; there is no search tool, and Bash
// is the thing we removed so that isolation is a property of the tool list rather than of
// the model's behaviour. Discovery therefore happens here, where it is deterministic.
//
// The Manifest is left out: an Agent reads state/previous-findings.json and nothing else in
// state/, which is also enforced by a deny rule.
func Files(dir string) ([]string, error) {
	var out []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") && path != dir {
				return fs.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		if strings.HasPrefix(d.Name(), ".") || rel == filepath.Join("state", "manifest.json") {
			return nil
		}
		out = append(out, filepath.ToSlash(rel))
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	slices.Sort(out)
	return out, nil
}
