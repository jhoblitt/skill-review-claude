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

## Axes

Load only the axes the request calls for. Paths are relative to this
skill's directory.

- **structure** — is every rule stated exactly once? →
  `references/drift.md`
- **shape** — does the control flow match the dependencies? →
  `references/concurrency.md`

A pre-PR gate or a repo-wide audit runs both. A targeted request loads
one: "is this rule stated anywhere else?" is structure, "does this
parallelize?" is shape.

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

## Repo-wide audit

On request ("audit the repo for drift", "audit the repo's
procedures"), run both censuses over every rule and every procedure in
the repo rather than a diff — same classes, same verification, same
contract. When a rule keeps drifting across rounds, recommend a
behavior-pinning eval: an eval tests the behavior wherever the prose
lives, and is the only rendering that cannot drift silently.
