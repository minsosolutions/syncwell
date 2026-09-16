# Reporting and follow-up

**Document:** `product/reporting-and-followup.md`
**Version:** 3.3 · **Last updated:** 2026-09-08 · **Owner:** Product
**Applies to:** Grantly 2026.9 and later

Follow-up is what happens after money is paid: the grantee reports on how it was used, the tenant
approves or rejects the report, and in some cases money is paid back.

## Report types

| Type | When | Configurable |
| --- | --- | --- |
| Interim report | At a date set in the payment plan | Optional per programme since 2026.6 |
| Final report | After the project end date | Always required |
| Financial appendix | With the final report | Optional per programme |

A programme defines its own report form, in the same form builder used for application forms. Report
forms can be changed between programmes but not while a report is open for submission.

## Report states

| State | Set by | Next |
| --- | --- | --- |
| `pending` | System, when the report becomes due | `submitted` |
| `submitted` | Grantee | `approved`, `returned` |
| `returned` | Caseworker | `submitted` |
| `approved` | Caseworker | — |
| `overdue` | System, after the deadline | `submitted` |

A report in `overdue` still accepts a submission. Grantly does not block a grantee from reporting
late; it records that the report was late in the `submitted_late` flag.

## Reminders

Reminders are sent by email to the contact on the application.

| Reminder | Default timing | Configurable |
| --- | --- | --- |
| Upcoming report | 14 days before deadline | 7, 14 or 30 days since 2026.9 |
| Deadline today | On the deadline | On/off |
| Overdue | 7 and 30 days after | On/off |

> **Note**
> Reminders go to the contact person recorded on the application, which is not necessarily the
> portal account that submitted it. If a grantee says they never received a reminder, check the
> contact field before checking the mail log.

Reminders are sent at 07:00 Europe/Stockholm. A reminder whose deadline falls on a Saturday or Sunday
is sent on the preceding Friday.

## Approving a report

A caseworker opens the report, reviews the answers and attachments, and either approves it or returns
it with a reason. Approving the final report:

- Releases the final instalment in a split payment plan. See [Disbursement](disbursement.md).
- Moves the grant to `paid` if no instalments remain.
- Closes the follow-up for that grant.

Returning a report reopens it for the grantee and resets the deadline to a date the caseworker
chooses. There is no limit on how many times a report can be returned.

## Repayment

A repayment request is recorded against the grant when the tenant decides that some or all of the
money must be paid back — typically after a rejected final report or an unused grant.

| Field | Description |
| --- | --- |
| `repayment_id` | Identifier |
| `amount` | Amount to be repaid, up to the amount paid out |
| `reason` | Free text, required |
| `due_date` | Date the repayment is due |
| `state` | `requested`, `partially_settled`, `settled`, `written_off` |

Repayments are tracked, not collected. Grantly produces a repayment notice as a PDF; the actual
invoice is issued by the tenant's finance system. Settlements are recorded by hand or through
`POST /v1/repayments/{id}/settlements`.

A grant with an outstanding repayment is flagged on the applicant's record. The flag is visible to
caseworkers when the same applicant applies to a later programme, which is the usual reason tenants
ask for it.

## Statistics

The follow-up dashboard shows, per programme:

- Reports due in the next 30 days
- Overdue reports, by age
- Reports returned more than once
- Outstanding repayments and their total

Figures include anonymised historical records, so a total does not drop when a retention period
elapses. See [Data retention](../legal/data-retention.md).

## Related documents

- [Disbursement](disbursement.md)
- [Portal for applicants](portal-for-applicants.md)
- [Bulk export](../runbooks/bulk-export.md)
