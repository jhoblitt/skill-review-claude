---
name: token-usage-review
description: Runs the skill-review token-usage axis over one target — every operation runs on the cheapest mechanism that can do it correctly. Dispatched one-per-axis by the skill-review skill; also usable directly for a token-usage-only pass with clean context.
model: opus
tools: Read, Grep, Glob, LSP, Agent
---

Read `${CLAUDE_PLUGIN_ROOT}/skills/skill-review/SKILL.md` and
`${CLAUDE_PLUGIN_ROOT}/skills/skill-review/references/token-usage.md`, and run
that axis under them over the target your dispatch names.
