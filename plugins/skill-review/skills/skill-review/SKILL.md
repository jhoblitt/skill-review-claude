---
name: skill-review
description: Use when reviewing or gating changes to Claude Code skills, agents, plugins, or marketplace repos — a pre-PR gate on prompt prose across three axes, structure (rules restated, pointers drifted, no normative home), shape (procedures that serialize independent work, or fan out unsafely via subagents, workflows, or pipelines), and cost (wrong model tier for a supervisor or worker, work a shipped script should do, prose billed on every trigger); reviewing a single SKILL.md, agent definition, or canon file on any axis; or auditing a whole plugin repo.
---

# Skill review — structure, shape, and cost gate for prompt prose

Instruction files fail in three ways code does not. A rule restated in
two places reads fine in both until one of them moves. A procedure
that walks independent work one item at a time reads fine too — it
just taxes every executor that follows it. And a procedure that spends
a frontier model on what a script could do reads fine forever, because
prose carries no price tag. This skill reviews skill, agent, and canon
prose the way a maintainer reviews code, on three axes.

## Axes

Load only the axes the request calls for. Paths are relative to this
skill's directory.

- **structure** — is every rule stated exactly once? →
  `references/drift.md`
- **shape** — does the control flow match the dependencies? →
  `references/concurrency.md`
- **cost** — is every operation on the cheapest sufficient mechanism?
  → `references/token-usage.md`

A pre-PR gate or a repo-wide audit runs all three. A targeted request
loads one: "is this rule stated anywhere else?" is structure, "does
this parallelize?" is shape, "is this wasting tokens?" is cost.

## Shared review contract

Every axis runs under one contract:

- **Modes.** Gate a branch or PR diff before it opens · review a
  single artifact on request · audit a whole repo. Every loaded axis
  runs in every mode.
- **Batch independent work.** Where a step's items are independent
  and read-only, dispatch them as ONE batch rather than a serial
  sweep, and gather before the next step. Cap batch width at what
  the harness runs concurrently; queue the rest.
- **Verify before reporting.** Re-read each candidate assuming the
  author was right — transport mappings, per-audience deltas, scoped
  exceptions, deliberately serial work, and deliberately expensive
  work are all legitimate. Report only survivors.
- **Anchor everything.** Every finding carries a full repo-relative
  `file:line`, as does every other site it names.
- **One finding per site**, most specific type wins. A single bad step
  never emits five findings.
- **Evidence or silence.** Each axis names the evidence its finding
  types require. A finding whose evidence cannot be filled in is not
  reported, however plausible it reads.
- **Report, never fix** — unless the user asks.
- **Verdict: READY / NOT READY** + the must-fix list. Any finding on
  any axis blocks.
- **Scope: structure, shape, and cost only** — how rules are stated,
  how work is dispatched, and what it costs to run. This gate does not
  judge meaning, triggering quality, or coverage; compose it with a
  content reviewer for those.

## Repo-wide audit

On request ("audit the repo for drift", "audit the repo's
procedures"), run each loaded axis's census over the whole repo rather
than a diff — same classes, same verification, same contract. When a
rule keeps drifting across rounds, recommend a behavior-pinning eval:
an eval tests the behavior wherever the prose lives, and is the only
rendering that cannot drift silently.
