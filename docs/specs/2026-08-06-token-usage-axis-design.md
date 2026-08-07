# Design — token usage review axis

**Status:** accepted, 2026-08-06. **Non-normative** — see `AGENTS.md`,
"Design records", for what that means and what this file may hold. A
dated record of why the token-usage axis is shaped the way it is.

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

## Non-goals

Judging meaning, triggering quality, or coverage. Auto-fixing
findings. Measuring actual token consumption — this axis reads
instruction prose and the mechanisms it prescribes, not telemetry.

## Release

Lands as `feat:` (minor). `plugin.json` `version` is never touched in
a PR; semantic-release owns that field.
