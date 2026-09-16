package collect

import "testing"

func TestSegmentRejectsPathEscapes(t *testing.T) {
	for _, bad := range []string{"../../etc/passwd", "..", "a/b", `a\b`, ".hidden", "", "cust-x\x00"} {
		if _, err := segment("channel", bad); err == nil {
			t.Fatalf("segment(%q) was accepted; it can escape the data directory", bad)
		}
	}
	for _, ok := range []string{"cust-nordstad-support", "syncwell-internt", "a.b_c-1"} {
		if _, err := segment("channel", ok); err != nil {
			t.Fatalf("segment(%q) rejected: %v", ok, err)
		}
	}
}
