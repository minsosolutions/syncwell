# Source: meeting transcripts

Auto-generated transcripts of customer meetings — steering groups, migration syncs, ad-hoc
calls. They arrive as a **file drop**, not through an API, which makes this the most awkward
source in the kit and the most rewarding one.

Read one before you write any code. They are not minutes. People talk over each other, someone
joins late, a decision gets made verbally in half a sentence and nobody writes it down. That
last property is exactly why a transcript is worth reading across: things are true in here
that are true nowhere else.

## Routing rule

**Attendee email domain.** Frontmatter carries the attendee list; strip the vendor's own domain
(`minso.se`) and what remains should be exactly one customer's domain. Zero remaining, or two,
means quarantine.

An internal meeting where Minso discusses a customer has no customer attendee, so it will not
route. Decide deliberately whether that is right — there is real signal in those, and no
deterministic rule that safely claims it.

## The mock

There is no API. Files sit in this directory and a collector reads them from disk. That is
also how it works in production for most transcript tools: a scheduled export drops files into
object storage.

Frontmatter is deliberately flat — no nested YAML — so you can parse it with a few lines of
string handling rather than pulling in a YAML dependency:

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

## Writing the real thing

Whichever tool you actually use, the collector shape is the same: list recordings since the
last run, fetch the transcript, write it down with its attendee list intact. The attendee list
is the part you must not lose — it is the routing key.

- [Granola](https://www.granola.ai/) — no public API at time of writing; export or a shared
  folder is the practical route
- [Fireflies.ai GraphQL API](https://docs.fireflies.ai/) — `transcripts` query returns
  participants, sentences and speaker labels
- [Microsoft Graph: onlineMeetings transcripts](https://learn.microsoft.com/en-us/graph/api/resources/calltranscript)
  — Teams transcripts as VTT; attendees come from the meeting resource, not the transcript
- [Zoom Cloud Recording API](https://developers.zoom.us/docs/api/rest/reference/zoom-api/methods/#operation/recordingGet)
  — transcript as a recording file type
- [Google Meet REST API](https://developers.google.com/meet/api/guides/transcripts) —
  `conferenceRecords.transcripts`, with structured entries per speaker

Practical notes: transcripts are long, so summarising per meeting before the Agent reads them
is usually worth it — but summarise with citations back to the transcript, or you lose the
evidence trail. Speaker labels are unreliable and often just "Speaker 1". Diarisation mangles
proper nouns, so do not match on exact names.

## On disk

`<YYYY-MM-DD>-<slug>.md` — glob that pattern rather than `*.md`, or your collector will try to route this README. See [docs/DATA-FORMATS.md](../../../docs/DATA-FORMATS.md).
