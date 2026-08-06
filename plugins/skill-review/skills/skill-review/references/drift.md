# Axis 1 — duplication and drift

Runs under the shared review contract in `SKILL.md`, which owns modes,
batching, verification, anchoring, one-finding-per-site, the verdict,
and the scope boundary.

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

## Procedure

1. **Rule census.** From the diff, list every rule, contract field,
   threshold, vocabulary token, and procedure step the change adds,
   edits, or moves.
2. **Rendering hunt.** Dispatch the hunts as one batch per the
   shared contract — one searcher per rule or rule-batch, each
   applying step 3 to what it finds. Each searches the WHOLE repo
   by distinctive phrases, numbers, and token names, not exact
   strings: drifted copies no longer match exactly.
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

Then verify per the shared review contract.

## Finding block

```text
<file:line> — <rule>: RESTATED | DRIFTED | ORPHANED | HOMELESS
  renderings: <full paths of every other rendering>
  normative home: <where the one statement should live>
  fix: <one line — usually "keep one, point the rest">
```
