# Design — token usage review axis

**Status:** accepted, 2026-08-06. **Non-normative.** A dated record of
why the token-usage axis is shaped the way it is. The live rules live
in `plugins/skill-review/skills/skill-review/` — the shared contract in
`SKILL.md`, this axis in `references/token-usage.md`. When this
document and those disagree, they win and this is stale history.

Follows [2026-08-06-concurrency-axis-design.md](2026-08-06-concurrency-axis-design.md),
whose deferred router decision this design executes.

## Problem

Skills cost tokens in ways their prose never signals. A procedure that
spends a frontier model on what a script could do reads fine forever;
so does one that supervises a fan-out with a model too weak to
decompose the work, which is worse — a bad decomposition wastes every
worker it dispatches, so the saving is multiplied into a loss.

Neither existing axis catches this. Axis 1 asks whether a rule is
stated once. Axis 2 asks whether control flow matches dependencies —
correctness of shape, not its price.

## Decisions

**Axis 2's Sizing cluster migrates here.** `HEAVY`, `OVERSIZED`,
`UNTIERED`, and `COSTBLIND` were cost findings living in the shape
axis. Leaving them there while adding a cost axis would put model
tiering under two owners — a DRIFTED finding against this repo's own
axis 1, and a guarantee the two would diverge on first edit. Axis 2
drops to eight types and becomes purely shape correctness; axis 3 owns
all cost.

**The router refactor lands with this axis.** The concurrency design
deferred it with an explicit trigger: revisit if a third axis arrives.
It has. Three rules, three procedures, and fourteen more finding types
would have taken `SKILL.md` past 300 lines — billed on every trigger,
before any work happens, including triggers that turn out to be the
wrong skill. That is this axis's own `MONOLITH` finding, and the gate
would have failed its new axis on sight. `SKILL.md` becomes a router
owning the shared contract; each axis becomes a reference loaded on
demand.

**Every finding still blocks; precision scales instead of severity.**
Three axes bring the vocabulary to 26 blocking types, up from four.
Rather than introduce the advisory tier rejected for axis 2, each type
names the evidence it requires, and a finding whose evidence cannot be
filled in is not reported. The rule lives once in the shared contract
(evidence-or-silence); the per-type requirements live with their axis.
The judgment-heavy types get the strictest bars: `WEAKSUPER` must name
the decomposition decision the tier cannot make and the worker count
at risk; `OVERTIER` must name the deterministic check that backstops
going cheaper, and reports nothing without one.

## The rule

**Cheapest sufficient mechanism.** Every operation runs on the
cheapest mechanism that can do it correctly: code before model, small
model before large, targeted read before whole file, load-on-demand
before load-always. Cost is justified only by capability actually
required.

*Sufficient* is load-bearing. An under-tiered supervisor is not cheap,
it is insufficient — which is why the rule catches under-tiering and
over-tiering with one statement instead of two.

Corollaries: cheapness upstream is amplified downstream · cost paid
per invocation dominates cost paid per run · deterministic output
belongs in code · cheap gates precede expensive work · judgment is not
a cost centre.

## Finding vocabulary

Model class: `WEAKSUPER`, `OVERTIER`, `TIEREDJUDGE`, `UNTIERED`.

Code versus model: `REGEN` (agent made to write code the skill could
ship), `INCONTEXT` (bulk data through context a script could reduce),
`MODELDET` (model doing deterministic work a tool does better).

Footprint: `MONOLITH`. Context economy: `VERBOSERET`, `WIDEREAD`.
Ordering: `LATEGATE`. Orchestration cost: `HEAVY`, `OVERSIZED`,
`COSTBLIND`.

## Layout

```text
skills/skill-review/
  SKILL.md                    router: intro, axis table, shared
                              contract, repo-wide audit
  references/drift.md         axis 1
  references/concurrency.md   axis 2
  references/token-usage.md   axis 3
```

`SKILL.md` drops from 9,075 to 3,624 bytes — a 60% cut in
always-billed content while total content grows to 15,570 across four
files. A structure-only request now loads the router plus one
reference instead of every axis.

The shared contract gains one bullet, evidence-or-silence, and updates
three renderings from two axes to three (modes, verdict, scope
boundary). Each reference opens by pointing at the contract rather
than restating it.

## Non-goals

Judging meaning, triggering quality, or coverage. Auto-fixing
findings. Measuring actual token consumption — this axis reads
instruction prose and the mechanisms it prescribes, not telemetry.

## Release

Lands as `feat:` (minor). `plugin.json` `version` is never touched in
a PR; semantic-release owns that field.
