# Syncwell — workshop kit

Minso runs **Grantly**, a grant management platform, for three customers: a municipality, a
private foundation and a national agency. One project manager covers all three.

What that person needs to know is never in one place. A customer mentions a date change in a
Tuesday meeting and nobody writes it down. A ticket gets closed as solved and three weeks
later someone says in Slack that it still happens. An email asks a direct question and no
answer is ever sent. Each source on its own looks fine. The problem only exists **between**
them.

Build something that reads across sources and surfaces what needs attention — and then does
something about it.

This repo is a **starting point**, not a solution. It gives you realistic data, a running
mock of every external system, and one worked example of each half of the problem. What you
build on top is the point of the afternoon.

## What "done" looks like in 2.5 hours

You will not build all of it. Done, for this afternoon, is a thing that genuinely works
across **two or three** sources rather than a sketch that covers five.

Concretely, by the end you can run it in front of the room and it:

1. **Collects** two or three sources into per-customer directories, with every record either
   routed to exactly one customer or quarantined with a stated reason. Nothing is dropped,
   and nothing is guessed.
2. **Surfaces at least three things that need attention** for one customer, each with
   evidence naming the specific records it came from — and **at least one of them must rest
   on two or more different sources**. A finding visible in a single source does not count;
   that is the part that does not need you.
3. **Writes an Output** — a markdown report at minimum, and for at least one item a Linear
   issue created through the mock API.
4. **Has an answer for the second run.** Run it twice. Say what happens to what the first run
   created, and why that is the behaviour you chose.

There are specific things planted in this data, findable only by crossing sources. Your
facilitator knows what they are and will tell you how many you found.

## Where to focus

**Spend your time on the reading-across, not on the plumbing.** The mock API, the data, one
collector and one output writer are already here so that you do not have to build them.

Worth your time:

- **The routing rule.** It is four lines of code and it is the whole isolation story. Look at
  what lands in `data/unrouted/` and decide what you would actually do about it.
- **What the Agent is given.** One customer's directory, plus `data/references/`. Does it read
  everything every run? Does it need to?
- **Evidence.** A finding without a traceable source is a guess with confidence. Make the
  Agent cite records.
- **The second run.** There is already a run from 2026-09-09 in here, and a human has edited
  Linear since. That is where the interesting problem is.

Not worth your time today:

- Authentication, retries, rate limits, webhooks. The mock has none of these on purpose.
- A spreadsheet writer. Markdown is enough to judge the output.
- Covering all four sources shallowly.
- The cross-customer view — that is the stretch goal, and only after one customer works.

## A rough shape for the time

| | |
| --- | --- |
| 0:00–0:20 | Get it running. Read some of the actual data — pick one customer and skim a week of it. |
| 0:20–1:00 | Collectors for your chosen sources. Copy the Slack one. |
| 1:00–2:00 | The Agent, and one Output. |
| 2:00–2:20 | Run it a second time. Deal with what happens. |
| 2:20–2:30 | Show the room. |

## Ground rules

- **Separation by customer is deterministic.** Decided by code over stable metadata, never by
  a prompt and never by a model. An Agent gets one directory.
- **Quarantine, do not guess.** Anything that does not resolve to exactly one customer goes to
  `data/unrouted/` with a reason.
- **Cite your evidence.** Every finding names the records behind it.

Everything else is yours. If you think the shape sketched below is wrong, build it
differently and tell us why — that is a better outcome than agreeing with us.

## Quick start

```bash
docker compose up          # mock API on :8099, Mailpit on :8025
```

If a port is already taken on your machine, drop a `.env` next to `docker-compose.yml` with
`MOCKAPI_PORT=9099` (or `MAILPIT_SMTP_PORT` / `MAILPIT_UI_PORT`) rather than editing the compose file.

No Docker? The mock API is a single Go binary with no dependencies:

```bash
go run ./cmd/mockapi        # serves :8099 from ./data
```

Then collect Slack into per-customer directories and look at what happened:

```bash
go run ./cmd/collect -source slack
```

```
routed      cust-nordstad-support -> nordstad (18 messages)
...
quarantined syncwell-internt (8 messages): channel "syncwell-internt" has no customer prefix
```

Every endpoint needs `Authorization: Bearer <anything>`:

```bash
curl -H 'Authorization: Bearer dev' 'localhost:8099/zammad/api/v1/tickets?page=1&per_page=25'
curl -H 'Authorization: Bearer dev' localhost:8099/        # route index
```

## The shape

```
 Mock API / file drops          Collector              Customer Directory        Agent          Outputs
 ─────────────────────          ─────────              ──────────────────        ─────          ───────
 all customers, mixed    →   deterministic routing  →   one customer only    →   reads one  →   report.md
                             channel / org / domain                              directory      Linear issue
                                     ↓                                                          spreadsheet
                                 unrouted/                                                      email
```

Three kinds of thing, and it is worth keeping them apart:

- **Collectors** pull from a Source and write into exactly one Customer Directory. They are
  ordinary deterministic scripts. No model runs here.
- **References** (`data/references/`) are shared product documentation, owned by no customer
  and readable on every customer's behalf.
- **Outputs** are what the Agent produces or maintains — a markdown report, a spreadsheet, a
  Linear issue, an email. Some are internal, some a customer sees.

### The one hard constraint

**Customer separation is structural, not prompted.** A Collector decides which customer a
record belongs to using deterministic code over stable metadata — a Slack channel name, a
Zammad organization, a sender domain. The Agent is then handed a single directory and has no
path to any other. If separation depends on an instruction in a prompt, it is not separation.

A record that does not resolve to exactly **one** customer goes to `data/unrouted/` with a
reason. It is never dropped, and never guessed at. See
[ADR-0001](./docs/adr/0001-isolation-lives-in-the-routing-rule.md) for why it works this way.

## What ships, and what you write

| | Ships | You write |
| --- | --- | --- |
| Sources | Slack, Zammad (HTTP), transcripts, email (file drops) | — |
| Collectors | Slack — [`internal/collect/slack.go`](./internal/collect/slack.go) | Zammad, transcripts, email |
| References | 14 documents | — |
| Outputs | Markdown — [`internal/output/markdown.go`](./internal/output/markdown.go) | Spreadsheet, Linear, email |
| The Agent | nothing | all of it |

Each source directory has a `README.md` with the real vendor's API docs and the routing rule
for that source. Copy the Slack collector's shape; it is about 150 lines.

The reference Output writer takes JSON on stdin and renders markdown, which keeps the Agent's
judgement separate from the rendering:

```bash
echo '{"customer":"nordstad","run_at":"2026-09-16T09:00:00Z","items":[
  {"title":"Launch date moved","severity":"high","summary":"...","sources":["transcript:2026-08-26"]}]}' \
  | go run ./cmd/report
```

## Things we have deliberately not decided

These are genuinely open. We have opinions, not answers, and the data is arranged so you
will run into all of them.

1. **How does the Agent know what it created?** There is already a run from 2026-09-09 in
   here. Its Linear issues carry a marker in the description, and each customer has a
   `state/manifest.json` listing what that run made. **Since then a human has been in Linear.**
   The two signals do not agree. Which one do you trust, and what do you do when they clash?
2. **What happens to a human's edits?** If someone reprioritises an issue you created, does
   the next run respect it, revert it, or notice and ask?
3. **Where does a human have to approve?** Some actions are safe to take. Emailing a customer
   probably is not. Where is the line, and how does the Agent wait?
4. **How much configuration is real?** `config/syncwell.json` holds the routing rules today.
   Per-source prompts, per-customer prompts, per-customer tone — we suspect these are needed,
   but we have not built them and might be wrong.
5. **Cross-customer work.** A second Agent reading every customer's Outputs to find common
   themes. Strictly a stretch goal, and the isolation constraint above is exactly what makes
   it a separate thing rather than a flag on the first Agent.

## Layout

```
cmd/mockapi         mock Slack, Zammad and Linear APIs
cmd/collect         reference Slack collector
cmd/report          reference markdown output writer
config/             customers, domains, routing rules
data/sources/       what the vendor sees: all customers, mixed together
data/references/    shared product documentation
data/customers/     collector output, one directory per customer, plus state/ and outputs/
data/unrouted/      quarantine
data/linear/        the mock Linear workspace
docs/DATA-FORMATS.md   exact shape of every file
CONTEXT.md          the vocabulary; use these words
```

The data is invented. Nothing in it comes from a real customer.
