# Source: Zammad

Zammad is the support desk. Tickets come from customer staff — a grants manager, an IT
coordinator — not from grant applicants. Applicants appear only as the subject of a ticket
("applicant 2026-0412 reports…"), never as the author.

## Routing rule

**Organization, falling back to email domain.** Zammad has a native `organization_id` on both
tickets and users, which is the right thing to route on when it is set. It is not always set:
someone writes in from a personal address and the desk never assigns them, and that ticket
arrives with `organization_id: null`.

So: use `organization_id` when present, otherwise map the requester's email domain via
`config/syncwell.json`, otherwise quarantine.

Two cases in this data are worth your attention:

- **Utvecklingsverket uses two domains.** Zammad's `domain` field holds one value per
  organization, but their architect writes from `ruv.se`. The config knows about both; Zammad
  does not. That asymmetry is why the fallback exists.
- **One ticket has no organization at all** and is obviously, in its text, from a customer you
  can name. The deterministic rule cannot see that. Quarantine it anyway — then decide what a
  human should do with the quarantine queue, because that is the real design question.

## The mock

| | |
| --- | --- |
| `GET /zammad/api/v1/tickets` | `?page=&per_page=` — **no articles** |
| `GET /zammad/api/v1/tickets/{id}` | full ticket, articles embedded |
| `GET /zammad/api/v1/ticket_articles/by_ticket/{id}` | the articles |
| `GET /zammad/api/v1/organizations` · `/users` · `/groups` | `?page=&per_page=` |

Send `Authorization: Bearer <anything>`. Responses are bare JSON arrays, as Zammad's are.

The list endpoint deliberately omits `articles`, exactly as the real index action does — the
ticket title tells you almost nothing, and the conversation is where the content is. Fetching
30 tickets means 31 requests. Page size defaults to **10** and maxes at 100, which is Zammad's
real ceiling.

Deviations: no search endpoint, no `expand=true`, no state or group filters, no attachments,
no triggers or webhooks.

## Writing the real thing

- [API intro and authentication](https://docs.zammad.org/en/latest/api/intro.html) — HTTP token
  auth, `Authorization: Token token=<token>`, not `Bearer`
- [Ticket API](https://docs.zammad.org/en/latest/api/ticket/index.html) ·
  [Ticket articles](https://docs.zammad.org/en/latest/api/ticket/articles.html)
- [Organizations](https://docs.zammad.org/en/latest/api/organization.html) ·
  [Users](https://docs.zammad.org/en/latest/api/user.html)
- [Search](https://docs.zammad.org/en/latest/api/ticket/index.html#search) — for a real
  collector, `?query=updated_at:>now-1d` beats paging everything every run

Real-world notes: article bodies may be `text/html` and need stripping before a model reads
them; `internal: true` articles are agent-only notes that the customer never saw, and whether
those belong in a customer-facing Output is a decision you should make on purpose.

## On disk

`organizations.json`, `users.json`, `groups.json`, `tickets/<id>.json` with articles embedded.
See [docs/DATA-FORMATS.md](../../../docs/DATA-FORMATS.md).
