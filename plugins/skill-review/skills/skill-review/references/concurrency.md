# Axis 2 — concurrency

Runs under the shared review contract in `SKILL.md`, which owns modes,
batching, verification, anchoring, one-finding-per-site,
evidence-or-silence, the verdict, and the scope boundary.

This axis judges whether a shape is CORRECT for the dependencies. What
that shape costs — orchestration weight, agent count, model tier,
sizing levers — belongs to `references/token-usage.md`.

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
   unbounded sets need a stated cap, and Workflow-class orchestration
   needs the explicit user opt-in it requires.
5. **Closure pass.** Every fan-out has a gather step that reconciles
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
where each item could flow through its stages independently) ·
**NOSYNTH** (fan-out with no gather step).

```text
<file:line> — <step>: <TYPE>
  items: <the set, and the evidence they are (or are not) independent>
  shape now: <serial | batched fan-out | pipeline | barrier | workflow>
  shape wanted: <same vocabulary>
  fix: <one line>
```

The `items:` line is this axis's evidence bar: a finding that cannot
name its set and its independence evidence is not reported.
