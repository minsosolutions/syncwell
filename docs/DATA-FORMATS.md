# Data formats

`data/sources/` is the multi-tenant backing store — what the Vendor sees before any routing.
`data/customers/<slug>/` is Collector output: one Customer's records only.
`data/unrouted/` is Quarantine.

## Slack (`data/sources/slack/`) — mirrors a real Slack export

- `users.json` — `[{id, name, real_name, is_bot, profile:{email, title}}]`
- `channels.json` — `[{id, name, created, is_shared, topic:{value}, purpose:{value}, members:[user_id]}]`
- `<channel-name>/<YYYY-MM-DD>.json` — `[{type:"message", user, text, ts, thread_ts?, reply_count?}]`

`ts` is a Unix epoch with microseconds as a string, e.g. `"1755846000.001200"`. Threaded replies carry
`thread_ts` equal to the parent's `ts` and live in the same daily file.

## Zammad (`data/sources/zammad/`) — mirrors the Zammad REST objects

- `organizations.json` — `[{id, name, domain, domain_assignment, note, active}]`
- `users.json` — `[{id, email, firstname, lastname, organization_id, active}]`
- `groups.json` — `[{id, name}]`
- `tickets/<id>.json` — `{id, number, title, group_id, state, priority, customer_id, organization_id,
  created_at, updated_at, close_at, articles:[{id, ticket_id, from, to, subject, body, content_type,
  type, sender, internal, created_at}]}`

`organization_id` may be `null` — that is the Quarantine case.

## Transcripts (`data/sources/transcripts/`) — file drop, no API

`<YYYY-MM-DD>-<slug>.md`, flat key/value frontmatter (no nesting, comma-separated lists):

```
---
title: Monthly steering — Nordstad kommun
date: 2026-08-26
duration_minutes: 45
source: granola
attendees: anna.lind@nordstad.se, per.ek@nordstad.se, robin@minso.se
---

**Anna Lind:** ...
```

## Email (`data/sources/email/`) — file drop, already converted to markdown

`<YYYY-MM-DD>-<slug>.md`:

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

Body in markdown.
```

## Linear (`data/linear/`) — written through the mock GraphQL API

- `projects.json`, `issues.json`

An Issue created by a Run carries a Provenance Marker as an HTML comment at the end of its
description: `<!-- syncwell:run=<ISO8601> customer=<slug> item=<id> -->`

## Customer Directory (`data/customers/<slug>/`)

```
slack/  zammad/  transcripts/  email/     # Collector output, same formats as above
state/manifest.json                        # what the last Run created
outputs/                                   # on-disk Outputs (report.md, *.xlsx)
```
