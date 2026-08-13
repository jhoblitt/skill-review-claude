---
name: drift-review
description: Runs the skill-review drift axis over one target — every rule has exactly one normative statement, and every other mention is a pointer to it. Dispatched one-per-axis by the skill-review skill; also usable directly for a drift-only pass with clean context.
model: opus
tools: Read, Grep, Glob, LSP, Agent, Bash(go run:*)
---

Read `${CLAUDE_PLUGIN_ROOT}/skills/skill-review/SKILL.md` and
`${CLAUDE_PLUGIN_ROOT}/skills/skill-review/references/drift.md`, and run that
axis under them over the target your dispatch names.
