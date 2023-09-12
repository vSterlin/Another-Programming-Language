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

func (p *Parser) parseBooleanExpr() Expr {
	val, err := strconv.ParseBool(p.current().Value)
	if err != nil {
		return nil
	}
	p.next()

	return &BooleanExpr{Val: val}
}

func (p *Parser) parseIdentifierExpr() Expr {
	name := p.current().Value
	p.next()
	return &IdentifierExpr{Name: name}
}

func (p *Parser) parseParenExpr() Expr {
	p.next()
	val := p.parseExpr()
	curr := p.current()
	if curr.Type != RPAREN {
		//	TODO: return error
		fmt.Println("BOOO!")
	}
	// eat the RPAREN
	p.next()

	return val
}

func (p *Parser) parsePrimaryExpr() Expr {
	switch p.current().Type {
	case IDENTIFIER:
		return p.parseIdentifierExpr()
	case NUMBER:
		return p.parseNumberExpr()
	case BOOLEAN:
		return p.parseBooleanExpr()

	case LPAREN:
		return p.parseParenExpr()
	}

	return nil
}

// additiveOperator ::= '+' | '-'
// additiveExpression ::= multiplicativeExpression (additiveOperator multiplicativeExpression)*
func (p *Parser) parseAdditiveExpr() Expr {
	lhs := p.parseMultiplicativeExpr()

	for p.pos < p.len && (p.current().Type == ADD || p.current().Type == SUB) {
		curr := p.current()
		p.next()
		rhs := p.parseMultiplicativeExpr()
		switch curr.Type {
		case ADD:
			lhs = &BinaryExpr{Op: "+", Lhs: lhs, Rhs: rhs}
		case SUB:
			lhs = &BinaryExpr{Op: "-", Lhs: lhs, Rhs: rhs}
		default:
			return lhs
		}
	}
	return lhs
}

// multiplicativeOperator ::= '*' | '/'
// multiplicativeExpression ::= primaryExpression (multiplicativeOperator primaryExpression)*
func (p *Parser) parseMultiplicativeExpr() Expr {
	lhs := p.parsePrimaryExpr()

	for p.pos < p.len && (p.current().Type == MUL || p.current().Type == DIV) {
		curr := p.current()
		p.next()
		rhs := p.parsePrimaryExpr()
		switch curr.Type {
		case MUL:
			lhs = &BinaryExpr{Op: "*", Lhs: lhs, Rhs: rhs}
		case DIV:
			lhs = &BinaryExpr{Op: "/", Lhs: lhs, Rhs: rhs}
		default:
			return lhs
		}
	}
	return lhs
}

// expression ::= primaryExpression | additiveExpression
func (p *Parser) parseExpr() Expr {
	return p.parseAdditiveExpr()
}

// statement ::= expression
func (p *Parser) parseStmt() Stmt {
	ex := p.parseExpr()
	return &ExprStmt{Expr: ex}
}

// program ::= statement*
func (p *Parser) parseProgram() *Program {
	var stmts []Stmt
	for p.pos < p.len {
		stmts = append(stmts, p.parseStmt())
	}
	return &Program{Stmts: stmts}
}

// helper functions
func (p *Parser) current() *Token {
	return p.tokens[p.pos]
}

func (p *Parser) next() {
	p.pos++
}
