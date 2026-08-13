# Design — one subagent per review axis

**Status:** accepted, 2026-08-13. **Non-normative** — see `AGENTS.md`,
"Design records", for what that means and what this file may hold. A
dated record of why the axes are dispatched the way they are.

## Problem

Four axes ran in one context. The shared contract had them share a
census too: one subagent per artifact returned every loaded axis's
inventory for it, and classification, verification, and the verdict
all happened in the dispatcher.

That shape has two costs, and the second is the one that bites. Four
axes in one context spend wall-clock as a sum — measured against this
gate's own budget, the security axis alone took 24 minutes, so the
other three are pure addition on top of it. And a census serving four
axes has to return four inventory shapes at once — a grant census, a
work-item census, an operation census — each of which only one
consumer reads in full. The dispatcher then holds all four references
and all four inventories at peak.

## Decisions

**One agent per axis, running its axis end to end.** Each loaded axis
gets its own subagent, which reads the contract and that axis's
reference, censuses for itself, classifies, verifies, and returns its
findings. Rejected: a hybrid keeping the shared census under per-axis
judgment agents — it preserves the read savings but keeps the
four-shaped return contract that couples the axes, which is the defect
rather than the cost. Also rejected: splitting only on direct
invocation and leaving the gate on the shared census, which buys a
second shape to keep in sync forever, in the mode with the least
budget to spare.

**The duplicated reads are accepted, not overlooked.** Three axes now
read each artifact for themselves where one census read it once. That
is this repository's own WIDEREAD, and it is the price of four
contexts that neither compete for one budget nor blend four
vocabularies; wall-clock becomes the slowest axis rather than the sum
of all four. Recorded here because the saving it gives up was the
stated reason for the shared census one release earlier — a later
reviewer restoring that census should know it is reopening a decision,
not fixing an oversight. The companion defect, DISPATCH, is avoided by
construction: the agents go out in one message.

**Agent definitions rather than dispatch prose.** The axis agents ship
as files, so the model tier and the tool roster each have one
normative home in frontmatter instead of living in prose an executor
may or may not honour. Their bodies name the two contract files and
stop; anything more would be four copies of an axis, which is the
worst case the drift axis has. Rejected: inline dispatch instructions,
which leave nothing reusable outside the skill and no place to narrow
capability.

**Frontier tier, on the axis's own authority.** Every agent pins the
top tier — not as a preference, but because `references/token-usage.md`
already exempts review, verification, and adjudication from cost
cutting, and an axis is all three. The tier a supervisor needs is the
hardest decision it makes, and each agent makes its axis's hardest
one.

**Rosters narrower than the union, and the one shell scoped to a
single verb.** Three axes need no shell and no longer hold one. The
drift axis does need one for `tools/dupscan`, and the security axis
demands the narrowest grant the harness can express — which it can:
shipped first-party agents declare `Agent(claude-security:explore)` in
`tools` frontmatter. The obstacle was the recipe, not the syntax. A
build-then-run pair reached the binary through a per-run
`${TMPDIR:-/tmp}` path that no prefix specifier binds, and the first
draft granted a bare shell and recorded that as unavoidable. It was
not: `go run -C <dir> . -min 12 <paths>` changes directory before
module resolution — the very reason the recipe built first — and
collapses the step to one verb. `references/drift.md` now says that,
and the grant is `Bash(go run:*)`.

**The dispatch brief is part of the contract, and the diff reaches one
axis.** The first draft said each agent runs its axis and returns, and
never said what it was handed. Under the gate that is a break rather
than an omission: the pre-change diff exists only in the dispatcher's
prompt and the shell is denied, so an agent briefed with a target
alone cannot run a census that begins "from the diff". The contract
now fixes the brief. Only drift censuses from pre-change state — the
other three read full artifact text and say so — so only drift's brief
carries the diff. Shipping it to all four was this repository's own
WIDEREAD, and at the shipped cap four verbatim copies in one message
overrun any output ceiling besides. That copy travels fenced on an
unforgeable marker with a treat-as-data instruction beside it: a fence
built around this contract does not survive into a fresh context, and
the agent an injection would steer is the one that reports steering.

**The dispatcher is not pinned, and the tier inversion is the accepted
loss.** Four opus agents under a supervisor on the ambient default
inverts the tiering rule, so the gate was changed to pass `--model` —
and changed back. This gate fails open, and `claude -p` exits non-zero
on a tier the install cannot serve; every such exit folds into the
pass-through. A plan without that tier, or a deployment naming models
by full ID, would have switched the gate off silently and permanently,
which is the one failure a quality gate must not have.
`--fallback-model` covers an overloaded model, not an unavailable one.
A supervisor sometimes below its workers is the smaller loss than a
gate nobody notices is gone; the agents keep their pins in
frontmatter, where a bad value costs one axis instead of every review.
The same pass stopped handing the dispatcher all four axis references,
which the agents now read for themselves — the peak-context cost this
design exists to remove was otherwise still paid at the top.

**One test holds the places that have to agree.** The dispatch
identifiers, the four agent files behind them, and the tier those
files pin are read together by nothing: `claude plugin validate`
passes with `agents/` deleted outright, so a rename would drop an axis
from every review while the gate still returned READY. A Go test now
resolves each identifier to its file, checks the name it declares, and
fails on a split tier.

**Reconciliation collapses by fix, not by precedence.** Four
independent agents can each fire on one line, so the gather step needs
a rule the single shared context used to supply implicitly. Rejected:
a fixed axis ranking — mechanical, but it discards a genuine second
defect whose remedy the surviving finding does not deliver, which
makes NOT READY under-report what must actually change. Also rejected:
one finding per site per axis, which re-admits the noise the rule
exists to prevent. The test promoted to the contract is the one
`references/security.md` already applied within its own axis: two
fixes are two findings.

## Non-goals

Narrowing the gate's own roster. The gap it opens in the drift axis is
no longer unrecorded: `references/drift.md` states the degraded pass,
which is where a reviewer looks for it. Rewriting
`2026-08-13-security-axis-design.md`,
whose census decision this supersedes: records are history, and the
live rules win. Judging meaning, triggering quality, or coverage.
Auto-fixing findings.

## Release

Lands as `feat:` (minor). `plugin.json` `version` is never touched in
a PR; semantic-release owns that field.
