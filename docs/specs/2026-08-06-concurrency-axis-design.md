# Design — concurrency review axis

**Status:** accepted, 2026-08-06. **Non-normative.** A dated record of
why the concurrency axis is shaped the way it is. The live rules live
in `plugins/skill-review/skills/skill-review/` — the shared contract in
`SKILL.md`, this axis in `references/concurrency.md`. When this
document and those disagree, they win and this is stale history, not a
finding.

## Problem

`skill-review` ships one axis: a duplication and drift gate over
instruction prose, built on the one-normative-home rule. It answers
"is this rule stated once?" and nothing else.

It has no answer for the other way instruction files fail silently. A
procedure that walks independent work one item at a time reads
perfectly well — it just taxes every executor that follows it,
forever, and nothing in the prose signals the cost. The inverse fails
worse: a procedure that fans out over work whose items share state or
depend on each other's output corrupts results, and reads exactly as
fine.

## Decisions

**Placement — two axis sections in one `SKILL.md`.** Rejected: a
sibling skill (duplicates the shared review contract across two
files — a finding this repo's own gate would raise), and a router
with `references/` per axis (three files to carry two axes). Revisit
the router only if a third axis arrives.

**Direction — bidirectional, including cost and scale.** The axis
judges missed concurrency, unsafe concurrency, and concurrency sized
wrong. A one-directional "parallelize more" reviewer produces harmful
advice: told to fan out file edits without isolation, an executor
corrupts trees.

**Teeth — every finding blocks.** Any finding on either axis flips the
verdict to NOT READY. Rejected: a correctness/advisory severity split.

The accepted risk is an everything-blocks gate crying wolf on arguable
design opinions. It is contained by precision rather than by severity:
a finding must name the specific item set and the evidence for or
against its independence, and each site emits exactly one finding
(most specific type wins). A vague "this could probably be parallel"
does not clear the bar and is not reported.

## The rule

**Shape matches the work.** A procedure's control shape must match the
dependency structure of the work it describes. Independent items fan
out; dependent items serialize; per-item stages pipeline; only
cross-item reasoning earns a barrier. Serializing independent work and
parallelizing dependent work are the same defect — shape mismatched to
dependencies.

Corollaries:

- Independence is a property of the work, not the wording. Prove it
  before recommending fan-out: name the items, name what they do not
  share.
- Fan-out that mutates shared state requires isolation (worktree, or
  provably disjoint paths). Without it the shape is wrong even when
  the items are logically independent.
- A barrier is earned only by cross-item reasoning — dedup across the
  full set, merge, early exit on a total. Needing to flatten, map, or
  filter is not a barrier.
- Batching is part of the shape. Subagents dispatched across turns run
  serially whatever the prose intends.
- Orchestration weight is part of the shape. A fleet for two lookups
  is as wrong as a serial loop over fifty.

## The procedure

Six steps, sharing the drift axis's census → classify → verify spine:

1. **Work-item census** over the FULL text of each artifact, never the
   diff hunk — a procedure's shape is invisible in a hunk. List every
   step operating over a set, its expected size, and whether the
   artifact bounds it.
2. **Independence test** per set: shared mutable state, item N
   consuming item N−1's output, order affecting the result. Record the
   verdict and the evidence. A set whose independence cannot be
   demonstrated is correctly serial and is not a finding.
3. **Shape classification** — name the shape prescribed and the shape
   the dependencies want, from one vocabulary: serial, batched
   fan-out, pipeline, barrier, workflow. A mismatch is a candidate.
4. **Safety pass** over every prescribed fan-out: writers need
   isolation or disjoint paths, order-dependent work must not fan out,
   unbounded sets need a stated cap.
5. **Sizing pass**: agent count against independent parts, deliberate
   per-agent model choice, orchestration weight against the work,
   Workflow's explicit-opt-in requirement, and whether the artifact
   lets the executor scale fan-out to the situation.
6. **Closure pass**: every fan-out has a gather step that reconciles
   results.

Verification is not a seventh step. It moves to the shared review
contract that both axes point at.

## Finding vocabulary

Safety — the shape corrupts or runs away: `UNSAFE` (parallel writers,
no isolation), `ORDERED` (fan-out over order-dependent work),
`UNCAPPED` (unbounded set, no stated cap), `OPTIN` (Workflow-class
orchestration without the explicit user opt-in it demands).

Missed concurrency: `SERIAL` (provably independent items walked one at
a time), `DISPATCH` (subagents spread across turns instead of one
message, so they run serially anyway), `BARRIER` (collect-all where
each item could flow through its stages independently).

Sizing: `HEAVY` (orchestration heavier than the work), `OVERSIZED`
(more agents than independent parts), `UNTIERED` (no deliberate model
choice where subtask difficulty plainly differs), `COSTBLIND` (cost
scales with input and the executor is given no way to size it),
`NOSYNTH` (fan-out with no gather step).

`UNCAPPED` and `COSTBLIND` are the pair most easily confused.
`UNCAPPED` is fan-out width the artifact never bounds; `COSTBLIND` is
a bounded fan-out the executor is given no lever to size down. Where
both apply, `UNCAPPED` wins.

## Report format

Structurally parallel to the drift finding block. The mandatory
`items:` line is the precision bar made structural — a finding cannot
be filed without naming the set and its independence evidence.

```text
<file:line> — <step>: <TYPE>
  items: <the set, and the evidence they are (or are not) independent>
  shape now: <serial | batched fan-out | pipeline | barrier | workflow>
  shape wanted: <same vocabulary>
  fix: <one line>
```

## Document layout

`SKILL.md` gains a shared review contract that owns everything
cross-axis — modes, verify-before-reporting, `file:line` anchors,
report-never-fix, one-finding-per-site, the verdict, and the scope
boundary — followed by an Axis 1 and an Axis 2 section, each owning
its rule, procedure, and finding block.

Three deletions carry as much weight as the additions, and each exists
so the change does not violate the rule the skill enforces:

- `## Findings and verdict` dissolves; each finding block moves under
  its axis, the verdict moves to the shared contract.
- The drift procedure's step 6 (verify before reporting) becomes a
  pointer to the shared contract. Drift goes six steps to five.
- The intro drops its contract clause (anchors, verify, report-never-
  fix), which moves to the shared contract. The intro states the
  problem only.

## Scope boundary and its renderings

The boundary widens from "structure findings only" to structure and
shape: how rules are stated, and how work is dispatched. It still does
not judge meaning, triggering quality, or coverage.

Its normative home is the shared review contract in `SKILL.md`. The
`README.md` intro and Scope section, `plugin.json`, and both
`marketplace.json` descriptions are audience renderings, legitimate
under the one-normative-home rule but DRIFTED if the boundary moves
without them. They are rewritten, not converted to pointers.

## Non-goals

Judging meaning, triggering quality, or coverage. Auto-fixing
findings. Reviewing executable code for concurrency bugs — this axis
reads instruction prose and the procedures it prescribes.

## Release

Lands as `feat:` (minor). `plugin.json` `version` is never touched in
a PR; semantic-release owns that field.
