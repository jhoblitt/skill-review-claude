# Design — security review axis

**Status:** accepted, 2026-08-13. **Non-normative** — see `AGENTS.md`,
"Design records", for what that means and what this file may hold. A
dated record of why the security axis is shaped the way it is.

## Problem

Three axes ask how a rule is stated, how work is dispatched, and what
it costs. None asks what the content a procedure ingests can make it
do. A skill that pipes an issue body into a prompt unfenced reads
fine; one that hands the step reading it a shell and an outward
channel reads fine too — right up until someone authors the content.
This repository's own gate hardened itself against exactly this (an
unforgeable `=` fence, a treat-as-data instruction, a deny-listed
roster inherited by subagents) with no axis to say why.

## Decisions

**One rule covering steering AND blast radius.** The request named
prompt injection in particular but not exclusively. Rejected: an
injection-only axis — fences fail, so a review that checks fences
without checking what a slipped instruction could reach certifies
half a defense. The rule is data-carries-no-authority.

**Trust is provenance, not proximate authorship.** Adversarial review
of the first draft found one root under four holes: classifying an
ingress by its immediate author lets hostility launder through a
subagent's summary, a file the skill itself wrote, or a re-fetch
after a gate approved different bytes. Taint propagation — output as
untrusted as the least trusted input, gates bound to the bytes they
approved — closes all four with one corollary.

**Manifests are in scope, and the gate's fire set grows to match.**
Capability grants live in `allowed-tools` frontmatter, hook and MCP
configuration, and permission prose; the census reads them or the
grant findings are unfindable. `.claude-plugin/*.json` already fired
the gate; `hooks/*.json` and `.mcp.json` did not — a branch changing
only hook or MCP configuration, the canonical grant-widening change,
opened its PR unreviewed — so this change adds both, and funds the
fourth axis by raising the hook budget and review timeout.

**Steering findings turn on the act, not only the instruction.**
Adversarial review found the first draft blessing the attack it was
built for: a fenced, tool-free reader classifying an issue into an
author-enumerated set, with a dispatcher mapping one value to closing
the issue, satisfied every rule while an injection flipped the
classification. Selection is informing only while the outcomes stay
reversible in band; a gate is required past that, and the gate itself
had to be defined as criteria a TOOL applies — a check the reading
step performs falls to whatever defeated the act.

**Legitimacy carve-outs must name the finding type they rescue.** A
"deliberately fail-open designs" exemption added to the shared
contract in the first round was cut in the second: no finding type on
any axis condemns fail-open, so the clause could only cull true
positives — and an attacker triggers the failure it blesses.

**Boundary with token usage.** A wide roster is this axis's finding
only when untrusted input is in reach of it; model tiering stays
under one owner and capability under another, and no finding type
migrates.

**One census across axes.** Four axes each dispatching a full-text
census of the same artifact is the repo's own WIDEREAD and DISPATCH
verdict against itself. The shared contract now owns a single census
batch that returns every loaded axis's inventory.

**Every finding blocks; precision by evidence bars, per precedent.**
EXFIL carries the strictest bar — ingress, non-public data, and an
outward channel reaching an audience the data has not, each leg
named, taint-joined legs counting — because a leg merely implied is
where the false positives live, and a channel returning data to its
own audience discloses nothing. Its defeats are written narrowly for
the same reason: a template excuses a post only where every slot is
shape-constrained, since a free-text slot carries whatever an
injection wrote.

**The gate reads its own output as data.** The axis condemns relaying
attacker-influenceable text in instruction position, which is what
the hook did with the review's findings and observations. Both now
travel fenced, on the same unforgeable marker the review's input
uses, and the prompt now states the invariant that makes that marker
unforgeable from the changed-path list as well as from the diff.

**The census is two passes, because one was unaffordable.** Running
the axis against real installed skills put step 1 alone at roughly
twelve minutes and 113k tokens on a plugin with hooks and tools —
past the whole gate's budget before the other five steps began. A
cheap pass over the artifact and its manifests now bounds an
expensive one, which reads the code behind a grant only where the
first pass found untrusted input in reach of it. Measured afterwards
on `rook-code-review`, that split kept 93% of the target's Go
unopened.

**The hook is a bounded pass, and the timeout is set from the
harness's ceiling rather than from a ratio.** A full run of this axis
against that skill took 24 minutes and 356k tokens — no timeout a
blocking `PreToolUse` hook can hold, since the user waits behind it.
The hook budget moves to the harness's own 600-second default for a
command hook, the review to 540 of it, and the depth beyond that is
reached by invoking the skill directly. The prior 240/180 pair was
neither measured nor at the ceiling.

**Named work is not a promotion, and an inherited roster is not an
over-grant.** The same field test fired EXEC on every implementation
agent told to run its own repository's tests, and OVERGRANT on every
skill that declares no `allowed-tools` — including this one, which
would have made the gate block itself. Both are carve-outs now: an
axis whose false positives are universal reports nothing a maintainer
will read.

## Non-goals

Auditing shipped executable code for vulnerabilities — this axis
reads instruction prose and the grants that arm it; compose it with a
code security reviewer for the code itself. Judging meaning,
triggering quality, or coverage. Auto-fixing findings. Designing
sandboxes or permission systems.

## Release

Lands as `feat:` (minor). `plugin.json` `version` is never touched in
a PR; semantic-release owns that field.
