# Axis 3 — token usage

Runs under the shared review contract in `SKILL.md`, which owns modes,
batching, per-axis dispatch, verification, anchoring,
one-finding-per-site, evidence-or-silence, the verdict, and the scope
boundary.

This axis judges what a procedure COSTS to run. Whether its shape is
correct for the dependencies belongs to `references/concurrency.md`.

## The cheapest-sufficient-mechanism rule

Every operation runs on the cheapest mechanism that can do it
CORRECTLY: code before model, small model before large, targeted read
before whole file, load-on-demand before load-always. Cost is
justified only by capability actually required. Corollaries:

- Cheapness upstream is amplified downstream. An under-tiered
  supervisor is the most expensive error on this axis: a bad
  decomposition wastes every worker it dispatches, so the saving is
  multiplied into a loss. Tier the supervisor to the hardest decision
  it must make, not to the average one.
- Cost paid per invocation dominates cost paid per run. Prose resident
  in a skill's entry file is billed on every trigger, before any work
  happens; a wasteful step is billed only when reached.
- A round trip is a unit of cost. Every tool call re-bills the context
  it has accumulated, so a run of calls that one call could carry pays
  for the split whether its steps are independent or strictly ordered.
  Granularity is a cost question even where the serial shape is right.
- Deterministic output belongs in code. If a machine can check the
  result, a model should not be generating it — and bulk data should
  reach the model already reduced, or not at all.
- Cheap gates precede expensive work.
- Judgment is not a cost centre. Review, verification, refutation, and
  adjudication stay at full capability; capability cuts there cost
  more than they save.

## Procedure

1. **Operation census** over the full text of each artifact: every
   step that spends tokens — model calls, subagent dispatches, tool
   calls, data transformations — plus the artifact's own resident
   prose, which every trigger pays for. Dispatch the censuses as one
   batch per the shared contract, one artifact each.
2. **Mechanism test** per operation: name the cheapest mechanism that
   does it correctly — code · tool · small model · large model — and
   the specific capability that forces anything more. A requirement
   you cannot name is a requirement the operation does not have.
3. **Tier check** on every model choice: supervisors and decomposers
   against the hardest decision they must make, workers against how
   deterministically their output can be checked, and judgment work
   against the carve-out that exempts it.
4. **Code check**: work the artifact makes an agent implement that it
   could ship, bulk data routed through context that a script could
   reduce outside it, and deterministic work handed to a model.
5. **Footprint check**: content resident in the entry file that only
   some invocations need, and could load on demand instead.
6. **Granularity check**: runs of tool calls the artifact prescribes —
   a command per file, a script per item, a read per path — where one
   call carries the run and no step between them branches on what the
   last returned.
7. **Flow check**: what subagents are told to return, what they are
   told to read, and whether cheap gates precede the expensive work
   they could short-circuit.

Then verify per the shared review contract.

## Finding block

Model class: **WEAKSUPER** (supervisor or decomposer tiered below the
decisions it must make) · **OVERTIER** (mechanical work on a frontier
model) · **TIEREDJUDGE** (judgment, verification, or adjudication
tiered down) · **UNTIERED** (no deliberate model choice where subtask
difficulty plainly differs).

Code versus model: **REGEN** (the agent is made to write code the
skill could ship) · **INCONTEXT** (bulk data routed through model
context that a script could reduce outside it) · **MODELDET** (a model
doing deterministic work a tool does better).

Footprint: **MONOLITH** (content billed on every trigger that only
some invocations need).

Context economy: **VERBOSERET** (subagent return value unconstrained —
prose where a schema would do, quoted blocks where an anchor would) ·
**WIDEREAD** (whole-file reads where targeted reads suffice, or one
context redundantly shipped to N agents) · **CHATTY** (a run of tool
calls where one call carries the whole run).

Ordering: **LATEGATE** (a cheap check placed after the expensive work
it could have short-circuited).

Orchestration cost: **HEAVY** (orchestration heavier than the work) ·
**OVERSIZED** (more agents than independent parts) · **COSTBLIND**
(cost scales with input and the executor is given no lever to size
it).

```text
<file:line> — <operation>: <TYPE>
  mechanism now: <code | tool | small model | large model | prose>[ ×N]
  mechanism wanted: <same vocabulary>[ ×N]
  evidence: <this type's required evidence, below>
  fix: <one line>
```

The `×N` count is carried only where a finding turns on how many calls
an operation takes rather than which mechanism runs it.

### Required evidence

A finding whose evidence line cannot be filled in is not reported.

- **WEAKSUPER** — the specific decomposition decision the tier cannot
  make, and how many workers its output puts at risk.
- **OVERTIER** — the deterministic check that backstops the cheaper
  tier. No backstop, no finding.
- **TIEREDJUDGE** — the judgment being delegated and what a wrong call
  costs downstream.
- **REGEN**, **MODELDET** — what makes the output machine-checkable.
- **INCONTEXT** — the volume routed through context, and where it
  could be reduced instead.
- **MONOLITH** — which triggers pay for the content without needing
  it.
- **VERBOSERET**, **WIDEREAD** — what is shipped that the consumer
  never reads.
- **CHATTY** — the calls that combine into one, and what makes their
  intermediate results unnecessary to inspect.
- **LATEGATE** — the expensive work the cheap check short-circuits,
  and how often it would.
- **UNTIERED** — the subtasks whose difficulty differs, and by what.
- **HEAVY**, **OVERSIZED** — the work the orchestration exceeds, and
  the independent-part count it exceeds it by.
- **COSTBLIND** — the input the cost scales with, and the lever the
  executor is missing.
