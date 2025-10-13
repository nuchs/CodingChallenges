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
			want: jp.Token{jp.STRING, ""},
		},
		{
			desc: "String",
			data: "\"bacon egg\"",
			want: jp.Token{jp.STRING, "bacon egg"},
		},
		{
			desc: "Special characters",
			data: "\"{}[]():null true false\"",
			want: jp.Token{jp.STRING, "{}[]():null true false"},
		},
		{
			desc: "Quotes",
			data: "\"\\\"arrgh\\\"\"",
			want: jp.Token{jp.STRING, "\\\"arrgh\\\""},
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
			want: jp.Token{jp.NUM, "0"},
		},
		{
			desc: "Positive int",
			data: "123",
			want: jp.Token{jp.NUM, "123"},
		},
		{
			desc: "Negative int",
			data: "-123",
			want: jp.Token{jp.NUM, "-123"},
		},
		{
			desc: "Positive small float",
			data: "0.456",
			want: jp.Token{jp.NUM, "0.456"},
		},
		{
			desc: "Negative small float",
			data: "-0.78901",
			want: jp.Token{jp.NUM, "-0.78901"},
		},
		{
			desc: "Positive big float",
			data: "123.456",
			want: jp.Token{jp.NUM, "123.456"},
		},
		{
			desc: "Negative big float",
			data: "-999.78901",
			want: jp.Token{jp.NUM, "-999.78901"},
		},
		{
			desc: "Big positive e",
			data: "2E+2",
			want: jp.Token{jp.NUM, "2E+2"},
		},
		{
			desc: "Small positive e",
			data: "2e+2",
			want: jp.Token{jp.NUM, "2e+2"},
		},
		{
			desc: "Big negative e",
			data: "2E-2",
			want: jp.Token{jp.NUM, "2E-2"},
		},
		{
			desc: "Small negative e",
			data: "2e-2",
			want: jp.Token{jp.NUM, "2e-2"},
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
			desc: "dash",
			data: "-a",
			err:  "'-' must be followed by a digit",
		},
		{
			desc: "dot",
			data: ".a",
			err:  "'.' must be followed by a digit",
		},
		{
			desc: "e",
			data: "1eaa",
			err:  "exponent must be followed by a sign",
		},
		{
			desc: "E",
			data: "1E11",
			err:  "exponent must be followed by a sign",
		},
		{
			desc: "e too short",
			data: "1e",
			err:  "truncated exponent",
		},
		{
			desc: "E too short",
			data: "1E1",
			err:  "truncated exponent",
		},
		{
			desc: "e-",
			data: "1e-x",
			err:  "exponent must have a value",
		},
		{
			desc: "e+",
			data: "1e+a",
			err:  "exponent must have a value",
		},
		{
			desc: "E-",
			data: "1E-b",
			err:  "exponent must have a value",
		},
		{
			desc: "E+",
			data: "1E+z",
			err:  "exponent must have a value",
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
		{desc: "Open brace", data: "{", want: []jp.TokenType{jp.LBRACE, jp.EOF}},
		{desc: "Close brace", data: "}", want: []jp.TokenType{jp.RBRACE, jp.EOF}},
		{desc: "Open bracket", data: "[", want: []jp.TokenType{jp.LBRCKT, jp.EOF}},
		{desc: "Close bracket", data: "]", want: []jp.TokenType{jp.RBRCKT, jp.EOF}},
		{desc: "Colon", data: ":", want: []jp.TokenType{jp.COLON, jp.EOF}},
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
