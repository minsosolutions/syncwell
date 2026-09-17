package state

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Proposal is a Gated Action a Run prepared and did not take. It lives in the Customer
// Directory beside the Manifest, so the next Run sees what is already pending or suppressed
// on its behalf — unlike Drift reconciliation, which is our bookkeeping and stays out.
//
// Hash binds the approval to the artifact byte for byte: a human approves this subject, this
// body, these recipients, never a decision that is rendered afterwards. See ADR-0003.
type Proposal struct {
	Hash       string   `json:"hash"`
	Customer   string   `json:"customer"`
	FindingKey string   `json:"finding_key"`
	Kind       string   `json:"kind"` // email
	To         []string `json:"to"`
	Subject    string   `json:"subject"`
	Body       string   `json:"body"`
	RunAt      string   `json:"run_at"`
	State      string   `json:"state"` // pending | applied | rejected
	DecidedAt  string   `json:"decided_at,omitempty"`
	Note       string   `json:"note,omitempty"`
}

// NewProposal computes the Proposal's identity from its content, so re-proposing the same
// artifact is the same Proposal and changing one word is a different one.
func NewProposal(customer, key, kind, subject, body string, to []string, runAt string) Proposal {
	p := Proposal{
		Hash: "", Customer: customer, FindingKey: key, Kind: kind,
		To: to, Subject: subject, Body: body, RunAt: runAt, State: "pending",
	}
	sum := sha256.Sum256([]byte(strings.Join([]string{customer, key, kind, subject, body, strings.Join(to, ",")}, "\x00")))
	p.Hash = hex.EncodeToString(sum[:])[:16]
	return p
}

func proposalsDir(customerDir string) string { return filepath.Join(customerDir, "state", "proposals") }

func (p Proposal) Save(customerDir string) error {
	if err := os.MkdirAll(proposalsDir(customerDir), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(proposalsDir(customerDir), p.Hash+".json"), append(b, '\n'), 0o644)
}

func LoadProposals(customerDir string) ([]Proposal, error) {
	entries, err := os.ReadDir(proposalsDir(customerDir))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var out []Proposal
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(proposalsDir(customerDir), e.Name()))
		if err != nil {
			return nil, err
		}
		var p Proposal
		if err := json.Unmarshal(b, &p); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].RunAt+out[i].Hash < out[j].RunAt+out[j].Hash })
	return out, nil
}

func LoadProposal(customerDir, hash string) (Proposal, error) {
	b, err := os.ReadFile(filepath.Join(proposalsDir(customerDir), hash+".json"))
	if err != nil {
		return Proposal{}, err
	}
	var p Proposal
	return p, json.Unmarshal(b, &p)
}

// Supersede drops the pending Proposals this Customer holds for a Finding that are not the
// artifact just prepared. A decided Proposal is never touched: it is the audit trail.
func Supersede(customerDir, key, keepHash string) error {
	existing, err := LoadProposals(customerDir)
	if err != nil {
		return err
	}
	for _, p := range existing {
		if p.FindingKey == key && p.State == "pending" && p.Hash != keepHash {
			if err := os.Remove(filepath.Join(proposalsDir(customerDir), p.Hash+".json")); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Proposal) Decide(state, note string) {
	p.State = state
	p.Note = note
	p.DecidedAt = time.Now().UTC().Format(time.RFC3339)
}
