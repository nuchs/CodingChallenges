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
		fmt.Fprintf(os.Stderr, "Failed to load json: %s\n", err)
		os.Exit(1)
	}
	defer closeSource(src)

	p := NewParser(src)
	if err := p.Parse(); err != nil {
		fmt.Printf("Bad JSON: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Good JSON")
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

func closeSource(src io.Closer) {
	if err := src.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to close json source: %s", err)
	}
}
