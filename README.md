# skill-review-claude

A [Claude Code plugin marketplace](https://code.claude.com/docs/en/plugin-marketplaces)
for skill and plugin authors. One plugin, `skill-review`, carrying a
three-axis review gate for instruction prose: skills, agents, canon and
reference files, marketplace repos.

Instruction files fail in three ways code does not, and each axis
enforces one rule against one of them:

| axis | rule | what it catches |
| --- | --- | --- |
| **drift** | every rule has exactly one normative statement; every other mention is a pointer to it | a rule restated in a second place, two renderings that have drifted apart, a pointer whose target moved, a rule with no normative home |
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
Anything it notices outside those three lands in a separate observations
section that leaves the verdict alone. It does not judge meaning,
triggering quality, or coverage; compose it with a content reviewer to
cover those properly. Findings are reported, never auto-fixed, and any
finding on any axis blocks the verdict — nothing else does.

## Automatic gate

The plugin ships a `PreToolUse` hook that runs the gate before
`gh pr create` — but only when the branch actually changes instruction
prose. Everything else exits in a few milliseconds, with no model
involvement and no tokens spent.

The gate ships as Go source and compiles itself on first use, caching
the binary under the plugin's data directory, where it survives plugin
updates until the source changes. The one-time build costs a couple of
seconds on the call that triggers it; every call after runs the cached
binary.

It fires when the diff against the base branch touches a `SKILL.md`, a
`CLAUDE.md` or `AGENTS.md`, an `.md` file directly inside `agents/`,
`commands/`, or `references/`, or a `.claude-plugin/*.json` manifest.
When it does fire, the review runs in a separate `claude -p` process with the
shell, the file writers, network access, and MCP tools denied outright, so it
reads your branch rather than editing it. Your session receives the verdict
rather than the whole review: `NOT READY` blocks the command and reports the
findings; `READY` passes the command through, and hands on any observations
the review filed.

The gate fails open. A missing `claude` or Go toolchain, a failed
build, a timeout, an undiscoverable base branch, or any other surprise
lets the PR through — a broken gate must never block work.

| variable | effect |
| --- | --- |
| `SKILL_REVIEW_GATE=off` | skip the gate for that command |
| `SKILL_REVIEW_GATE_TIMEOUT` | seconds before the review gives up (default 180) |
| `SKILL_REVIEW_GATE_DIFF_CAP` | prose diff lines handed to the review (default 3000) |

Keep the timeout under the hook's own budget — the `timeout` field in
`hooks.json`, which ships with the plugin and is replaced on update. Past it the
harness kills the hook first, and because the gate fails open, a raised timeout
buys no review at all. The headroom over the default is modest, so a much larger
cap may not fit however you tune the two together.

A branch whose prose diff runs past the cap is reviewed on the truncated diff
plus the post-change files, so the pre-change half of what was cut goes unread.
The review still returns a verdict either way, `READY` included — raise the cap
for a bulk prose rewrite you want graded in full, and the timeout with it, within
the ceiling above.

Hooks load at session start, so run `/reload-plugins` or restart after
installing or updating.

## Development

Validate after changes:

```sh
claude plugin validate .
```

Before editing prose, read [`AGENTS.md`](AGENTS.md). It is the register of
which files may render the review contract for another audience, what each may
carry, and what obliges them to be rewritten when a rule moves — this file's
Scope section among them.

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
