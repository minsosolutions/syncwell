package collect

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/minsosolutions/syncwell/internal/routing"
)

func transcriptConfig() *routing.Config {
	c := &routing.Config{Customers: []routing.Customer{
		{Slug: "nordstad", Domains: []string{"nordstad.se"}},
		{Slug: "ebbahus", Domains: []string{"ebbahus.se"}},
	}}
	c.Vendor.Domain = "minso.se"
	return c
}

func drop(t *testing.T, dataDir, name, attendees string) {
	t.Helper()
	body := "---\ntitle: Test\ndate: 2026-08-26\nattendees: " + attendees + "\n---\n\n**Anna Lind:** hello.\n"
	path := filepath.Join(dataDir, "sources", "transcripts", name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestTranscriptsRoutesAndQuarantines(t *testing.T) {
	dataDir := t.TempDir()
	drop(t, dataDir, "2026-08-26-nordstad-styrgrupp.md", "anna.lind@nordstad.se, robin@minso.se")
	drop(t, dataDir, "2026-09-01-internt.md", "robin@minso.se, mia.berg@minso.se")
	// Neither of these is a Source Record; the glob must not try to route them.
	drop(t, dataDir, "README.md", "anna.lind@nordstad.se")
	drop(t, dataDir, "notes.md", "anna.lind@nordstad.se")

	rep, err := Transcripts(transcriptConfig(), dataDir)
	if err != nil {
		t.Fatal(err)
	}

	routed := filepath.Join(dataDir, "customers", "nordstad", "transcripts", "2026-08-26-nordstad-styrgrupp.md")
	got, err := os.ReadFile(routed)
	if err != nil {
		t.Fatalf("transcript did not reach the customer directory: %v", err)
	}
	if !strings.Contains(string(got), "attendees: anna.lind@nordstad.se") {
		t.Fatal("attendee list lost in transit; it is the routing key")
	}

	if _, err := os.Stat(filepath.Join(dataDir, "unrouted", "transcripts", "2026-09-01-internt.md")); err != nil {
		t.Fatalf("internal meeting was not quarantined: %v", err)
	}
	reason, err := os.ReadFile(filepath.Join(dataDir, "unrouted", "transcripts", "2026-09-01-internt.reason.txt"))
	if err != nil {
		t.Fatalf("quarantined without a stated reason: %v", err)
	}
	if len(strings.TrimSpace(string(reason))) == 0 {
		t.Fatal("quarantine reason is empty")
	}

	for _, stray := range []string{"README.md", "notes.md"} {
		for _, dir := range []string{filepath.Join("customers", "nordstad", "transcripts"), filepath.Join("unrouted", "transcripts")} {
			if _, err := os.Stat(filepath.Join(dataDir, dir, stray)); err == nil {
				t.Fatalf("%s was collected as a transcript", stray)
			}
		}
	}

	if len(rep.Routed) != 1 || len(rep.Quarantined) != 1 {
		t.Fatalf("report says routed=%v quarantined=%v", rep.Routed, rep.Quarantined)
	}
}

// A transcript whose frontmatter cannot be read has no routing key. It must quarantine, not
// look like an internal meeting and not silently reach a Customer.
func TestTranscriptsWithoutUsableFrontmatter(t *testing.T) {
	dataDir := t.TempDir()
	write := func(name, body string) {
		t.Helper()
		path := filepath.Join(dataDir, "sources", "transcripts", name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("2026-08-26-no-attendees.md", "---\ntitle: Test\ndate: 2026-08-26\n---\n\n**Anna:** hello.\n")
	// The word only appears in what someone said. Frontmatter ends at the closing ---.
	write("2026-08-27-spoken.md", "---\ntitle: Test\n---\n\n**Anna:** attendees: anna.lind@nordstad.se were all there.\n")

	rep, err := Transcripts(transcriptConfig(), dataDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Routed) != 0 {
		t.Fatalf("routed %v on unreadable frontmatter", rep.Routed)
	}
	if len(rep.Quarantined) != 2 {
		t.Fatalf("quarantined %v, want both files", rep.Quarantined)
	}
	// The reason is the whole value of Quarantine: "no attendee list" is a broken export to
	// go fix, "vendor-internal meeting" is a working export nobody needs to look at.
	for _, line := range rep.Quarantined {
		if !strings.Contains(line, "attendees") || strings.Contains(line, "vendor-internal") {
			t.Fatalf("misreported as an internal meeting: %s", line)
		}
	}
}
