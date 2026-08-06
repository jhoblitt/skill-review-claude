# Axis 2 — concurrency

Runs under the shared review contract in `SKILL.md`, which owns modes,
batching, verification, anchoring, one-finding-per-site, the verdict,
and the scope boundary.

## The shape-matches-the-work rule

A procedure's control shape must match the dependency structure of the
work it describes. Independent items fan out; dependent items
serialize; per-item stages pipeline; only cross-item reasoning earns a
barrier. Serializing independent work and parallelizing dependent work
are the SAME defect — shape mismatched to dependencies. Corollaries:

- Independence is a property of the work, not the wording. Prove it
  before recommending fan-out: name the items, name what they do not
  share.
- Fan-out that mutates shared state requires isolation (a worktree, or
  provably disjoint paths). Without it the shape is wrong even when
  the items are logically independent.
- A barrier is earned only by cross-item reasoning — a dedup across
  the full set, a merge, an early exit on a total. Needing to flatten,
  map, or filter is not a barrier.
- Batching is part of the shape. Subagents dispatched across turns run
  serially whatever the prose intends.
- Orchestration weight is part of the shape. A fleet for two lookups
  is as wrong as a serial loop over fifty.

## Procedure

1. **Work-item census** over the FULL text of each artifact, never the
   diff hunk — a procedure's shape is invisible in a hunk. Dispatch
   the censuses as one batch per the shared contract, one artifact
   each. Each lists every step that operates over a set (files,
   modules, findings, PRs, review dimensions, candidate designs,
   test cases), its expected size, and whether the artifact bounds
   it.
2. **Independence test** per set: shared mutable state? item N
   consuming item N−1's output? order affecting the result? Record the
   verdict AND the evidence. A set whose independence you cannot
   demonstrate is correctly serial and is not a finding.
3. **Shape classification.** Name the shape prescribed and the shape
   the dependencies want, from one vocabulary — serial · batched
   fan-out · pipeline · barrier · workflow. A mismatch is a candidate.
4. **Safety pass** over every prescribed fan-out: writers need
   isolation or disjoint paths, order-dependent work must not fan out,
   unbounded sets need a stated cap.
5. **Sizing pass**: agent count against independent parts, deliberate
   per-agent model choice, orchestration weight against the work, the
   explicit user opt-in Workflow-class orchestration requires, and
   whether the artifact lets the executor scale fan-out to the
   situation.
6. **Closure pass.** Every fan-out has a gather step that reconciles
   results.

Then verify per the shared review contract.

## Finding block

Safety — the shape corrupts or runs away: **UNSAFE** (parallel
writers, no isolation) · **ORDERED** (fan-out over order-dependent
work) · **UNCAPPED** (unbounded set, no stated cap) · **OPTIN**
(Workflow-class orchestration without the explicit user opt-in it
requires).

Missed concurrency: **SERIAL** (provably independent items walked one
at a time) · **DISPATCH** (subagents spread across turns instead of
one message, so they run serially anyway) · **BARRIER** (collect-all
where each item could flow through its stages independently).

Sizing: **HEAVY** (orchestration heavier than the work) ·
**OVERSIZED** (more agents than independent parts) · **UNTIERED** (no
deliberate model choice where subtask difficulty plainly differs) ·
**COSTBLIND** (cost scales with input, executor given no lever to size
it) · **NOSYNTH** (fan-out with no gather step). UNCAPPED is width the
artifact never bounds; COSTBLIND is a bounded fan-out with no lever to
size down — where both apply, UNCAPPED wins.

```text
<file:line> — <step>: <TYPE>
  items: <the set, and the evidence they are (or are not) independent>
  shape now: <serial | batched fan-out | pipeline | barrier | workflow>
  shape wanted: <same vocabulary>
  fix: <one line>
```

The `items:` line is mandatory: a finding that cannot name its set and
its independence evidence does not clear the bar and is not reported.
