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

func temp() {
	l := NewLexer("1123 xxx")

	tokens, _ := l.GetTokens()

	p := NewParser(tokens)

	_ = p.parseNumberExpr()
	expVar := p.parseVariableExpr()

	fmt.Printf("exp: %#v\n", expVar)
}

// helper functions
func (p *Parser) current() *Token {
	return p.tokens[p.pos]
}

func (p *Parser) next() {
	p.pos++

}
