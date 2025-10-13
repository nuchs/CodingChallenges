package main

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	LBRACE = "{"
	RBRACE = "}"
	LPAREN = "("
	RPAREN = ")"
	LBRCKT = "["
	RBRCKT = "]"
	COLON  = ":"

	IDENT  = "IDENT"
	STRING = "STRING"
	NUM    = "NUM"
	NULL   = "NULL"
	TRUE   = "TRUE"
	FALSE  = "FALSE"

	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"
)

func newTokenFromRune(tt TokenType, r rune) Token {
	return Token{tt, string(r)}
}

func newTokenFromString(tt TokenType, s string) Token {
	return Token{tt, s}
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
