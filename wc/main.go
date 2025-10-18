// Solution for CodingChallenge 1 - build a wc clone
// https://codingchallenges.fyi/challenges/challenge-wc
package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	spec, err := LoadSpec(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	res := NewResults(len(spec.Sources))
	for _, src := range spec.Sources {
		count, err := countSrc(src)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to count source %q: %v", src, err)
			os.Exit(2)
		}
		res.Add(count)
	}

	res.Print(spec)
}

func countSrc(src string) (Counts, error) {
	var empty Counts
	rd := os.Stdin
	if src != "stdin" {
		f, err := os.Open(src)
		if err != nil {
			return empty, fmt.Errorf("failed to open source for reading: %w", err)
		}
		rd = f
	}
	defer closeSource(src, rd)

	counts, err := Count(rd, src)
	if err != nil {
		return empty, err
	}

	return counts, nil
}

func closeSource(src string, rd io.Closer) {
	if err := rd.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to close source %+v: %s", src, err)
	}
}
