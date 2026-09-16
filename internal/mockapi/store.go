package mockapi

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// Store reads the JSON backing files on every request rather than caching them, so a
// facilitator can edit data/ by hand mid-workshop and see it immediately.
// ponytail: whole-file read/write under one mutex; fine for a few thousand records.
type Store struct {
	dir string
	mu  sync.Mutex
}

func NewStore(dir string) *Store { return &Store{dir: dir} }

type obj = map[string]any

// readList reads a file holding a JSON array of objects.
func (s *Store) readList(rel string) ([]obj, error) {
	b, err := os.ReadFile(filepath.Join(s.dir, rel))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var out []obj
	return out, json.Unmarshal(b, &out)
}

// readObjectDir reads a directory of *.json files, each holding a single object.
func (s *Store) readObjectDir(rel string) ([]obj, error) {
	paths, err := filepath.Glob(filepath.Join(s.dir, rel, "*.json"))
	if err != nil {
		return nil, err
	}
	sort.Strings(paths)
	out := make([]obj, 0, len(paths))
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		var o obj
		if err := json.Unmarshal(b, &o); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, nil
}

// readGlobConcat reads every *.json array under a directory and concatenates them in
// filename order. Slack stores one file per day, so filename order is date order.
func (s *Store) readGlobConcat(rel string) ([]obj, error) {
	paths, err := filepath.Glob(filepath.Join(s.dir, rel, "*.json"))
	if err != nil {
		return nil, err
	}
	sort.Strings(paths)
	var out []obj
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		var day []obj
		if err := json.Unmarshal(b, &day); err != nil {
			return nil, err
		}
		out = append(out, day...)
	}
	return out, nil
}

func (s *Store) writeList(rel string, v []obj) error {
	path := filepath.Join(s.dir, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}
