---
title: Delivery sync — Nordstad kommun
date: 2026-09-09
duration_minutes: 41
source: granola
attendees: anna.lind@nordstad.se, per.ek@nordstad.se, karin.sjo@nordstad.se, robin@minso.se, jonas.frid@minso.se
---

**Per Ek:** — no but the thing is they repaved it and then two weeks later they dug it up again for the fibre.

**Robin Bauhn:** That's universal, I think. Hi Karin.

**Karin Sjö:** Hi. I can only stay half.

**Robin Bauhn:** Okay, then let's do your bits first. Jonas is here too.

**Jonas Frid:** Hej.

**Karin Sjö:** Good, because I want to talk about the uploads.

**Jonas Frid:** Yeah. So — I looked at the application numbers you sent. There were, I think, eleven that I could actually correlate.

**Karin Sjö:** And?

**Jonas Frid:** In nine of them there is no upload at all on our side. Nothing arrived. The form submission arrived, the attachment did not.

**Karin Sjö:** So they didn't attach it.

**Jonas Frid:** That's one reading. The other reading is the upload was rejected before it got to the part that logs. And the reason I lean that way is — in three of them there's a rejected request in the edge logs about ninety seconds before the submission. Big body. And the response is a four-oh-one.

**Per Ek:** Four-oh-one is auth.

**Jonas Frid:** Yes.

**Per Ek:** So their session is dead but the form still works?

**Jonas Frid:** The form page is already loaded, so yes, it looks fine to them. The token that the upload endpoint checks is a different, shorter-lived one. If they sit on the page for a long time —

**Karin Sjö:** Which is exactly what happens. They sit there. They phone us, they read the questions, they go and find the file —

**Jonas Frid:** Right.

**Anna Lind:** Sorry, I was on mute. Is this — does this mean applications have been assessed without their budgets?

**Karin Sjö:** Yes. That's what I've been saying since July.

**Anna Lind:** Right.

**Jonas Frid:** I want to be careful, I haven't proven it. Eleven cases, three with a matching log line, and the logs only go back thirty days so I can't go further.

**Robin Bauhn:** But it's enough to raise as a defect.

**Jonas Frid:** It's enough for me to raise it, yeah. I'll write it up.

**Robin Bauhn:** Do that today if you can.

**Jonas Frid:** Mm.

**Karin Sjö:** And meanwhile? Because the autumn round is —

**Anna Lind:** Soon.

**Karin Sjö:** Soon. What do I tell people?

**Jonas Frid:** Tell them to reload the page before they attach. That will work. It's stupid but it will work.

**Karin Sjö:** I can't put that on the website.

**Anna Lind:** You can put "if your file does not appear in the summary, reload and try again". That's not a lie.

**Karin Sjö:** Fine.

**Robin Bauhn:** And we'll do a proper fix. Jonas, is it a big fix?

**Jonas Frid:** The fix is small. The testing is not small.

**Robin Bauhn:** Okay.

**Karin Sjö:** I have to jump. Was there anything else for me?

**Robin Bauhn:** The export — did the release help?

**Karin Sjö:** Oh. I haven't tried it in the day. I've just kept doing the night thing.

**Robin Bauhn:** [laughs] Try it in the day.

**Karin Sjö:** I'll try it in the day. Bye, everyone.

**Per Ek:** Bye.

**Robin Bauhn:** Okay. Per, the migration.

**Per Ek:** So — Registraturen came back. Decisions and the original application, not the appendices. Which is what Anna said in the first place.

**Anna Lind:** I did say.

**Per Ek:** You did say. So the volume is about a third of what we thought and we don't need to buy the media transfer from Sedermera — from the old vendor, sorry, I never remember the name.

**Robin Bauhn:** Which one is it?

**Per Ek:** Civitas Nord? Something.

**Anna Lind:** It doesn't matter.

**Robin Bauhn:** Okay, so we do the migration without the appendices, which means the dry run scope shrinks. Jonas, does the cutover window still look right?

**Jonas Frid:** Should be shorter, if anything. I'd still want the full window booked.

**Robin Bauhn:** Fine, we keep the window as planned.

**Per Ek:** I've told drift to keep the Friday.

**Robin Bauhn:** Good.

**Anna Lind:** Can I ask about the retention question again? The one I raised at the steering.

**Robin Bauhn:** The rejected applications one.

**Anna Lind:** Yes. Our jurist has now been asked by the kommunstyrelse to give a note, and she will not write anything that she cannot source. She needs it from you, in writing.

**Robin Bauhn:** Understood. Let me get Sara to pull the official line and we'll send it.

**Anna Lind:** When?

**Robin Bauhn:** I'll come back to you on when.

**Anna Lind:** [pause] Okay.

**Robin Bauhn:** Anything else? Per?

**Per Ek:** Nothing from me.

**Anna Lind:** No. Thank you.

**Robin Bauhn:** Thanks all.

**Jonas Frid:** Hej —
