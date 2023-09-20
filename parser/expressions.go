package parser

import (
	"fmt"
	"language/ast"
	. "language/lexer"
	"strconv"
)

func (p *Parser) parseNumberExpr() (ast.Expr, error) {
	val, err := strconv.Atoi(p.current().Value)
	if err != nil {
		return nil, NewParserError(p.pos, fmt.Sprintf("expected number, got %s", p.current().Value))
	}
	p.next()

	return &ast.NumberExpr{Val: val}, nil
}

func (p *Parser) parseBooleanExpr() (ast.Expr, error) {
	val, err := strconv.ParseBool(p.current().Value)
	if err != nil {
		return nil, NewParserError(p.pos, fmt.Sprintf("expected boolean, got %s", p.current().Value))
	}
	p.next()

	return &ast.BooleanExpr{Val: val}, nil
}

// TODO: review if some error handling is needed
func (p *Parser) parseStringExpr() (ast.Expr, error) {
	val := p.current().Value
	p.next()

	return &ast.StringExpr{Val: val}, nil
}

// TODO: review if some error handling is needed
func (p *Parser) parseIdentifierExpr() (ast.Expr, error) {
	name := p.current().Value
	p.next()
	return &ast.IdentifierExpr{Name: name}, nil
}

func (p *Parser) parseParenExpr() (ast.Expr, error) {
	p.next()
	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if err := p.consume(RPAREN); err != nil {
		return nil, err
	}
	return val, nil
}

// arrayExpression ::= '[' (expression (',' expression)*)? ']';
func (p *Parser) parseArrayExpr() (ast.Expr, error) {
	if err := p.consume(LBRACK); err != nil {
		return nil, err
	}
	exprs := []ast.Expr{}
	for p.pos < p.len && p.current().Type != RBRACK {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, expr)
		if p.current().Type == COMMA {
			p.next()
		}

	}
	if err := p.consume(RBRACK); err != nil {
		return nil, err
	}

	return &ast.ArrayExpr{Elements: exprs}, nil

}

// primaryExpression ::= identifier | number | boolean | string | '(' expression ')' | callExpression | arrayExpression;
func (p *Parser) parsePrimaryExpr() (ast.Expr, error) {

	switch p.current().Type {

	case IDENTIFIER:
		// if followed by LPAREN, then it's a call expression
		if p.pos+1 < p.len && p.tokens[p.pos+1].Type == LPAREN {
			return p.parseCallExpr()
		} else if p.pos+1 < p.len && p.tokens[p.pos+1].Type == LBRACK {
			return p.parseSliceExpr()
		} else {
			return p.parseIdentifierExpr()
		}
	case NUMBER:
		return p.parseNumberExpr()
	case BOOLEAN:
		return p.parseBooleanExpr()
	case STRING:
		return p.parseStringExpr()
	case LPAREN:
		return p.parseParenExpr()
	case LBRACK:
		return p.parseArrayExpr()

	}

	return nil, NewParserError(p.pos, fmt.Sprintf("expected primary expression, got %d", p.current().Type))
}

// callExpression ::= identifier '(' (identifier (',' identifier)*)? ')';
func (p *Parser) parseCallExpr() (ast.Expr, error) {

	calleeId, err := p.parseIdentifierExpr()
	if err != nil {
		return nil, err
	}
	if err := p.consume(LPAREN); err != nil {
		return nil, err
	}

	args := []ast.Expr{}
	for !p.isEnd() && p.current().Type != RPAREN {
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		if p.current().Type == COMMA {
			p.next()
		}
	}
	if err := p.consume(RPAREN); err != nil {
		return nil, err
	}

	return &ast.CallExpr{
		Callee: calleeId.(*ast.IdentifierExpr),
		Args:   args,
	}, nil

}

// sliceExpression ::= identifier '[' expression ':' expression ' (':' expression)?]';
func (p *Parser) parseSliceExpr() (ast.Expr, error) {
	var err error
	id, err := p.parseIdentifierExpr()
	if err != nil {
		return nil, err
	}

	if err = p.consume(LBRACK); err != nil {
		return nil, err
	}

	low, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if err = p.consume(COLON); err != nil {
		return nil, err
	}
	high, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	var step ast.Expr
	if p.current().Type == COLON {
		p.next()
		step, err = p.parseExpr()
		if err != nil {
			return nil, err
		}

	}
	if err = p.consume(RBRACK); err != nil {
		return nil, err
	}

	return &ast.SliceExpr{Id: id.(*ast.IdentifierExpr), Low: low, High: high, Step: step}, nil

}

// equalityExpression ::= relationalExpression (equalityOperator relationalExpression)*;
func (p *Parser) parseEqualityExpr() (ast.Expr, error) {
	lhs, err := p.parseRelationalExpr()

	if err != nil {
		return nil, err
	}

	for !p.isEnd() {
		if p.tokenTypeEqual(p.current().Type, EQ, NEQ) {
			curr := p.current()
			p.next()

			rhs, err := p.parseRelationalExpr()
			if err != nil {
				return nil, err
			}
			switch curr.Type {
			case EQ:
				lhs = &ast.BinaryExpr{Op: "==", Lhs: lhs, Rhs: rhs}
			case NEQ:
				lhs = &ast.BinaryExpr{Op: "!=", Lhs: lhs, Rhs: rhs}
			}
		} else {
			return lhs, nil
		}
	}
	return lhs, nil
}

// relationalExpression ::= additiveExpression (relationalOperator additiveExpression)*;
func (p *Parser) parseRelationalExpr() (ast.Expr, error) {
	lhs, err := p.parseAdditiveExpr()

	if err != nil {
		return nil, err
	}

	for !p.isEnd() {

		if p.tokenTypeEqual(p.current().Type, LT, GT, LTE, GTE) {

			curr := p.current()
			p.next()
			rhs, err := p.parseAdditiveExpr()
			if err != nil {
				return nil, err
			}
			switch curr.Type {
			case LT:
				lhs = &ast.BinaryExpr{Op: "<", Lhs: lhs, Rhs: rhs}
			case GT:
				lhs = &ast.BinaryExpr{Op: ">", Lhs: lhs, Rhs: rhs}
			case LTE:
				lhs = &ast.BinaryExpr{Op: "<=", Lhs: lhs, Rhs: rhs}
			case GTE:
				lhs = &ast.BinaryExpr{Op: ">=", Lhs: lhs, Rhs: rhs}
			}
		} else {
			return lhs, nil
		}
	}
	return lhs, nil

}

// additiveOperator ::= '+' | '-';
// additiveExpression ::= multiplicativeExpression (additiveOperator multiplicativeExpression)*;
func (p *Parser) parseAdditiveExpr() (ast.Expr, error) {
	lhs, err := p.parseMultiplicativeExpr()

	if err != nil {
		return nil, err
	}
	for p.pos < p.len && (p.current().Type == ADD || p.current().Type == SUB) {
		curr := p.current()
		p.next()
		rhs, err := p.parseMultiplicativeExpr()
		if err != nil {
			return nil, err
		}
		switch curr.Type {
		case ADD:
			lhs = &ast.BinaryExpr{Op: "+", Lhs: lhs, Rhs: rhs}
		case SUB:
			lhs = &ast.BinaryExpr{Op: "-", Lhs: lhs, Rhs: rhs}
		default:
			return lhs, nil
		}
	}
	return lhs, nil
}

// multiplicativeOperator ::= '*' | '/' | '**' | '%';
// multiplicativeExpression ::= primaryExpression (multiplicativeOperator primaryExpression)*;
func (p *Parser) parseMultiplicativeExpr() (ast.Expr, error) {
	lhs, err := p.parsePrimaryExpr()

	if err != nil {
		return nil, err
	}
	for p.pos < p.len {

		curr := p.current()
		if p.tokenTypeEqual(curr.Type, MUL, DIV, POW, MOD) {

			p.next()
			rhs, err := p.parsePrimaryExpr()
			if err != nil {
				return nil, err
			}
			switch curr.Type {
			case MUL:
				lhs = &ast.BinaryExpr{Op: ast.MUL, Lhs: lhs, Rhs: rhs}
			case DIV:
				lhs = &ast.BinaryExpr{Op: ast.DIV, Lhs: lhs, Rhs: rhs}
			case POW:
				lhs = &ast.BinaryExpr{Op: ast.POW, Lhs: lhs, Rhs: rhs}
			case MOD:
				lhs = &ast.BinaryExpr{Op: ast.MOD, Lhs: lhs, Rhs: rhs}
			default:
				return lhs, nil
			}
		} else {
			return lhs, nil
		}
	}
	return lhs, nil
}

// logicAndExpression ::= equalityExpression (logicAndOperator equalityExpression)*;
func (p *Parser) parseAndExpr() (ast.Expr, error) {
	lhs, err := p.parseEqualityExpr()

	if err != nil {
		return nil, err
	}
	for p.pos < p.len && (p.current().Type == AND) {
		p.next()
		rhs, err := p.parseEqualityExpr()
		if err != nil {
			return nil, err
		}
		lhs = &ast.LogicalExpr{Op: ast.AND, Lhs: lhs, Rhs: rhs}

	}
	return lhs, nil
}

// logicOrExpression ::= logicAndExpression (logicOrOperator logicAndExpression)*;
func (p *Parser) parseOrExpr() (ast.Expr, error) {
	lhs, err := p.parseAndExpr()

	if err != nil {
		return nil, err
	}
	for p.pos < p.len && (p.current().Type == OR) {

		p.next()
		rhs, err := p.parseAndExpr()
		if err != nil {
			return nil, err
		}
		lhs = &ast.LogicalExpr{Op: ast.OR, Lhs: lhs, Rhs: rhs}

	}

	return lhs, nil
}

// expression ::= logicOrExpression;
func (p *Parser) parseExpr() (ast.Expr, error) {
	return p.parseOrExpr()
}
