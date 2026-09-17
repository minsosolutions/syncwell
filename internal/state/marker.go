// Package state holds what a Customer Directory remembers between Runs: the Manifest, the
// Provenance Marker that pairs with it, and the Proposals awaiting a human.
//
// Nothing here consults a model. The Manifest and the Marker are the two traces a Run leaves
// and they are read with different rules — see
// docs/adr/0002-authorship-and-stewardship-are-separate-signals.md.
package state

import (
	"fmt"
	"regexp"
	"strings"
)

// Marker is the Provenance Marker an Output carries: the Run that wrote it, the Customer it
// belongs to, and the Finding Key it stands for.
type Marker struct {
	Run      string
	Customer string
	Item     string
}

func (m Marker) String() string {
	return fmt.Sprintf("<!-- syncwell:run=%s customer=%s item=%s -->", m.Run, m.Customer, m.Item)
}

var markerRE = regexp.MustCompile(`<!--\s*syncwell:run=(\S+)\s+customer=(\S+)\s+item=(\S+)\s*-->`)

// FindMarker reads the Marker out of an Output's body. Its absence is a claim in its own
// right: a human who strips it has withdrawn Stewardship.
func FindMarker(body string) (Marker, bool) {
	m := markerRE.FindStringSubmatch(body)
	if m == nil {
		return Marker{}, false
	}
	return Marker{Run: m[1], Customer: m[2], Item: m[3]}, true
}

// Stamp replaces whatever Marker a body carries with this one, leaving the prose alone.
func Stamp(body string, m Marker) string {
	body = strings.TrimRight(markerRE.ReplaceAllString(body, ""), " \t\n")
	return body + "\n\n" + m.String()
}

// Prose is an Output's body without its Marker, which is what a human's edit is compared
// against: re-stamping must not read as a human rewrite.
func Prose(body string) string {
	return strings.TrimRight(markerRE.ReplaceAllString(body, ""), " \t\n")
}
