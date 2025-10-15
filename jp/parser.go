package main

import (
	"fmt"
	"io"
)

type Parser struct {
	lx  Lexer
	tok Token
}

func NewParser(src io.Reader) Parser {
	lx := NewLexer(src)
	return Parser{
		lx:  lx,
		tok: lx.NextToken(),
	}
}

func (p *Parser) Parse() error {
	if err := p.parseExpression(); err != nil {
		return fmt.Errorf("Parse failure: %w", err)
	}

	if p.tok.Type != EOF {
		return fmt.Errorf("additional top level token: %q", p.tok.Literal)
	}

	return nil
}

func (p *Parser) readToken() {
	p.tok = p.lx.NextToken()
}

func (p *Parser) parseExpression() error {
	switch p.tok.Type {
	case LBRACE:
		return p.parseObject()
	case LBRCKT:
		return p.parseArray()
	case NULL:
	case STRING:
	case NUM:
	case TRUE, FALSE:
	default:
		return fmt.Errorf("invalid expression, unexpected token: %q", p.tok.Literal)
	}

	p.readToken()

	return nil
}

func (p *Parser) parseArray() error {
	p.readToken()
	if p.tok.Type == RBRCKT {
		return nil
	}
	if err := p.parseExpression(); err != nil {
		return fmt.Errorf("bad expression in array: %w", err)
	}

	for p.tok.Type == COMMA {
		p.readToken()
		if err := p.parseExpression(); err != nil {
			return fmt.Errorf("bad expression in array: %w", err)
		}
	}

	return nil
}

func (p *Parser) parseObject() error {
	p.readToken()
	if p.tok.Type == RBRACE {
		return nil
	}
	if p.tok.Type != IDENT {
		return fmt.Errorf("expected identifier in object found %s", p.tok)
	}

	return nil
}
