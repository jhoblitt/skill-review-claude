---
name: concurrency-review
description: Runs the skill-review concurrency axis over one target — a procedure's control flow matches the dependency structure of the work it describes. Dispatched one-per-axis by the skill-review skill; also usable directly for a concurrency-only pass with clean context.
model: opus
tools: Read, Grep, Glob, LSP, Agent
---

Read `${CLAUDE_PLUGIN_ROOT}/skills/skill-review/SKILL.md` and
`${CLAUDE_PLUGIN_ROOT}/skills/skill-review/references/concurrency.md`, and run
that axis under them over the target your dispatch names.
