---
title: Support check-in — Nordstad kommun
date: 2026-08-12
duration_minutes: 38
source: granola
attendees: karin.sjo@nordstad.se, per.ek@nordstad.se, sara.holm@minso.se, jonas.frid@minso.se
---

**Sara Holm:** — so it's just us until Karin gets here, she said she'd be five minutes late.

**Per Ek:** She's in the other building today I think.

**Sara Holm:** How's the heat over there?

**Per Ek:** Terrible. The server room is fine and the people are dying, which tells you something about priorities.

**Sara Holm:** [laughs] Okay, while we wait — the ticket list. You have three open. The bulk export one, the one about the e-mail templates, and then the SSO certificate which is really yours, not ours.

**Per Ek:** The certificate is on me, yeah, I have to get it from Sitewise — from, sorry, from our — the identity people.

**Sara Holm:** No rush on that one, it doesn't expire until November.

**Jonas Frid:** The bulk export one I looked at last week. So what's happening is when you export the full round, all statuses, it's building the whole thing in memory before it streams. And when you cross a certain —

**Per Ek:** How big is a certain.

**Jonas Frid:** It depends on the attachments. Roughly, when the row count goes past eleven, twelve thousand it gets slow and then the gateway kills it at sixty seconds.

**Per Ek:** So it's not going to work for the autumn round.

**Jonas Frid:** Not the way it is now. I've got a change that chunks it. It's not shipped.

**Sara Holm:** That's ticket ten-oh-four-two.

**Jonas Frid:** Right.

**Karin Sjö:** Hi, sorry, sorry. Parking. Hi.

**Sara Holm:** No problem. We're on the bulk export.

**Karin Sjö:** Oh — that one. Yeah, that one is annoying.

**Jonas Frid:** So the fix chunks the export. I want to say it'll be in the release at the end of the month.

**Sara Holm:** Can we — should we leave the ticket open until it ships, or —

**Jonas Frid:** Doesn't matter to me.

**Sara Holm:** I'll leave it. [pause] Okay, Karin, you had a thing about the phone calls.

**Karin Sjö:** Yes. So I've now started writing them down, because you asked. Since the — since like the twenty-eighth of July I have nineteen.

**Per Ek:** Nineteen calls about the same thing?

**Karin Sjö:** Nineteen calls where the person says they attached their budget and it isn't there. And I said last time I thought it was them not pressing the button. I'm now less sure.

**Jonas Frid:** What changed?

**Karin Sjö:** Two of them were from the same förening, and the second time the person did it with me on the phone. Step by step. And it still wasn't there afterwards.

**Jonas Frid:** Hm.

**Karin Sjö:** She picked the file, it showed the name, it showed the — there's a little bar —

**Jonas Frid:** Progress bar.

**Karin Sjö:** Yeah. It went all the way. And then she submitted and we both looked and it was gone.

**Jonas Frid:** Was it a big file?

**Karin Sjö:** It was a budget, so — it's an Excel with pictures in it, I think. They put the logo in. They all put the logo in.

**Per Ek:** Everything has a logo.

**Karin Sjö:** And she had been sitting with the form for a long time before, because she called me first, and then we talked for — I don't know, half an hour — and then she filled it in.

**Jonas Frid:** Okay.

**Sara Holm:** Is that useful?

**Jonas Frid:** Maybe. I'd want the — if you have the time of the submission I can look at what the backend saw.

**Karin Sjö:** I have the application number. I don't have the time.

**Jonas Frid:** The number is enough.

**Karin Sjö:** I'll send it over today.

**Jonas Frid:** And if it happens again, ask them how long they'd had the page open. That's — I have a suspicion but I don't want to say it out loud and be wrong.

**Sara Holm:** Say it.

**Jonas Frid:** No, I'll look first.

**Karin Sjö:** [laughs] Fine.

**Per Ek:** Is this the same as the thing where people complain the PDF preview is blank?

**Jonas Frid:** That's different. That's the font thing, that's fixed.

**Per Ek:** Okay.

**Sara Holm:** Karin, in the meantime — do you want us to put something in the applicant help text? Like, if your file doesn't appear, do X.

**Karin Sjö:** Honestly what would help is if the system told them. Right now there's no error. It just — the file isn't there and nobody is told anything. Not them, not us.

**Sara Holm:** Yeah.

**Karin Sjö:** I only find out when they phone.

**Per Ek:** And some of them presumably don't phone.

**Karin Sjö:** Right. That's the part I don't like.

**Sara Holm:** Okay. Jonas, can you look at it this week, and then we'll decide whether it's a ticket or a defect or —

**Jonas Frid:** I can look. I can't promise this week, I have the — the Utvecklingsverket thing.

**Sara Holm:** Next week then.

**Jonas Frid:** Next week.

**Sara Holm:** Good. Anything else? Per?

**Per Ek:** No. Well — the steering meeting on the twenty-sixth, is that still —

**Sara Holm:** As far as I know. Robin owns that one.

**Per Ek:** Fine.

**Karin Sjö:** I'm supposed to be there too apparently.

**Sara Holm:** Lucky you. Okay, thanks all.

**Jonas Frid:** Bye.

**Karin Sjö:** Bye —
