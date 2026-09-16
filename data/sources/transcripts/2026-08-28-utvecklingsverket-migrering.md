---
title: Migration working meeting — Regionalt Utvecklingsverket
date: 2026-08-28
duration_minutes: 51
source: granola
attendees: hakan.nyberg@utvecklingsverket.se, tobias.ahl@ruv.se, elin.sund@utvecklingsverket.se, jonas.frid@minso.se, robin@minso.se
---

**Jonas Frid:** — I can share, hold on. Can you see the sheet?

**Tobias Ahl:** I see a spreadsheet, yes.

**Jonas Frid:** Good. That's the mapping as it stands.

**Håkan Nyberg:** Elin is joining but she said she'd be fifteen late, she's in the other thing.

**Robin Bauhn:** Fine, we'll do the technical part first.

**Håkan Nyberg:** How was the dry run prep?

**Jonas Frid:** Mixed. The extract loads. The volumes match — three hundred and eleven thousand cases, which is what you said. Where it gets messy is the attachments, which is the usual, and the — honestly the worst part is the free text fields that have HTML in them from the old editor.

**Tobias Ahl:** From the old old editor. There were two.

**Jonas Frid:** There were two, and they use different tags, and one of them isn't valid HTML.

**Tobias Ahl:** That would be the one from 2014.

**Jonas Frid:** It would.

**Robin Bauhn:** Does it block the dry run?

**Jonas Frid:** No. It makes some records look ugly. It's a cleanup pass, I can do it, I just don't want anyone opening a record in the dry run output, seeing angle brackets, and filing it as a disaster.

**Håkan Nyberg:** They will do exactly that.

**Jonas Frid:** I know.

**Håkan Nyberg:** I'll warn them.

**Robin Bauhn:** Put it in the note that goes with the output.

**Håkan Nyberg:** I'll put it in the note.

**Jonas Frid:** The eleven unmapped fields — Tobias, you did most of them.

**Tobias Ahl:** I did nine. Two I can't answer alone. One is "handläggarkod", which I think is dead, but Britta says it's used in the statistics, and Britta is usually right.

**Jonas Frid:** Then we carry it.

**Tobias Ahl:** Carry it.

**Jonas Frid:** And the other?

**Tobias Ahl:** The other is the one I've been raising since — I put it in the channel a week ago and nobody replied, which is fine, everyone was away.

**Robin Bauhn:** Which one?

**Tobias Ahl:** Case identifiers. Whether the legacy case numbers survive the migration as real identifiers or become a historical field.

**Robin Bauhn:** Right, we parked this at the end of July as well.

**Tobias Ahl:** We parked it in July and then I raised it again on the twentieth and now I'm raising it here.

**Håkan Nyberg:** Remind me of the shape of it.

**Tobias Ahl:** The old number is DNR-year-unit-sequence. Everyone uses it. It is on paper decisions, it is in other authorities' systems, it is in letters we sent in 2019. If Grantly issues its own number and the old one becomes a searchable field, then every time someone quotes an old number we have to do a lookup rather than a direct — it works, but it is a different world.

**Jonas Frid:** And technically — I checked after last time — we can store it and index it and search it. We cannot make it the primary key. The format doesn't fit and the sequence has gaps.

**Tobias Ahl:** Which I expected.

**Håkan Nyberg:** So the decision is: do we accept a new primary number with the old one searchable, or do we do something clever.

**Jonas Frid:** I would strongly prefer you don't do something clever.

**Håkan Nyberg:** [laughs] Noted.

**Tobias Ahl:** The problem is it isn't only a technical decision. It affects the handläggare, it affects the letters, it affects what we tell other authorities. That's not mine to decide.

**Håkan Nyberg:** No. It's mine, probably, or it's the styrgrupp's.

**Tobias Ahl:** And the dry run is in ten days.

**Jonas Frid:** The dry run doesn't need it decided. The cutover does.

**Tobias Ahl:** How late can it be decided?

**Jonas Frid:** If you want the old numbers visible in the dry run output, I need to know before the — before next Friday. If you're happy for the dry run to show new numbers and we discuss the presentation later, then it can wait.

**Håkan Nyberg:** Let's not take it now. Half the people who should be in the room aren't. Let's take it next time.

**Tobias Ahl:** Next time is the — tenth?

**Håkan Nyberg:** The tenth. Put it first on the agenda.

**Robin Bauhn:** Okay. Moving on then —

**Elin Sund:** Sorry. Hello. I'm here.

**Håkan Nyberg:** Welcome. We just finished the identifiers.

**Elin Sund:** I don't need the identifiers. I need the agreement.

**Robin Bauhn:** It's with our legal. I chased on Tuesday.

**Elin Sund:** And what did they say?

**Robin Bauhn:** They said this week or early next.

**Elin Sund:** Robin, I want to be direct, because I think in these meetings we are all very polite and then nothing moves. We are processing personal data under an agreement that does not describe what we are actually doing. That is not a risk I can carry, and it is not a risk I can carry quietly, because my job is to not carry it quietly.

**Robin Bauhn:** Understood.

**Elin Sund:** I will put it in writing so that it is documented on my side.

**Robin Bauhn:** That's fair. Please do.

**Håkan Nyberg:** Elin, is there a date attached to this, so we all know?

**Elin Sund:** I will put a date in the letter.

**Håkan Nyberg:** Fine.

**Jonas Frid:** Is that — should I keep going on the migration or is that in question?

**Elin Sund:** Keep going. I am not stopping anything today.

**Jonas Frid:** Okay.

**Robin Bauhn:** Let's leave it there. Next meeting the tenth, identifiers first.

**Håkan Nyberg:** Thank you all.

**Tobias Ahl:** Hej —
