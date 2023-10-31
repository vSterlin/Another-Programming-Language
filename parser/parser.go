package parser

import (
	"fmt"
	"language/ast"
	. "language/lexer"
)

type Parser struct {
	tokens []*Token
	pos    int
	len    int
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
		len:    len(tokens),
	}
}

// program ::= statement*;
func (p *Parser) ParseProgram() (*ast.Program, error) {
	var stmts []ast.Stmt
	for p.pos < p.len {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	return &ast.Program{Stmts: stmts}, nil
}

// helper functions
func (p *Parser) current() *Token {
	return p.tokens[p.pos]
}

func (p *Parser) next() {
	p.pos++
}

func (p *Parser) peek() *Token {
	return p.tokens[p.pos+1]
}

func (p *Parser) peek2() *Token {
	return p.tokens[p.pos+2]
}

func (p *Parser) peek3() *Token {
	return p.tokens[p.pos+3]
}

func (p *Parser) isEnd() bool {
	return p.pos >= p.len
}

func (p *Parser) tokenTypeEqual(actual TokenType, expectedType ...TokenType) bool {
	for _, expected := range expectedType {
		if actual == expected {
			return true
		}
	}
	return false
}

func (p *Parser) consume(tokType TokenType) error {

	curr := p.current()
	if curr.Type != tokType {
		fmt.Println("consumeError!")
		return NewParserError(p.pos, fmt.Sprintf("expected %d, got %d", tokType, curr.Type))
	} else {
		p.next()
		return nil
	}
}
