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
	row int
	col int
}

func NewLexer(src io.Reader) Lexer {
	lx := Lexer{
		src: bufio.NewReader(src),
	}
	lx.row = 1
	lx.readRune()

	return lx
}

func (lx *Lexer) NextToken() Token {
	lx.skipWhitespace()

	if lx.err == io.EOF {
		return NewTokenFromString(EOF, "", lx.row, lx.col)
	}
	if lx.err != nil {
		return NewTokenFromString(
			ILLEGAL,
			fmt.Sprintf("bad token %q: %s", lx.c, lx.err),
			lx.row,
			lx.col,
		)
	}

	var tok Token
	switch lx.c {
	case '{':
		tok = NewTokenFromRune(LBRACE, lx.c, lx.row, lx.col)
	case '}':
		tok = NewTokenFromRune(RBRACE, lx.c, lx.row, lx.col)
	case '[':
		tok = NewTokenFromRune(LBRCKT, lx.c, lx.row, lx.col)
	case ']':
		tok = NewTokenFromRune(RBRCKT, lx.c, lx.row, lx.col)
	case ':':
		tok = NewTokenFromRune(COLON, lx.c, lx.row, lx.col)
	case ',':
		tok = NewTokenFromRune(COMMA, lx.c, lx.row, lx.col)
	case '"':
		lx.readRune()
		tok = NewTokenFromString(STRING, lx.readString(), lx.row, lx.col)
	default:
		if unicode.IsLetter(lx.c) || lx.c == '_' {
			ident := lx.readIdentifier()
			tok = NewTokenFromString(lookupIdentifier(ident), ident, lx.row, lx.col)
		} else if num, err := lx.readNumber(); err == nil {
			tok = NewTokenFromString(NUM, num, lx.row, lx.col)
		} else {
			tok = NewTokenFromString(
				ILLEGAL,
				fmt.Sprintf("unrecognised token: %v", err),
				lx.row,
				lx.col,
			)
		}
	}

	lx.readRune()

	return tok
}

func (lx *Lexer) readRune() {
	lx.c, _, lx.err = lx.src.ReadRune()
	if lx.c == '\n' {
		lx.row++
		lx.col = 1
	} else {
		lx.col++
	}
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

	if err := lx.readIntegralPart(&buf); err != nil {
		return "", err
	}
	if lx.c == '.' {
		if err := lx.readFractionalPart(&buf); err != nil {
			return "", err
		}
	}
	if lx.c == 'e' || lx.c == 'E' {
		if err := lx.readExponent(&buf); err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}

func (lx *Lexer) readIntegralPart(buf *strings.Builder) error {
	if lx.c == '-' {
		next, err := lx.src.Peek(1)
		if err != nil {
			return fmt.Errorf("readNumber - failed to peek after '-': %w", err)
		}
		if !unicode.IsDigit(rune(next[0])) {
			return fmt.Errorf("'-' must be followed by a digit")
		}
		buf.WriteRune(lx.c)
		lx.readRune()
	}

	lx.readDigits(buf)
	return nil
}

func (lx *Lexer) readFractionalPart(buf *strings.Builder) error {
	next, err := lx.src.Peek(1)
	if err != nil {
		return fmt.Errorf("readNumber - failed to peek after '.': %w", err)
	}
	if !unicode.IsDigit(rune(next[0])) {
		return fmt.Errorf("'.' must be followed by a digit")
	}
	buf.WriteRune(lx.c)
	lx.readRune()
	lx.readDigits(buf)

	return nil
}

func (lx *Lexer) readExponent(buf *strings.Builder) error {
	next, err := lx.src.Peek(2)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("truncated exponent")
		}
		return fmt.Errorf(
			"readNumber - failed to peek after '%c': %w",
			lx.c,
			err,
		)
	}
	sign, value := rune(next[0]), rune(next[1])
	if sign != '+' && sign != '-' {
		return errors.New("exponent must be followed by a sign")
	}
	if !unicode.IsDigit(value) {
		return errors.New("exponent must have a value")
	}
	buf.WriteRune(lx.c)
	lx.readRune()
	buf.WriteRune(lx.c)
	lx.readRune()
	lx.readDigits(buf)

	return nil
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
