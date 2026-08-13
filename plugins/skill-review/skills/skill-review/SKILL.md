---
name: skill-review
description: Use when reviewing or gating changes to Claude Code skills, agents, plugins, or marketplace repos — a pre-PR gate on prompt prose across four axes, drift (rules restated, pointers stale, no normative home), concurrency (procedures that serialize independent work, or fan out unsafely via subagents, workflows, or pipelines), token usage (wrong model tier for a supervisor or worker, work a shipped script should do, prose billed on every trigger), and security (prompt-injection surface — untrusted content entering context unfenced or obeyed as instructions, capability or outward channels in reach of hostile input, secrets on wide channels); reviewing a single SKILL.md, agent definition, or canon file on any axis; or auditing a whole plugin repo.
---

# Skill review — drift, concurrency, token-usage, and security gate

Instruction files fail in four ways code does not. A rule restated in
two places reads fine in both until one of them moves. A procedure
that walks independent work one item at a time reads fine too — it
just taxes every executor that follows it. A procedure that spends a
frontier model on what a script could do reads fine forever, because
prose carries no price tag. And a procedure that reads content an
attacker can author, with capability that work never needed in reach,
reads fine right up until the attacker writes it. This skill reviews
skill, agent, and canon prose the way a maintainer reviews code, on
four axes.

## Axes

Load only the axes the request calls for. Paths are relative to this
skill's directory.

- **drift** — is every rule stated exactly once? →
  `references/drift.md`
- **concurrency** — does the control flow match the dependencies? →
  `references/concurrency.md`
- **token usage** — is every operation on the cheapest sufficient
  mechanism? → `references/token-usage.md`
- **security** — can anything the procedure reads steer it, or reach
  beyond it? → `references/security.md`

A pre-PR gate or a repo-wide audit runs all four. A targeted request
loads one: "is this rule stated anywhere else?" is drift, "does this
parallelize?" is concurrency, "is this wasting tokens?" is token
usage, "could a hostile issue body steer this?" is security.

## Shared review contract

Every axis runs under one contract:

- **Modes.** Gate a branch or PR diff before it opens · review a
  single artifact on request · audit a whole repo. Every loaded axis
  runs in every mode.
- **Batch independent work.** Where a step's items are independent
  and read-only, dispatch them as ONE batch rather than a serial
  sweep, and pipeline each result into the next step. Gather only
  where the next step needs the whole set. Cap batch width at what
  the harness runs concurrently; queue the rest. Axes that
  census artifacts share one: a single subagent per artifact returns
  each loaded axis's inventory for it.
- **Verify before reporting.** Re-read each candidate assuming the
  author was right — transport mappings, per-audience deltas, scoped
  exceptions, deliberately serial work, deliberately expensive work,
  and capability a named work requires are all legitimate. Report
  only survivors.
- **Anchor everything.** Every finding carries a full repo-relative
  `file:line`, as does every other site it names.
- **One finding per site**, most specific type wins. A single bad step
  never emits five findings.
- **Evidence or silence.** Each axis names the evidence its finding
  types require. A finding whose evidence cannot be filled in is not
  reported, however plausible it reads.
- **Report, never fix** — unless the user asks.
- **Verdict: READY / NOT READY** + the must-fix list. Any finding on
  any axis blocks, and nothing else does. Anything you notice outside
  the four axes belongs in a separate observations section that leaves
  the verdict alone — a gate whose blocking set the reviewer can widen
  has no stopping condition.
- **Scope: drift, concurrency, token usage, and security only** — how
  rules are stated, how work is dispatched, what it costs to run, and
  what its input can make it do. This gate does not judge meaning,
  triggering quality, or coverage; compose it with a content reviewer
  to cover those properly.

## Repo-wide audit

On request ("audit the repo for drift", "audit the repo's
procedures"), run each loaded axis's census over the whole repo rather
than a diff — same classes, same verification, same contract. When a
rule keeps drifting across rounds, recommend a behavior-pinning eval:
an eval tests the behavior wherever the prose lives, and is the only
rendering that cannot drift silently.
