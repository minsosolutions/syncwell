# Application workflow

**Document:** `product/application-workflow.md`
**Version:** 5.2 · **Last updated:** 2026-08-14 · **Owner:** Product
**Applies to:** Grantly 2026.8 and later

An application moves through a fixed set of states from draft to decision. Programmes can configure
deadlines, required fields and whether returns to the applicant are allowed, but they cannot add or
remove states.

## States

| State | Set by | Applicant can edit | Next states |
| --- | --- | --- | --- |
| `draft` | Applicant | Yes | `submitted`, `deleted` |
| `submitted` | Applicant | No | `under_review`, `returned`, `withdrawn` |
| `returned` | Caseworker | Yes | `submitted`, `withdrawn` |
| `under_review` | Caseworker | No | `on_hold`, `approved`, `rejected` |
| `on_hold` | Caseworker | No | `under_review`, `rejected` |
| `approved` | Decision maker | No | `paying`, `cancelled` |
| `rejected` | Decision maker | No | — (terminal) |
| `withdrawn` | Applicant | No | — (terminal) |
| `cancelled` | Caseworker | No | — (terminal) |

`paying` and the states after it are covered in [Disbursement](disbursement.md).

> **Note**
> `on_hold` was added in release 2026.8. Tenants that were using a custom label on `under_review` to
> mean "waiting for the applicant" should migrate those cases to `on_hold`; the two states are
> reported separately in statistics.

## Fields set on transition

| Transition | Fields written |
| --- | --- |
| → `submitted` | `submitted_at`, `submission_reference` |
| → `returned` | `returned_at`, `return_reason` (required, free text) |
| → `under_review` | `assigned_caseworker_id`, `review_round_id` |
| → `on_hold` | `hold_reason`, `hold_until` (optional) |
| → `approved` | `decision_date`, `amount_granted`, `decision_maker_id` |
| → `rejected` | `decision_date`, `rejection_reason_code`, `decision_maker_id` |
| → `withdrawn` | `withdrawn_at` |

`decision_date` is the date the decision is formally issued, not the date it was entered into
Grantly. It can be backdated by up to 90 days by a user with the **Decision maker** role. This
matters: retention periods for rejected applications are counted from `decision_date`. See
[Data retention](../legal/data-retention.md).

## Deadlines

A programme has an opening date and a closing date. After the closing date:

- Drafts can no longer be submitted. The submit button is disabled and the portal shows the deadline.
- Submitted applications are unaffected.
- Returned applications can still be resubmitted if the caseworker sets an individual deadline on the
  return. Without an individual deadline, a return issued after the programme closes cannot be
  resubmitted — check this before returning a late application.

Deadlines are enforced at the second, in the programme's timezone. An application submitted at
23:59:59 on the closing date is accepted; 00:00:00 the following day is not.

## Withdrawal

An applicant can withdraw a `submitted` or `returned` application from the portal up until the
programme deadline. After the deadline, only a caseworker can withdraw on the applicant's behalf,
recording who requested it.

Withdrawal is not deletion. A withdrawn application keeps its answers and attachments and is retained
under the withdrawn-application period in [Data retention](../legal/data-retention.md).

## Bulk actions

Caseworkers with the **Programme administrator** role can perform these actions on a selection:

| Action | Limit per action | Notes |
| --- | --- | --- |
| Assign caseworker | 500 applications | |
| Move to `under_review` | 500 applications | Only from `submitted` |
| Generate decision letters | 500 applications | Added in 2026.7 |
| Export | See [Bulk export](../runbooks/bulk-export.md) | |

> **Warning**
> Bulk decisions are not reversible from the interface. An application moved to `rejected` in error
> must be corrected by Support, and the correction is visible in the audit log.

## Audit log

Every state transition writes an audit entry with the actor, the timestamp in UTC, the old and new
state, and any reason text. Audit entries cannot be edited or deleted by tenant users and are
retained for five years.

## Related documents

- [Review and scoring](review-and-scoring.md)
- [Portal for applicants](portal-for-applicants.md)
- [Attachments and file uploads](attachments-and-file-uploads.md)
- [API overview](../integrations/api-overview.md)
