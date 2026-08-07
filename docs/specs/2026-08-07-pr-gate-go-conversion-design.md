# Design — pr-gate hook as a Go binary

**Status:** accepted, 2026-08-07. **Non-normative** — see `AGENTS.md`,
"Design records", for what that means and what this file may hold. A
dated record of why the PR gate became a compiled binary and how it
gets built without an install step.

## Problem

`hooks/pr-gate.sh` is 209 lines of bash that fire on every Bash call
in every session with the plugin enabled. Most of its weight is
defense against the language it is written in: a builtin prefilter
exists to dodge a `jq` fork, the prompt travels by here-string because
argv has a ceiling, NUL-delimited path plumbing guards quoting, and
`awk`/`grep` render every parse as another process. A compiled binary
does the entire gate in one exec with a real JSON parser and no
external dependencies — `jq` stops being a requirement.

What a compiled artifact lacks is a way to exist: the plugin system
has no install-time lifecycle. The manifest schema carries no
install/setup/build key (confirmed against the live schema and every
manifest in the local plugin cache), and the `Setup` hook event fires
only under explicit `--init`-family flags. Whatever builds the binary
must run from the hook path itself.

## Decisions

**The entry point does not move.** `hooks/pr-gate.sh` keeps its name
and becomes a thin wrapper; `hooks.json` is untouched. The gate logic
ports to Go under `hooks/pr-gate/` — `go 1.26`, stdlib only, so the
build needs no module downloads and no vendor tree.

**Build on first use, cached in `${CLAUDE_PLUGIN_DATA}`.** The wrapper
builds the binary when it is missing or older than any of `go.mod`,
`*.go`, and execs it otherwise. The cache is the documented per-plugin
data directory, falling back to `${XDG_CACHE_HOME:-~/.cache}/skill-review-claude`
when unset. `${CLAUDE_PLUGIN_ROOT}` is rejected as a cache: it is a
per-version directory that moves on update and is reapable mid-session.
Staleness by mtime makes updates rebuild for free — a fresh version
clone carries fresh mtimes. Builds serialize under `flock` because
PreToolUse hooks run in parallel. No SessionStart pre-build: the one
~2s compile lands on a call that already spends minutes in review.

**The wrapper is thin, with two deliberate exceptions.** All gate
logic — guards included — lives in Go. Except: the `SKILL_REVIEW_GATE=off`
kill switch is also checked in the wrapper, because the escape hatch
for a broken gate must not depend on the build that may be what broke;
and a failed build writes a marker holding a source fingerprint, and
further attempts are skipped while the fingerprint matches and the
marker is under an hour old. Without the marker a broken build would
retry on every Bash call in the session; the fingerprint half re-arms
immediately when the source changes, the expiry half re-arms after
transient toolchain or network failures. Success removes the marker.

**Fail-open extends to the new failure class.** The doctrine at the
top of the bash gate carries over verbatim: a broken gate must never
brick `gh pr create`. No Go toolchain, a failed compile, an unwritable
cache — each exits 0 silently, exactly as a missing `jq` does today.

**The port is behavior-identical, not a redesign.** Same environment
contract (`SKILL_REVIEW_GATE`, `_ACTIVE`, `_TIMEOUT`, `_DIFF_CAP`),
same git invocations flag-for-flag, same fire-set regex, same
byte-identical prompt including the `=` fence, same verdict and
Observations handling, same exit codes. The `timeout` command becomes
a SIGTERM-first `cmd.Cancel` with a kill delay; process-wide `cd`
becomes per-subprocess `Dir`, with the review still run from the repo
top. The load-bearing comments move with the code they defend.

**Tests stay end-to-end, plus unit tests for the pure core.**
`tests/pr-gate-test.sh` keeps driving `bash "$GATE"` through the real
stdin contract against stub `claude` binaries — it pre-builds the
binary once into `CLAUDE_PLUGIN_DATA="$WORK/data"` so the PATH-pinned
cases never need `go`, and gains wrapper cases: cold build, rebuild on
source change, cooldown honored and re-armed, missing toolchain fails
open. `gate_test.go` covers the extracted pure functions (fire-set,
command regex, cap fallback, Observations extraction).

**Lint machinery lands with the code.** `validate.yml` gains a pinned
`setup-go` reading the module's `go.mod` via `go-version-file`, then
`gofmt -l` (fail on output), `go vet`, `golangci-lint` (pinned, default
linter set, config file only when a finding needs silencing), and
`go test` before the e2e suite. `shellcheck` covers the wrapper and
the test suite. Dependabot gains a `gomod` entry for the module.

## Non-goals

Prebuilt or downloaded binaries. `go run` (leaves the gate needing the
toolchain forever instead of once). Windows beyond today's support —
the entry point is already bash. Any change to what fires the gate,
what the review is handed, or how verdicts read.

## Release

Lands as `feat:` (minor). `plugin.json` `version` is never touched in
a PR; semantic-release owns that field.
