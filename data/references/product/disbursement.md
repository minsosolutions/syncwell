# Disbursement

**Document:** `product/disbursement.md`
**Version:** 4.1 · **Last updated:** 2026-09-08 · **Owner:** Product
**Applies to:** Grantly 2026.9 and later

Disbursement covers everything after a grant is approved: the payment plan, the payment runs that
generate payment files, and the reconciliation of what was actually paid.

Grantly does not move money. It produces payment files that the tenant's finance system or bank
processes, and it records what finance reports back.

## States after approval

| State | Meaning |
| --- | --- |
| `paying` | At least one instalment is scheduled or paid |
| `paid` | All instalments settled |
| `suspended` | Payments stopped, decision stands |
| `repaid` | Grant repaid in full, see [Reporting and follow-up](reporting-and-followup.md) |
| `cancelled` | Decision revoked, unpaid instalments voided |

## Payment plans

A payment plan is created when an application is approved. The default plan is a single instalment of
the full granted amount, payable on the first payment run after the decision.

| Plan type | Instalments | Available since |
| --- | --- | --- |
| Single | 1 | — |
| Split | 2, on approval and on final report approval | — |
| Monthly | Up to 24 | — |
| Quarterly | Up to 8 | 2026.8 |
| Custom | Up to 24, dates and amounts set by hand | — |

Instalment amounts must sum exactly to `amount_granted`. Grantly distributes an uneven amount by
putting the remainder in the **final** instalment, so a 10 000 SEK grant over three instalments
becomes 3 333 / 3 333 / 3 334.

> **Note**
> Before release 2026.9 the remainder was rounded off the last instalment rather than added to it,
> which left plans short by up to one krona. Existing plans were not corrected automatically; check
> plans created before 2026-09-08 if the totals do not reconcile.

## Payment runs

A payment run collects every instalment due on or before its run date and produces one payment file.

| Field | Description |
| --- | --- |
| `payment_run_id` | Identifier |
| `run_date` | Date the file is generated |
| `scope` | Whole tenant, or a single programme (since 2026.5) |
| `state` | `draft`, `locked`, `exported`, `reconciled` |
| `total_amount` | Sum of included instalments |

Runs are created manually or on a schedule. The recommended schedule is weekly. A run in `draft` can
have instalments added or removed; locking it freezes the contents and allows the file to be
generated.

| Limit | Value |
| --- | --- |
| Instalments per run | 20 000 |
| Payment runs per tenant per day | 5 |
| Retention of payment files | 7 years, statutory — see [Data retention](../legal/data-retention.md) |

## Payment file formats

| Format | Use |
| --- | --- |
| ISO 20022 `pain.001.001.09` | Default, bank transfers |
| Bankgirot Leverantörsbetalningar | Legacy, SEK only |
| CSV | Manual import into a finance system |

The bank account on each payment is taken from the applicant's record at the time the run is locked,
not at the time of the decision. A bank account changed between decision and payment takes effect on
the next run.

## Reconciliation

Finance reports back which payments cleared. Reconciliation can be done three ways:

1. Upload a bank return file (`pain.002` or the bank's CSV).
2. Mark instalments as paid in the interface.
3. `POST /v1/payment_runs/{id}/reconciliation` — see
   [API overview](../integrations/api-overview.md).

An instalment that is not reconciled within 30 days of its run date is flagged. Flagged instalments
appear on the finance dashboard and are the usual sign that a payment file was generated but never
processed.

> **Warning**
> A payment run that has been exported cannot be deleted. If a file was generated in error, cancel
> the individual instalments and issue a new run. Deleting the file outside Grantly does not change
> what Grantly believes was paid.

## Bookkeeping

Approved, paid and repaid amounts are exported nightly to the tenant's bookkeeping system where that
integration is configured. See [Bookkeeping export](../integrations/bookkeeping-export.md).

## Related documents

- [Application workflow](application-workflow.md)
- [Reporting and follow-up](reporting-and-followup.md)
- [Bookkeeping export](../integrations/bookkeeping-export.md)
