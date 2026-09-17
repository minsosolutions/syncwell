package collect

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func mail(t *testing.T, dataDir, name, from, to, cc string) {
	t.Helper()
	body := "---\nfrom: " + from + "\nto: " + to + "\ncc: " + cc +
		"\nsubject: Test\ndate: 2026-09-11T10:22:00+02:00\nmessage_id: <a1@x.se>\nin_reply_to:\n---\n\nHej.\n"
	path := filepath.Join(dataDir, "sources", "email", name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestEmailRoutesOnEveryPeer(t *testing.T) {
	dataDir := t.TempDir()
	// Display-name form, and a peer who appears only in cc.
	mail(t, dataDir, "2026-09-11-nordstad-fraga.md", "Anna Lind <anna.lind@nordstad.se>", "robin@minso.se", "per.ek@nordstad.se")
	// Q1: a customer contact writing from a personal address.
	mail(t, dataDir, "2026-09-05-styrelsematerial-fraga.md", "Gustav Ebbe <g.ebbe.privat@gmail.com>", "robin@minso.se", "")
	// cc drags a second customer into the thread.
	mail(t, dataDir, "2026-09-12-korsad-trad.md", "Anna Lind <anna.lind@nordstad.se>", "robin@minso.se", "Nina Dahl <nina.dahl@ebbahus.se>")
	mail(t, dataDir, "README.md", "Anna Lind <anna.lind@nordstad.se>", "robin@minso.se", "")

	rep, err := Email(transcriptConfig(), dataDir)
	if err != nil {
		t.Fatal(err)
	}

	routed := filepath.Join(dataDir, "customers", "nordstad", "email", "2026-09-11-nordstad-fraga.md")
	got, err := os.ReadFile(routed)
	if err != nil {
		t.Fatalf("mail did not reach the customer directory: %v", err)
	}
	if !strings.Contains(string(got), "message_id:") {
		t.Fatal("headers lost in transit; threading needs message_id")
	}

	for _, q := range []struct{ name, wants string }{
		{"2026-09-05-styrelsematerial-fraga", "gmail.com"},
		{"2026-09-12-korsad-trad", "span"},
	} {
		reason, err := os.ReadFile(filepath.Join(dataDir, "unrouted", "email", q.name+".reason.txt"))
		if err != nil {
			t.Fatalf("%s: quarantined without a stated reason: %v", q.name, err)
		}
		if !strings.Contains(string(reason), q.wants) {
			t.Fatalf("%s: reason %q does not say %q", q.name, reason, q.wants)
		}
		if _, err := os.Stat(filepath.Join(dataDir, "unrouted", "email", q.name+".md")); err != nil {
			t.Fatalf("%s: reason kept but the mail was dropped: %v", q.name, err)
		}
	}

	// The leaked-in customer must not also get a copy.
	if _, err := os.Stat(filepath.Join(dataDir, "customers", "ebbahus", "email")); err == nil {
		t.Fatal("cross-customer thread reached a customer directory")
	}

	if len(rep.Routed) != 1 || len(rep.Quarantined) != 2 {
		t.Fatalf("report says routed=%v quarantined=%v", rep.Routed, rep.Quarantined)
	}
}
