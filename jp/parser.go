package main

import (
	"fmt"
	"io"
)

type Parser struct {
	lx    Lexer
	tok   Token
	debug bool
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
		return fmt.Errorf("additional top level token: %q", p.tok)
	}

	return nil
}

func (p *Parser) readToken() {
	p.tok = p.lx.NextToken()
	if p.debug {
		fmt.Println(p.tok)
	}
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
		return fmt.Errorf("invalid expression, unexpected token: %q", p.tok)
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

	if p.tok.Type != RBRCKT {
		return fmt.Errorf("malformed array, expected ']', got '%s'", p.tok)
	}

	return nil
}

func (p *Parser) parseObject() error {
	p.readToken()
	if p.tok.Type == RBRACE {
		return nil
	}

	err := p.readKV()
	if err != nil {
		return fmt.Errorf("failed to rad object kv: %w", err)
	}

	for p.tok.Type == COMMA {
		p.readToken()
		err := p.readKV()
		if err != nil {
			return fmt.Errorf("failed to rad object kv: %w", err)
		}
	}

	return nil
}

func (p *Parser) readKV() error {
	if p.tok.Type != IDENT {
		return fmt.Errorf("expected identifier in object found %s", p.tok)
	}
	p.readToken()
	if p.tok.Type != COLON {
		return fmt.Errorf("expected ':' in object found %s", p.tok)
	}
	p.readToken()
	if err := p.parseExpression(); err != nil {
		return fmt.Errorf("bad expression in object: %w", err)
	}

	return nil
}
