# Repository conventions

## Where the rules live

`plugins/skill-review/skills/skill-review/SKILL.md` owns the shared review
contract — its bullets are the list, and this register does not keep a second
copy of them. Each of the four axis rules lives in its own reference beside
its corollaries: `references/drift.md`, `references/concurrency.md`,
`references/token-usage.md`, `references/security.md`.

`plugins/skill-review/agents/*.md` carry no rules in their bodies. Each names
the two files above that its axis runs under and nothing else, so prose there
that explains an axis is a finding rather than a rendering. Their frontmatter
is the one home for that agent's own tier and tool roster, and its
`description` renders the axis under the table below.

Rationale stays with the rule it explains, wherever that rule lives.

## Sanctioned audience renderings

Six places re-explain some of that for a different audience instead of
pointing at it. Under `references/drift.md`'s corollary that is a restatement,
so this is a deliberate exception to the one-normative-home rule rather than an
application of it: a marketplace listing cannot cite a contract at a reader who
has not installed it, and a dispatcher deciding whether to load the skill has
not read it either. The exception is bounded by the obligations below,
and each rendering carries only what this table says:

| rendering | may carry |
| --- | --- |
| `README.md`, intro and axes table | the axis names, each rule in a line, what each catches, why these failures survive ordinary review, the modes, and the anchoring rule |
| `README.md`, "How a review runs" diagram | the modes, the axis names, the one-agent-per-axis dispatch, the verify-before-reporting step, the anchoring rule, and the verdict rule, as the shape of one pass |
| `README.md`, "Scope" section | the scope boundary, the verdict rule, and report-never-fix, stated as outcomes a user sees |
| `plugins/skill-review/.claude-plugin/plugin.json`, `description` | the axis names, the four rule names, and the modes |
| `.claude-plugin/marketplace.json`, both `description` fields | the axis names and the modes |
| `plugins/skill-review/skills/skill-review/SKILL.md`, frontmatter `description` | the axis names, what each catches, and the modes |
| `plugins/skill-review/agents/*.md`, frontmatter `description` | the axis name, its rule in a line, and how the agent is dispatched |

A rendering left behind when a rule it carries moves is DRIFTED, and it must be
rewritten in the same change that moves the rule. A rendering that reaches past
its row is RESTATED. Extend the table in the same change that extends a
rendering, or the next reviewer has to guess where the licence ends.

## Design records

`docs/specs/*.md` record how a decision was reached, on the date it was reached.
They are non-normative: where one disagrees with the live rules, the live rules
win and the record is history rather than a finding. They may state what was
decided and what was rejected. They may not restate what a decision produced —
a restatement there drifts, and nothing is allowed to catch it.

That exemption is why a design record may hold a rule this register does not
sanction.

The register governs renderings of the review contract, and nothing else
outside the table may carry one. Rules that are not the contract's — how this
repository releases, say — are normal prose with a home of their own, and the
table has no claim on them.
