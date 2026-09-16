# Review and scoring

**Document:** `product/review-and-scoring.md`
**Version:** 3.8 · **Last updated:** 2026-09-10 · **Owner:** Product
**Applies to:** Grantly 2026.9 and later

Review happens in rounds. A round groups the applications to a programme, the reviewers assigned to
them and the scoring model in force. A programme can have several rounds; an application belongs to
exactly one round at a time.

## Review rounds

| Field | Description |
| --- | --- |
| `review_round_id` | Identifier |
| `name` | Free text, shown to reviewers |
| `opens_at` / `closes_at` | When reviewers can enter scores |
| `scoring_model_id` | Model in force for this round |
| `state` | `planned`, `open`, `closed`, `decided` |

A round moves to `closed` when it passes `closes_at` or when an administrator closes it manually.
Scores cannot be changed in a `closed` round without reopening it, which is logged.

## Assigning reviewers

Reviewers can be assigned in three ways:

1. **Per application** — select applications and assign a reviewer to them.
2. **Per programme** — a reviewer assigned to a programme sees every application in every open round
   of that programme. Added in release 2026.9.
3. **By import** — upload a CSV of `application_reference,reviewer_email` pairs, up to 2 000 rows.

A single application can have up to 10 reviewers. A reviewer can be assigned to at most 500
applications across all tenants they have access to; above that the assignment is rejected and the
administrator is told which reviewer is over the limit.

> **Note**
> Assigning a reviewer to a programme does not retroactively assign them to rounds that are already
> `closed`. To include a reviewer in a closed round, reopen it.

## Scoring models

A scoring model is a list of criteria. Each criterion has a name, a scale and, since release 2026.6,
a weight.

| Property | Values |
| --- | --- |
| Scale | `1-3`, `1-5`, `1-10`, `pass_fail` |
| Weight | Integer 1–5 (default 1) |
| Comment | `optional`, `required`, `hidden` |

The total score is the weighted mean of the numeric criteria, rounded to one decimal. `pass_fail`
criteria do not contribute to the total; a `fail` on any such criterion flags the application in the
list view but does not block a decision.

A scoring model cannot be changed once a round using it has opened. Create a new model and a new
round instead.

## Conflict of interest

A reviewer can declare a conflict on an individual application. On declaring:

- The reviewer loses access to the application immediately.
- Any score they had entered is removed from the total and retained in the audit log.
- The application is flagged for reassignment in the administrator's task list.

Conflicts are declared by the reviewer, not assigned by an administrator. An administrator can
remove a reviewer from an application, but that is recorded as a reassignment, not a conflict.

## Reviewer visibility

| Reviewer sees | Default |
| --- | --- |
| Applicant name and organisation | Yes |
| Applicant contact details | No |
| Other reviewers' scores | No, until the round is `closed` |
| Other reviewers' comments | No, until the round is `closed` |
| Attachments | Yes, if `available` |

Anonymous review can be enabled per programme, which hides applicant name and organisation from
reviewers. It cannot be enabled after a round has opened.

## From scores to decision

Scores are advisory. The decision is made by a user with the **Decision maker** role and is entered
against the application, not against the round. Grantly does not compute decisions from scores and
does not rank applications automatically; the list can be sorted by total score, which is as far as
it goes.

When the decision is entered, the application moves to `approved` or `rejected` as described in
[Application workflow](application-workflow.md#states). A rejection requires a
`rejection_reason_code` from the tenant's own list of codes.

Decision protocols can be generated per round as a PDF containing the applications, the scores, the
decisions and the decision maker. The protocol is generated once and stored; regenerating it creates
a new version and keeps the old one.

## Related documents

- [Application workflow](application-workflow.md)
- [Disbursement](disbursement.md)
- [SSO and authentication](../integrations/sso-and-authentication.md)
