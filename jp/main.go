// Solution for CodingChallenge 2 - build JSON parser
// https://codingchallenges.fyi/challenges/challenge-json-parser
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"unicode"
)

type TokenType int

const (
	WhiteSpace TokenType = iota
	OpenBrace
	CloseBrace
)

var tokenNames = []string{
	"WhiteSpace",
	"OpenBrace",
	"CloseBrace",
}

func (tt TokenType) String() string {
	return tokenNames[tt]
}

type Token struct {
	Type   TokenType
	Value  rune
	Line   int
	Column int
}

func main() {
	fmt.Println("Hello JP")
	source, err := openSource()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load json: %s", err)
		os.Exit(1)
	}
	defer func() {
		if err := source.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close json source: %s", err)
		}
	}()

	rd := bufio.NewReader(source)
	ts := make([]Token, 1)
	row := 1
	col := 1
	for {
		r, _, err := rd.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "Unable to tokenise input: %s", err)
			os.Exit(1)
		}

		switch {
		case r == '{':
			ts = append(ts, Token{OpenBrace, '{', row, col})
		case r == '}':
			ts = append(ts, Token{CloseBrace, '}', row, col})
		case unicode.IsSpace(r):
			ts = append(ts, Token{WhiteSpace, r, row, col})
			if r == '\n' {
				row++
				col = 0
			}
		default:
			fmt.Fprintf(os.Stderr, "Invalid token %q @ (row %v, col %v)", r, row, col)
			os.Exit(1)
		}
		col++
	}
	fmt.Println(ts)
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
