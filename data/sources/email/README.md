# Source: email

Mail between Minso and customer staff, already converted from MIME to markdown. Like
transcripts, this is a **file drop** rather than a live API.

Email is where things get asked directly and formally — a deadline, a written answer someone
needs for their own governance. It is also where things go to die: a question arrives, nobody
picks it up, and there is no ticket or channel message to show for it. Absence is signal here,
and absence is hard to see in any single source.

## Routing rule

**Peer email domain.** Take the addresses in `from`, `to` and `cc`, drop the vendor's own
domain (`minso.se`), and map what remains through `config/syncwell.json`. Exactly one customer
routes; zero or two quarantines.

Two things in this data will exercise it:

- A customer contact writing from a **personal address**. The rule cannot route it, and it
  should not try. Quarantine and move on — then decide what a human does with that queue.
- A customer whose people use **two domains**. The config knows both. This is why the rule maps
  domains through configuration rather than assuming one domain per customer.

Note that `cc` can pull two customers into one thread. Handle it, or you will leak.

## The mock

No API. Files in this directory, read from disk. Frontmatter is flat key/value:

```
---
from: Anna Lind <anna.lind@nordstad.se>
to: robin@minso.se
cc:
subject: Retention for rejected applications
date: 2026-09-11T10:22:00+02:00
message_id: <a1b2@nordstad.se>
in_reply_to:
---
```

Threads are reconstructed from `message_id` and `in_reply_to`. A thread with no reply is a
legitimate and interesting state.

## Sending email

Outbound goes to **Mailpit**, which speaks real SMTP on `localhost:1025` and shows you the
resulting inbox at <http://localhost:8025>. Nothing leaves the machine.

```
SMTP host localhost, port 1025, no auth, no TLS
```

## Writing the real thing

- [Gmail API: `users.messages.list`](https://developers.google.com/gmail/api/reference/rest/v1/users.messages/list)
  — `q` supports Gmail search syntax; use `historyId` for incremental sync rather than
  re-listing
- [Microsoft Graph: `messages`](https://learn.microsoft.com/en-us/graph/api/user-list-messages)
  — and [delta query](https://learn.microsoft.com/en-us/graph/delta-query-messages) for
  incremental sync
- [IMAP RFC 3501](https://datatracker.ietf.org/doc/html/rfc3501) — if you are stuck with it,
  use UIDs and `UIDVALIDITY`, never sequence numbers
- [JMAP](https://jmap.io/) — if your provider supports it, easier than both

Practical notes: real mail is MIME, so you will be picking `text/plain` over a
`multipart/alternative`, or converting HTML. Quoted reply chains mean the same text appears in
every message of a thread — strip it or the Agent reads the same complaint five times. Threading
by `References`/`In-Reply-To` is more reliable than by subject.

## On disk

`<YYYY-MM-DD>-<slug>.md` — glob that pattern rather than `*.md`, or your collector will try to route this README. See [docs/DATA-FORMATS.md](../../../docs/DATA-FORMATS.md).
