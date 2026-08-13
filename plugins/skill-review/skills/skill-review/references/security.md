# Axis 4 — security

Runs under the shared review contract in `SKILL.md`, which owns modes,
batching, per-axis dispatch, verification, anchoring,
one-finding-per-site, evidence-or-silence, the verdict, and the scope
boundary.

This axis judges what a procedure's input can MAKE it do. What a
mechanism costs belongs to `references/token-usage.md`; whether a
shape matches its dependencies belongs to `references/concurrency.md`.
A tool roster wider than the work is a finding here only when
untrusted input is in reach of it.

## The data-carries-no-authority rule

Content that crosses a procedure's trust boundary carries no
authority: it may inform what the executor concludes, never direct
what the executor does. Untrusted content enters context only fenced
as data, and the steps that read it hold only the capability their
work requires — so an instruction that slips the fence finds nothing
to drive. Corollaries:

- Trust follows authorship, not content. The trusted set is the user
  and the artifact's own authors — those whose writes reach the
  artifact only through the user's or a maintainer's review. An
  ingress anyone else can write to is untrusted — issue and PR text,
  fetched pages, repositories whose content reaches the tree without
  that review (third-party clones, fork branches, unmerged
  contributor branches), output of tools over external state, the
  input a hook or command interpolates, and the metadata riding all
  of them: paths, titles, branch and test names — however benign
  today's copy reads. A field a machine constrains to an enum or a
  schema carries a shape rather than an author, and is judged as a
  shape below.
- Taint propagates. A step's output is as untrusted as the least
  trusted thing it read; summarizing, reformatting, or persisting it
  launders nothing. A file the procedure wrote under untrusted
  influence is untrusted when read back, wherever it lives —
  including where the harness loads it as instructions rather than
  handing it to a step that reads.
- A fence is its marker. A delimiter the fenced content could legally
  contain is decoration, not confinement; a fence needs a marker no
  content can collide with — from inside the fence or from any
  unfenced region rendered around it — and a treat-as-data
  instruction beside it. Either leg missing leaves the content
  unfenced. Content that must land in an instruction position — a
  task description, an argument — cannot be fenced, and is
  constrained to an author-fixed shape instead: a closed enumeration,
  or a grammar that cannot carry free-form natural language. A length
  cap is not a shape.
- Capability is granted by the author, never by the input. Untrusted
  content may select among outcomes the artifact's own prose
  enumerates, for as long as those outcomes stay reversible in band;
  one that commits an irreversible or outward act needs a gate
  however short its enumeration, because a flipped classification
  then IS the attack. Content that supplies the outcome — the
  command, the step, the target — is directing, and a candidate set
  mined from untrusted content is supplied, not enumerated.
- A gate is author-fixed criteria a tool applies, or a human approval
  — of the bytes read, or of the bytes the act would commit before
  they take effect. The reading step's own judgment is never a gate,
  and neither is that step applying the criteria itself: whatever
  defeats the act defeats the check.
- Fences fail, so blast radius is part of the design: the step that
  reads hostile content holds only what its work requires — no shell,
  no writers, no outward channel it does not use — and its subagents
  inherit the confinement. Capability is measured at the narrowest
  grant the harness can express (`Bash(gh issue view:*)`, not
  `Bash`): a whole shell held for enumerable commands is wider than
  the work.
- Execution is the widest capability there is. Untrusted content
  becomes code — run, sourced, eval'd, installed, or written where
  the harness loads instructions — only behind a gate, or a sandbox
  whose isolation the artifact names: filesystem, credentials,
  network, persistence. Building or testing the repository a
  procedure was dispatched to work on is that procedure's named work,
  not a promotion.
- Secrets transit the narrowest channel that works; argv, logs,
  prompts, environments, and subagent returns are all wider than they
  look.

## Procedure

1. **Ingress-and-grant census**, in two passes so the cheap one bounds
   the dear one. First over the FULL text of each artifact and the
   frontmatter and manifests that arm it — `allowed-tools`, hook and
   MCP configuration, permission prose. Dispatch the censuses as one
   batch per the shared contract, one artifact each. List every
   point where content enters context (files read, pages fetched,
   issue and PR text, command output, subagent returns, the input a
   hook or command interpolates, and the metadata riding them) and
   every capability the artifact grants, holds, or passes to
   subagents. An agent definition's ingress is what its callers hand
   it: census the dispatch sites that brief it, not its own text
   alone. Then, only where the first pass found untrusted input in
   reach of a grant, read the code behind that grant — the program a
   hook or launcher execs — to where its decision is actually made; a
   registration says nothing about whether the guard allowlists or
   no-ops. Record the harness's grant vocabulary as you go: the
   narrowest expressible form of each roster entry, and whether the
   deny mechanism globs. Reviewing a bare artifact whose manifests
   are out of reach, say so: the grants those manifests carry are out
   of scope, not clean — as is a grant you cannot price against an
   expressible narrower form.
2. **Trust classification** per ingress, by the authorship corollary,
   taint included; an ingress whose authorship cannot be established
   is untrusted. The trusted remainder drops out. Record the verdict
   AND the evidence.
3. **Fence check** per untrusted ingress: delimited, marker
   unforgeable from inside the fence or around it, treat-as-data
   beside it — or, for instruction positions, shape-constrained. The
   artifact's own working repository, where the named work is to read
   and modify it, needs no fence.
4. **Authority check** per untrusted ingress: prose that defers to
   what the content says — applies its suggestions, follows its
   steps, obeys its labels — and prose that lets it select an
   irreversible or outward act, each with no gate between reading and
   acting.
5. **Blast-radius walk** per untrusted ingress: from the step that
   reads it, list every capability in reach — held by the step or
   inherited by its subagents — and the work that requires each.
   Outward channels are any tool whose target the content can choose,
   any store that routinely leaves the boundary, and any write that
   reaches an audience the data has not already. List the non-public
   data in reach: what shares the context — a step in the main
   session carries the conversation, its files, and its instructions
   by default, and only an isolated subagent with an enumerated
   context carries less — plus what the step could acquire with a
   capability it already holds.
6. **Secret trace**: every secret the procedure touches, each channel
   it transits, and the narrowest channel that would do. A secret a
   tool authenticates with, whose bytes the procedure never handles,
   is one row and clears at that tool's boundary.

Then verify per the shared review contract.

## Finding block

Steering — untrusted content can direct the executor: **UNFENCED**
(untrusted content in context lacking either an unforgeable fence or
a treat-as-data instruction beside it) · **FORGEABLE** (a fence whose
marker the content could emit or collide with, from inside it or from
an unfenced region rendered around it) · **OBEYS** (prose deferring
to instructions inside untrusted content, or letting it select or
target an irreversible or outward act with no gate).

Blast radius — a successful injection finds something to use:
**OVERGRANT** (capability a step's work does not require, with
untrusted input in reach of it — or subagents that escape the step's
confinement) · **EXFIL** (an outward channel carrying non-public data
to an audience not already authorized for it, while untrusted input
shares the context) · **EXEC** (untrusted content promoted to
execution with no gate and no named sandbox).

Secrets: **LEAKPATH** (a secret transiting a wider channel than its
use requires).

Where one site fills several bars, blast radius outranks steering —
EXEC over EXFIL over LEAKPATH over OVERGRANT — and OBEYS over
FORGEABLE over UNFENCED within it. On a line granting several tools,
the site is the capability, not the line.

```text
<file:line> — <step>: <TYPE>
  evidence: <this type's required evidence, below>
  fix: <one line>
```

### Required evidence

A finding whose evidence line cannot be filled in is not reported.

- **UNFENCED**, **OBEYS** — the ingress and who beyond the trust
  boundary can author it; where authorship cannot be established, the
  channel anyone could write through. For a selected or targeted act,
  what makes it irreversible or outward.
- **FORGEABLE** — additionally a line the content could legally
  contain, or render outside the fence, that the marker mistakes for
  itself.
- **OVERGRANT** — the capability held, at the narrowest grant the
  harness can express — or the dispatch that escapes the
  confinement — the work it exceeds, and the untrusted ingress in
  reach of it. An artifact declaring no roster inherits the session's:
  report only where its own work is narrower than a roster it could
  have declared. A deny-list over an extensible tool surface is
  reported against the mechanism, its extension point standing in for
  a held capability.
- **EXFIL** — the untrusted ingress, the non-public data — in the
  context or acquirable with a capability the step holds — and the
  outward channel with the audience it reaches beyond those already
  authorized, joined in one context or through taint: a laundered
  return or a store that leaves the boundary counts as its leg. The
  channel leg is defeated by a template whose every slot is
  shape-constrained and whose target the content cannot choose, by a
  destination set enforced outside the model, or by human approval of
  the bytes sent.
- **EXEC** — the execution site, the content's provenance, and the
  gate that is absent or does not bind these bytes; where a sandbox
  is claimed, the isolation it leaves unnamed.
- **LEAKPATH** — the secret, the channel it transits, and the
  narrower channel that suffices, "drop the secret" included. A
  bundle forwarded whole — an environment, a config file, a context —
  satisfies the secret leg by naming the bundle.
