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
		return NewTokenFromString(EOF, "", lx.row)
	}
	if lx.err != nil {
		return NewTokenFromString(
			ILLEGAL,
			fmt.Sprintf("bad token %q: %s", lx.c, lx.err),
			lx.row,
		)
	}

	var tok Token
	switch lx.c {
	case '{':
		tok = NewTokenFromRune(LBRACE, lx.c, lx.row)
	case '}':
		tok = NewTokenFromRune(RBRACE, lx.c, lx.row)
	case '[':
		tok = NewTokenFromRune(LBRCKT, lx.c, lx.row)
	case ']':
		tok = NewTokenFromRune(RBRCKT, lx.c, lx.row)
	case ':':
		tok = NewTokenFromRune(COLON, lx.c, lx.row)
	case ',':
		tok = NewTokenFromRune(COMMA, lx.c, lx.row)
	case '"':
		str, err := lx.readString()
		if err != nil {
			tok = NewTokenFromString(
				ILLEGAL,
				fmt.Sprintf("bad string: %s", err),
				lx.row,
			)
			break
		}
		tok = NewTokenFromString(STRING, str, lx.row)
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		num, err := lx.readNumber()
		if err != nil {
			tok = NewTokenFromString(
				ILLEGAL,
				fmt.Sprintf("bad number: %s", err),
				lx.row,
			)
			break
		}
		tok = NewTokenFromString(NUM, num, lx.row)
	default:
		if unicode.IsLetter(lx.c) || lx.c == '_' {
			ident := lx.readIdentifier()
			tok = NewTokenFromString(lookupIdentifier(ident), ident, lx.row)
		} else {
			tok = NewTokenFromString(
				ILLEGAL,
				fmt.Sprintf("unrecognised token: %v", string(lx.c)),
				lx.row,
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

func (lx *Lexer) peek(num int) ([]rune, error) {
	next, err := lx.src.Peek(num)
	if err != nil && err != io.EOF {
		return nil, err
	}
	runes := make([]rune, 0, num)

	for _, n := range next {
		runes = append(runes, rune(n))
	}

	return runes, err
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
	if err := lx.readFractionalPart(&buf); err != nil {
		return "", err
	}
	if err := lx.readExponent(&buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (lx *Lexer) readIntegralPart(buf *strings.Builder) error {
	next, err := lx.peek(1)
	if lx.c == '-' {
		if err != nil {
			if errors.Is(err, io.EOF) {
				return errors.New("truncated integral part")
			}
			return fmt.Errorf("readNumber - failed to peek after '-': %w", err)
		}
		if !unicode.IsDigit(next[0]) {
			return fmt.Errorf("'-' must be followed by a digit")
		}
	} else if lx.c == '0' && err != io.EOF && unicode.IsDigit(next[0]) {
		return errors.New("numbers cannot lead with zero")
	}
	buf.WriteRune(lx.c)

	lx.readDigits(buf)
	return nil
}

func (lx *Lexer) readFractionalPart(buf *strings.Builder) error {
	next, err := lx.peek(2)
	if len(next) == 0 || next[0] != '.' {
		return nil
	}
	if err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("truncated fractional part")
		}
		return fmt.Errorf("readNumber - failed to peek after '.': %w", err)
	}
	if !unicode.IsDigit(next[1]) {
		return fmt.Errorf("'.' must be followed by a digit")
	}
	lx.readRune()
	buf.WriteRune(lx.c)
	lx.readDigits(buf)

	return nil
}

func (lx *Lexer) readExponent(buf *strings.Builder) error {
	next, err := lx.peek(3)
	length := len(next)
	switch {
	// There is no exponential part
	case length == 0, next[0] != 'e' && next[0] != 'E':
		return nil
	// We start the exponential part but don't have a value
	case length == 1:
		return errors.New("truncated exponent")
	// Valid, unsigned exponetial part e.g. e2, E42, etc
	case unicode.IsDigit(next[1]):
		lx.readRune()
		buf.WriteRune(lx.c)
	// The 'e' is followed by an invalid character
	case next[1] != '+' && next[1] != '-':
		return errors.New("exponent must be followed by a sign or digit")
	// We have a signed exponetial part (e.g. e+, e-) but either there is no
	// numerical value after it
	case length == 2 || !unicode.IsDigit(next[2]):
		return errors.New("signed exponent must be followed by a digit")
	case err != nil && err != io.EOF:
		return fmt.Errorf(
			"readExponent - failed to peek after '%c': %w",
			lx.c,
			err,
		)
	// Valid signed exponential part e.g. e+23, E-123
	default:
		lx.readRune()
		buf.WriteRune(lx.c)
		lx.readRune()
		buf.WriteRune(lx.c)
	}

	lx.readDigits(buf)

	return nil
}

func (lx *Lexer) readDigits(buf *strings.Builder) {
	next, err := lx.peek(1)
	for !errors.Is(err, io.EOF) && unicode.IsDigit(next[0]) {
		lx.readRune()
		buf.WriteRune(lx.c)
		next, err = lx.peek(1)
	}
}

func (lx *Lexer) readString() (string, error) {
	lx.readRune()
	var buf strings.Builder
	esc := false

	for esc || lx.c != '"' {
		switch {
		case lx.err != nil:
			return "", errors.New("unterminated string")
		case !esc && lx.c < 0x20:
			return "", fmt.Errorf("control character 0x%x in stream", lx.c)
		case esc && lx.c == 'u':
			next, err := lx.peek(4)
			if err != nil || !isHex(next) {
				return "", errors.New("invalid unicode sequence")
			}
		case esc && lx.c == '"':
		case esc && lx.c == 'b':
		case esc && lx.c == 'f':
		case esc && lx.c == 'n':
		case esc && lx.c == 'r':
		case esc && lx.c == 't':
		case esc && lx.c == '\\':
		case esc && lx.c == '/':
		case esc:
			return "", fmt.Errorf("invalid escape sequence: \\%c", lx.c)
		}
		buf.WriteRune(lx.c)
		esc = !esc && lx.c == '\\'
		lx.readRune()
	}

	return buf.String(), nil
}

func isHex(seq []rune) bool {
	for _, r := range seq {
		switch r {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		case 'a', 'A', 'b', 'B', 'c', 'C', 'd', 'D', 'e', 'E', 'f', 'F':
		default:
			return false
		}
	}
	return true
}

func (lx *Lexer) readIdentifier() string {
	var buf strings.Builder
	buf.WriteRune(lx.c)

	next, err := lx.peek(1)
	for !errors.Is(err, io.EOF) && unicode.IsLetter(next[0]) || lx.c == '_' || unicode.IsDigit(lx.c) {
		lx.readRune()
		buf.WriteRune(lx.c)
		next, err = lx.peek(1)
	}

	return buf.String()
}
