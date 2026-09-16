---
title: Integration working session — Ebbahus Stiftelse
date: 2026-08-11
duration_minutes: 47
source: granola
attendees: nina.dahl@ebbahus.se, jonas.frid@minso.se, mia.berg@minso.se
---

**Jonas Frid:** — yeah I'm back, I got back Sunday.

**Nina Dahl:** Where were you?

**Jonas Frid:** Gotland. It rained for six days and then it was beautiful on the day we left.

**Nina Dahl:** That is the law.

**Mia Berg:** Okay. So this is the deep dive on the bookkeeping export. Jonas has been in it since Monday.

**Nina Dahl:** Good, because Ulrika asked me yesterday and I had nothing to tell her.

**Jonas Frid:** So — first, the cost centre thing you described. You were right. It is the split-year grants. When a grant is booked across two fiscal years the second row is generated from the schedule and not from the decision, and the schedule rows don't carry the cost centre. That's a one-line fix, honestly.

**Nina Dahl:** Oh, good.

**Jonas Frid:** It's not shipped yet. I found something else while I was in there and I'm — it's a bit tangled.

**Nina Dahl:** Something worse?

**Jonas Frid:** Not worse. Different. The export job retries on failure, and the retry doesn't deduplicate. So if the first attempt writes half the file and then dies, the retry writes the whole file, and you get some rows twice.

**Nina Dahl:** Wait. Has that happened?

**Jonas Frid:** Once that I can see, in May. Three rows.

**Nina Dahl:** Three rows in May. [pause] I'd have to ask Ulrika whether she caught that.

**Jonas Frid:** It's a small number.

**Nina Dahl:** It's a small number of payments to real organisations who either got paid twice or — no, sorry, it's the bookkeeping, not the actual payment.

**Jonas Frid:** It's the bookkeeping. The money is fine. The ledger might not be.

**Nina Dahl:** Right. Okay. That's less alarming.

**Mia Berg:** Still needs fixing.

**Jonas Frid:** Yes. So what I want to do is rework the job so it writes to a staging table and then commits once, which kills both problems. That's a couple of days.

**Nina Dahl:** And in the meantime?

**Jonas Frid:** In the meantime it runs as it is. It works most nights.

**Nina Dahl:** Most nights.

**Jonas Frid:** Most nights.

**Mia Berg:** When can you start on it?

**Jonas Frid:** I can start next week. I've got the Nordstad export thing and then the — the migrering for the agency, which is eating a lot.

**Mia Berg:** Okay, so realistically end of August.

**Jonas Frid:** Realistically end of August, beginning of September.

**Nina Dahl:** As long as it's before the November payment run.

**Jonas Frid:** It will be.

**Nina Dahl:** Can I have that in the ticket? So I can show Gustav.

**Mia Berg:** I'll update the ticket.

**Nina Dahl:** Thank you.

**Jonas Frid:** While I have you — when I'm working on this I sometimes need to stop the nightly job so it doesn't run on top of me. Is that a problem for you?

**Nina Dahl:** For a night? No. Ulrika does the reconciliation weekly, not daily.

**Jonas Frid:** Okay, good. I'll try to keep it to when I'm actually in there.

**Nina Dahl:** Just tell me if it's more than a couple of days, because then the file dates get strange and she notices.

**Jonas Frid:** Sure.

**Mia Berg:** Should we write that down as a — like, a rule? Notify Nina if the job is off for more than two nights?

**Nina Dahl:** That would be sensible.

**Jonas Frid:** Fine by me.

**Mia Berg:** Okay.

**Nina Dahl:** The other thing on my list is the autumn call. You said mid-August for changes.

**Mia Berg:** I did.

**Nina Dahl:** I went through the test form on Friday. Two things. The intresseanmälan asks for the organisation number before it asks for the organisation name, which confuses people, and the confirmation email says "your application" when it isn't an application yet, it's an expression of interest.

**Mia Berg:** Both doable. The field order is config, the email text is config.

**Nina Dahl:** Then that's all I have.

**Mia Berg:** I'll do them today and ping you.

**Nina Dahl:** Perfect.

**Jonas Frid:** Can I ask about the documents thing Gustav mentioned last time? Someone said applicants were losing documents.

**Nina Dahl:** Oh — the PDF thing. Gustav thinks it's our template.

**Jonas Frid:** Is it always PDFs?

**Nina Dahl:** I don't know. I wasn't the one taking the calls, that was Emil, and Emil is here one day a week.

**Jonas Frid:** Okay.

**Nina Dahl:** It was two applicants out of two hundred. I honestly haven't thought about it since.

**Jonas Frid:** Fair enough.

**Mia Berg:** Anything else? I have five minutes then I have another —

**Nina Dahl:** No, I think that's it. Thank you both, this was actually useful, usually I come out of these more confused.

**Jonas Frid:** [laughs] High praise.

**Mia Berg:** Bye Nina.

**Nina Dahl:** Hej då —
