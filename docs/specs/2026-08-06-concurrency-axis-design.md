# Design — concurrency review axis

**Status:** accepted, 2026-08-06. **Non-normative** — see `AGENTS.md`,
"Design records", for what that means and what this file may hold. A
dated record of why the concurrency axis is shaped the way it is.

## Problem

`skill-review` ships one axis: a duplication and drift gate over
instruction prose, built on the one-normative-home rule. It answers
"is this rule stated once?" and nothing else.

It has no answer for the other way instruction files fail silently. A
procedure that walks independent work one item at a time reads
perfectly well — it just taxes every executor that follows it,
forever, and nothing in the prose signals the cost. The inverse fails
worse: a procedure that fans out over work whose items share state or
depend on each other's output corrupts results, and reads exactly as
fine.

## Decisions

**Placement — two axis sections in one `SKILL.md`.** Rejected: a
sibling skill (duplicates the shared review contract across two
files — a finding this repo's own gate would raise), and a router
with `references/` per axis (three files to carry two axes). Revisit
the router only if a third axis arrives.

**Direction — bidirectional, including cost and scale.** The axis
judges missed concurrency, unsafe concurrency, and concurrency sized
wrong. A one-directional "parallelize more" reviewer produces harmful
advice: told to fan out file edits without isolation, an executor
corrupts trees.

**Teeth — every finding blocks.** Any finding on either axis flips the
verdict to NOT READY. Rejected: a correctness/advisory severity split.

The accepted risk is an everything-blocks gate crying wolf on arguable
design opinions. It is contained by precision rather than by severity — the evidence
bars and one-finding-per-site, which the live rules own.

## Scope boundary and its renderings

The boundary widened from "structure findings only" to structure and
shape: how rules are stated, and how work is dispatched. It still does
not judge meaning, triggering quality, or coverage.

This design also settled that a few files may render the boundary for
another audience rather than point at it, and that such a rendering is
rewritten rather than converted to a pointer. Which files those are, and
what each may carry, is `AGENTS.md`'s register — not this document.

## Non-goals

Judging meaning, triggering quality, or coverage. Auto-fixing
findings. Reviewing executable code for concurrency bugs — this axis
reads instruction prose and the procedures it prescribes.

## Release

Lands as `feat:` (minor). `plugin.json` `version` is never touched in
a PR; semantic-release owns that field.
