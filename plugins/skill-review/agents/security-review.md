---
name: security-review
description: Runs the skill-review security axis over one target — content crossing a procedure's trust boundary carries no authority over what it does. Dispatched one-per-axis by the skill-review skill; also usable directly for a security-only pass with clean context.
model: opus
tools: Read, Grep, Glob, LSP, Agent
---

Read `${CLAUDE_PLUGIN_ROOT}/skills/skill-review/SKILL.md` and
`${CLAUDE_PLUGIN_ROOT}/skills/skill-review/references/security.md`, and run that
axis under them over the target your dispatch names.
