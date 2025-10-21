package main_test

import (
	"reflect"
	"strings"
	"testing"

	jp "github.com/nuchs/ccjp"
)

func TestStringTokens(t *testing.T) {
	testCases := []struct {
		desc string
		data string
		want jp.Token
	}{
		{
			desc: "Empty string",
			data: "\"\"",
			want: jp.NewTokenFromString(jp.STRING, "", 1),
		},
		{
			desc: "String",
			data: "\"bacon egg\"",
			want: jp.NewTokenFromString(jp.STRING, "bacon egg", 1),
		},
		{
			desc: "Special characters",
			data: "\"{}[]():null true false\"",
			want: jp.NewTokenFromString(jp.STRING, "{}[]():null true false", 1),
		},
		{
			desc: "Escapes",
			data: `"\"\b\f\r\n\t\/\\\u0123\uaAfF"`,
			want: jp.NewTokenFromString(jp.STRING, `\"\b\f\r\n\t\/\\\u0123\uaAfF`, 1),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			lx := jp.NewLexer(strings.NewReader(tC.data))
			got := lx.NextToken()
			if !reflect.DeepEqual(got, tC.want) {
				t.Fatalf("Bad token: got %+v, want %+v", got, tC.want)
			}
		})
	}
}

func TestNumberTokens(t *testing.T) {
	testCases := []struct {
		desc string
		data string
		want jp.Token
	}{
		{
			desc: "zero",
			data: "0",
			want: jp.NewTokenFromString(jp.NUM, "0", 1),
		},
		{
			desc: "Positive int",
			data: "123",
			want: jp.NewTokenFromString(jp.NUM, "123", 1),
		},
		{
			desc: "Negative int",
			data: "-123",
			want: jp.NewTokenFromString(jp.NUM, "-123", 1),
		},
		{
			desc: "Positive small float",
			data: "0.456",
			want: jp.NewTokenFromString(jp.NUM, "0.456", 1),
		},
		{
			desc: "Negative small float",
			data: "-0.78901",
			want: jp.NewTokenFromString(jp.NUM, "-0.78901", 1),
		},
		{
			desc: "Positive big float",
			data: "123.456",
			want: jp.NewTokenFromString(jp.NUM, "123.456", 1),
		},
		{
			desc: "Negative big float",
			data: "-999.78901",
			want: jp.NewTokenFromString(jp.NUM, "-999.78901", 1),
		},
		{
			desc: "Big e",
			data: "2E23",
			want: jp.NewTokenFromString(jp.NUM, "2E23", 1),
		},
		{
			desc: "Small e",
			data: "3e4",
			want: jp.NewTokenFromString(jp.NUM, "3e4", 1),
		},
		{
			desc: "Big positive e",
			data: "2E+2",
			want: jp.NewTokenFromString(jp.NUM, "2E+2", 1),
		},
		{
			desc: "Small positive e",
			data: "2e+2",
			want: jp.NewTokenFromString(jp.NUM, "2e+2", 1),
		},
		{
			desc: "Big negative e",
			data: "2E-2",
			want: jp.NewTokenFromString(jp.NUM, "2E-2", 1),
		},
		{
			desc: "Small negative e",
			data: "2e-2",
			want: jp.NewTokenFromString(jp.NUM, "2e-2", 1),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			lx := jp.NewLexer(strings.NewReader(tC.data))
			got := lx.NextToken()
			if !reflect.DeepEqual(got, tC.want) {
				t.Fatalf("Bad token: got %+v, want %+v", got, tC.want)
			}
		})
	}
}

func TestBadTokens(t *testing.T) {
	testCases := []struct {
		desc string
		data string
		err  string
	}{
		{
			desc: "leading zero",
			data: "0123",
			err:  "numbers cannot lead with zero",
		},
		{
			desc: "truncated dash",
			data: "-",
			err:  "truncated integral part",
		},
		{
			desc: "dash non number",
			data: "-a",
			err:  "'-' must be followed by a digit",
		},
		{
			desc: "truncated dot",
			data: "1.",
			err:  "truncated fractional part",
		},
		{
			desc: "dot non number",
			data: "0.a",
			err:  "'.' must be followed by a digit",
		},
		{
			desc: "bare dot",
			data: ".1",
			err:  "unrecognised token: .",
		},
		{
			desc: "truncated e",
			data: "1e",
			err:  "truncated exponent",
		},
		{
			desc: "truncated E",
			data: "1E",
			err:  "truncated exponent",
		},
		{
			desc: "e- non number",
			data: "1e-x",
			err:  "signed exponent must be followed by a digit",
		},
		{
			desc: "e+ non number",
			data: "1e+a",
			err:  "signed exponent must be followed by a digit",
		},
		{
			desc: "E- non number",
			data: "1E-b",
			err:  "signed exponent must be followed by a digit",
		},
		{
			desc: "E+ non number",
			data: "1E+z",
			err:  "signed exponent must be followed by a digit",
		},
		{
			desc: "Unterminated string",
			data: "\"blah",
			err:  "unterminated string",
		},
		{
			desc: "bad esc sequence",
			data: `"what's the \q word?"`,
			err:  `invalid escape sequence: \q`,
		},
		{
			desc: "line break",
			data: "\"blah\nblah\"",
			err:  "control character 0xa in stream",
		},
		{
			desc: "unicode too short",
			data: `"\u111"`,
			err:  "invalid unicode sequence",
		},
		{
			desc: "unicode invalid",
			data: `"\u1X23"`,
			err:  "invalid unicode sequence",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			lx := jp.NewLexer(strings.NewReader(tC.data))
			got := lx.NextToken()
			if got.Type != jp.ILLEGAL {
				t.Fatalf("Wrong token type: got %q, want \"ILLEGAL\"", got.Type)
			}
			if !strings.HasSuffix(got.Literal, tC.err) {
				t.Fatalf("Wrong error: got %q, want %q", got.Literal, tC.err)
			}
		})
	}
}

func TestSimpleTokens(t *testing.T) {
	testCases := []struct {
		desc string
		data string
		want []jp.TokenType
	}{
		{desc: "empty", data: "", want: []jp.TokenType{jp.EOF}},
		{desc: "Identifier", data: "bob", want: []jp.TokenType{jp.IDENT, jp.EOF}},
		{desc: "Open brace", data: "{", want: []jp.TokenType{jp.LBRACE, jp.EOF}},
		{desc: "Close brace", data: "}", want: []jp.TokenType{jp.RBRACE, jp.EOF}},
		{desc: "Open bracket", data: "[", want: []jp.TokenType{jp.LBRCKT, jp.EOF}},
		{desc: "Close bracket", data: "]", want: []jp.TokenType{jp.RBRCKT, jp.EOF}},
		{desc: "Colon", data: ":", want: []jp.TokenType{jp.COLON, jp.EOF}},
		{desc: "Comma", data: ",", want: []jp.TokenType{jp.COMMA, jp.EOF}},
		{desc: "Null", data: "null", want: []jp.TokenType{jp.NULL, jp.EOF}},
		{desc: "True", data: "true", want: []jp.TokenType{jp.TRUE, jp.EOF}},
		{desc: "False", data: "false", want: []jp.TokenType{jp.FALSE, jp.EOF}},
		{desc: "Skip Whitespace", data: " \n\r\t", want: []jp.TokenType{jp.EOF}},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			lx := jp.NewLexer(strings.NewReader(tC.data))
			got := readAll(&lx)
			if !reflect.DeepEqual(got, tC.want) {
				t.Fatalf("Bad tokenisation: got %+v, want %+v", got, tC.want)
			}
		})
	}
}

func readAll(lx *jp.Lexer) []jp.TokenType {
	tt := []jp.TokenType{}

	for {
		tok := lx.NextToken()
		tt = append(tt, tok.Type)
		if tok.Type == jp.EOF || tok.Type == jp.ILLEGAL {
			break
		}
	}

	return tt
}
