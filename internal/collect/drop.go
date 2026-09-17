package collect

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/minsosolutions/syncwell/internal/routing"
)

// dropGlob matches <YYYY-MM-DD>-<slug>.md and nothing else, so the README in a drop directory
// and any stray note beside it are never mistaken for a Source Record.
const dropGlob = "[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]-*.md"

// collectDrop walks a file-drop Source: no API, just files a scheduled export left in
// data/sources/<source>/. Each file is routed by its own rule and copied byte-for-byte into
// one Customer Directory, or quarantined with the reason beside it. The file is never
// rewritten — its headers are the routing key and its body is the evidence a finding cites.
//
// Transcripts and email share this walk. The path guards in it are what keep one Customer's
// records out of another's directory, and a second copy of them is a copy that can drift.
func collectDrop(dataDir, source string, route func(body string) (*routing.Customer, error)) (Report, error) {
	var rep Report
	files, err := filepath.Glob(filepath.Join(dataDir, "sources", source, dropGlob))
	if err != nil {
		return rep, err
	}

	for _, file := range files {
		body, err := os.ReadFile(file)
		if err != nil {
			return rep, err
		}
		name, err := segment(source, strings.TrimSuffix(filepath.Base(file), ".md"))
		if err != nil {
			return rep, err
		}

		cust, routeErr := route(string(body))
		if routeErr != nil {
			var u *routing.Unroutable
			if !errors.As(routeErr, &u) {
				return rep, routeErr
			}
			if err := quarantine(dataDir, source, name, body, u.Reason); err != nil {
				return rep, err
			}
			rep.Quarantined = append(rep.Quarantined, fmt.Sprintf("%s.md: %s", name, u.Reason))
			continue
		}

		slug, err := segment("customer slug", cust.Slug)
		if err != nil {
			return rep, err
		}
		if err := writeFile(filepath.Join(dataDir, "customers", slug, source, name+".md"), body); err != nil {
			return rep, err
		}
		rep.Routed = append(rep.Routed, fmt.Sprintf("%s.md -> %s", name, cust.Slug))
	}
	return rep, nil
}

// quarantine keeps the record byte-for-byte and states why it is here in a file beside it.
func quarantine(dataDir, source, name string, body []byte, reason string) error {
	dir := filepath.Join(dataDir, "unrouted", source)
	if err := writeFile(filepath.Join(dir, name+".md"), body); err != nil {
		return err
	}
	return writeFile(filepath.Join(dir, name+".reason.txt"), []byte(reason+"\n"))
}

// frontmatter reads one key out of the flat key/value header. The format is deliberately not
// YAML (see the source READMEs), so neither is the parser. The scan stops at the closing ---
// so nothing in the body — a line someone said, a quoted reply — can be read as a header.
func frontmatter(body, key string) string {
	lines := strings.Split(body, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return ""
	}
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "---" {
			return ""
		}
		if rest, ok := strings.CutPrefix(line, key+":"); ok {
			return strings.TrimSpace(rest)
		}
	}
	return ""
}

func writeFile(path string, body []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, body, 0o644)
}
