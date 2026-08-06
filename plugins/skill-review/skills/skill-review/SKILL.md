---
name: skill-review
description: Use when reviewing or gating changes to Claude Code skills, agents, plugins, or marketplace repos — a pre-PR gate on prompt prose across two axes, duplication/drift (restated rules, drifted pointers, missing normative homes) and concurrency (procedures that serialize independent work, or fan out unsafely or wastefully via subagents, workflows, or pipelines); reviewing a single SKILL.md, agent definition, or canon file on either axis; or auditing a whole plugin repo.
---

# Skill review — structure and shape gate for prompt prose

Instruction files fail in two ways code does not. A rule restated in
two places reads fine in both until one of them moves, and nothing
fails loudly. A procedure that walks independent work one item at a
time reads fine too — it just taxes every executor that follows it,
forever, and the prose never signals the cost. This skill reviews
skill, agent, and canon prose the way a maintainer reviews code, on
two axes: **structure** — how rules are stated — and **shape** — how
work is dispatched.

## Shared review contract

Both axes run under one contract:

- **Modes.** Gate a branch or PR diff before it opens · review a
  single artifact on request · audit a whole repo. Both axes run in
  every mode.
- **Batch independent work.** Where a step's items are independent
  and read-only, dispatch them as ONE batch rather than a serial
  sweep, and gather before the next step. Cap batch width at what
  the harness runs concurrently; queue the rest.
- **Verify before reporting.** Re-read each candidate assuming the
  author was right — transport mappings, per-audience deltas, scoped
  exceptions, and deliberately serial work are all legitimate. Report
  only survivors.
- **Anchor everything.** Every finding carries a full repo-relative
  `file:line`, as does every other site it names.
- **One finding per site**, most specific type wins. A single bad step
  never emits five findings.
- **Report, never fix** — unless the user asks.
- **Verdict: READY / NOT READY** + the must-fix list. Any finding on
  either axis blocks.
- **Scope: structure and shape only** — how rules are stated and how
  work is dispatched. This gate does not judge meaning, triggering
  quality, or coverage; compose it with a content reviewer for those.

## Axis 1 — duplication and drift

### The one-normative-home rule

Every rule has exactly ONE normative statement; every other mention is
a POINTER that names the rule and its home without re-explaining it. A
pointer fails loudly when its target changes; a restatement fails
silently. Corollaries:

- Procedures are executed by reference ("run `<file>`'s steps 1–4"),
  never paraphrased at the call site.
- Contracts — field lists, verdict vocabularies, caps, thresholds —
  live where they are defined; consumers cite them, transport-map
  them, or render them for an audience, and say which they are doing.
- A rendering that owns audience-specific content (a JSON field
  mapping, a per-mode delta, a scoped exception) is not a restatement;
  a rendering that re-explains the rule is, whatever else it adds.

### Drift procedure

1. **Rule census.** From the diff, list every rule, contract field,
   threshold, vocabulary token, and procedure step the change adds,
   edits, or moves.
2. **Rendering hunt.** Dispatch the hunts as one batch per the
   shared contract — one searcher per rule or rule-batch, each
   applying step 3 to what it finds. Each searches the WHOLE repo
   by distinctive phrases, numbers, and token names, not exact
   strings: drifted copies no longer match exactly.
3. **Classify every rendering**: NORMATIVE (the one home) · POINTER
   (names rule and home, no re-explanation) · RESTATED (re-explains —
   a finding) · DRIFTED (contradicts another rendering — a finding
   citing both) · ORPHANED (a pointer whose target moved or no longer
   states the rule — a finding).
4. **Homeless-rule check.** A rule the diff introduces with no
   normative home is a finding; name the home it should get.
5. **Consistency sweep** over the drift-prone classes: caps and
   thresholds, verdict and state vocabularies, field lists across
   JSON/contract/prose renderings, cost estimates, and counts ("all
   seven components").

Then verify per the shared review contract.

### Drift finding block

```text
<file:line> — <rule>: RESTATED | DRIFTED | ORPHANED | HOMELESS
  renderings: <full paths of every other rendering>
  normative home: <where the one statement should live>
  fix: <one line — usually "keep one, point the rest">
```

## Axis 2 — concurrency

### The shape-matches-the-work rule

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

### Concurrency procedure

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

### Concurrency finding block

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

## Repo-wide audit

On request ("audit the repo for drift", "audit the repo's
procedures"), run both censuses over every rule and every procedure in
the repo rather than a diff — same classes, same verification, same
contract. When a rule keeps drifting across rounds, recommend a
behavior-pinning eval: an eval tests the behavior wherever the prose
lives, and is the only rendering that cannot drift silently.
