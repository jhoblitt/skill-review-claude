package main

import (
	"regexp"
	"sort"
	"strings"
)

// Fenced blocks are dropped whole. A duplicated code sample is not a
// duplicated rule, and this axis judges prose; leaving fences in makes every
// document that shows the same command look like a restatement.
var fenceRe = regexp.MustCompile("(?s)```.*?```")

// Markdown punctuation carries no meaning for this comparison, and keeping it
// would split `**rule**` from `rule`. Path characters stay: `references/x.md`
// is one token, and a pointer that names its target is exactly what should
// survive normalization.
var (
	markupRe = regexp.MustCompile("[`*_>#|\\[\\]()]")
	splitRe  = regexp.MustCompile(`[^\w/.\-]+`)
)

// Normalize reduces a document to comparable tokens. Case, wrapping and
// markdown emphasis all vary between two copies of the same sentence, and none
// of them changes what the sentence says. Folding them is what lets this find
// a restatement a line-oriented search cannot see, because the copy wraps at a
// different column.
func Normalize(text string) []string {
	text = fenceRe.ReplaceAllString(text, " ")
	text = markupRe.ReplaceAllString(text, " ")
	text = strings.ToLower(text)
	fields := splitRe.Split(text, -1)
	kept := fields[:0]
	for _, w := range fields {
		if w != "" {
			kept = append(kept, w)
		}
	}
	return kept
}

// Pair is one file pair sharing a contiguous run of tokens.
type Pair struct {
	A, B   string
	Run    int    // length in tokens of the longest shared contiguous run
	Sample string // that run, rendered from the document it was measured in
}

type best struct {
	run   int
	src   string // document the run was measured in; start indexes THIS one
	start int
}

// Scan reports, for every pair of documents, the longest contiguous run of
// tokens they share, keeping pairs at or above min.
//
// The run is tracked per PARTNER rather than per document. A document whose
// n-gram at i is shared with B and at i+1 with C has no run of two with
// either, and a single counter would credit one to both — inflating exactly
// the long-verbatim-copy signal this exists to rank by.
func Scan(docs map[string][]string, n, min int) []Pair {
	if n < 1 {
		n = 1
	}

	index := map[string]map[string]bool{}
	for name, words := range docs {
		for i := 0; i+n <= len(words); i++ {
			g := strings.Join(words[i:i+n], " ")
			if index[g] == nil {
				index[g] = map[string]bool{}
			}
			index[g][name] = true
		}
	}

	// Keyed by the ordered pair so each pair is reported once, whichever
	// document is walked first.
	bests := map[[2]string]best{}

	for name, words := range docs {
		run := map[string]int{}   // partner -> current contiguous run
		seen := map[string]bool{} // partners sharing the CURRENT n-gram
		for i := 0; i+n <= len(words); i++ {
			g := strings.Join(words[i:i+n], " ")
			for p := range seen {
				delete(seen, p)
			}
			for p := range index[g] {
				if p != name {
					seen[p] = true
				}
			}
			for p := range run {
				if !seen[p] {
					delete(run, p)
				}
			}
			for p := range seen {
				run[p]++
				length := run[p] + n - 1
				key := [2]string{name, p}
				if key[0] > key[1] {
					key[0], key[1] = key[1], key[0]
				}
				if b, ok := bests[key]; !ok || length > b.run {
					bests[key] = best{run: length, src: name, start: i - run[p] + 1}
				}
			}
		}
	}

	var out []Pair
	for key, b := range bests {
		if b.run < min {
			continue
		}
		src := docs[b.src]
		out = append(out, Pair{
			A:      key[0],
			B:      key[1],
			Run:    b.run,
			Sample: strings.Join(src[b.start:b.start+b.run], " "),
		})
	}
	// Longest first, ties broken by name, so two runs over an unchanged tree
	// produce output that diffs cleanly.
	sort.Slice(out, func(i, j int) bool {
		if out[i].Run != out[j].Run {
			return out[i].Run > out[j].Run
		}
		if out[i].A != out[j].A {
			return out[i].A < out[j].A
		}
		return out[i].B < out[j].B
	})
	return out
}
