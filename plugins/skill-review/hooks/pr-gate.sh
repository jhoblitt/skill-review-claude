#!/usr/bin/env bash
# Build-on-first-use launcher for the Go gate in hooks/pr-gate/.
#
# Fails OPEN, same doctrine as the gate itself: a broken build must never
# brick `gh pr create`, or anything else — this launcher runs on every Bash
# call in the session. Correctness under parallel hooks comes from the
# per-PID temp plus atomic mv; flock, where present, only dedupes the work.
#
# The kill switch is checked here as well as in Go: the escape hatch for a
# broken gate cannot live only inside the artifact whose build may be what
# broke.
set -uo pipefail

[ "${SKILL_REVIEW_GATE:-on}" = "off" ] && exit 0

src="${CLAUDE_PLUGIN_ROOT:-}/hooks/pr-gate"
[ -n "${CLAUDE_PLUGIN_ROOT:-}" ] && [ -d "$src" ] || exit 0

# The binary lives in the per-plugin data directory, never in
# CLAUDE_PLUGIN_ROOT: the cache is a per-version directory that moves on
# every update and can be reaped mid-session. Data outlives versions, and a
# fresh version clone carries fresh mtimes, so the staleness check below
# rebuilds after updates on its own.
data="${CLAUDE_PLUGIN_DATA:-${XDG_CACHE_HOME:-$HOME/.cache}/skill-review-claude}"
bin="$data/pr-gate"

stale() {
  [ ! -x "$bin" ] && return 0
  [ -n "$(find "$src" -maxdepth 1 \( -name '*.go' -o -name go.mod \) \
    -newer "$bin" -print -quit 2>/dev/null)" ]
}

# A failed build must not retry on every Bash call. The marker pins the
# source fingerprint it failed on: a source change re-arms immediately (the
# dev iterating on a fix), expiry re-arms after transient toolchain or
# network trouble. Success removes it.
build() {
  local marker="$bin.buildfail" fp
  fp=$(cat "$src"/go.mod "$src"/*.go 2>/dev/null | cksum)
  if [ -f "$marker" ] && [ "$(cat "$marker" 2>/dev/null)" = "$fp" ] &&
    [ -z "$(find "$marker" -mmin +60 -print 2>/dev/null)" ]; then
    return 1
  fi
  if (cd "$src" && go build -o "$bin.tmp.$$" .) >/dev/null 2>&1 &&
    mv -f "$bin.tmp.$$" "$bin" 2>/dev/null; then
    rm -f "$marker"
    return 0
  fi
  rm -f "$bin.tmp.$$"
  printf '%s' "$fp" >"$marker" 2>/dev/null
  return 1
}

if stale; then
  command -v go >/dev/null 2>&1 || exit 0
  mkdir -p "$data" 2>/dev/null || exit 0
  if command -v flock >/dev/null 2>&1; then
    exec 9>"$bin.lock" || exit 0
    flock 9 || exit 0
    if stale; then
      build || exit 0
    fi
    exec 9>&-
  else
    build || exit 0
  fi
fi

exec "$bin"
