# Grantly shared reference documentation

**Document:** `README.md`
**Version:** 4.4 · **Last updated:** 2026-09-10 · **Applies to:** all Grantly tenants

This is the shared reference library for **Grantly**, Minso's grant management platform. Everything
here applies to every tenant. It describes standard product behaviour, standard limits and the
standard procedures Minso staff follow.

> **Note**
> Nothing in this library is customer-specific. Where a tenant has a contractual override — a longer
> retention period, a raised quota, a bespoke integration — the override is recorded in that
> customer's own agreement and takes precedence over anything written here. If you cannot find a
> customer's override, assume the documented default and say so explicitly when you answer.

## How to use this library

1. Find the document that owns the topic (see the index below). Each topic has exactly one owner.
2. Check the **Last updated** line at the top of the document.
3. Check [changelog.md](changelog.md) for anything that has shipped since that date.
4. When answering a customer, cite the document and its version.

## Index

### Product

| Document | Covers |
| --- | --- |
| [Application workflow](product/application-workflow.md) | Application states, transitions, deadlines, withdrawal |
| [Review and scoring](product/review-and-scoring.md) | Review rounds, scoring models, conflict of interest, decisions |
| [Disbursement](product/disbursement.md) | Payment plans, instalments, payment files, reconciliation |
| [Reporting and follow-up](product/reporting-and-followup.md) | Interim and final reports, reminders, repayment |
| [Portal for applicants](product/portal-for-applicants.md) | Applicant accounts, drafts, sessions, notifications |
| [Attachments and file uploads](product/attachments-and-file-uploads.md) | File types, size limits, virus scanning, upload behaviour |

### Integrations

| Document | Covers |
| --- | --- |
| [API overview](integrations/api-overview.md) | REST API, authentication, pagination, rate limits, webhooks |
| [SSO and authentication](integrations/sso-and-authentication.md) | SAML, OIDC, BankID, session lifetimes, staff roles |
| [Bookkeeping export](integrations/bookkeeping-export.md) | Nightly ledger export, file formats, failure handling |

### Legal

| Document | Covers |
| --- | --- |
| [Data retention](legal/data-retention.md) | Retention periods, anonymisation, erasure requests |
| [Processing agreement](legal/processing-agreement.md) | DPA structure, sub-processors, transfers, audit rights |

### Runbooks (Minso staff)

| Document | Covers |
| --- | --- |
| [Bulk export](runbooks/bulk-export.md) | Running and troubleshooting large exports on behalf of a tenant |
| [Incident escalation](runbooks/incident-escalation.md) | Severity levels, on-call, customer communication |

## Release naming

Grantly releases are named `YYYY.M` — for example `2026.9` for the September 2026 release. Releases
ship on the second Tuesday of the month unless stated otherwise in the changelog. Hotfixes are named
`YYYY.M.N` and are listed under their parent release.

## Document conventions

- **States** are written in `code` and are always lowercase with underscores, matching the API.
- **Limits** given as "per tenant" apply to the whole tenant, not per programme or per user.
- **Note** callouts are behaviour that surprises people. **Warning** callouts are destructive or
  irreversible.
- Dates are ISO 8601 (`YYYY-MM-DD`). Times are Europe/Stockholm unless a field says UTC.

## Ownership

| Area | Owner |
| --- | --- |
| Product documents | Product, reviewed by Support |
| Integration documents | Engineering |
| Legal documents | Legal, reviewed by the DPO |
| Runbooks | Support |

Corrections go to the documentation owner listed above. Do not edit a customer answer into a shared
document; if a tenant needs different wording, it belongs in their agreement.
