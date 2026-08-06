#!/usr/bin/env bash
# PreToolUse gate: run the skill-review axes before a prose PR opens.
#
# Fails OPEN. Every unexpected condition exits 0 silently, because a broken
# gate must never brick `gh pr create`. The only non-zero exit is a verdict
# of NOT READY from a review that actually ran.
#
# Matchers match tool NAMES, so this runs on every Bash call in the session.
# The bail-outs below are ordered cheapest-first for that reason.
set -uo pipefail

# The review subprocess loads this plugin too.
[ -n "${SKILL_REVIEW_GATE_ACTIVE:-}" ] && exit 0
[ "${SKILL_REVIEW_GATE:-on}" = "off" ] && exit 0

command -v jq >/dev/null 2>&1 || exit 0
input=$(cat)

# Builtin prefilter before paying for jq. Deliberately loose: it only has to
# be free of false negatives, and anything it lets through hits the real
# check below.
case "$input" in *gh*pr*) ;; *) exit 0 ;; esac

[ "$(printf '%s' "$input" | jq -r '.tool_name // empty')" = "Bash" ] || exit 0
cmd=$(printf '%s' "$input" | jq -r '.tool_input.command // empty')
printf '%s' "$cmd" | grep -Eq '(^|[;&|[:space:]])gh[[:space:]]+pr[[:space:]]+create' || exit 0

cwd=$(printf '%s' "$input" | jq -r '.cwd // empty')
[ -n "$cwd" ] && cd "$cwd" 2>/dev/null || exit 0
git rev-parse --git-dir >/dev/null 2>&1 || exit 0

base=""
for candidate in \
  "$(git symbolic-ref --quiet --short refs/remotes/origin/HEAD 2>/dev/null)" \
  origin/main origin/master main master; do
  [ -n "$candidate" ] || continue
  if git rev-parse --verify --quiet "$candidate" >/dev/null 2>&1; then
    base="$candidate"
    break
  fi
done
[ -n "$base" ] || exit 0

changed=$(git diff --name-only "$base...HEAD" 2>/dev/null) || exit 0
prose=$(printf '%s\n' "$changed" | grep -E \
  -e '(^|/)SKILL\.md$' \
  -e '(^|/)(CLAUDE|AGENTS)\.md$' \
  -e '(^|/)(agents|commands|references)/[^/]+\.md$' \
  -e '(^|/)\.claude-plugin/[^/]+\.json$' || true)
[ -n "$prose" ] || exit 0

command -v claude >/dev/null 2>&1 || exit 0
root="${CLAUDE_PLUGIN_ROOT:-}"
[ -n "$root" ] && [ -r "$root/skills/skill-review/SKILL.md" ] || exit 0

read -r -d '' prompt <<EOF
Gate this branch before its pull request opens.

Read these four files first and follow them exactly. They are the contract;
do not substitute your own review criteria.

  $root/skills/skill-review/SKILL.md
  $root/skills/skill-review/references/drift.md
  $root/skills/skill-review/references/concurrency.md
  $root/skills/skill-review/references/token-usage.md

Review the diff of $base...HEAD in $cwd. These changed files are instruction
prose and are what triggered this gate:

$prose

Report survivors in each axis's finding block format, with full
repo-relative file:line anchors. Then end your reply with exactly one line,
nothing after it:

VERDICT: READY
or
VERDICT: NOT READY
EOF

# Judgment work — deliberately not tiered down. See the token-usage axis.
out=$(SKILL_REVIEW_GATE_ACTIVE=1 \
  timeout "${SKILL_REVIEW_GATE_TIMEOUT:-180}" \
  claude -p "$prompt" \
  --permission-mode dontAsk \
  --allowed-tools Read Grep Glob "Bash(git diff:*)" "Bash(git log:*)" "Bash(git show:*)" \
  2>/dev/null) || exit 0

printf '%s' "$out" | grep -q '^VERDICT: NOT READY[[:space:]]*$' || exit 0

{
  echo "skill-review gate: NOT READY."
  echo
  echo "The branch changes instruction prose and the review found unresolved"
  echo "findings. Report the findings below to the user and fix them before"
  echo "opening the PR."
  echo
  printf '%s\n' "$out"
  echo
  echo "To bypass deliberately, re-run with SKILL_REVIEW_GATE=off in the"
  echo "environment."
} >&2
exit 2
