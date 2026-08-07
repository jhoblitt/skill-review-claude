// dupscan — stage one of the drift axis's rendering hunt, done mechanically.
//
// The hunt has to find copies that no longer match exactly, which is why
// drift.md asks for distinctive phrases rather than exact strings. Done by
// hand that means grep, and grep is line-oriented: a rule restated with
// different line breaks is invisible to it, however verbatim the words. This
// normalizes wrapping away and reports the longest contiguous run of tokens
// each pair of documents shares, longest first.
//
// It classifies nothing. Every pair it prints is a CANDIDATE for step 3 of
// that procedure — a pointer that names its target shares tokens with the
// target too, and is not a finding.
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type excludes []string

func (e *excludes) String() string { return strings.Join(*e, ",") }

func (e *excludes) Set(v string) error {
	*e = append(*e, v)
	return nil
}

func run() int {
	n := flag.Int("n", 8, "n-gram size in tokens; the granularity of a shared run")
	min := flag.Int("min", 12, "report pairs sharing at least this many contiguous tokens")
	var skip excludes
	flag.Var(&skip, "exclude", "skip paths containing this substring (repeatable)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"usage: dupscan [-n 8] [-min 12] [-exclude SUBSTR] <dir-or-file>...\n\n"+
				"Reports markdown file pairs sharing a contiguous run of prose.\n"+
				"Candidates for the drift axis to classify, never findings.\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		return 2
	}

	docs := map[string][]string{}
	for _, root := range flag.Args() {
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil
			}
			for _, s := range skip {
				if strings.Contains(path, s) {
					return nil
				}
			}
			b, err := os.ReadFile(path) //nolint:gosec // paths come from the caller's own argv
			if err != nil {
				return err
			}
			docs[filepath.ToSlash(path)] = Normalize(string(b))
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "dupscan: %v\n", err)
			return 2
		}
	}

	if len(docs) < 2 {
		fmt.Fprintf(os.Stderr, "dupscan: need at least two markdown files, found %d\n", len(docs))
		return 2
	}

	pairs := Scan(docs, *n, *min)
	if len(pairs) == 0 {
		fmt.Printf("no pair shares %d+ contiguous tokens across %d files\n", *min, len(docs))
		return 0
	}

	for _, p := range pairs {
		fmt.Printf("[%d tokens] %s\n            %s\n", p.Run, p.A, p.B)
		sample := p.Sample
		if len(sample) > 200 {
			sample = sample[:200] + "…"
		}
		fmt.Printf("   %q\n\n", sample)
	}
	fmt.Printf("%d candidate pair(s) across %d files — classify per references/drift.md\n",
		len(pairs), len(docs))
	return 0
}

func main() {
	os.Exit(run())
}
