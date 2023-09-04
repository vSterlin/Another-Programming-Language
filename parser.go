package main

import (
	"fmt"
	"strconv"
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

func (p *Parser) parseNumberExpr() Expr {
	val, err := strconv.Atoi(p.current().Value)
	if err != nil {
		return nil
	}
	p.next()

	return &NumberExpr{Val: val}
}

func (p *Parser) parseVariableExpr() Expr {
	name := p.current().Value
	p.next()

	return &VariableExpr{Name: name}
}

func (p *Parser) parsePrimaryExpr() Expr {
	switch p.current().Type {
	case IDENTIFIER:
		return p.parseVariableExpr()
	case NUMBER:
		return p.parseNumberExpr()
	}

	return nil
}

func (p *Parser) parseAdditiveExpr() Expr {
	lhs := p.parsePrimaryExpr()

	for p.pos < p.len && (p.current().Type == ADD || p.current().Type == SUB) {
		switch p.current().Type {
		case ADD:
			// to account for operator
			p.next()
			rhs := p.parsePrimaryExpr()
			lhs = &BinaryExpr{Op: "+", Lhs: lhs, Rhs: rhs}
		case SUB:
			p.next()
			rhs := p.parsePrimaryExpr()
			lhs = &BinaryExpr{Op: "-", Lhs: lhs, Rhs: rhs}
		default:
			return lhs
		}
	}
	return lhs
}

func (p *Parser) parseExpr() Expr {
	return p.parsePrimaryExpr()
}

func temp() {
	l := NewLexer("1 + 77")

	tokens, _ := l.GetTokens()

	p := NewParser(tokens)

	ex := p.parseAdditiveExpr()

	fmt.Printf("ex: %#v\n", ex)

}

// helper functions
func (p *Parser) current() *Token {
	return p.tokens[p.pos]
}

func (p *Parser) next() {
	p.pos++

}
