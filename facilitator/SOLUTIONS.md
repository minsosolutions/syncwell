# Facilitator notes — what is seeded in the data

**Do not hand this to participants.** It is the answer key and the generation spec.

Window: **2026-07-22 → 2026-09-16**. "Today" is 2026-09-16. Last Run was **2026-09-09**.

## Cast

**Vendor — Minso** (`minso.se`), sells **Grantly**, a multi-tenant grant management platform.

| Person | Email | Role |
| --- | --- | --- |
| Robin Bauhn | robin@minso.se | Project manager (the Syncwell user) |
| Sara Holm | sara.holm@minso.se | Support lead |
| Jonas Frid | jonas.frid@minso.se | Backend engineer |
| Mia Berg | mia.berg@minso.se | Customer success |

**Customers**

| Slug | Name | Character | Domains | Slack channels |
| --- | --- | --- | --- | --- |
| `nordstad` | Nordstad kommun | Large municipality; high volume, operational noise | `nordstad.se` | `#cust-nordstad-support`, `#cust-nordstad-leverans` |
| `ebbahus` | Ebbahus Stiftelse | Small private foundation; one bespoke integration, high-touch | `ebbahus.se` | `#cust-ebbahus-support` |
| `utvecklingsverket` | Regionalt Utvecklingsverket | National agency; mid-migration, compliance-heavy | `utvecklingsverket.se`, `ruv.se` | `#cust-utvecklingsverket-support`, `#cust-utvecklingsverket-migrering` |

People: Anna Lind (grants manager) `anna.lind@nordstad.se`, Per Ek (IT) `per.ek@nordstad.se`,
Karin Sjö (caseworker) `karin.sjo@nordstad.se`; Gustav Ebbe (director) `gustav.ebbe@ebbahus.se`,
Nina Dahl (program officer) `nina.dahl@ebbahus.se`; Håkan Nyberg (programme director)
`hakan.nyberg@utvecklingsverket.se`, Elin Sund (DPO) `elin.sund@utvecklingsverket.se`,
Tobias Ahl (architect) `tobias.ahl@ruv.se` — note the second domain, a deliberate routing test.

Internal Vendor-only channel: `#syncwell-internt` (no `cust-` prefix).

## Seeded at-risk items

Each is visible **only by crossing two or more Sources**. An Agent reading one Source cannot find them.

### Nordstad kommun

**N1 — Launch date moved; the plan still says the old one.**
Transcript `2026-08-26` steering: Anna says the public launch moved **Oct 1 → Sep 22** to match a
council decision. The Linear project from the last Run still reads Oct 1, and Zammad ticket on data
migration is scheduled for **Sep 25** — after the new launch.
_Sources: transcript + Linear/manifest + Zammad._

**N2 — A "solved" ticket that was never solved.**
Zammad ticket **#10042** (bulk export times out) closed `solved` on 2026-08-19. On 2026-09-08 in
`#cust-nordstad-support` Karin writes that the export still times out and they now "just run it at
night". Nobody reopened it.
_Sources: Zammad state + Slack._

**N3 — An escalation with no answer anywhere.**
Email from Anna, 2026-09-11, asking for a **written** answer on GDPR retention for rejected
applications, flagging a council deadline. There is no reply in email, Slack or Zammad.
_Sources: email + the **absence** of a response elsewhere — the hardest one to detect._

### Ebbahus Stiftelse

**E1 — Bespoke integration silently switched off.**
`#cust-ebbahus-support`, 2026-08-05: Jonas disables the nightly bookkeeping export job "while
debugging". Never re-enabled. Email from Nina, 2026-09-09: September disbursements are missing from
their ledger. Nobody connects the two.
_Sources: Slack + email._

**E2 — Renewal at risk, tracked nowhere.**
Transcript `2026-09-02`: Gustav mentions the board "will review whether to continue" at the October
board meeting. A Zammad ticket complains about response times. No renewal item exists in any Output.
_Sources: transcript + Zammad._

**E3 — Upload defect, blamed on themselves.** See the cross-customer theme below.

### Regionalt Utvecklingsverket

**U1 — Migration blocked on a decision nobody owns.**
`#cust-utvecklingsverket-migrering`, 2026-08-20: Tobias asks whether legacy case IDs are preserved.
Transcript `2026-08-28` raises it and defers it "to next time". Silence since. The migration cannot
proceed without it.
_Sources: Slack + transcript._

**U2 — Compliance deadline against a stale Reference.**
Email from Elin (DPO), 2026-08-29: they need an updated processing agreement before **Oct 1** or they
must suspend processing. The Reference doc `references/legal/data-retention.md` still states a
**24-month** retention period; Zammad ticket and `references/changelog.md` show it changed to **36
months** in the July release. An Agent that cites the Reference gives the customer a wrong answer.
_Sources: email + Reference + Zammad. This is the payoff for wiring up References._

**U3 — Upload defect as an audit risk.** See below.

## Cross-customer theme (second iteration)

One root cause — **attachment uploads over ~10 MB fail silently when the applicant's session has been
idle**, an expired token on the upload endpoint — reported in three vocabularies that share no
keywords:

- **Nordstad**: applicants phone in saying "the budget file disappeared" (framed as applicant confusion, high volume)
- **Ebbahus**: Gustav says the portal "loses documents on submit" and assumes their own PDF export is at fault
- **Utvecklingsverket**: framed as compliance — "we cannot prove the applicant submitted before the deadline"

Only an Agent reading across Customers sees these are one defect.

## Quarantine cases (seeded, 3)

| # | Record | Why it does not route |
| --- | --- | --- |
| Q1 | Email from `g.ebbe.privat@gmail.com` about the board deck | Personal domain; sender domain maps to no Customer |
| Q2 | Zammad ticket from `karin.sjo@hotmail.com` | Personal address, `organization_id` is `null` |
| Q3 | `#syncwell-internt` message comparing Nordstad's and Utvecklingsverket's migration approaches | Channel has no `cust-` prefix, and the message concerns **two** Customers |

Q3 is the important one: any rule that resolves it to a single Customer is a leak.

## Last Run (2026-09-09) and seeded Drift

Each Customer has `state/manifest.json` listing the Linear issues that Run created. Since then a human
has been in Linear:

| Customer | Issue | What the human did | Which signal breaks |
| --- | --- | --- | --- |
| `nordstad` | `SYN-14` | **Deleted** it | In the Manifest, absent from Linear |
| `ebbahus` | `SYN-21` | **Retitled** and raised priority | Marker intact, title no longer matches the Manifest |
| `utvecklingsverket` | `SYN-33` | **Hand-created** an issue | In Linear, no Marker, not in the Manifest |
| `utvecklingsverket` | `SYN-31` | Edited the description and **stripped the Marker** | Manifest claims it; the issue itself denies it |

`SYN-31` is the designed conflict: **Provenance Marker and Manifest disagree**, and the kit takes no
position on which wins. That is the open question.
