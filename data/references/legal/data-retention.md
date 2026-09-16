# Data retention

**Document:** `legal/data-retention.md`
**Version:** 3.1 · **Last updated:** 2026-02-11 · **Owner:** Legal, reviewed by the DPO
**Applies to:** all Grantly tenants · **Review cycle:** annual

This document states the retention periods Grantly applies to personal data by default, how data is
anonymised at the end of a period, and how erasure requests are handled. It is the authoritative
source for retention questions and is the document Minso cites in the processing agreement. See
[Processing agreement](processing-agreement.md) for the contractual framework around it.

## Roles

The tenant (the grant-giving organisation) is the **controller**. Minso is the **processor**.
Retention periods are therefore set by the tenant; the periods below are the product defaults that
apply unless the tenant has agreed a different period in writing. Minso does not delete or extend
data outside these periods without a documented instruction from the controller.

## Default retention periods

| Data category | Retention period | Counted from | End-of-period action |
| --- | --- | --- | --- |
| Draft application, never submitted | 12 months | Last edit | Hard delete |
| **Rejected application** | **24 months** | Date the rejection decision is issued | Anonymisation |
| Withdrawn application | 24 months | Date of withdrawal | Anonymisation |
| Approved application, grant paid | 10 years | Date of the final disbursement | Anonymisation |
| Approved application, grant repaid | 10 years | Date the repayment is settled | Anonymisation |
| Final and interim reports | Follows the parent application | — | Anonymisation |
| Attachments | Follows the parent application | — | Hard delete |
| Applicant portal account, no applications | 24 months | Last sign-in | Hard delete |
| Staff user account | 24 months | Account deactivation | Hard delete |
| Audit log | 5 years | Event timestamp | Hard delete |
| Payment files and bookkeeping exports | 7 years | File generation | Hard delete, statutory |
| Support correspondence | 3 years | Ticket closure | Hard delete |
| Backups | 35 days | Backup creation | Rolling overwrite |

Rejected applications are retained for **24 months** from the date the rejection decision is issued.
At the end of that period the application, its answers and its attachments are anonymised: the
applicant's name, personal identity number, contact details, organisation contact details and all
free-text fields that may contain personal data are removed, and the record is reduced to the
statistical fields listed under [Anonymisation](#anonymisation).

The 24-month period is chosen so that a rejected applicant can refer back to a decision across two
subsequent application rounds, and so that the tenant can answer questions about an earlier decision
within the normal complaints window. Tenants subject to a longer statutory audit period should agree
a longer period in their processing agreement.

> **Note**
> The period is counted from the **decision date**, not from the application date or the date the
> programme closed. An application rejected on 2026-01-15 is anonymised on 2028-01-15, regardless of
> when it was submitted.

## Anonymisation

Anonymisation is irreversible. After anonymisation, the following fields remain:

| Field | Kept |
| --- | --- |
| `application_id` | Yes |
| `programme_id` | Yes |
| `decision` | Yes |
| `decision_date` | Yes |
| `amount_applied_for` | Yes |
| `amount_granted` | Yes |
| `municipality_code` | Yes |
| `category` | Yes |
| Applicant name, identifiers, contact details | No |
| Free-text answers and attachments | No |

Anonymised records remain in reports and statistics, so historical totals do not change when a
retention period elapses.

> **Warning**
> Anonymisation cannot be undone and cannot be reversed from a backup after 35 days. If a tenant
> needs a record preserved beyond its period — for example because it is subject to an ongoing
> appeal — a legal hold must be placed on the application **before** the period elapses.

## Legal hold

A tenant administrator can place a legal hold on an individual application or on a whole programme.
A record under legal hold is excluded from anonymisation and from hard deletion until the hold is
lifted. Holds are recorded in the audit log with the user who placed them and the stated reason.

## How the job runs

The retention job runs nightly at 02:30 Europe/Stockholm per tenant. It processes records whose
period elapsed on or before the previous day. A tenant administrator can see the next scheduled
anonymisation date on each application under **Administration → Retention**.

If the job fails, it retries the following night. A record is never partially anonymised; each record
is processed in a single transaction.

## Erasure requests

An applicant may ask the controller to erase their data. The request goes to the tenant, not to
Minso. Once the tenant approves it, a tenant administrator triggers erasure from the application
record; Minso does not act on erasure requests received directly from applicants.

Erasure is refused where the controller has a legal obligation to keep the record — in practice for
approved applications where public money has been paid out, and for payment files and bookkeeping
exports, which are kept for seven years under the Swedish Accounting Act.

## Related documents

- [Processing agreement](processing-agreement.md)
- [Application workflow](../product/application-workflow.md)
- [Attachments and file uploads](../product/attachments-and-file-uploads.md)
