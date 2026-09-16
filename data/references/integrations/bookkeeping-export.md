# Bookkeeping export

**Document:** `integrations/bookkeeping-export.md`
**Version:** 3.2 · **Last updated:** 2026-07-14 · **Owner:** Engineering
**Applies to:** Grantly 2026.7 and later

The bookkeeping export is an optional nightly job that writes a tenant's grant and payment
transactions to a file the tenant's accounting system can import. It is configured per tenant and is
off by default.

This is an export, not a two-way integration. Grantly never reads back from the accounting system;
reconciliation of actual payments is a separate flow described in
[Disbursement](../product/disbursement.md#reconciliation).

## Schedule

| Property | Value |
| --- | --- |
| Runs | Nightly at 01:15 Europe/Stockholm |
| Window covered | Transactions with a booking date of the previous day |
| Typical duration | Under 2 minutes |
| Maximum duration before timeout | 30 minutes |
| Catch-up | One missed night is included in the next run |

A run that produces no transactions still writes an empty file with a header, so a missing file means
the job did not run, not that there was nothing to export.

> **Note**
> Since release 2026.7 each run writes a summary line to the tenant job log under
> **Administration → Jobs**: run time, transaction count, file name and outcome. This is the first
> place to look when a customer says transactions are missing.

## Transaction types exported

| Type | Triggered by |
| --- | --- |
| `commitment` | Application approved, full granted amount |
| `disbursement` | Instalment included in an exported payment run |
| `repayment` | Repayment settlement recorded |
| `cancellation` | Approved grant cancelled, unpaid amount released |
| `adjustment` | Manual correction by a finance user |

Each transaction carries the account codes configured in the account mapping. A transaction whose
programme has no mapping is written to the tenant's fallback account and flagged in the job log.

## File formats

| Format | Notes |
| --- | --- |
| SIE 4 | Default for Swedish accounting systems |
| CSV | Column set fixed, semicolon separated, UTF-8 with BOM |
| Custom | Per-tenant template, maintained by Minso |

> **Warning**
> A custom template is bespoke code owned by Minso, not configuration. Changes to it are a
> development task and are released with a normal Grantly release, not on request the same day.

## Delivery

| Method | Notes |
| --- | --- |
| SFTP push | Minso pushes to the tenant's server; key-based authentication only |
| SFTP pull | Tenant collects from Minso; files kept 90 days |
| Webhook | `payment_run.exported` only, see [API overview](api-overview.md) |

Files are named `{tenant}_{type}_{YYYYMMDD}.{ext}`. File names are stable; if a run is repeated, the
new file gets a `_r2` suffix rather than overwriting.

## Failure handling

1. The job retries once, 15 minutes after a failure.
2. If the retry fails, the run is marked `failed` in the job log and an alert is raised to Minso
   support.
3. Failed runs are **not** retried the following night automatically. The missed window is included
   in the next successful run only if the job is re-enabled and re-run.

| Failure | Usual cause |
| --- | --- |
| SFTP connection refused | Tenant firewall or rotated host key |
| `mapping_missing` | New programme without account codes |
| `file_too_large` | Over 100 000 transactions in one run |
| `job_disabled` | The job was switched off and not switched back on |

> **Note**
> The job can be disabled per tenant — for example while investigating a mapping problem. A disabled
> job produces no file and no alert, because nothing has failed. Re-enabling it is a manual step;
> check **Administration → Jobs** shows the job as `enabled` before closing an investigation.

## Limits

| Limit | Value |
| --- | --- |
| Transactions per run | 100 000 |
| File size | 250 MB |
| SFTP retention on Minso side | 90 days |
| Export file retention in Grantly | 7 years, statutory — see [Data retention](../legal/data-retention.md) |

## Reconciling a discrepancy

When a tenant reports that amounts do not match:

1. Check the job log for the period. Confirm a run exists for every night.
2. Compare the run's transaction count with the number of exported payment runs in the same period.
3. Check for transactions written to the fallback account.
4. Check whether any payment run was locked but never exported; those produce no `disbursement`
   rows.
5. If a run is missing entirely, confirm the job is enabled before re-running it.

## Related documents

- [Disbursement](../product/disbursement.md)
- [API overview](api-overview.md)
- [Incident escalation](../runbooks/incident-escalation.md)
