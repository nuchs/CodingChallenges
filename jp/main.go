// Solution for CodingChallenge 2 - build JSON parser
// https://codingchallenges.fyi/challenges/challenge-json-parser
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	src, err := openSource()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load json: %s", err)
		os.Exit(1)
	}
	defer closeSource(src)

	_ = NewLexer(src)
}

func openSource() (io.ReadCloser, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat stdin: %w", err)
	}
	if (info.Mode() & os.ModeCharDevice) == 0 {
		return os.Stdin, nil
	}

	if len(os.Args) < 2 {
		return nil, errors.New("no json source provided")
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", os.Args[1], err)
	}

	return f, nil
}

func closeSource(source io.Closer) {
	if err := source.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to close json source: %s", err)
	}
}
