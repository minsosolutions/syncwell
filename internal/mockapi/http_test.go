package mockapi

import "testing"

func items(n int) []obj {
	out := make([]obj, n)
	for i := range out {
		out[i] = obj{"i": i}
	}
	return out
}

// Walking the cursor must visit every item exactly once and stop. If this breaks,
// collectors silently lose data, which is the one bug this mock exists to provoke
// deliberately rather than by accident.
func TestCursorPageWalksEverything(t *testing.T) {
	all := items(47)
	seen := 0
	cursor := ""
	for pages := 0; ; pages++ {
		if pages > 100 {
			t.Fatal("cursor never terminated")
		}
		page, next := cursorPage(all, cursor, 10, 20, 200)
		for _, it := range page {
			if it["i"] != seen {
				t.Fatalf("page item %v out of order, want %d", it["i"], seen)
			}
			seen++
		}
		if next == "" {
			break
		}
		cursor = next
	}
	if seen != 47 {
		t.Fatalf("walked %d items, want 47", seen)
	}
}

func TestCursorPageDefaultsAndClamps(t *testing.T) {
	page, _ := cursorPage(items(100), "", 0, 20, 200)
	if len(page) != 20 {
		t.Fatalf("limit=0 gave %d, want the default 20", len(page))
	}
	page, _ = cursorPage(items(500), "", 9999, 20, 200)
	if len(page) != 200 {
		t.Fatalf("oversized limit gave %d, want the max 200", len(page))
	}
}

func TestOffsetPageIsOneBased(t *testing.T) {
	first := offsetPage(items(25), 1, 10, 10, 100)
	if len(first) != 10 || first[0]["i"] != 0 {
		t.Fatalf("page 1 = %v", first)
	}
	// Zammad clients often omit ?page=; treat that as page 1, not page 0.
	implicit := offsetPage(items(25), 0, 10, 10, 100)
	if len(implicit) != 10 || implicit[0]["i"] != 0 {
		t.Fatalf("page 0 = %v, want same as page 1", implicit)
	}
	last := offsetPage(items(25), 3, 10, 10, 100)
	if len(last) != 5 || last[0]["i"] != 20 {
		t.Fatalf("page 3 = %v", last)
	}
	if past := offsetPage(items(25), 9, 10, 10, 100); len(past) != 0 {
		t.Fatalf("page past the end = %v, want empty", past)
	}
}
