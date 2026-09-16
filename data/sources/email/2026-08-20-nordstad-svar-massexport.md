---
from: Sara Holm <sara.holm@minso.se>
to: karin.sjo@nordstad.se
cc: jonas.frid@minso.se
subject: Re: Bulk export still times out (ticket 10042)
date: 2026-08-20T09:18:00+02:00
message_id: <20260820091800.3301@minso.se>
in_reply_to: <20260818142700.7745@nordstad.se>
---

Hi Karin,

Sorry for the slow reply — I was out Tuesday.

Two things:

**Short term.** Run the export outside working hours. The timeout is caused by the export building the whole result in memory before it streams, and it is much more likely to get through when the system is not busy. Not elegant, but it works today.

**The fix.** Jonas has chunked the export so it streams incrementally. That ships in the release on 31 August. After that the size of the round should not matter.

I have marked 10042 as solved on the basis that the fix is merged and there is a workaround in the meantime — shout if you disagree and I will reopen it.

Best,
Sara

--
Sara Holm
Support Lead, Minso AB
sara.holm@minso.se

> Following up on 10042. I tried the full-round export three times today:
> - 09:12 — failed after about a minute, white page
> [...]
> I am chasing whether there is anything I can do this week, because the
> sammanställning for the nämnd is due on the 27th.
