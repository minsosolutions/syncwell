# Source: Slack

Customer conversation happens in **Slack Connect** channels — one shared channel per customer
per topic, with Minso people and that customer's people in it. Internal Minso channels live in
the same workspace and must not be collected as customer data.

## Routing rule

**Channel name prefix.** A channel named `cust-<slug>-<topic>` belongs to that customer;
`config/syncwell.json` holds the prefixes. Anything without a matching prefix is quarantined.

This is a naming convention, which means it is only as good as the discipline that maintains
it. It is still the right rule: it is visible, it is checkable, and when someone creates
`#cust-nordstsd-support` the record quarantines loudly rather than landing somewhere wrong.

The quarantine case worth looking at is `#syncwell-internt`. Read what is in there — a message
can be about two customers at once, and no rule should resolve that to one of them.

## The mock

| | |
| --- | --- |
| `GET /slack/api/conversations.list` | `?cursor=&limit=` |
| `GET /slack/api/conversations.history` | `?channel=&cursor=&limit=` — newest first |
| `GET /slack/api/users.list` | `?cursor=&limit=` |
| `POST /slack/api/chat.postMessage` | `{"channel":"...","text":"..."}` |

Send `Authorization: Bearer <anything>`. Responses carry Slack's real envelope:
`{"ok": true, ...}` on success, `{"ok": false, "error": "channel_not_found"}` on failure.

**Pagination is where this bites.** `response_metadata.next_cursor` is an opaque string, and
it is empty only on the last page. The default `limit` here is **20** — smaller than Slack's
real default, deliberately — so a collector that reads one page and stops gets most of a
channel and none of the warning. `internal/collect/slack.go` walks it properly; copy that loop.

Deviations from real Slack: no rate limits, no `oldest`/`latest` filters, no scopes, no
webhooks or Events API, and any non-empty token works.

## Writing the real thing

- [`conversations.list`](https://api.slack.com/methods/conversations.list) ·
  [`conversations.history`](https://api.slack.com/methods/conversations.history) ·
  [`users.list`](https://api.slack.com/methods/users.list) ·
  [`chat.postMessage`](https://api.slack.com/methods/chat.postMessage)
- [Pagination](https://api.slack.com/docs/pagination) — cursor-based, same shape as here
- [Rate limits](https://api.slack.com/docs/rate-limits) — tiered per method; you will need backoff
- [Slack Connect](https://api.slack.com/apis/connect) — shared channels behave differently from internal ones
- Scopes for a bot token: `channels:read`, `channels:history`, `users:read`,
  and `chat:write` if you post. Shared channels also need `groups:read` in some setups.

Two things the mock does not teach you. Real message text is
[mrkdwn](https://api.slack.com/reference/surfaces/formatting), where user mentions arrive as
`<@U123>` and need resolving against `users.list`. And thread replies are **not** returned by
`conversations.history` — you need
[`conversations.replies`](https://api.slack.com/methods/conversations.replies) per parent, so
budget a call per thread.

## On disk

`channels.json`, `users.json`, and `<channel-name>/<YYYY-MM-DD>.json` — one file per day, the
shape a real Slack export uses. See [docs/DATA-FORMATS.md](../../../docs/DATA-FORMATS.md).
