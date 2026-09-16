# Attachments and file uploads

**Document:** `product/attachments-and-file-uploads.md`
**Version:** 2.6 · **Last updated:** 2026-08-14 · **Owner:** Product
**Applies to:** Grantly 2026.8 and later

Applicants attach documents from the portal; caseworkers attach documents from the back office. Both
use the same upload service, limits and virus scanning.

## Accepted file types

| Category | Extensions |
| --- | --- |
| Documents | `.pdf`, `.doc`, `.docx`, `.odt`, `.rtf` |
| Spreadsheets | `.xls`, `.xlsx`, `.ods`, `.csv` |
| Images | `.jpg`, `.jpeg`, `.png`, `.tif` |
| Archives | `.zip` (contents are scanned, not expanded) |

The type check is on the file's content type, not only the extension. A `.pdf` that is not a PDF is
rejected with `unsupported_file_type`. Extension casing is ignored.

Executables, scripts and macro-enabled Office formats are rejected.

## Limits and quotas

| Limit | Value | Scope |
| --- | --- | --- |
| Maximum file size | **25 MB** | Per file |
| Maximum total attachments | 200 MB | Per application |
| Maximum number of attachments | 50 | Per application |
| Maximum filename length | 200 characters | Per file |
| Upload request timeout | 300 seconds | Per file |
| Concurrent uploads | 3 | Per applicant session |

A programme can lower the per-file limit but cannot raise it above 25 MB. The per-application total
is fixed.

> **Note**
> The 25 MB limit is checked in the browser before the upload starts and again by the upload service.
> A file over the limit is rejected immediately with `file_too_large` and the applicant sees the
> limit and the file's actual size in the error message.

## Upload behaviour

1. The applicant selects a file. The browser checks the extension and size.
2. The portal requests a short-lived upload ticket from the API. The ticket is bound to the
   applicant's session and to the target application.
3. The file is uploaded in chunks. A progress bar shows the percentage uploaded.
4. On completion the file is queued for virus scanning and appears in the attachment list with the
   state `scanning`.
5. When the scan passes, the state becomes `available`. If the scan fails, the state becomes
   `quarantined` and the applicant is notified by email.

Large files are uploaded in 5 MB chunks. If an individual chunk fails, the client retries that chunk
up to three times before reporting `upload_failed` to the applicant.

### Sessions and upload tickets

Uploads are performed against the applicant's active portal session. Applicant sessions have an idle
timeout of 30 minutes; see
[Portal for applicants](portal-for-applicants.md#sessions-and-timeouts).

The upload ticket obtained in step 2 is valid for 30 minutes and is refreshed automatically by the
portal while an upload is in progress, so an upload that takes longer than the remaining session time
completes normally.

> **Note**
> If the applicant's session has expired before the upload starts, the portal detects this when the
> upload ticket is requested, signs the applicant in again and then starts the upload. The applicant
> does not lose the selected file and does not need to reselect it.

An upload is never discarded without the applicant being told. Every upload ends in one of four
visible outcomes: the file appears in the attachment list, the applicant sees an error message, the
applicant is asked to sign in again, or the file appears as `quarantined`. There is no path in which
an attachment is accepted by the browser and then silently dropped.

## Attachment states

| State | Meaning | Visible to applicant | Visible to caseworker |
| --- | --- | --- | --- |
| `uploading` | Chunks in transit | Yes | No |
| `scanning` | Queued for or undergoing virus scan | Yes | Yes, marked pending |
| `available` | Scanned and usable | Yes | Yes |
| `quarantined` | Scan failed | Yes, cannot download | Yes, cannot download |
| `deleted` | Removed before submission | No | No |

Attachments can be deleted while the application is in `draft`. After submission they are immutable;
a replacement is added as a new file. See [Application workflow](application-workflow.md#states).

## Virus scanning

Scanning is asynchronous and normally completes within 60 seconds. A quarantined file is retained for 30 days and then hard deleted. The application record keeps a note
that a file was quarantined, including the filename and the timestamp.

## Retention

Attachments follow the retention period of the application they belong to and are hard deleted, not
anonymised, when that period elapses. See [Data retention](../legal/data-retention.md).

## Downloading

Caseworkers can download attachments individually or as a ZIP for the whole application. A ZIP over
500 MB is not generated in the browser; use [Bulk export](../runbooks/bulk-export.md).

## Common errors

| Error code | Cause | What the applicant should do |
| --- | --- | --- |
| `file_too_large` | File over 25 MB | Reduce the file size or split it |
| `unsupported_file_type` | Content type not on the accepted list | Convert to PDF |
| `attachment_limit_reached` | 50 files or 200 MB reached | Remove a file first |
| `upload_failed` | Three chunk retries failed | Retry; check the network connection |
| `session_expired` | Session ended and could not be renewed | Sign in again and retry the upload |
