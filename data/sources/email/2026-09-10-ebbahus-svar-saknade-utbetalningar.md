---
from: Sara Holm <sara.holm@minso.se>
to: nina.dahl@ebbahus.se
cc: robin@minso.se, jonas.frid@minso.se
subject: Re: September disbursements missing from our ledger
date: 2026-09-10T08:33:00+02:00
message_id: <20260910083300.2214@minso.se>
in_reply_to: <20260909114700.9964@ebbahus.se>
---

Hi Nina,

Thanks — and yes, please send the file listing, that will let us line up exactly where it stops from both ends.

I have logged this as a case and assigned it to Jonas. He is in the migration work for another customer this week, so realistically he looks at it Thursday or Friday. If that is too late for the quarter, say so and I will push it up.

Your first guess is the one I would check too. The staging table rework touches the file writing, so if it shipped and the destination path moved, the symptom would look like this. I do not know off the top of my head whether it has actually shipped yet — Jonas will know.

I will come back to you by Friday even if the answer is "still looking".

Best,
Sara

--
Sara Holm
Support Lead, Minso AB
sara.holm@minso.se

> Ulrika (our accountant) came to me this morning. There is nothing from
> September in the bookkeeping. [...] The last export file she has is dated
> 4 August.
> [...]
> My best guess is that it is connected to the rework Jonas described on
> 11 August — the staging table change.
