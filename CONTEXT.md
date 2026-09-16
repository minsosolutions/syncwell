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

**Output**:
An artifact an Agent creates or maintains, either on disk (report, spreadsheet) or in an external system (Linear issue, email).
_Avoid_: artifact, deliverable, target

**Provenance Marker**:
A machine-readable stamp carried by an Output identifying the Run that created it.
_Avoid_: tag, fingerprint

**Manifest**:
The record a Run leaves of the Outputs it created, held in the Customer Directory.
_Avoid_: state file, ledger

**Drift**:
Divergence between what a Manifest claims exists and what the Output system actually contains — a human deleted, edited, or hand-created an Output between Runs.
_Avoid_: conflict, staleness
