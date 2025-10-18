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

	for _, src := range spec.Sources {
		if err := countSrc(spec, src); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to count source %q: %v", src, err)
			os.Exit(2)
		}
	}
}

func countSrc(spec Spec, src string) error {
	rd, err := openSource(src)
	if err != nil {
		return err
	}
	defer closeSource(src, rd)

	counts, err := Count(rd)
	if err != nil {
		return err
	}

	fmt.Println(counts.Format(spec, src))

	return nil
}

func openSource(src string) (io.ReadCloser, error) {
	if src == "stdin" {
		return os.Stdin, nil
	}

	f, err := os.Open(src)
	if err != nil {
		return nil, fmt.Errorf("failed to open source for reading: %w", err)
	}
	return f, nil
}

func closeSource(src string, rd io.Closer) {
	if err := rd.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to close source %+v: %s", src, err)
	}
}
