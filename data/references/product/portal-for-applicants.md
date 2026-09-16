# Portal for applicants

**Document:** `product/portal-for-applicants.md`
**Version:** 4.0 · **Last updated:** 2026-09-08 · **Owner:** Product
**Applies to:** Grantly 2026.9 and later

The applicant portal is the public-facing side of Grantly. Applicants create an account, find open
programmes, fill in an application, attach documents, submit, and later report on the grant.

Each tenant has its own portal on its own address, with the tenant's logo and colours. Applicants
have no visibility of other tenants.

## Accounts

| Sign-in method | Identifies | Typical use |
| --- | --- | --- |
| Email and password | Email address | Individuals, small associations |
| BankID | Personal identity number | Swedish individuals |
| Organisation account | Organisation number | Associations and companies |

An organisation account can have several members. Any member can see and edit every application
belonging to the organisation; there are no per-application permissions inside an organisation. This
surprises people, so it is worth stating when a customer asks.

Account details are covered in
[SSO and authentication](../integrations/sso-and-authentication.md).

## Finding a programme

The portal lists programmes that are open, with the opening and closing dates in the applicant's own
timezone. Archived programmes are not listed. A programme can be unlisted and reachable only by
direct link, which tenants use for invitation-only calls.

## Drafts

An application is a draft until it is submitted.

| Behaviour | Detail |
| --- | --- |
| Autosave | Every 30 seconds and on field blur |
| Manual save | Always available |
| Draft list sorting | Last edited first, since 2026.9 |
| Drafts per applicant | 20 open drafts |
| Draft retention | 12 months from last edit — see [Data retention](../legal/data-retention.md) |

> **Note**
> Autosave saves answers, not attachments. An attachment is stored when its upload completes, which
> is a separate operation. See [Attachments and file uploads](attachments-and-file-uploads.md).

## Sessions and timeouts

| Setting | Value | Configurable |
| --- | --- | --- |
| Idle timeout | 30 minutes | No |
| Absolute session lifetime | 12 hours | No |
| Warning before idle timeout | 2 minutes | No |
| Sessions per account | 5 concurrent | No |

Activity means a request to the server. Typing into a form field is activity only because autosave
fires; a long stretch of reading or of preparing a document in another application is not.

Two minutes before the idle timeout, the portal shows a dialogue offering to extend the session. If
the applicant does not respond, the session ends and they are returned to the sign-in page. After
signing in again they are taken back to the page they were on, and autosaved answers are intact.

> **Note**
> A session that ends mid-task never discards work that the portal has confirmed as saved. Answers
> are protected by autosave; uploads in progress are protected by the upload ticket refresh described
> in [Attachments and file uploads](attachments-and-file-uploads.md#sessions-and-upload-tickets).

## Submitting

On submit the portal validates every required field, checks that required attachments are present and
`available`, and shows a summary for confirmation. After confirmation the applicant receives an email
with a submission reference, added in release 2026.8.

A submitted application cannot be edited. The applicant can withdraw it, or the caseworker can return
it for changes. See [Application workflow](application-workflow.md).

## Notifications

| Event | Email | In portal |
| --- | --- | --- |
| Application submitted | Yes | Yes |
| Application returned | Yes | Yes |
| Decision issued | Yes | Yes |
| Attachment quarantined | Yes | Yes |
| Report due | Yes | Yes |
| Report approved | Yes | Yes |

Applicants cannot turn off notifications about decisions or reports. They can turn off programme
announcements.

## Accessibility

The portal is built to WCAG 2.2 AA. Contrast and focus indicators were reworked in release 2026.5. An
accessibility statement is published per tenant; Minso supplies the technical part of it and the
tenant supplies the contact and complaints route.

## Languages

Swedish and English are supported. The language follows the applicant's browser and can be changed in
the account menu. Form content is written by the tenant and is not translated by Grantly; a form
authored only in Swedish stays in Swedish.

## Related documents

- [Application workflow](application-workflow.md)
- [Attachments and file uploads](attachments-and-file-uploads.md)
- [Reporting and follow-up](reporting-and-followup.md)
