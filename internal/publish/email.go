package publish

import (
	"fmt"
	"net/smtp"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/minsosolutions/syncwell/internal/output"
	"github.com/minsosolutions/syncwell/internal/routing"
	"github.com/minsosolutions/syncwell/internal/state"
)

// Email crossing to the Customer is a Gated Action: a Run prepares it and never sends it.
// The Run does not wait either — it leaves Proposals in the Customer Directory and finishes.
// See docs/adr/0003-approval-gates-on-blast-radius.md.

var addressRE = regexp.MustCompile(`[\w.+-]+@[\w.-]+\.[A-Za-z]{2,}`)

// Recipients derives who an email may go to, in Go, from the Customer's own records. An
// address qualifies only if it appears in a Source Record collected for this Customer and
// matches one of that Customer's configured domains. A model that picks recipients can pick
// the wrong Customer's contact, which is the leak the Routing Rule prevents one step earlier.
func Recipients(customerDir string, cust *routing.Customer, f output.Finding) ([]string, error) {
	cited, err := addressesIn(citedFiles(customerDir, f))
	if err != nil {
		return nil, err
	}
	found := filterByDomain(cited, cust.Domains)
	if len(found) > 0 {
		return found, nil
	}
	// No Customer address in the cited records — the Finding may rest on chat and tickets.
	// Fall back to the Customer's correspondents, still from their own records only.
	all, err := addressesIn(mailFiles(customerDir))
	if err != nil {
		return nil, err
	}
	return filterByDomain(all, cust.Domains), nil
}

func citedFiles(customerDir string, f output.Finding) []string {
	var out []string
	for _, e := range f.Evidence {
		if e.Kind == "record" && e.Path != "" {
			out = append(out, filepath.Join(customerDir, filepath.Clean(e.Path)))
		}
	}
	return out
}

func mailFiles(customerDir string) []string {
	files, _ := filepath.Glob(filepath.Join(customerDir, "email", "*.md"))
	return files
}

func addressesIn(paths []string) ([]string, error) {
	var out []string
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		for _, line := range strings.Split(unescapeJSON(string(b)), "\n") {
			if isIdentifierLine(line) {
				continue
			}
			for _, a := range addressRE.FindAllString(line, -1) {
				out = append(out, strings.ToLower(a))
			}
		}
	}
	return out, nil
}

var jsonEscapeRE = regexp.MustCompile(`\\u[0-9a-fA-F]{4}`)

// unescapeJSON turns \uXXXX back into the character it stands for before addresses are read
// out of a file. Zammad stores `Gustav Ebbe \u003cgustav.ebbe@ebbahus.se\u003e`, and a regex
// over the raw bytes derives the recipient `u003cgustav.ebbe@ebbahus.se` — an address that
// does not exist, on a domain that does.
func unescapeJSON(s string) string {
	if !strings.Contains(s, `\u`) {
		return s
	}
	return jsonEscapeRE.ReplaceAllStringFunc(s, func(esc string) string {
		n, err := strconv.ParseUint(esc[2:], 16, 32)
		if err != nil {
			return esc
		}
		return string(rune(n))
	})
}

// isIdentifierLine skips the headers that carry a message id rather than a person. A
// message id is shaped like an address and lives on the Customer's own domain, so regexing
// a mail file without this posts to <20260909114700.9964@ebbahus.se>.
func isIdentifierLine(line string) bool {
	l := strings.ToLower(strings.TrimLeft(line, " \t\""))
	for _, h := range []string{"message_id", "message-id", "in_reply_to", "in-reply-to", "references"} {
		if strings.HasPrefix(l, h) {
			return true
		}
	}
	return false
}

func filterByDomain(addresses, domains []string) []string {
	var out []string
	for _, a := range addresses {
		at := strings.LastIndex(a, "@")
		if at < 0 {
			continue
		}
		for _, d := range domains {
			if strings.EqualFold(a[at+1:], d) {
				out = append(out, a)
				break
			}
		}
	}
	slices.Sort(out)
	return slices.Compact(out)
}

// Propose writes the Gated Actions this Run prepared. A Finding whose Proposal a human has
// already rejected is not proposed again: a Suppression is obeyed until the human revokes it.
func Propose(customerDir string, cust *routing.Customer, m *state.Manifest, r output.Report, runAt string) ([]state.Proposal, []Action, error) {
	var proposals []state.Proposal
	var actions []Action
	for _, f := range r.Items {
		if f.Proposal == nil || f.Proposal.Kind != "email" {
			continue
		}
		if m.Suppressed(f.Key, "email") {
			actions = append(actions, Action{"skipped", "", f.Key, "email suppressed: a human rejected it"})
			continue
		}
		to, err := Recipients(customerDir, cust, f)
		if err != nil {
			return nil, nil, err
		}
		if len(to) == 0 {
			actions = append(actions, Action{"skipped", "", f.Key, "no recipient derivable from this customer's records"})
			continue
		}
		p := state.NewProposal(cust.Slug, f.Key, "email", f.Proposal.Subject, f.Proposal.Body, to, runAt)
		if m.Applied(p.Hash) {
			actions = append(actions, Action{"skipped", "", f.Key, "this exact email has already been sent"})
			continue
		}
		if err := p.Save(customerDir); err != nil {
			return nil, nil, err
		}
		// A reworded draft is a different artifact, so the approval that was pending on the
		// old wording must not carry over to it.
		if err := state.Supersede(customerDir, f.Key, p.Hash); err != nil {
			return nil, nil, err
		}
		proposals = append(proposals, p)
		actions = append(actions, Action{"proposed", p.Hash[:8], f.Key, fmt.Sprintf("email to %s, awaiting approval", strings.Join(to, ", "))})
	}
	return proposals, actions, nil
}

// Send hands the approved artifact to the mail server byte for byte. The Provenance Marker
// rides along in the footer: an Output a Customer reads still says who made it.
func Send(addr, from string, p state.Proposal) error {
	body := p.Body + "\n\n-- \nSent by Syncwell on behalf of Minso.\n" +
		state.Marker{Run: p.RunAt, Customer: p.Customer, Item: p.FindingKey}.String() + "\n"
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n%s",
		from, strings.Join(p.To, ", "), p.Subject, strings.ReplaceAll(body, "\n", "\r\n"))
	return smtp.SendMail(addr, nil, from, p.To, []byte(msg))
}

// Approve applies a Proposal through the same writer path a Run would use, and records it in
// the Manifest: every later Run asks "have we already told them?".
func Approve(customerDir, smtpAddr, from, hash string, m *state.Manifest) (state.Proposal, error) {
	p, err := state.LoadProposal(customerDir, hash)
	if err != nil {
		return p, err
	}
	if p.State != "pending" {
		return p, fmt.Errorf("proposal %s is %s, not pending", hash, p.State)
	}
	if err := Send(smtpAddr, from, p); err != nil {
		return p, fmt.Errorf("sending to %s: %w", strings.Join(p.To, ", "), err)
	}
	p.Decide("applied", "")
	if err := p.Save(customerDir); err != nil {
		return p, err
	}
	m.Record(state.Action{Key: p.FindingKey, Kind: p.Kind, State: "applied", Hash: p.Hash, To: p.To, Subject: p.Subject})
	return p, m.Save(customerDir)
}

// Reject turns a refusal into a Suppression. It is remembered against the Finding and the
// action, never expires on its own, and is revoked by the human who made it.
func Reject(customerDir, hash, note string, m *state.Manifest) (state.Proposal, error) {
	p, err := state.LoadProposal(customerDir, hash)
	if err != nil {
		return p, err
	}
	p.Decide("rejected", note)
	if err := p.Save(customerDir); err != nil {
		return p, err
	}
	m.Record(state.Action{Key: p.FindingKey, Kind: p.Kind, State: "suppressed", Hash: p.Hash, To: p.To, Subject: p.Subject})
	return p, m.Save(customerDir)
}

// Revoke lifts a Suppression, which is the other half of obeying a human gesture: a decision
// a human can make and not unmake is not a decision, it is a trap.
func Revoke(customerDir, hash string, m *state.Manifest) error {
	kept := m.Actions[:0]
	for _, a := range m.Actions {
		if a.Hash == hash && a.State == "suppressed" {
			continue
		}
		kept = append(kept, a)
	}
	m.Actions = kept
	if p, err := state.LoadProposal(customerDir, hash); err == nil {
		p.State = "pending"
		p.DecidedAt = ""
		if err := p.Save(customerDir); err != nil {
			return err
		}
	}
	return m.Save(customerDir)
}
