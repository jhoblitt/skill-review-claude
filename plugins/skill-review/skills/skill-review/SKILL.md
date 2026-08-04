---
name: skill-review
description: Use when reviewing or gating changes to Claude Code skills, agents, plugins, or marketplace repos — a pre-PR duplication/drift gate on prompt prose; reviewing a SKILL.md, agent definition, or canon/reference file for restated rules, drifted pointers, or missing normative homes; or auditing a whole plugin repo for instruction drift.
---

# Skill review — duplication and drift gate for prompt prose

Instruction files drift in a way code does not: a rule restated in two
places reads fine in both until one of them moves, and nothing fails
loudly. This skill reviews skill, agent, and canon prose the way a
maintainer reviews code — findings verified before reporting, anchored
to full repo-relative `file:line` paths, and REPORTED, never fixed
unless the user asks.

## The one-normative-home rule

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

## Pre-PR gate procedure

Run on a plugin/skill repo branch (or PR diff) before it opens:

1. **Rule census.** From the diff, list every rule, contract field,
   threshold, vocabulary token, and procedure step the change adds,
   edits, or moves.
2. **Rendering hunt.** For each, search the WHOLE repo for its other
   renderings — by distinctive phrases, numbers, and token names, not
   exact strings: drifted copies no longer match exactly.
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
6. **Verify before reporting.** Re-read each candidate assuming the
   author was right — transport mappings, per-audience deltas, and
   scoped exceptions are legitimate. Report only survivors, each with
   full repo-relative `file:line` for every rendering involved.

## Findings and verdict

```text
<file:line> — <rule>: RESTATED | DRIFTED | ORPHANED | HOMELESS
  renderings: <full paths of every other rendering>
  normative home: <where the one statement should live>
  fix: <one line — usually "keep one, point the rest">
```

Verdict: **READY** / **NOT READY** + the must-fix list. Structure
findings only — this gate does not judge meaning, triggering quality,
or coverage; compose it with a content reviewer for those.

## Repo-wide audit

On request ("audit the repo for drift"), run the census over every
rule in the repo rather than a diff — same classes, same verification,
same contract. When a rule keeps drifting across rounds, recommend a
behavior-pinning eval: an eval tests the behavior wherever the prose
lives, and is the only rendering that cannot drift silently.
