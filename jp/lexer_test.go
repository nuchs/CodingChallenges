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
			want: jp.NewTokenFromString(jp.STRING, "", 1, 2),
		},
		{
			desc: "String",
			data: "\"bacon egg\"",
			want: jp.NewTokenFromString(jp.STRING, "bacon egg", 1, 11),
		},
		{
			desc: "Special characters",
			data: "\"{}[]():null true false\"",
			want: jp.NewTokenFromString(jp.STRING, "{}[]():null true false", 1, 24),
		},
		{
			desc: "Quotes",
			data: "\"\\\"arrgh\\\"\"",
			want: jp.NewTokenFromString(jp.STRING, "\\\"arrgh\\\"", 1, 11),
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
			want: jp.NewTokenFromString(jp.NUM, "0", 1, 1),
		},
		{
			desc: "Positive int",
			data: "123",
			want: jp.NewTokenFromString(jp.NUM, "123", 1, 3),
		},
		{
			desc: "Negative int",
			data: "-123",
			want: jp.NewTokenFromString(jp.NUM, "-123", 1, 4),
		},
		{
			desc: "Positive small float",
			data: "0.456",
			want: jp.NewTokenFromString(jp.NUM, "0.456", 1, 5),
		},
		{
			desc: "Negative small float",
			data: "-0.78901",
			want: jp.NewTokenFromString(jp.NUM, "-0.78901", 1, 8),
		},
		{
			desc: "Positive big float",
			data: "123.456",
			want: jp.NewTokenFromString(jp.NUM, "123.456", 1, 7),
		},
		{
			desc: "Negative big float",
			data: "-999.78901",
			want: jp.NewTokenFromString(jp.NUM, "-999.78901", 1, 10),
		},
		{
			desc: "Big positive e",
			data: "2E+2",
			want: jp.NewTokenFromString(jp.NUM, "2E+2", 1, 4),
		},
		{
			desc: "Small positive e",
			data: "2e+2",
			want: jp.NewTokenFromString(jp.NUM, "2e+2", 1, 4),
		},
		{
			desc: "Big negative e",
			data: "2E-2",
			want: jp.NewTokenFromString(jp.NUM, "2E-2", 1, 4),
		},
		{
			desc: "Small negative e",
			data: "2e-2",
			want: jp.NewTokenFromString(jp.NUM, "2e-2", 1, 4),
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
