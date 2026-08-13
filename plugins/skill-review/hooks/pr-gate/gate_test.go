package main

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// The same fire-set the e2e suite asserts through the shipped artifact. Kept
// in both places on purpose: this table localizes a regex regression to a
// function, the suite proves the wiring.
func TestIsProse(t *testing.T) {
	fires := []string{
		"SKILL.md",
		"plugins/p/skills/s/SKILL.md",
		"CLAUDE.md",
		"AGENTS.md",
		"references/drift.md",
		"plugins/p/skills/s/references/drift.md",
		"agents/reviewer.md",
		"commands/review.md",
		".claude-plugin/plugin.json",
		"hooks/hooks.json",
		"plugins/p/hooks/hooks.json",
		".mcp.json",
		"plugins/p/.mcp.json",
	}
	silent := []string{
		"references/axes/drift.md",
		"agents/team/reviewer.md",
		"references/schema.json",
		"hooks/team/extra.json",
		"mcp.json",
		"commands/review.sh",
		"README.md",
		"docs/specs/design.md",
	}
	for _, p := range fires {
		if !isProse(p) {
			t.Errorf("isProse(%q) = false, want true", p)
		}
	}
	for _, p := range silent {
		if isProse(p) {
			t.Errorf("isProse(%q) = true, want false", p)
		}
	}
}

func TestIsGhPrCreate(t *testing.T) {
	yes := []string{
		"gh pr create",
		"gh pr create --draft",
		"cd x && gh  pr   create",
		"foo;gh pr create",
		"foo | gh pr create",
		"git push && gh\tpr\tcreate --fill",
	}
	no := []string{
		"gh pr list",
		"gh prcreate",
		"mygh pr create",
		"echo 'gh_pr_create'",
		"gh pr view --web",
	}
	for _, c := range yes {
		if !isGhPrCreate(c) {
			t.Errorf("isGhPrCreate(%q) = false, want true", c)
		}
	}
	for _, c := range no {
		if isGhPrCreate(c) {
			t.Errorf("isGhPrCreate(%q) = true, want false", c)
		}
	}
}

// The unusable values the bash gate enumerated beside its fallback: digits
// only, positive, or the default wins.
func TestDiffCapFrom(t *testing.T) {
	cases := map[string]int{
		"":                     3000,
		"0":                    3000,
		"-5":                   3000,
		"+5":                   3000,
		"nope":                 3000,
		"5":                    5,
		"007":                  7,
		"99999999999999999999": 3000,
	}
	for in, want := range cases {
		if got := diffCapFrom(in); got != want {
			t.Errorf("diffCapFrom(%q) = %d, want %d", in, got, want)
		}
	}
}

func TestTruncateDiff(t *testing.T) {
	diff := "l1\nl2\nl3\nl4"
	kept, cut := truncateDiff(diff, 2)
	if !cut || kept != "l1\nl2" {
		t.Errorf("truncateDiff cap 2 = (%q, %v), want (\"l1\\nl2\", true)", kept, cut)
	}
	kept, cut = truncateDiff(diff, 4)
	if cut || kept != diff {
		t.Errorf("truncateDiff at exactly the cap = (%q, %v), want unchanged, false", kept, cut)
	}
	kept, cut = truncateDiff(diff, 3000)
	if cut || kept != diff {
		t.Errorf("truncateDiff under the cap = (%q, %v), want unchanged, false", kept, cut)
	}
}

func TestTruncationNotice(t *testing.T) {
	n := truncationNotice(5, "/repo/top")
	if !strings.HasPrefix(n, "\n") {
		t.Error("notice must start with a newline: it follows the closing fence line")
	}
	if !strings.Contains(n, "stops at 5 lines") {
		t.Errorf("notice missing the cap: %q", n)
	}
	if !strings.Contains(n, "readable under /repo/top.") {
		t.Errorf("notice missing the root: %q", n)
	}
	if strings.HasSuffix(n, "\n") {
		t.Error("notice must not end with a newline: the template supplies it")
	}
}

func TestObservations(t *testing.T) {
	out := strings.Join([]string{
		"## Axis 1 — duplication and drift",
		"",
		"No survivors.",
		"",
		"## Observations",
		"",
		"- a note that does not block",
		"- another",
		"",
		"VERDICT: READY",
	}, "\n")
	got := observations(out)
	// Trailing newlines are stripped, matching what the bash gate's $(awk ...)
	// substitution handed to jq; interior blank lines survive.
	want := "\n- a note that does not block\n- another"
	if got != want {
		t.Errorf("observations = %q, want %q", got, want)
	}

	if got := observations("## Axis 1\n\nnothing\n\nVERDICT: READY"); got != "" {
		t.Errorf("no Observations heading must yield nothing, got %q", got)
	}

	// A later heading closes the section; its lines must not leak.
	out = "## Observations\n- kept\n## Axis 9\n- dropped\nVERDICT: READY"
	if got := observations(out); got != "- kept" {
		t.Errorf("heading after Observations must close it, got %q", got)
	}
}

func TestHasNotReady(t *testing.T) {
	if !hasNotReady("prose\nVERDICT: NOT READY") {
		t.Error("plain NOT READY line not detected")
	}
	if !hasNotReady("VERDICT: NOT READY  \nafter") {
		t.Error("trailing spaces on the verdict line not tolerated")
	}
	if hasNotReady("VERDICT: READY") {
		t.Error("READY misread as NOT READY")
	}
	if hasNotReady("a line saying VERDICT: NOT READY mid-sentence") {
		t.Error("mid-line text misread as a verdict")
	}
	if hasNotReady("XVERDICT: NOT READY") {
		t.Error("prefixed text misread as a verdict")
	}
}

func TestBuildPrompt(t *testing.T) {
	p := buildPrompt(promptParams{
		Root:       "/plugroot",
		Base:       "origin/main",
		Top:        "/repo",
		ProseList:  "  a/SKILL.md\n",
		ProseDiff:  "diff --git a/a/SKILL.md b/a/SKILL.md",
		Truncation: "",
	})
	for _, must := range []string{
		"Gate this branch before its pull request opens.",
		"/plugroot/skills/skill-review/SKILL.md",
		"the diff of origin/main...HEAD below",
		"post-change files under /repo directly",
		"\n\n  a/SKILL.md\n\nTheir diff follows",
		"every path in\nthe list above is indented",
		"===== BEGIN DIFF =====\ndiff --git a/a/SKILL.md b/a/SKILL.md\n===== END DIFF =====\n\n",
		"VERDICT: READY\nor\nVERDICT: NOT READY",
	} {
		if !strings.Contains(p, must) {
			t.Errorf("prompt missing %q", must)
		}
	}
	// The dispatcher is handed the contract alone; each axis agent reads its own
	// reference, which is the peak-context saving the per-axis split exists for.
	// Asserted on the template, not on p: a branch that changes a reference puts
	// that path in the prose list and the diff, where it belongs.
	for _, absent := range []string{
		"references/drift.md", "references/concurrency.md",
		"references/token-usage.md", "references/security.md",
	} {
		if strings.Contains(promptTemplate, absent) {
			t.Errorf("prompt template still hands the dispatcher %q", absent)
		}
	}
	if strings.HasSuffix(p, "\n") {
		t.Error("prompt must not end with a newline; the caller's stdin write adds it")
	}

	// With a truncation notice the fence line is followed by the notice, then
	// the blank line the template already carries.
	p = buildPrompt(promptParams{
		Root: "/r", Base: "b", Top: "/t", ProseList: "  x\n", ProseDiff: "d",
		Truncation: truncationNotice(5, "/t"),
	})
	if !strings.Contains(p, "===== END DIFF =====\n\nThat diff stops at 5 lines") {
		t.Error("truncation notice not seated directly after the closing fence")
	}
}

// The roster is asserted by property, not by a second copy of the list: a
// duplicate here would drift from the one that ships, which is the failure
// this repository exists to catch. What matters is that the review keeps what
// it reads and fans out with, loses each class of tool that could escape it,
// and that the two sets never overlap — a name in both is denied, so an
// allow-list entry that quietly stops working would otherwise go unnoticed.
func TestReviewRoster(t *testing.T) {
	allowed := map[string]bool{}
	for _, a := range allowedTools {
		allowed[a] = true
	}
	denied := map[string]bool{}
	for _, d := range deniedTools {
		denied[d] = true
	}

	for _, a := range allowedTools {
		if denied[a] {
			t.Errorf("%q is both granted and denied; the denial wins", a)
		}
	}

	// Read, Grep and Glob are the review's only route to the tree with the
	// shell denied; Task is the dispatch every axis runs in.
	for _, need := range []string{"Read", "Grep", "Glob", "Task"} {
		if !allowed[need] {
			t.Errorf("review cannot work without %q granted", need)
		}
	}

	// One representative per class of escape, named in the finding vocabulary
	// the axes use: a shell, a writer, the network, persistence past the
	// process, reach outside it, and loading more capability.
	for _, hazard := range []string{
		"Bash", "Write", "WebFetch", "CronCreate", "EnterWorktree",
		"SendMessage", "Workflow", "Skill", "ToolSearch",
	} {
		if !denied[hazard] {
			t.Errorf("%q reaches the review; it must be denied", hazard)
		}
	}
}

// The flags carry the roster only if each list lands under its own flag. This
// parses argv back the way the e2e suite does, so a misplaced flag fails here
// rather than silently granting the deny list.
func TestReviewArgs(t *testing.T) {
	args := reviewArgs()

	section := func(flag string) []string {
		var out []string
		for i, a := range args {
			if a != flag {
				continue
			}
			for _, v := range args[i+1:] {
				if strings.HasPrefix(v, "--") {
					break
				}
				out = append(out, v)
			}
			return out
		}
		return nil
	}

	got := strings.Join(section("--allowed-tools"), " ")
	if got != strings.Join(allowedTools, " ") {
		t.Errorf("--allowed-tools carried %q", got)
	}
	got = strings.Join(section("--disallowed-tools"), " ")
	if got != strings.Join(deniedTools, " ") {
		t.Errorf("--disallowed-tools carried %q", got)
	}
	if len(args) == 0 || args[0] != "-p" {
		t.Error("review must run headless: -p has to lead the argv")
	}
	if strings.Join(section("--permission-mode"), " ") != "dontAsk" {
		t.Error("review must not be able to prompt; the hook has no one to ask")
	}
	// The gate fails open, so a tier the install cannot serve would disable it
	// silently rather than degrade it.
	if len(section("--model")) != 0 {
		t.Error("review must not pin a model; an unavailable tier would fail the gate open")
	}
}

func frontmatter(body, key string) string {
	m := regexp.MustCompile(`(?m)^` + key + `: *(.+)$`).FindStringSubmatch(body)
	if m == nil {
		return ""
	}
	return strings.TrimSpace(m[1])
}

// The dispatch identifiers in SKILL.md, the agent files they name, and the tier
// those files pin have to agree, and nothing else reads them together:
// `claude plugin validate` passes with agents/ deleted outright, so a renamed
// or missing agent would silently drop an axis from every review.
func TestAxisAgentsMatchContract(t *testing.T) {
	const root = "../../"

	contract, err := os.ReadFile(root + "skills/skill-review/SKILL.md")
	if err != nil {
		t.Fatalf("reading the contract: %v", err)
	}

	dispatched := map[string]bool{}
	for _, m := range regexp.MustCompile(`skill-review:([a-z-]+)`).
		FindAllStringSubmatch(string(contract), -1) {
		dispatched[m[1]] = true
	}
	if len(dispatched) != 4 {
		t.Fatalf("contract dispatches %d axis agents, want 4: %v", len(dispatched), dispatched)
	}

	tiers := map[string]string{}
	for name := range dispatched {
		body, err := os.ReadFile(root + "agents/" + name + ".md")
		if err != nil {
			t.Errorf("contract dispatches %q with no agent file behind it: %v", name, err)
			continue
		}
		if got := frontmatter(string(body), "name"); got != name {
			t.Errorf("agents/%s.md declares name %q; that dispatch cannot resolve", name, got)
		}
		tiers[name] = frontmatter(string(body), "model")
	}

	for name, tier := range tiers {
		if tier == "" {
			t.Errorf("agents/%s.md pins no model", name)
		} else if tier != tiers["drift-review"] {
			t.Errorf("agents/%s.md pins %q while drift-review pins %q; the axes must not split tiers",
				name, tier, tiers["drift-review"])
		}
	}
}
