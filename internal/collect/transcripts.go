package collect

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/minsosolutions/syncwell/internal/routing"
)

// transcriptGlob matches <YYYY-MM-DD>-<slug>.md and nothing else, so the README and any
// stray note in the drop directory are never mistaken for a Source Record.
const transcriptGlob = "[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]-*.md"

// Transcripts collects the meeting transcript file drop. There is no API: a scheduled export
// leaves files in data/sources/transcripts/ and this reads them from disk. Each file is routed
// by its attendee list and copied verbatim into one Customer Directory, or quarantined with
// the reason beside it. The file is never rewritten — the attendee list is the routing key and
// the transcript is the evidence.
func Transcripts(cfg *routing.Config, dataDir string) (Report, error) {
	var rep Report
	srcDir := filepath.Join(dataDir, "sources", "transcripts")
	files, err := filepath.Glob(filepath.Join(srcDir, transcriptGlob))
	if err != nil {
		return rep, err
	}

	for _, file := range files {
		body, err := os.ReadFile(file)
		if err != nil {
			return rep, err
		}
		name, err := segment("transcript", strings.TrimSuffix(filepath.Base(file), ".md"))
		if err != nil {
			return rep, err
		}

		list := attendees(string(body))
		if len(list) == 0 {
			if err := quarantine(dataDir, name, body, "frontmatter has no attendees; nothing to route on"); err != nil {
				return rep, err
			}
			rep.Quarantined = append(rep.Quarantined, name+".md: frontmatter has no attendees; nothing to route on")
			continue
		}

		cust, routeErr := cfg.ByAttendeeDomains(list)
		if routeErr != nil {
			var u *routing.Unroutable
			if !errors.As(routeErr, &u) {
				return rep, routeErr
			}
			if err := quarantine(dataDir, name, body, u.Reason); err != nil {
				return rep, err
			}
			rep.Quarantined = append(rep.Quarantined, fmt.Sprintf("%s.md: %s", name, u.Reason))
			continue
		}

		slug, err := segment("customer slug", cust.Slug)
		if err != nil {
			return rep, err
		}
		if err := writeFile(filepath.Join(dataDir, "customers", slug, "transcripts", name+".md"), body); err != nil {
			return rep, err
		}
		rep.Routed = append(rep.Routed, fmt.Sprintf("%s.md -> %s", name, cust.Slug))
	}
	return rep, nil
}

// attendees reads the flat frontmatter's comma-separated attendee list. The format is
// deliberately not YAML (see data/sources/transcripts/README.md), so neither is the parser.
// The scan stops at the closing --- so nothing anyone says in the meeting can be read as
// frontmatter.
func attendees(body string) []string {
	lines := strings.Split(body, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return nil
	}
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "---" {
			return nil
		}
		rest, ok := strings.CutPrefix(line, "attendees:")
		if !ok {
			continue
		}
		var out []string
		for _, a := range strings.Split(rest, ",") {
			if a = strings.TrimSpace(a); a != "" {
				out = append(out, a)
			}
		}
		return out
	}
	return nil
}

// quarantine keeps the transcript byte-for-byte and states why it is here in a file beside it.
func quarantine(dataDir, name string, body []byte, reason string) error {
	dir := filepath.Join(dataDir, "unrouted", "transcripts")
	if err := writeFile(filepath.Join(dir, name+".md"), body); err != nil {
		return err
	}
	return writeFile(filepath.Join(dir, name+".reason.txt"), []byte(reason+"\n"))
}

func writeFile(path string, body []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, body, 0o644)
}
