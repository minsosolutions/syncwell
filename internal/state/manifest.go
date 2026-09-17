package state

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"time"
)

// Snapshot is the Written Snapshot: what a Run last wrote to a Maintained Output. Without it
// a human's retitle cannot be told from our own previous write.
type Snapshot struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Priority    int    `json:"priority,omitempty"`
}

// Output is one Maintained Output across Runs. Stewardship is deliberately not stored: it is
// carried by the Marker in the Output itself and read fresh every Run (ADR-0002).
type Output struct {
	Kind        string   `json:"kind"` // linear_issue | linear_project | file
	Key         string   `json:"key"`  // the Finding Key this Output stands for
	ID          string   `json:"id,omitempty"`
	Identifier  string   `json:"identifier,omitempty"`
	Path        string   `json:"path,omitempty"` // kind=file, relative to the Customer Directory
	AuthoredRun string   `json:"authored_run"`
	LastRun     string   `json:"last_run,omitempty"`
	Records     []string `json:"records,omitempty"` // "<path>#<locator>" of the Evidence behind it
	Snapshot    Snapshot `json:"snapshot,omitempty"`
	// HumanFields are the fields a human has diverged from our Snapshot. Once a field is
	// here it belongs to the human: later Runs comment rather than overwrite (#9).
	HumanFields []string `json:"human_fields,omitempty"`
	// AgedRun is the Run that last found no Finding for this Output. It is stamped so the
	// "not surfaced again" comment is made once, not on every Run that follows.
	AgedRun string `json:"aged_run,omitempty"`
	// Deleted records that a human removed this Output. It is a Suppression, not a gap to
	// fill: the Finding is never re-filed under this key.
	Deleted bool `json:"deleted,omitempty"`
}

// Action is a Gated Action that reached a decision: applied, or suppressed by a rejection.
// Recorded even though an email has no identity to steward, because every later Run asks
// "have we already told them?" — and because an audit trail is half the reason to gate.
type Action struct {
	Key       string   `json:"finding_key"`
	Kind      string   `json:"kind"`  // email
	State     string   `json:"state"` // applied | suppressed
	Hash      string   `json:"hash"`
	At        string   `json:"at"`
	To        []string `json:"to,omitempty"`
	Subject   string   `json:"subject,omitempty"`
	DecidedBy string   `json:"decided_by,omitempty"`
}

type Manifest struct {
	Customer  string   `json:"customer"`
	UpdatedAt string   `json:"updated_at"`
	Outputs   []Output `json:"outputs"`
	Actions   []Action `json:"actions"`
}

func manifestPath(customerDir string) string {
	return filepath.Join(customerDir, "state", "manifest.json")
}

// LoadManifest reads the Manifest, returning an empty one for a Customer that has never run.
func LoadManifest(customerDir, customer string) (*Manifest, error) {
	b, err := os.ReadFile(manifestPath(customerDir))
	if errors.Is(err, fs.ErrNotExist) {
		return &Manifest{Customer: customer}, nil
	}
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	if m.Customer == "" {
		m.Customer = customer
	}
	return &m, nil
}

// Save rewrites the Manifest whole.
// ponytail: last writer wins. A Run and an approval in the admin UI touching one Customer at
// the same moment can lose one of the two writes; a per-Customer file lock is the upgrade if
// Runs ever overlap with a human's clicking.
func (m *Manifest) Save(customerDir string) error {
	m.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(manifestPath(customerDir)), 0o755); err != nil {
		return err
	}
	return os.WriteFile(manifestPath(customerDir), append(b, '\n'), 0o644)
}

// ByKey finds the Output a Finding Key claims, deleted entries included: a Suppression is
// found so that it can be honoured.
func (m *Manifest) ByKey(kind, key string) *Output {
	for i := range m.Outputs {
		if m.Outputs[i].Kind == kind && m.Outputs[i].Key == key {
			return &m.Outputs[i]
		}
	}
	return nil
}

func (m *Manifest) ByID(id string) *Output {
	for i := range m.Outputs {
		if m.Outputs[i].ID == id {
			return &m.Outputs[i]
		}
	}
	return nil
}

func (m *Manifest) Put(o Output) *Output {
	if existing := m.ByID(o.ID); existing != nil && o.ID != "" {
		*existing = o
		return existing
	}
	m.Outputs = append(m.Outputs, o)
	return &m.Outputs[len(m.Outputs)-1]
}

// Suppressed is true when a human has already rejected this action for this Finding. A
// rejection never expires on its own; it is revoked by the human who made it.
func (m *Manifest) Suppressed(key, kind string) bool {
	for _, a := range m.Actions {
		if a.Key == key && a.Kind == kind && a.State == "suppressed" {
			return true
		}
	}
	return false
}

// Applied is true when this exact artifact has already been sent.
func (m *Manifest) Applied(hash string) bool {
	return slices.ContainsFunc(m.Actions, func(a Action) bool { return a.Hash == hash && a.State == "applied" })
}

func (m *Manifest) Record(a Action) {
	a.At = time.Now().UTC().Format(time.RFC3339)
	m.Actions = append(m.Actions, a)
}

// SharesRecord decides whether a Finding is the same Finding as the one stored under this
// key. The Agent proposes the key; it is honoured only when the new Finding still stands on
// a record the old one cited (#10). A key alone is a name, not an identity.
func (o *Output) SharesRecord(records []string) bool {
	for _, r := range records {
		if slices.Contains(o.Records, r) {
			return true
		}
	}
	return false
}
