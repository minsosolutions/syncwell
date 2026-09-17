package agent

import (
	"fmt"
	"strings"
)

// Prompt is what one Run asks of the Agent. It names the Customer rather than letting the
// Agent infer one: the working directory is already the answer, and an Agent asked to
// attest its own identity from the records will refuse.
func Prompt(customer, runAt string, working, references []string) string {
	return fmt.Sprintf(`You are reading the working set for one customer, %q, on behalf of the project manager who looks after them. Set "customer" to %q and "run_at" to %q in your answer.

You can read files, and that is all: there is no search tool and no shell, so every readable path is listed
for you below. Read is the only way in, and any path not listed is not yours to read.

Files collected for this customer, relative to your working directory:
%s

Shared product documentation, readable at these paths:
%s

Your job is to surface what needs attention. Expect to find several things: this is a live customer with
five weeks of unread history, and a run that reports nothing has almost certainly not looked hard enough.

What matters most is what is only visible by crossing sources — a decision made in a meeting that no ticket
reflects, a ticket closed as solved that chat says still happens, a question asked and never answered, a
change announced in one place and contradicted in another. Prefer those. A finding visible in a single
record is worth reporting only when it is serious.

Work efficiently: read files in parallel batches rather than one at a time, several per turn. Read all of
the customer's own records — there are not many. Read a reference document only to check a specific claim
against it.

How to write a finding:
- Quote evidence word for word from the file you cite, so the reader can check it. Copy the text; do not
  paraphrase it or tidy it up.
- Cite paths relative to your working directory, with a locator: the Slack ts, the email message_id, the
  transcript line number, the Zammad article id.
- If the finding is that nobody answered something, cite the record that went unanswered and set
  awaiting_reply on it. Do not assert the silence yourself; it is checked and stated for you.
- Reuse a key from state/previous-findings.json when you are surfacing the same concern again.
- Where writing to the customer is the obvious next step, draft it in proposal. Never name recipients; they
  are derived from the records. A human approves every word before anything is sent.`,
		customer, customer, runAt, list(working), list(references))
}

func list(paths []string) string {
	if len(paths) == 0 {
		return "(none)"
	}
	return "- " + strings.Join(paths, "\n- ")
}
