---
from: Tobias Ahl <tobias.ahl@ruv.se>
to: jonas.frid@minso.se
cc:
subject: Field mapping sheet — the eleven unknowns
date: 2026-08-04T13:14:00+02:00
message_id: <20260804131400.5507@ruv.se>
in_reply_to:
---

Jonas,

I have been through your sheet. Answers to nine of the eleven below. Two I cannot answer alone.

**Answered:**

- `arendetyp_kod` — maps to your `case_type`. The values 4 and 7 are both "regionalt projektstöd"; 7 replaced 4 in 2018 and nobody migrated the old rows. Map both.
- `handlaggare_sign` — three-letter initials, free text, not reliable. Map to a text field, do not try to resolve to a user.
- `beslutsdatum` — beware: for cases before 2016 this is the date of the protocol, not the decision. Same field, different meaning.
- `belopp_beviljat` — öre, not kronor, in the old register. Kronor in the case system. Yes, really.
- `medfinansiering_pct` — can be null, can be 0, and these mean different things.
- `avslutsorsak` — code list attached in my previous mail, 14 values, 3 unused since 2019.
- `kommunkod` — standard SCB codes. Some are historical municipalities that no longer exist.
- `projektperiod_start` / `projektperiod_slut` — straightforward.
- `sekretessmarkering` — boolean in practice, integer in the schema. Anything non-zero means the case is restricted and must not appear in exports.

**Cannot answer alone:**

- `handlaggarkod` — I think it is dead. Britta says it feeds the statistics. Britta is usually right, so assume it is live until she tells me otherwise.
- Case identifiers. This is the one I keep raising. Whether the legacy DNR survives as a real identifier or becomes a field. It is not a technical question and I should not be the one answering it.

The `sekretessmarkering` one I would put in bold in your own notes. If those rows leak into an export we have a genuine incident, not an inconvenience.

/Tobias

Tobias Ahl
IT-arkitekt
Regionalt Utvecklingsverket (RUV)
tobias.ahl@ruv.se
