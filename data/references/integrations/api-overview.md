# API overview

**Document:** `integrations/api-overview.md`
**Version:** 6.1 · **Last updated:** 2026-09-08 · **Owner:** Engineering
**Applies to:** Grantly API `v1`

The Grantly REST API gives a tenant programmatic access to its own data. It is the supported way to
integrate Grantly with a finance system, a case management system or a data warehouse.

Base URL: `https://{tenant}.grantly.se/api/v1`

## Versioning

| Version | Status |
| --- | --- |
| `v1` | Current |
| `v0` | Retired 2026-05-12, see [changelog](../changelog.md) |

Breaking changes get a new major version. Additive changes — new fields, new optional parameters, new
endpoints — ship in `v1` without notice, so clients must ignore unknown fields.

## Authentication

API access uses OAuth 2.0 client credentials. A tenant administrator creates an API client under
**Administration → Integrations** and receives a client ID and secret. The secret is shown once.

```
POST /api/v1/oauth/token
grant_type=client_credentials&scope=applications:read payments:write
```

Access tokens are valid for **60 minutes**. Refresh by requesting a new token; there is no refresh
token in this flow.

| Scope | Grants |
| --- | --- |
| `applications:read` | Applications, answers, attachment metadata |
| `applications:write` | State transitions, caseworker assignment |
| `attachments:read` | Attachment content |
| `payments:read` | Payment plans, runs, files |
| `payments:write` | Reconciliation, repayment settlements |
| `reports:read` | Interim and final reports |
| `admin:read` | Users, programmes, audit log |

Scopes are per client. There is no user context; everything an API client does is attributed to the
client in the audit log.

## Core endpoints

| Method | Path | Notes |
| --- | --- | --- |
| `GET` | `/applications` | Filters: `programme_id`, `state`, `updated_since` |
| `GET` | `/applications/{id}` | Includes answers |
| `POST` | `/applications/{id}/transitions` | Body: `{ "to": "under_review" }` |
| `GET` | `/applications/{id}/attachments` | Metadata only |
| `GET` | `/attachments/{id}/content` | Redirects to a signed URL, valid 5 minutes |
| `GET` | `/programmes` | |
| `GET` | `/payment_runs` | |
| `POST` | `/payment_runs/{id}/reconciliation` | See [Disbursement](../product/disbursement.md) |
| `GET` | `/reports` | |
| `POST` | `/repayments/{id}/settlements` | |
| `GET` | `/audit_log` | `admin:read`, max 90 days per query |

`updated_since` takes an ISO 8601 timestamp and was documented in release 2026.9. It is the
recommended way to run an incremental sync.

## Pagination

Cursor-based. Responses include `next_cursor` when more data exists.

```
GET /applications?limit=100&cursor=eyJpZCI6...
```

| Parameter | Default | Maximum |
| --- | --- | --- |
| `limit` | 50 | 200 |

Do not construct cursors. They are opaque and their format changes without notice.

## Rate limits

| Limit | Value |
| --- | --- |
| Requests per minute, per client | 600 |
| Requests per minute, per tenant | 1 200 |
| Concurrent requests per client | 20 |
| Attachment content downloads per minute | 60 |

Every response carries `X-RateLimit-Remaining` and `X-RateLimit-Reset` (added in 2026.7). Exceeding
a limit returns `429` with a `Retry-After` header in seconds.

> **Note**
> Rate limits are not raised per tenant. A job that needs more throughput should use `updated_since`
> rather than re-reading everything, or use [Bulk export](../runbooks/bulk-export.md) for full
> extracts.

## Errors

| Status | Meaning |
| --- | --- |
| `400` | Malformed request |
| `401` | Missing or expired token |
| `403` | Token valid, scope insufficient |
| `404` | Not found, or not visible to this client |
| `409` | State conflict, e.g. an invalid transition |
| `422` | Validation failed, `errors` lists the fields |
| `429` | Rate limited |
| `503` | Maintenance, retry with backoff |

Error bodies are `{ "code": "...", "message": "...", "errors": [...] }`. Match on `code`, not on
`message`.

## Webhooks

A tenant can register up to 10 webhook endpoints.

| Event | Fires on |
| --- | --- |
| `application.submitted` | Submission |
| `application.state_changed` | Any transition |
| `decision.issued` | Approval or rejection |
| `payment_run.exported` | Payment file generated |
| `report.submitted` | Report submitted |

Deliveries are signed with `X-Grantly-Signature` (HMAC-SHA256 over the raw body). Retries run at 1,
5, 15 and 60 minutes; after four failures the delivery is dropped and the endpoint is disabled after
100 consecutive failures. Webhooks are at-least-once — handlers must be idempotent on `event_id`.

## Related documents

- [SSO and authentication](sso-and-authentication.md)
- [Bookkeeping export](bookkeeping-export.md)
- [Application workflow](../product/application-workflow.md)
