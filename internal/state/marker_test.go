package state

import "testing"

func TestStampIsIdempotentAndReadable(t *testing.T) {
	m := Marker{Run: "2026-09-17T06:00:00Z", Customer: "ebbahus", Item: "export-off"}
	body := Stamp("Nightly export has been off since 5 August.", m)
	body = Stamp(body, Marker{Run: "2026-09-18T06:00:00Z", Customer: "ebbahus", Item: "export-off"})

	if got := Prose(body); got != "Nightly export has been off since 5 August." {
		t.Fatalf("prose changed under re-stamping: %q", got)
	}
	got, ok := FindMarker(body)
	if !ok || got.Run != "2026-09-18T06:00:00Z" || got.Item != "export-off" {
		t.Fatalf("second stamp did not replace the first: %+v ok=%v", got, ok)
	}
}

func TestStrippedMarkerIsNotFound(t *testing.T) {
	// The seeded SYN-31 case: a human edited the description and took our Marker with it.
	if _, ok := FindMarker("Decision needed on legacy case IDs. Edited by HN."); ok {
		t.Fatal("found a marker in a body that has none")
	}
}
