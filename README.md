# skill-review-claude

A [Claude Code plugin marketplace](https://code.claude.com/docs/en/plugin-marketplaces)
for skill and plugin authors. One plugin, `skill-review`, carrying a
three-axis review gate for instruction prose: skills, agents, canon and
reference files, marketplace repos.

Instruction files fail in three ways code does not. A rule restated in
two places reads fine in both until one moves — the **structure** axis
enforces the one-normative-home rule and hunts restated, drifted, and
orphaned renderings. A procedure that walks independent work one item at
a time reads fine too — the **shape** axis checks that control flow
matches dependencies, flagging serialized independent work as readily as
unsafe fan-out via subagents, workflows, and pipelines. And a procedure
that spends a frontier model on what a script could do reads fine
forever — the **cost** axis holds every operation to the cheapest
mechanism that can do it correctly, from supervisor model tier down to
the prose billed on every trigger. All three run across a diff or a
whole repo and report maintainer-grade findings with full repo-relative
anchors.

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

The gate reports structure, shape, and cost findings — restated,
drifted, orphaned, and homeless rules; serialized or unsafe concurrency;
and operations running on a costlier mechanism than they need. It does
not judge meaning, triggering quality, or coverage; compose it with a
content reviewer for those. Findings are reported, never auto-fixed, and
any finding on any axis blocks the verdict.

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
