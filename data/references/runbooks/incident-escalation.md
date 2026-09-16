# Runbook: incident escalation

**Document:** `runbooks/incident-escalation.md`
**Version:** 4.3 · **Last updated:** 2026-09-10 · **Owner:** Support
**Audience:** Minso support, delivery and engineering staff

This runbook covers how an incident is classified, who is woken, and what the customer is told. It
applies to platform incidents affecting one or more tenants. It does not cover individual support
questions.

## Severity

| Severity | Definition | Examples |
| --- | --- | --- |
| **S1** | Platform unavailable, or data loss or exposure | Portal down, applications missing, cross-tenant data visible |
| **S2** | Core function unavailable for one or more tenants, no workaround | Cannot submit applications, payment files not generated, SSO broken |
| **S3** | Core function degraded, or a workaround exists | Export times out, slow list views, reminders delayed |
| **S4** | Cosmetic or isolated | Formatting error, one user affected |

A deadline makes severity worse, not better. An S3 the week before a programme deadline is treated as
an S2; say so explicitly in the incident channel rather than leaving it implied.

Anything that may be a personal data breach is an S1 until the DPO says otherwise, regardless of how
few records are involved. See [Processing agreement](../legal/processing-agreement.md#breach-notification).

## Response and update targets

| Severity | First response | Customer update cadence | Target resolution |
| --- | --- | --- | --- |
| S1 | 30 minutes, 24/7 | Hourly | 4 hours |
| S2 | 1 hour, office hours | Every 4 hours | 1 business day |
| S3 | 1 business day | Daily while open | Next release |
| S4 | 3 business days | On change | Backlog |

Office hours are 08:00–17:00 Europe/Stockholm, Monday to Friday. S1 is outside office hours as well.

## Roles during an incident

| Role | Responsibility |
| --- | --- |
| Incident lead | Owns the incident, decides severity, decides when it is closed |
| Engineer on call | Investigates and fixes |
| Communications | Writes customer updates; never the same person as the engineer |
| DPO | Consulted on anything touching personal data |

The incident lead is whoever declares the incident until they hand over explicitly. Handover is
written in the incident channel, not agreed verbally.

## Procedure

1. **Declare.** Open an incident channel named `inc-YYYYMMDD-short-name`. State severity, affected
   tenants and what is known.
2. **Classify.** Confirm whether one tenant or several are affected. Check the job log and the error
   rate dashboard before assuming it is tenant-specific.
3. **Notify.** Post the first customer update within the response target above, even if the update is
   only "we are investigating".
4. **Investigate and fix.** Keep a running timeline in the channel with timestamps.
5. **Confirm.** Verify with the affected tenant, not only with monitoring.
6. **Close.** State in the channel that the incident is closed and what the follow-up is.
7. **Review.** S1 and S2 get a written post-incident review within 5 business days.

> **Note**
> Step 2 is the one that is skipped. A symptom reported by one tenant is often present at others who
> have not noticed or have not reported it. Before treating something as tenant-specific, check
> whether the same error signature appears for other tenants — and check whether other tenants have
> reported the same thing in different words.

## Customer communication

- Say what is affected and what is not. A customer needs to know whether to tell their applicants.
- Give the next update time, and send it even when there is nothing new.
- Do not attribute cause until it is confirmed. "We are investigating" beats a wrong cause.
- Never tell a customer that the problem is on their side unless it has been demonstrated. An
  applicant-facing defect frequently arrives described as applicant error.

## Post-incident review

Covers: timeline, impact per tenant, root cause, what made it worse, what made it better, and
actions with owners and dates. Actions go on the engineering board with the incident ID. A review
without dated actions is not finished.

> **Warning**
> Do not mark a ticket as solved before the customer confirms. A fix deployed is not a fix verified.
> Where a ticket was closed and the symptom later recurs, reopen the original ticket rather than
> opening a new one, so the history stays in one place.

## Escalation contacts

| Need | Route |
| --- | --- |
| Engineering on call | On-call rota |
| Data protection | DPO |
| Customer relationship | Customer success owner for the tenant |
| Contractual or legal | Legal |

## Related documents

- [Processing agreement](../legal/processing-agreement.md)
- [Bulk export](bulk-export.md)
- [Bookkeeping export](../integrations/bookkeeping-export.md)
