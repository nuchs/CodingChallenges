package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Lexer struct {
	src *bufio.Reader
	c   rune
	err error
}

func NewLexer(src io.Reader) Lexer {
	lx := Lexer{
		src: bufio.NewReader(src),
	}
	lx.readRune()

	return lx
}

func (lx *Lexer) NextToken() Token {
	lx.skipWhitespace()

	if lx.err == io.EOF {
		return newTokenFromString(EOF, "")
	}
	if lx.err != nil {
		return newTokenFromString(
			ILLEGAL,
			fmt.Sprintf("bad token %q: %s", lx.c, lx.err),
		)
	}

	var tok Token
	switch lx.c {
	case '{':
		tok = newTokenFromRune(LBRACE, lx.c)
	case '}':
		tok = newTokenFromRune(RBRACE, lx.c)
	case '[':
		tok = newTokenFromRune(LBRCKT, lx.c)
	case ']':
		tok = newTokenFromRune(RBRCKT, lx.c)
	case ':':
		tok = newTokenFromRune(COLON, lx.c)
	case '"':
		lx.readRune()
		tok = newTokenFromString(STRING, lx.readString())
	default:
		if unicode.IsLetter(lx.c) || lx.c == '_' {
			ident := lx.readIdentifier()
			tok = newTokenFromString(lookupIdentifier(ident), ident)
		} else if num, err := lx.readNumber(); err == nil {
			tok = newTokenFromString(NUM, num)
		} else {
			tok = newTokenFromString(
				ILLEGAL,
				fmt.Sprintf("unrecognised token: %v", err),
			)
		}
	}

	lx.readRune()

	return tok
}

func (lx *Lexer) readRune() {
	lx.c, _, lx.err = lx.src.ReadRune()
}

func (lx *Lexer) skipWhitespace() {
	if lx.err != nil {
		return
	}
	for lx.c == ' ' || lx.c == '\t' || lx.c == '\n' || lx.c == '\r' {
		lx.readRune()
	}
}

func (lx *Lexer) readNumber() (string, error) {
	var buf strings.Builder

	if lx.c == '-' {
		next, err := lx.src.Peek(1)
		if err != nil {
			return "", fmt.Errorf("readNumber - failed to peek after '-': %w", err)
		}
		if !unicode.IsDigit(rune(next[0])) {
			return "", fmt.Errorf("'-' must be followed by a digit")
		}
		buf.WriteRune(lx.c)
		lx.readRune()
	}

	lx.readDigits(&buf)

	if lx.c == '.' {
		next, err := lx.src.Peek(1)
		if err != nil {
			return "", fmt.Errorf("readNumber - failed to peek after '.': %w", err)
		}
		if !unicode.IsDigit(rune(next[0])) {
			return "", fmt.Errorf("'.' must be followed by a digit")
		}
		buf.WriteRune(lx.c)
		lx.readRune()
		lx.readDigits(&buf)
	}

	if lx.c == 'e' || lx.c == 'E' {
		next, err := lx.src.Peek(2)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return "", errors.New("truncated exponent")
			}
			return "", fmt.Errorf(
				"readNumber - failed to peek after '%c': %w",
				lx.c,
				err,
			)
		}
		sign, value := rune(next[0]), rune(next[1])
		if sign != '+' && sign != '-' {
			return "", errors.New("exponent must be followed by a sign")
		}
		if !unicode.IsDigit(value) {
			return "", errors.New("exponent must have a value")
		}
		buf.WriteRune(lx.c)
		lx.readRune()
		buf.WriteRune(lx.c)
		lx.readRune()
		lx.readDigits(&buf)
	}

	return buf.String(), nil
}

func (lx *Lexer) readDigits(buf *strings.Builder) {
	for unicode.IsDigit(lx.c) {
		buf.WriteRune(lx.c)
		lx.readRune()
	}
}

func (lx *Lexer) readString() string {
	var buf strings.Builder
	esc := false

	for esc || lx.c != '"' {
		buf.WriteRune(lx.c)
		esc = !esc && lx.c == '\\'
		lx.readRune()
	}

	return buf.String()
}

func (lx *Lexer) readIdentifier() string {
	var buf strings.Builder

	for unicode.IsLetter(lx.c) || lx.c == '_' || unicode.IsDigit(lx.c) {
		buf.WriteRune(lx.c)
		lx.readRune()
	}

	return buf.String()
}
