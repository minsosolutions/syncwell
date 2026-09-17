# One prompt for every Source and Customer; config holds only what varies

The kit was built expecting per-Source and per-Customer prompts — `todo.md` assumes them and
README undecided #4 names them as suspected, unbuilt and possibly wrong. Now that Runs have
actually produced Findings, the suspicion can be settled from evidence rather than from
architecture taste.

**No Source needed its own prompt.** One prompt read Slack, email, transcripts and Zammad in a
single sweep, and the Findings that matter are the ones that cross them: a nightly job
disabled in chat against a ledger gap reported in mail, a ticket closed as solved against chat
saying it still happens. A per-Source prompt is a per-Source *pass*, and a pass that sees one
Source cannot find any of those. The thing we thought was configuration was the thing we were
building against.

**No Customer needed its own prompt either.** The variation between a municipality, a
foundation and an agency is entirely in their records — compliance language, volume, who
writes — and the Agent reads those records. A per-Customer prompt would restate in config
something already present in the Customer Directory, and would then have to be maintained
against it. A knob with the same value across all three Customers is a knob that should not
exist.

So `config/syncwell.json` holds only what a Run cannot read off the data: the Routing Rule keys
(domains, channel prefix, Zammad organization), each Customer's Linear project, and the two
endpoints a deployment has to be told about. Everything that differs *per Customer in kind*
lives in the Customer Directory, where it is data.

Two places take per-Source and per-Customer variation today without any of it being prompt
text, and the split is worth naming. The **file list** in the prompt is per-Customer and is
computed in Go, because the Agent has no search tool and discovery must be deterministic.
**Absence** is per-Source in code — only email has a thread relation, so only email can be
asked whether a reply exists (#16). Both are behaviour, not wording.

What would change our mind: a Source whose records need decoding before they can be judged —
a binary export, a language the Agent reads badly — or a Customer whose contract says
something the records never state. Neither is in this data. The first one to appear gets a
file in the Customer Directory, read by the Run and appended to the prompt; not a prompt
template per Source, which multiplies the passes we just argued against.
