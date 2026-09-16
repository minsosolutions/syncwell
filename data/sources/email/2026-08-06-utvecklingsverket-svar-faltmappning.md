---
from: Jonas Frid <jonas.frid@minso.se>
to: tobias.ahl@ruv.se
cc: robin@minso.se
subject: Re: Field mapping sheet — the eleven unknowns
date: 2026-08-06T15:41:00+02:00
message_id: <20260806154100.6612@minso.se>
in_reply_to: <20260804131400.5507@ruv.se>
---

Tobias,

This is the most useful mapping mail I have had from any customer, so thank you.

Notes back:

- `belopp_beviljat` in öre — found it, thank you, that would have been a very embarrassing dry run.
- `beslutsdatum` semantics change in 2016 — I will carry the original value and add a flag on pre-2016 rows so nobody later assumes they are the same thing.
- `sekretessmarkering` — agreed, this is the one that matters. I have made it a hard filter in the export path rather than a flag someone can forget to set. Non-zero means excluded, full stop.
- `handlaggarkod` — carrying it as-is until Britta rules.

On case identifiers: technically we can store the legacy DNR, index it and make it searchable. We cannot make it the primary key — the format has letters in positions our key does not allow, and the sequence has gaps, which our number series will not accept.

So the question that is left is not "can we", it is "what do you want the handläggare to see and quote". That is yours, not mine. I will raise it at the next meeting but somebody on your side needs to own the decision.

/Jonas

--
Jonas Frid
Backend Engineer, Minso AB
jonas.frid@minso.se

> **Cannot answer alone:**
> - handlaggarkod — I think it is dead. Britta says it feeds the statistics.
> - Case identifiers. This is the one I keep raising. [...] It is not a
>   technical question and I should not be the one answering it.
