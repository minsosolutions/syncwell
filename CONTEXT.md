# Syncwell

Syncwell pulls the scattered sources a project manager works from — chat, support tickets, meeting transcripts, email — into per-customer working sets, and maintains outputs from them. Customer separation is deterministic and structural, never a matter of prompting.

## Language

### Parties

**Vendor**:
Minso, the operator of Grantly. The party whose project managers Syncwell serves.
_Avoid_: supplier, us

**Customer**:
An organization that licenses Grantly. The unit of separation: every Source Record belongs to exactly one, or to none.
_Avoid_: tenant, client, account, org

**Applicant**:
A person applying for a grant through a Customer's Grantly instance. Appears only as the subject of a Source Record, never as its author.
_Avoid_: end user, applicant organization

**Grantly**:
The multi-tenant grant management platform the Vendor sells to Customers. The subject matter of every Source Record.

### Ingest

**Source**:
An external system that Source Records originate from: Slack, Zammad, meeting transcripts, email.
_Avoid_: integration, connector, channel

**Source Record**:
One item from a Source — a message, a ticket, a transcript, an email.
_Avoid_: document, item, event

**Collector**:
A deterministic script that pulls Source Records from one Source and writes each into exactly one Customer Directory, or into Quarantine.
_Avoid_: importer, ingester, sync

**Routing Rule**:
The deterministic mapping from a Source Record to a Customer — a Slack channel name, a Zammad organization, a sender domain. Evaluated by code, never inferred by a model.
_Avoid_: classifier, matcher, assignment

**Customer Directory**:
The on-disk directory holding every Source Record routed to one Customer. The complete and only data an Agent may read on that Customer's behalf.
_Avoid_: workspace, tenant folder

**Quarantine**:
Where a Source Record lands when its Routing Rule yields no single Customer. Records are quarantined with a stated reason, never dropped.
_Avoid_: dead letter, rejects, unknown

**Reference**:
Shared knowledge owned by no Customer — product documentation, runbooks, changelogs — readable on every Customer's behalf.
_Avoid_: knowledge base, context docs

### Acting

**Agent**:
The AI process that reads one Customer Directory plus References and maintains that Customer's Outputs.
_Avoid_: bot, assistant, LLM

**Run**:
One execution of an Agent on behalf of one Customer.
_Avoid_: job, sync, pass

**Finding**:
One thing needing a project manager's attention, held by an Agent's judgement and standing on its Evidence. Carries a key that identifies it across Runs.
_Avoid_: issue, item, insight, alert

**Finding Key**:
A Finding's identity across Runs: proposed by an Agent as a readable slug, and honoured only when the Finding still shares a cited record with the one stored under that key. Carried by the Provenance Marker and the Manifest alike.
_Avoid_: id, slug, fingerprint

**Evidence**:
What a Finding stands on: either a Source Record cited by path, locator and verbatim quote, or a stated absence — the thing that is not there, with where it was looked for and since when. An absence is verified in code, never taken on an Agent's word.
_Avoid_: citation, reference, source

**Output**:
An artifact an Agent creates or maintains, either on disk (report, spreadsheet) or in an external system (Linear issue, email).
_Avoid_: artifact, deliverable, target

**Provenance Marker**:
A machine-readable stamp carried by an Output identifying the Run that created it.
_Avoid_: tag, fingerprint

**Manifest**:
The Customer's cross-Run record of its Maintained Outputs: which ones a Run authored, whether we still steward them, and the Written Snapshot of each. Held in the Customer Directory.
_Avoid_: state file

**Written Snapshot**:
What a Run last wrote to a Maintained Output, recorded in the Manifest so a later Run can tell a human's edit from its own previous write.
_Avoid_: cached copy, last known state

**Authorship**:
The claim that a Run created an Output. A fact about the past: nothing done to the Output afterwards makes it untrue. Carried by the Manifest and the Provenance Marker alike.
_Avoid_: ownership

**Stewardship**:
The claim that Syncwell still maintains an Output. A fact about the present, revocable by a human, held per Output and carried by the Provenance Marker alone: removing the Marker withdraws it, restoring the Marker returns it.
_Avoid_: ownership, control

**Maintained Output**:
An Output with an identity that persists across Runs, reconciled rather than rewritten — a Linear issue or project.

**Snapshot Output**:
An Output that renders one Run and is regenerated whole by the next, holding no identity across Runs — a report, an email.

**Unmanaged Output**:
An artifact in a Customer's Output system that no Run authored. Seen during reconciliation and never written to, so that Syncwell does not file a duplicate beside a human's work.
_Avoid_: foreign issue, external

**Gated Action**:
An action a Run may not take unattended, because it crosses the boundary to the Customer or destroys work. Everything else is written by the Run itself.
_Avoid_: dangerous action, human-in-the-loop step

**Proposal**:
A Gated Action a Run has prepared and not taken: the exact artifact awaiting a human's approval, held in the Customer Directory. Approved, it is applied through the same writer path a Run uses; rejected, it becomes a Suppression.
_Avoid_: draft, pending action, request

**Suppression**:
A human gesture telling later Runs not to do something again — rejecting a Proposal, or deleting an Output a Run created — remembered against the Finding and action it concerned. Revocable by the human who made it.
_Avoid_: mute, ignore list, tombstone

**Drift**:
Divergence between what a Manifest claims and what the Output system actually contains — a human deleted, edited, or hand-created an Output between Runs. Detected by deterministic reconciliation, never by an Agent.
_Avoid_: conflict, staleness
