# skill-review-claude

A [Claude Code plugin marketplace](https://code.claude.com/docs/en/plugin-marketplaces)
for skill and plugin authors. One plugin, `skill-review`, carrying a
three-axis review gate for instruction prose: skills, agents, canon and
reference files, marketplace repos.

Instruction files fail in three ways code does not, and each axis
enforces one rule against one of them:

| axis | rule | what it catches |
| --- | --- | --- |
| **drift** | every rule has exactly one normative statement; every other mention is a pointer to it | a rule restated in a second place, two renderings that have drifted apart, a pointer whose target moved, a new rule with no normative home |
| **concurrency** | a procedure's control flow matches the dependency structure of the work it describes | independent work walked one item at a time, subagents spread across turns so they run serially anyway, a collect-all barrier that earns nothing — and unsafe fan-out: parallel writers without isolation, order-dependent work, unbounded sets with no cap |
| **token usage** | every operation runs on the cheapest mechanism that can do it correctly | a supervisor tiered below the decisions it must make, mechanical work on a frontier model, code the skill could ship but makes an agent rewrite, prose billed on every trigger that only some invocations need |

Every one of these reads fine on the page — a restatement fails only
once a copy moves, and prose carries no price tag at all — which is why
they survive the reviews code gets. All three axes run across a diff or
a whole repo and report maintainer-grade findings with full
repo-relative anchors.

## Install

Inside Claude Code:

```text
/plugin marketplace add jhoblitt/skill-review-claude
/plugin install skill-review@skill-review-claude
```

or from a shell:

```sh
claude plugin marketplace add jhoblitt/skill-review-claude
claude plugin install skill-review@skill-review-claude
```

To pick up updates later:

```sh
claude plugin marketplace update skill-review-claude   # refresh the index
claude plugin update skill-review@skill-review-claude  # install the new version
```

or inside Claude Code:

```text
/plugin marketplace update skill-review-claude
/plugin update skill-review@skill-review-claude
/reload-plugins
```

After the shell form, run `/reload-plugins` in running sessions — a
restart also works. The marketplace step alone only refreshes the
index; it does not update the installed plugin.

## Example prompts

- "gate this branch before I open the PR" (in a plugin repo)
- "review this SKILL.md for duplicated rules"
- "audit the repo for instruction drift"
- "is this new rule stated anywhere else?"
- "does this skill parallelize the work it could?"
- "is this fan-out safe, or are those agents writing the same files?"
- "is this skill wasting tokens?"
- "is a strong enough model supervising that fan-out?"

## Scope

The gate reports drift, concurrency, and token-usage findings —
restated, drifted, orphaned, and homeless rules; serialized or unsafe
fan-out; and operations running on a costlier mechanism than they need.
It does not judge meaning, triggering quality, or coverage; compose it
with a content reviewer for those. Findings are reported, never
auto-fixed, and any finding on any axis blocks the verdict.

## Automatic gate

The plugin ships a `PreToolUse` hook that runs the gate before
`gh pr create` — but only when the branch actually changes instruction
prose. Everything else exits in about 4 ms of shell, with no model
involvement and no tokens spent.

It fires when the diff against the base branch touches a `SKILL.md`, a
`CLAUDE.md` or `AGENTS.md`, anything under `agents/`, `commands/`, or
`references/`, or a `.claude-plugin/*.json` manifest. When it does fire,
the review runs in a separate `claude -p` process restricted to
read-only tools, so your session receives the verdict rather than the
whole review. `NOT READY` blocks the command and reports the findings;
`READY` is silent.

The gate fails open. A missing `claude` or `jq`, a timeout, an
undiscoverable base branch, or any other surprise lets the PR through —
a broken gate must never block work.

| variable | effect |
| --- | --- |
| `SKILL_REVIEW_GATE=off` | skip the gate for that command |
| `SKILL_REVIEW_GATE_TIMEOUT` | seconds before the review gives up (default 180) |

Hooks load at session start, so run `/reload-plugins` or restart after
installing or updating.

## Development

Validate after changes:

```sh
claude plugin validate .
```

Content changes land via PR. Commit messages follow
[Conventional Commits](https://www.conventionalcommits.org/) — commitlint
enforces this on every PR — and releasing is automated: on each merge to
`main`, [semantic-release](https://github.com/semantic-release/semantic-release)
derives the next version from the commit types in the changeset (`fix:`,
`docs:`, `refactor:`, `perf:` → patch; `feat:` → minor; a breaking change →
major; other types cut no release), writes it into the plugin manifest and
`CHANGELOG.md`, tags, and publishes a GitHub release. Never bump the
plugin version in a PR — the release commit owns that field.

## License

[Apache-2.0](LICENSE)
