---
title: Migration planning — Regionalt Utvecklingsverket
date: 2026-07-30
duration_minutes: 58
source: granola
attendees: hakan.nyberg@utvecklingsverket.se, tobias.ahl@ruv.se, robin@minso.se, jonas.frid@minso.se
---

**Håkan Nyberg:** — and then they moved the whole department to the third floor, so nobody can find anybody.

**Robin Bauhn:** Is Tobias joining?

**Håkan Nyberg:** He should. He's on the other network sometimes and then the invite goes to the wrong — he has two addresses, it's a whole thing.

**Robin Bauhn:** The ruv.se one?

**Håkan Nyberg:** The ruv.se one is the technical side. Historically it was a separate authority. It's supposed to be merged and it is not merged.

**Jonas Frid:** That explains some things about the extract files.

**Håkan Nyberg:** It explains almost everything about this agency.

**Tobias Ahl:** Hello. Sorry. I was in the wrong meeting.

**Robin Bauhn:** No problem. We were just talking about your two email addresses.

**Tobias Ahl:** Everyone does.

**Robin Bauhn:** Okay. So — the purpose today is to agree the migration approach and the sequence. Jonas has been through the extract.

**Jonas Frid:** I have. So, the shape of it: there are two source systems, not one. There's the case system, which is the big one, and there's the old grants register which I understand was frozen in —

**Tobias Ahl:** 2021. It's read-only. We keep it because of the arkivlag.

**Jonas Frid:** Right. And the plan we had was to migrate only the case system and leave the register alone.

**Håkan Nyberg:** That was the plan.

**Jonas Frid:** Does that still hold?

**Håkan Nyberg:** I believe so. Tobias?

**Tobias Ahl:** With one asterisk. There are cases in the live system that reference decisions in the frozen register. If we migrate the case and the reference points at a system nobody can reach from the new one, we've broken a chain that an auditor will one day walk.

**Håkan Nyberg:** How many?

**Tobias Ahl:** I'd have to count. Thousands, not tens of thousands.

**Robin Bauhn:** Can we carry the reference as a text field? So the chain is at least documented, even if it isn't clickable.

**Tobias Ahl:** Maybe. It depends what we do with identifiers in general, which is a bigger question and I'd rather take it separately.

**Robin Bauhn:** Okay, park that.

**Tobias Ahl:** Park it.

**Jonas Frid:** Second thing. Field mapping. I sent the sheet on Monday. There are eleven fields I couldn't map and four where I mapped them but I'm guessing.

**Tobias Ahl:** I saw it. I've done about half.

**Jonas Frid:** No rush, but the dry run needs it.

**Tobias Ahl:** When's the dry run?

**Robin Bauhn:** We had penciled the week of the seventh of September.

**Håkan Nyberg:** That's tight but it's possible. The thing that worries me is not the technology. It's that the verksamhet has to validate the output and the verksamhet is four people who also have day jobs.

**Robin Bauhn:** How long do they need?

**Håkan Nyberg:** Two weeks with a following wind.

**Robin Bauhn:** So the dry run output goes to them the — eleventh, twelfth, and they come back end of September.

**Håkan Nyberg:** Roughly.

**Robin Bauhn:** And then the real cutover in — October?

**Håkan Nyberg:** October or November. I'd like to say October. I suspect November.

**Jonas Frid:** From my side either is fine as long as it isn't December.

**Håkan Nyberg:** Nothing happens in December.

**Robin Bauhn:** Noted. Third thing — Elin asked me to raise the processing agreement. She's not here.

**Håkan Nyberg:** She's on leave until the eighteenth.

**Robin Bauhn:** Right. She wrote to me before she went saying the current agreement doesn't cover the new processing categories. I told her we'd look at it.

**Håkan Nyberg:** She will chase that. She is very good at chasing.

**Robin Bauhn:** I got that impression from the email.

**Håkan Nyberg:** [laughs] Take her seriously. It's the one thing that can actually stop this project.

**Robin Bauhn:** Understood.

**Tobias Ahl:** Can I go back to the identifiers for a moment, because I want to at least name it while everyone is here.

**Robin Bauhn:** Go.

**Tobias Ahl:** Our case numbers are meaningful. They encode the year and the unit. People know them by heart. Handläggare quote them on the phone. If the new system issues its own identifiers and the old number becomes a field somewhere, that is a change to how the entire organisation refers to its own work.

**Håkan Nyberg:** Mm.

**Tobias Ahl:** I don't know yet whether Grantly can keep them.

**Jonas Frid:** It can keep them as a secondary identifier. Whether it can use them as the primary — I'd have to check. I think not, because of the format.

**Tobias Ahl:** Then that's a decision somebody has to make.

**Håkan Nyberg:** Somebody meaning me.

**Tobias Ahl:** Somebody meaning not me.

**Håkan Nyberg:** [laughs] Fine. Put it on the list for the next one.

**Robin Bauhn:** I'll carry it forward.

**Håkan Nyberg:** Is there anything else? I have a departementsmöte.

**Robin Bauhn:** I think that's it. Jonas, you'll chase the field mapping.

**Jonas Frid:** I'll chase.

**Håkan Nyberg:** Thank you all. Tobias, stay a minute, I want to ask you about the other thing.

**Tobias Ahl:** Sure.

**Robin Bauhn:** We'll drop off then. Thanks —
