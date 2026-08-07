package main

import (
	"strings"
	"testing"
)

// The property the whole tool exists for: the same sentence wrapped at a
// different column is the same sentence. A line-oriented search sees two
// different lines and reports nothing, which is how a restatement survives
// review.
func TestNormalizeIsWrappingBlind(t *testing.T) {
	a := "Every rule has exactly ONE normative statement; every other\nmention is a pointer."
	b := "Every rule has exactly ONE\nnormative statement; every\nother mention is a pointer."
	if got, want := strings.Join(Normalize(a), " "), strings.Join(Normalize(b), " "); got != want {
		t.Fatalf("wrapping changed the tokens:\n  %q\n  %q", got, want)
	}
}

func TestNormalize(t *testing.T) {
	cases := []struct {
		name, in, want string
	}{
		{"fences dropped", "before\n```sh\ngh pr create --draft\n```\nafter", "before after"},
		{"emphasis stripped", "the **rule** and its `home`", "the rule and its home"},
		{"case folded", "NEVER Restated", "never restated"},
		{"paths kept whole", "see references/drift.md now", "see references/drift.md now"},
		{"punctuation split", "a rule, stated once; twice — no", "a rule stated once twice no"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := strings.Join(Normalize(c.in), " "); got != c.want {
				t.Fatalf("got %q, want %q", got, c.want)
			}
		})
	}
}

func words(s string) []string { return strings.Fields(s) }

func TestScanFindsLongestRun(t *testing.T) {
	shared := "a rule has exactly one normative statement and every other mention points at it"
	docs := map[string][]string{
		"a.md": words("preamble " + shared + " tail"),
		"b.md": words("different opening " + shared + " different close"),
	}
	got := Scan(docs, 4, 8)
	if len(got) != 1 {
		t.Fatalf("want 1 pair, got %d: %+v", len(got), got)
	}
	if got[0].Run != len(words(shared)) {
		t.Fatalf("run = %d, want %d", got[0].Run, len(words(shared)))
	}
	if got[0].Sample != shared {
		t.Fatalf("sample = %q, want %q", got[0].Sample, shared)
	}
	if got[0].A != "a.md" || got[0].B != "b.md" {
		t.Fatalf("pair = %s/%s, want a.md/b.md", got[0].A, got[0].B)
	}
}

// A single run counter per document credits one partner's streak to another.
// Here a.md shares one n-gram with b.md and a DIFFERENT one with c.md, back to
// back: neither pair has a run of two, and a shared counter would report one.
func TestScanTracksRunsPerPartner(t *testing.T) {
	docs := map[string][]string{
		"a.md": words("one two three four five six seven eight nine"),
		"b.md": words("zero one two three four zzz"),
		"c.md": words("yyy two three four five xxx"),
	}
	for _, p := range Scan(docs, 4, 1) {
		if p.Run > 5 {
			t.Fatalf("run %d for %s/%s exceeds what either pair shares: %q",
				p.Run, p.A, p.B, p.Sample)
		}
	}
}

// The sample is sliced from whichever document the winning run was measured
// in. Keying pairs by sorted name while indexing into the other document is
// the easy way to print text the pair does not contain.
func TestScanSampleComesFromASharedRun(t *testing.T) {
	shared := "the opening notice is the whole attribution and nothing else follows it"
	docs := map[string][]string{
		"zzz.md": words("padding padding padding padding " + shared),
		"aaa.md": words(shared + " trailing"),
	}
	got := Scan(docs, 4, 8)
	if len(got) != 1 {
		t.Fatalf("want 1 pair, got %d", len(got))
	}
	for _, doc := range []string{"zzz.md", "aaa.md"} {
		if !strings.Contains(strings.Join(docs[doc], " "), got[0].Sample) {
			t.Fatalf("sample %q is not present in %s", got[0].Sample, doc)
		}
	}
}

func TestScanHonorsMin(t *testing.T) {
	docs := map[string][]string{
		"a.md": words("alpha beta gamma delta epsilon zeta"),
		"b.md": words("alpha beta gamma delta epsilon eta"),
	}
	if got := Scan(docs, 3, 6); len(got) != 0 {
		t.Fatalf("min=6 should drop a 5-token run, got %+v", got)
	}
	if got := Scan(docs, 3, 5); len(got) != 1 {
		t.Fatalf("min=5 should keep a 5-token run, got %+v", got)
	}
}

func TestScanIgnoresSelfAndShortDocs(t *testing.T) {
	docs := map[string][]string{
		"a.md": words("repeated phrase repeated phrase repeated phrase"),
		"b.md": words("tiny"),
	}
	if got := Scan(docs, 4, 4); len(got) != 0 {
		t.Fatalf("a document must not pair with itself, got %+v", got)
	}
}
