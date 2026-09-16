# Grantly release changelog

**Document:** `changelog.md`
**Version:** rolling · **Last updated:** 2026-09-10 · **Applies to:** all Grantly tenants

Most recent release first. Releases ship on the second Tuesday of the month. Hotfixes are listed
under their parent release. Entries marked **Breaking** require action from the tenant before the
following release.

Only changes visible to customers are listed here. Internal refactoring, dependency bumps and
infrastructure work are tracked in the engineering board.

---

## 2026.9 — released 2026-09-08

- **Review:** reviewers can now be assigned to a whole programme instead of application by
  application. See [Review and scoring](product/review-and-scoring.md#assigning-reviewers).
- **Portal:** the applicant draft list is sorted by last edited instead of by created date.
- **API:** `GET /v1/applications` accepts `updated_since` as a filter. Undocumented before, now
  supported. See [API overview](integrations/api-overview.md).
- **Reporting:** final report reminders can be sent 7, 14 or 30 days before the deadline. Previously
  only 14 days.
- **Fix:** payment file generation no longer rounds instalment amounts down when the plan has an odd
  number of instalments.
- **Fix:** the programme copy function no longer copies archived form fields.

### 2026.9.1 — hotfix 2026-09-11

- **Fix:** Swedish characters in reviewer names were mangled in the PDF decision protocol.

---

## 2026.8 — released 2026-08-11

- **Workflow:** new `on_hold` state for applications parked pending clarification from the applicant.
  See [Application workflow](product/application-workflow.md#states).
- **Disbursement:** support for quarterly instalment plans.
- **SSO:** OIDC `acr_values` can be configured per tenant.
  See [SSO and authentication](integrations/sso-and-authentication.md).
- **Portal:** applicants receive an email confirmation with a submission reference on submit.
- **Attachments:** the accepted file type list now includes `.odt` and `.ods`.
- **Fix:** the caseworker inbox counter did not reset after bulk reassignment.
- **Fix:** CSV export used a semicolon separator in one place and a comma in another.

### 2026.8.2 — hotfix 2026-08-20

- **Fix:** scheduled nightly jobs that failed did not always appear in the tenant job log.

---

## 2026.7 — released 2026-07-14

- **Retention:** the default retention period for **rejected applications** changed from **24 months
  to 36 months**, to match the audit period most public-sector grant-givers are subject to. Existing
  tenants were migrated on the release date; tenants with a contractual override keep their agreed
  period. See [Data retention](legal/data-retention.md).
- **Workflow:** decision letters can be generated in bulk for a whole review round.
- **Bookkeeping export:** the nightly export job writes a per-run summary line to the tenant job log.
  See [Bookkeeping export](integrations/bookkeeping-export.md).
- **API:** rate limit headers (`X-RateLimit-Remaining`, `X-RateLimit-Reset`) added to all responses.
- **Portal:** the applicant session idle timeout is shown in the account menu.
- **Fix:** attachments with uppercase file extensions were rejected by the type check.
- **Fix:** programme budgets displayed with two decimals in the list view and none in the detail view.

---

## 2026.6 — released 2026-06-09

- **Review:** weighted scoring models. Criteria can be given a weight between 1 and 5.
- **Workflow:** applications can be withdrawn by the applicant up until the programme deadline.
- **Reporting:** interim reports can be made optional per programme.
- **Fix:** duplicate reminder emails when a deadline fell on a Sunday.
- **Fix:** the API returned `500` instead of `422` for a malformed `organisation_number`.

### 2026.6.1 — hotfix 2026-06-17

- **Fix:** SAML logout redirected to the wrong tenant when a user had access to two tenants.

---

## 2026.5 — released 2026-05-12

- **Breaking:** API version `v0` retired. All endpoints are `/v1`.
- **Disbursement:** payment files can be generated per programme instead of per payment run.
- **Portal:** improved contrast and focus indicators to meet WCAG 2.2 AA.
- **Attachments:** virus scanning moved to a queue; uploads are accepted immediately and quarantined
  if a scan later fails. See [Attachments and file uploads](product/attachments-and-file-uploads.md).
- **Fix:** the applicant portal did not show the programme deadline in the applicant's own timezone.

---

## 2026.4 — released 2026-04-14

- **Workflow:** free-text reason required when a caseworker returns an application to the applicant.
- **Reporting:** repayment requests can be recorded against a finished grant.
- **Runbooks:** bulk export moved from a manual database query to a supported tenant admin function.
  See [Bulk export](runbooks/bulk-export.md).
- **Fix:** archived programmes appeared in the applicant portal programme picker.
- **Fix:** timestamps in the audit log were written in local time instead of UTC.

---

## Older releases

Releases before 2026.4 are archived. Ask Support if you need an entry from the archive.
