package collect

import (
	"strings"

	"github.com/minsosolutions/syncwell/internal/routing"
)

// Transcripts collects the meeting transcript file drop and routes each meeting on who was in
// the room. There is no API: a scheduled export leaves files in data/sources/transcripts/.
func Transcripts(cfg *routing.Config, dataDir string) (Report, error) {
	return collectDrop(dataDir, "transcripts", func(body string) (*routing.Customer, error) {
		list := attendees(body)
		if len(list) == 0 {
			return nil, &routing.Unroutable{Reason: "frontmatter has no attendees; nothing to route on"}
		}
		return cfg.ByAttendeeDomains(list)
	})
}

// attendees splits the frontmatter's comma-separated attendee list.
func attendees(body string) []string {
	var out []string
	for _, a := range strings.Split(frontmatter(body, "attendees"), ",") {
		if a = strings.TrimSpace(a); a != "" {
			out = append(out, a)
		}
	}
	return out
}
