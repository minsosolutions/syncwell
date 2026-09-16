# Processing agreement

**Document:** `legal/processing-agreement.md`
**Version:** 4.0 · **Last updated:** 2026-07-14 · **Owner:** Legal, reviewed by the DPO
**Applies to:** all Grantly tenants

Every Grantly customer signs a data processing agreement (DPA) with Minso under Article 28 GDPR. This
document describes what the agreement covers, how it is versioned and how a tenant gets an updated
one. It is not itself the agreement.

## Roles

| Party | Role |
| --- | --- |
| The tenant (grant-giving organisation) | Controller |
| Minso AB | Processor |
| Minso's listed suppliers | Sub-processors |

Applicants are data subjects. Minso has no direct relationship with them and refers every data
subject request to the controller. See [Data retention](data-retention.md#erasure-requests).

## What the agreement covers

| Section | Content |
| --- | --- |
| 1 | Subject matter, duration, nature and purpose of processing |
| 2 | Categories of data subjects and personal data |
| 3 | Documented instructions and the limits on them |
| 4 | Confidentiality of Minso personnel |
| 5 | Technical and organisational measures (Annex 2) |
| 6 | Sub-processors and the objection procedure (Annex 3) |
| 7 | Assistance with data subject rights |
| 8 | Personal data breach notification |
| 9 | Deletion or return of data on termination |
| 10 | Audit and inspection rights |
| 11 | International transfers |

Retention periods are **not** set in the agreement body. They are set by the controller and default
to the product periods in [Data retention](data-retention.md); a tenant that needs different periods
records them in Annex 1 as a documented instruction.

## Template versions

| Template | Published | Status |
| --- | --- | --- |
| 3.0 | 2026-07-01 | Current |
| 2.4 | 2025-09-01 | Superseded, still in force where signed |
| 2.3 | 2025-01-15 | Superseded |

A signed agreement stays in force on the template version it was signed under until both parties sign
a new one. Minso does not unilaterally update a signed DPA.

> **Note**
> Template 3.0 restructured the annexes and moved the technical measures into Annex 2. A tenant
> moving from 2.4 to 3.0 signs a whole new agreement, not an amendment. Allow two to three weeks for
> a tenant's legal review; agencies with a formal review process routinely take longer.

## Requesting an updated agreement

1. The tenant's DPO or contract owner asks Minso Customer Success for the current template.
2. Minso sends template 3.0 with Annex 1 pre-filled from the tenant's configuration.
3. The tenant reviews, and returns any required changes to Annex 1.
4. Both parties sign electronically.
5. Minso files the signed agreement and records the template version on the tenant.

## Sub-processors

| Sub-processor | Purpose | Location |
| --- | --- | --- |
| Hosting provider | Application hosting and backups | Sweden |
| Managed database provider | Database operations | Sweden |
| Email delivery provider | Transactional email | EU |
| Virus scanning provider | Attachment scanning | EU |
| Support desk provider | Support ticket handling | EU |
| BankID provider | Applicant identification | Sweden |

The current list with company names and addresses is Annex 3 of the signed agreement. Minso notifies
controllers at least **30 days** before adding or replacing a sub-processor. A controller may object
in writing within the notice period; if the objection cannot be resolved, the controller may
terminate the affected service.

## International transfers

Grantly data is stored in Sweden. No personal data is transferred outside the EU/EEA under the
standard configuration. Where a sub-processor has support staff outside the EEA, access is governed
by standard contractual clauses and the transfer impact assessment supplied on request.

## Breach notification

Minso notifies the controller **without undue delay and at the latest within 24 hours** of becoming
aware of a personal data breach affecting that controller's data. The notification includes what is
known at the time and is followed by updates. Minso does not notify supervisory authorities or data
subjects; that is the controller's obligation. See
[Incident escalation](../runbooks/incident-escalation.md).

## Audit rights

The controller may audit once per calendar year, on 30 days' notice, and more often where required by
a supervisory authority. In the first instance Minso supplies its current third-party security report
and the completed security questionnaire; an on-site audit is available where the report is not
sufficient.

## Termination

On termination the controller chooses deletion or return of the data. The default is a full export in
machine-readable form within 30 days, followed by deletion within a further 30 days. Backups age out
on the normal backup cycle. See [Data retention](data-retention.md).

## Related documents

- [Data retention](data-retention.md)
- [SSO and authentication](../integrations/sso-and-authentication.md)
- [Incident escalation](../runbooks/incident-escalation.md)
