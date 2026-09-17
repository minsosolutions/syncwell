package collect

import (
	"fmt"
	netmail "net/mail"

	"github.com/minsosolutions/syncwell/internal/routing"
)

// Email collects the mail file drop — already converted from MIME to markdown — and routes
// each message on its peers: from, to and cc together, with the Vendor's own domain stripped.
// cc is the one that leaks: a thread copied to two Customers belongs to neither.
func Email(cfg *routing.Config, dataDir string) (Report, error) {
	return collectDrop(dataDir, "email", func(body string) (*routing.Customer, error) {
		peers, err := peers(body)
		if err != nil {
			return nil, err
		}
		return cfg.ByPeerDomains(peers)
	})
}

// peers reads every address off the from/to/cc headers. Addresses arrive in display-name form
// ("Anna Lind <anna.lind@nordstad.se>"), so net/mail does the parsing. A header we cannot
// parse quarantines rather than routing on the addresses we did manage to read: a half-read
// recipient list is how a Customer ends up in someone else's directory.
func peers(body string) ([]string, error) {
	var out []string
	for _, key := range []string{"from", "to", "cc"} {
		raw := frontmatter(body, key)
		if raw == "" {
			continue
		}
		list, err := netmail.ParseAddressList(raw)
		if err != nil {
			return nil, &routing.Unroutable{Reason: fmt.Sprintf("header %q is unparseable: %v", key, err)}
		}
		for _, a := range list {
			out = append(out, a.Address)
		}
	}
	if len(out) == 0 {
		return nil, &routing.Unroutable{Reason: "frontmatter has no from/to/cc; nothing to route on"}
	}
	return out, nil
}
