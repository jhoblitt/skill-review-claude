#!/usr/bin/env bash
# Tests for plugins/skill-review/hooks/pr-gate.sh.
#
# The gate must stay silent on everything that is not a prose PR, and block
# only on a NOT READY verdict. Every case drives the real script through its
# real stdin contract, against a stub `claude` so no test spends tokens.
#
# The fire-set cases assert through the script rather than re-matching its
# regex here. A copy of that pattern in this file would be a second rendering
# of the fire-set, free to drift from the one that ships — which is the bug
# this suite exists to prevent.

set -uo pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
GATE="$ROOT/plugins/skill-review/hooks/pr-gate.sh"
PLUGIN_ROOT="$ROOT/plugins/skill-review"

WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT

passed=0
failed=0

ok() {
  passed=$((passed + 1))
  printf '  ok      %s\n' "$1"
}

no() {
  failed=$((failed + 1))
  printf '  FAIL    %s\n          %s\n' "$1" "$2"
}

# A stub `claude` that reports the verdict it was built with.
stub_claude() {
  local dir="$WORK/bin-${1// /-}"
  mkdir -p "$dir"
  {
    echo '#!/usr/bin/env bash'
    printf 'echo "VERDICT: %s"\n' "$1"
  } >"$dir/claude"
  chmod +x "$dir/claude"
  printf '%s' "$dir"
}

# A throwaway repo whose feature branch touches exactly the given paths.
fixture() {
  local dir="$WORK/repo-$1"
  shift
  mkdir -p "$dir"
  {
    git -C "$dir" init -q -b main
    git -C "$dir" config user.email test@example.invalid
    git -C "$dir" config user.name test
    echo seed >"$dir/seed.txt"
    git -C "$dir" add -A
    git -C "$dir" commit -qm seed
    git -C "$dir" switch -qc feature
    local path
    for path in "$@"; do
      mkdir -p "$dir/$(dirname "$path")"
      echo changed >"$dir/$path"
    done
    git -C "$dir" add -A
    git -C "$dir" commit -qm change
  } >/dev/null 2>&1
  printf '%s' "$dir"
}

payload() {
  printf '{"tool_name":"%s","tool_input":{"command":"%s"},"cwd":"%s"}' "$1" "$2" "$3"
}

# PATH is pinned so `claude` is discoverable only when the stub provides it.
run_gate() {
  local json=$1 bin=$2
  shift 2
  printf '%s' "$json" |
    env -u SKILL_REVIEW_GATE -u SKILL_REVIEW_GATE_ACTIVE \
      PATH="$bin:/usr/bin:/bin" CLAUDE_PLUGIN_ROOT="$PLUGIN_ROOT" "$@" \
      bash "$GATE" >/dev/null 2>&1
}

expect() {
  local desc=$1 want=$2
  shift 2
  local got
  run_gate "$@"
  got=$?
  if [ "$got" -eq "$want" ]; then
    ok "$desc"
  else
    no "$desc" "expected exit $want, got $got"
  fi
}

READY=$(stub_claude READY)
BLOCK=$(stub_claude "NOT READY")
NO_CLAUDE="$WORK/bin-none"
mkdir -p "$NO_CLAUDE"

prose=$(fixture prose "plugins/p/skills/s/SKILL.md")
code=$(fixture code main.go Makefile)
plain="$WORK/not-a-repo"
mkdir -p "$plain"

echo "bail-outs — the gate stays silent and costs nothing"
expect "tool is not Bash" 0 \
  "$(payload Read 'gh pr create' "$prose")" "$BLOCK"
expect "Bash, but not a gh pr create" 0 \
  "$(payload Bash 'gh pr list' "$prose")" "$BLOCK"
expect "cwd is not a git repository" 0 \
  "$(payload Bash 'gh pr create' "$plain")" "$BLOCK"
expect "diff touches no instruction prose" 0 \
  "$(payload Bash 'gh pr create' "$code")" "$BLOCK"
expect "claude is not installed (fails open)" 0 \
  "$(payload Bash 'gh pr create' "$prose")" "$NO_CLAUDE"
expect "SKILL_REVIEW_GATE=off" 0 \
  "$(payload Bash 'gh pr create' "$prose")" "$BLOCK" SKILL_REVIEW_GATE=off
expect "SKILL_REVIEW_GATE_ACTIVE=1 (re-entry guard)" 0 \
  "$(payload Bash 'gh pr create' "$prose")" "$BLOCK" SKILL_REVIEW_GATE_ACTIVE=1

echo
echo "verdicts"
expect "prose PR, review says READY" 0 \
  "$(payload Bash 'gh pr create --draft' "$prose")" "$READY"
expect "prose PR, review says NOT READY" 2 \
  "$(payload Bash 'gh pr create --draft' "$prose")" "$BLOCK"

stderr=$(printf '%s' "$(payload Bash 'gh pr create' "$prose")" |
  env -u SKILL_REVIEW_GATE -u SKILL_REVIEW_GATE_ACTIVE \
    PATH="$BLOCK:/usr/bin:/bin" CLAUDE_PLUGIN_ROOT="$PLUGIN_ROOT" \
    bash "$GATE" 2>&1 >/dev/null)
case "$stderr" in
*"NOT READY"*"SKILL_REVIEW_GATE=off"*)
  ok "block message carries the verdict and the bypass" ;;
*)
  no "block message carries the verdict and the bypass" "stderr was: ${stderr:-<empty>}" ;;
esac

echo
echo "fire-set — which changed paths count as instruction prose"
fires() {
  local path=$1 repo
  repo=$(fixture "fire-${path//\//-}" "$path")
  expect "fires:  $path" 2 "$(payload Bash 'gh pr create' "$repo")" "$BLOCK"
}

silent() {
  local path=$1 repo
  repo=$(fixture "quiet-${path//\//-}" "$path")
  expect "silent: $path" 0 "$(payload Bash 'gh pr create' "$repo")" "$BLOCK"
}

fires SKILL.md
fires plugins/p/skills/s/SKILL.md
fires CLAUDE.md
fires AGENTS.md
fires references/drift.md
fires plugins/p/skills/s/references/drift.md
fires agents/reviewer.md
fires commands/review.md
fires .claude-plugin/plugin.json

silent references/axes/drift.md
silent agents/team/reviewer.md
silent references/schema.json
silent commands/review.sh
silent README.md
silent docs/specs/design.md

printf '\n%d passed, %d failed\n' "$passed" "$failed"
[ "$failed" -eq 0 ]
