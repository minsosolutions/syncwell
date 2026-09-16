# Runbook: bulk export

**Document:** `runbooks/bulk-export.md`
**Version:** 2.9 · **Last updated:** 2026-08-14 · **Owner:** Support
**Audience:** Minso support and delivery staff, and tenant administrators

Bulk export produces a full extract of a tenant's applications, answers, decisions and optionally
attachments. Tenants use it for archiving, for statutory reporting and for handing material to an
auditor.

Since release 2026.4 this is a supported tenant admin function. Do not run database queries against
a tenant to produce an export.

## When to use what

| Need | Use |
| --- | --- |
| Incremental sync into another system | API `updated_since`, see [API overview](../integrations/api-overview.md) |
| A spreadsheet of one programme | Interface export, list view → **Export** |
| Everything, including attachments | Bulk export, this runbook |
| Accounting transactions | [Bookkeeping export](../integrations/bookkeeping-export.md) |

## Running an export

1. Go to **Administration → Export**.
2. Choose the scope: whole tenant, one or more programmes, or a date range on `decision_date`.
3. Choose the contents:
   - Application data (always included)
   - Answers
   - Attachments
   - Review scores and comments
   - Audit log entries
4. Choose the format: CSV, JSON or XLSX. Attachments are always delivered as files inside a ZIP.
5. Start the export. It runs asynchronously; you can close the page.
6. When it completes, the requester receives an email with a download link.

Download links are valid for **7 days** and require sign-in. The export file is deleted after 14
days.

## Limits

| Limit | Value |
| --- | --- |
| Applications per export | 50 000 |
| Attachment volume per export | 20 GB |
| Output file size | 20 GB, split into 2 GB parts above that |
| Concurrent exports per tenant | 1 |
| Exports per tenant per day | 10 |
| Maximum run time before the job is abandoned | 4 hours |

> **Note**
> The concurrency limit is one. A second export started while the first is running is queued, not
> rejected. A tenant that says "the export button does nothing" usually has one already running —
> check **Administration → Jobs**.

## Expected run times

| Scope | Without attachments | With attachments |
| --- | --- | --- |
| Up to 1 000 applications | Under 2 minutes | 5–15 minutes |
| 1 000–10 000 | 5–20 minutes | 30–90 minutes |
| 10 000–50 000 | 30–90 minutes | 2–4 hours |

Run times are for a tenant on standard capacity and assume the export is not competing with the
nightly jobs. Exports started between 01:00 and 04:00 compete with retention, bookkeeping and backup
jobs and take longer.

> **Note**
> Recommend that large exports are started in the morning, not overnight. The overnight window looks
> quiet to the tenant but is the busiest window on the platform.

## Troubleshooting

| Symptom | Check | Action |
| --- | --- | --- |
| Export never completes | Job log for a 4-hour abandon | Narrow the scope, split by programme |
| Export completes, ZIP is short | Quarantined attachments | Expected; they are excluded and listed in `excluded.csv` |
| `too_many_applications` | Scope over 50 000 | Split by `decision_date` range |
| Email never arrives | Requester's address, spam filter | Download from **Administration → Export** instead |
| Download link expired | Older than 7 days | Re-run the export |
| Timeout above 10 000 applications | Whether attachments are included | Export data and attachments as two separate runs |

## Escalation

Escalate to Engineering when:

- An export abandons at 4 hours twice on the same narrowed scope.
- An export completes but the row count does not match the tenant's own count from the list view.
- The export job fails with a `500` rather than an abandon.

Include the export job ID, the tenant slug, the scope, the format and the job log lines. See
[Incident escalation](incident-escalation.md).

> **Warning**
> Do not close a ticket about a slow or failing export on the basis that a workaround was found. A
> narrowed scope or an off-hours run is a workaround, not a fix. Record the workaround on the ticket
> and keep the ticket open until the underlying run time is acceptable at the tenant's real volume.

## Related documents

- [API overview](../integrations/api-overview.md)
- [Attachments and file uploads](../product/attachments-and-file-uploads.md)
- [Data retention](../legal/data-retention.md)
