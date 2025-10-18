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
	return NewDebugParser(src, false)
}

func NewDebugParser(src io.Reader, dbg bool) Parser {
	p := Parser{
		lx:    NewLexer(src),
		debug: dbg,
	}

	p.readToken()

	return p
}

func (p *Parser) Parse() error {
	if err := p.parseExpression(); err != nil {
		return fmt.Errorf("Parse failure: %w", err)
	}

	if p.tok.Type != EOF {
		return fmt.Errorf("additional top level token: %s", p.tok)
	}

	return nil
}

func (p *Parser) readToken() {
	p.tok = p.lx.NextToken()
	p.log(p.tok.String())
}

func (p *Parser) log(msg string) {
	if p.debug {
		fmt.Println(msg)
	}
}

func (p *Parser) parseExpression() error {
	var err error
	switch p.tok.Type {
	case LBRACE:
		err = p.parseObject()
	case LBRCKT:
		err = p.parseArray()
	case NULL:
	case STRING:
	case NUM:
	case TRUE, FALSE:
	default:
		return fmt.Errorf("invalid expression, unexpected token: %s", p.tok)
	}

	p.readToken()

	return err
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
		return fmt.Errorf("failed to read object key/value: %w", err)
	}

	for p.tok.Type == COMMA {
		p.readToken()
		err := p.readKV()
		if err != nil {
			return fmt.Errorf("failed to read object key/value: %w", err)
		}
	}

	return nil
}

func (p *Parser) readKV() error {
	if p.tok.Type != STRING {
		return fmt.Errorf("expected key string in object found %s", p.tok)
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
