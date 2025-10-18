package main

import "fmt"

type TokenType string

type Token struct {
	Type     TokenType
	Literal  string
	row, col int
}

const (
	LBRACE = "{"
	RBRACE = "}"
	LPAREN = "("
	RPAREN = ")"
	LBRCKT = "["
	RBRCKT = "]"
	COLON  = ":"
	COMMA  = ","

	IDENT  = "IDENT"
	STRING = "STRING"
	NUM    = "NUM"
	NULL   = "NULL"
	TRUE   = "TRUE"
	FALSE  = "FALSE"

	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"
)

func NewTokenFromRune(tt TokenType, r rune, row, col int) Token {
	return Token{Type: tt, Literal: string(r), row: row, col: col}
}

func NewTokenFromString(tt TokenType, s string, row, col int) Token {
	return Token{Type: tt, Literal: s, row: row, col: col}
}

func (t Token) String() string {
	lit := ""
	if string(t.Type) != t.Literal {
		lit = fmt.Sprintf("(%s)", t.Literal)
	}
	return fmt.Sprintf("Line %d: %s\t%s", t.row, t.Type, lit)
}

var keywords = map[string]TokenType{
	"null":  NULL,
	"true":  TRUE,
	"false": FALSE,
}

func lookupIdentifier(id string) TokenType {
	if kw, ok := keywords[id]; ok {
		return kw
	}
	return IDENT
}
