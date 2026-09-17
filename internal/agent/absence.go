package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/minsosolutions/syncwell/internal/collect"
	"github.com/minsosolutions/syncwell/internal/output"
)

// ResolveAbsences turns each record an Agent flagged AwaitingReply into either a stated
// absence or a dead Finding. The Agent never writes the absence itself: a claim about
// everything is not a claim a model can make. See issue #16.
//
// Only email is answerable this way — a mail has a thread to follow. A record from a Source
// with no thread relation is left alone rather than guessed at.
func ResolveAbsences(r output.Report, customerDir string) (output.Report, []Drop) {
	var drops []Drop
	kept := r
	kept.Items = nil

	for _, f := range r.Items {
		reply, evidence, err := resolveFinding(f, customerDir)
		switch {
		case err != nil:
			drops = append(drops, Drop{Key: f.Key, Reason: err.Error()})
		case reply != "":
			drops = append(drops, Drop{Key: f.Key, Reason: "rests on silence, but " + reply + " replies to it"})
		default:
			f.Evidence = evidence
			kept.Items = append(kept.Items, f)
		}
	}
	return kept, drops
}

func resolveFinding(f output.Finding, customerDir string) (reply string, evidence []output.Evidence, err error) {
	evidence = f.Evidence
	for _, e := range f.Evidence {
		if !e.AwaitingReply || e.Source != "email" {
			continue
		}
		path, err := resolve(customerDir, e.Path)
		if err != nil {
			return "", nil, err
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return "", nil, fmt.Errorf("evidence cites %s, which cannot be read: %v", e.Path, err)
		}
		asked := mailOf(string(b))
		if asked.messageID == "" {
			return "", nil, fmt.Errorf("%s has no message_id, so nothing can be said about a reply to it", e.Path)
		}
		if id := replyTo(asked, filepath.Dir(path)); id != "" {
			return id, nil, nil
		}
		evidence = append(evidence, output.Evidence{
			Kind:     "absence",
			Claim:    "nothing replies to " + asked.messageID,
			Searched: []string{"email"},
			Since:    asked.date,
		})
	}
	return "", evidence, nil
}

type mail struct {
	messageID, inReplyTo, subject, date, from, to string
}

func mailOf(body string) mail {
	return mail{
		messageID: collect.Frontmatter(body, "message_id"),
		inReplyTo: collect.Frontmatter(body, "in_reply_to"),
		subject:   collect.Frontmatter(body, "subject"),
		date:      collect.Frontmatter(body, "date"),
		from:      collect.Frontmatter(body, "from"),
		to:        collect.Frontmatter(body, "to"),
	}
}

// replyTo returns the message id of a mail answering asked, or "". A reply is one that
// threads to it, or — for the mail clients that lose the header — a later mail between the
// same people carrying the same subject. Keyword overlap is deliberately not a reply.
func replyTo(asked mail, dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		b, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}
		m := mailOf(string(b))
		if m.messageID == "" || m.messageID == asked.messageID {
			continue
		}
		if m.inReplyTo == asked.messageID {
			return m.messageID
		}
		if m.date > asked.date && subjectOf(m.subject) == subjectOf(asked.subject) && sharesAPerson(m, asked) {
			return m.messageID
		}
	}
	return ""
}

// subjectOf strips the reply and forward prefixes a thread accumulates, in the languages
// this inbox actually sees.
func subjectOf(s string) string {
	s = strings.TrimSpace(s)
	for {
		lower := strings.ToLower(s)
		trimmed := false
		for _, p := range []string{"re:", "sv:", "vb:", "fw:", "fwd:"} {
			if strings.HasPrefix(lower, p) {
				s = strings.TrimSpace(s[len(p):])
				trimmed = true
				break
			}
		}
		if !trimmed {
			return strings.ToLower(s)
		}
	}
}

// sharesAPerson is true when the later mail involves whoever was asked. Addresses are
// matched loosely because headers carry display names around them.
func sharesAPerson(m, asked mail) bool {
	people := strings.ToLower(m.from + " " + m.to)
	for _, addr := range addresses(asked.to + " " + asked.from) {
		if strings.Contains(people, addr) {
			return true
		}
	}
	return false
}

func addresses(s string) []string {
	var out []string
	for _, f := range strings.FieldsFunc(strings.ToLower(s), func(r rune) bool {
		return r == ' ' || r == ',' || r == '<' || r == '>'
	}) {
		if strings.Contains(f, "@") {
			out = append(out, f)
		}
	}
	return out
}
