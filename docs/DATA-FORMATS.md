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

- `projects.json`, `issues.json`, `comments.json`

An Issue created by a Run carries a Provenance Marker as an HTML comment at the end of its
description: `<!-- syncwell:run=<ISO8601> customer=<slug> item=<Finding Key> -->`

The Marker is the only signal of **Stewardship**: a human who strips it has told us to stop
maintaining that issue, and a later Run reads that and obeys. A Run comments rather than
overwrites wherever a human has edited a field, and each comment carries its own Marker.

## Customer Directory (`data/customers/<slug>/`)

```
slack/  zammad/  transcripts/  email/     # Collector output, same formats as above
outputs/report.md                          # Snapshot Output, regenerated whole each Run
state/manifest.json                        # Maintained Outputs across Runs, and decided actions
state/previous-findings.json               # the only state file the Agent may read
state/proposals/<hash>.json                # Gated Actions awaiting, or past, a human's decision
state/runs/<run-at>.json                   # one Run's summary: findings, drift, actions
```

### `state/manifest.json`

```json
{
  "customer": "ebbahus",
  "updated_at": "2026-09-17T06:00:00Z",
  "outputs": [{
    "kind": "linear_issue",
    "key": "bookkeeping-export-disabled-since-5-august",
    "id": "iss_0021", "identifier": "SYN-21",
    "authored_run": "2026-09-09T06:00:00Z",
    "last_run": "2026-09-17T06:00:00Z",
    "records": ["email/2026-09-09-ebbahus-saknade-utbetalningar.md#<20260909114700.9964@ebbahus.se>"],
    "snapshot": {"title": "...", "description": "...", "priority": 2},
    "human_fields": ["title"],
    "aged_run": "", "deleted": false
  }],
  "actions": [{"finding_key": "...", "kind": "email", "state": "applied",
               "hash": "a3a3e4b8136b98cc", "at": "...", "to": ["nina.dahl@ebbahus.se"]}]
}
```

`snapshot` is the **Written Snapshot** — what the Run last wrote — so a later Run can tell a
human's edit from its own previous write. `records` is the Finding's identity material: a
Finding Key is honoured across Runs only when the new Finding still cites one of these.
`human_fields` are fields a human has taken; `deleted` is a Suppression. Stewardship is
deliberately **not** stored here: it lives in the Marker, in the Output itself.

`actions` is the audit trail of Gated Actions: `applied` means it was sent, `suppressed`
means a human rejected it and later Runs must not propose it again.

The Agent never reads this file — a deny rule and the prompt's file list both keep it out.

### `state/proposals/<hash>.json`

```json
{"hash": "a3a3e4b8136b98cc", "customer": "ebbahus", "finding_key": "...", "kind": "email",
 "to": ["nina.dahl@ebbahus.se"], "subject": "...", "body": "...",
 "run_at": "...", "state": "pending", "decided_at": ""}
```

`hash` is sha256 over customer, finding key, kind, subject, body and recipients: approval
binds to the artifact, so re-wording a draft produces a different Proposal and the approval
given to the old wording cannot send the new one. `to` is derived in Go from addresses that
appear in this Customer's own records and match their configured domains — never proposed by
the Agent.
