package parser

import (
	"fmt"
	"language/ast"
	. "language/lexer"
	"strconv"
	"strings"
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

func (p *Parser) parseStringExpr() (ast.Expr, error) {
	val := p.current().Value
	val = strings.Replace(val, "\\n", "\n", -1)
	p.next()

	return &ast.StringExpr{Val: val}, nil
}

func (p *Parser) parseIdentifierExpr() (*ast.IdentifierExpr, error) {
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

// arrowFunction ::= '(' (param (',' param)*)? ')' identifier '=>' expression | blockStatement ;
func (p *Parser) parseArrowFunc() (ast.Expr, error) {

	// eat LPAREN
	p.next()

	params := []*ast.Param{}

	for !p.isEnd() && p.current().Type != RPAREN {
		paramId, err := p.parseIdentifierExpr()
		if err != nil {
			return nil, err
		}

		paramType, err := p.parseTypeExpr()
		if err != nil {
			return nil, err
		}
		params = append(params, &ast.Param{Id: paramId, Type: paramType})
		if p.current().Type == COMMA {
			p.next()
		}
	}

	if err := p.consume(RPAREN); err != nil {
		return nil, err
	}

	var retType *ast.TypeExpr

	if p.current().Type == ARROW {
		// REVIEW: implicit void return, but maybe it should be
		// whatever the last statement returns
		retType = &ast.TypeExpr{Type: &ast.IdentifierExpr{Name: "void"}}
	} else {
		r, err := p.parseTypeExpr()
		if err != nil {
			return nil, err
		}
		retType = r
	}

	if err := p.consume(ARROW); err != nil {
		return nil, err
	}

	var body *ast.BlockStmt
	if p.current().Type == LBRACE {
		block, err := p.parseBlockStmt()
		if err != nil {
			return nil, err
		}

		body = block
	} else {
		// implicit return
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		body = &ast.BlockStmt{Stmts: []ast.Stmt{&ast.ReturnStmt{Arg: expr}}}
	}

	return &ast.ArrowFunc{Args: params, Body: body, ReturnType: retType}, nil

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

// primaryExpression ::= identifier | number | boolean | string | '(' expression ')' | arrayExpression | arrowFunction;
func (p *Parser) parsePrimaryExpr() (ast.Expr, error) {

	switch p.current().Type {

	case IDENTIFIER:
		return p.parseIdentifierExpr()

	case THIS:
		p.next()
		return &ast.ThisExpr{}, nil
	case NUMBER:
		return p.parseNumberExpr()
	case BOOLEAN:
		return p.parseBooleanExpr()
	case STRING:
		return p.parseStringExpr()
	case LPAREN:
		if (p.peek().Type == IDENTIFIER && p.peek2().Type == IDENTIFIER) ||
			(p.peek().Type == RPAREN && (p.peek2().Type == ARROW || p.peek3().Type == ARROW)) {
			return p.parseArrowFunc()
		} else {

			return p.parseParenExpr()
		}
	case LBRACK:
		return p.parseArrayExpr()

	}

	return nil, NewParserError(p.pos, fmt.Sprintf("expected primary expression, got %v", p.current().Type))
}

// updateExpression ::= primaryExpression ('++' | '--');
func (p *Parser) parseUpdateExpr() (ast.Expr, error) {

	id, err := p.parsePrimaryExpr()
	if err != nil {
		return nil, err
	}

	if p.isEnd() {
		return id, nil
	}

	curr := p.current().Type

	if !p.tokenTypeEqual(curr, INCR, DECR) {
		return id, nil
	}

	p.next()
	if curr == INCR {
		return &ast.UpdateExpr{Arg: id, Op: ast.INC}, nil
	} else {
		return &ast.UpdateExpr{Arg: id, Op: ast.DEC}, nil
	}
}

// unaryExpression ::= updateExpression | '!'unaryExpression | callExpression | sliceExpression;
func (p *Parser) parseUnaryExpr() (ast.Expr, error) {

	if p.current().Type == NOT {
		p.next()
		expr, err := p.parseUnaryExpr()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{Op: "!", Arg: expr}, nil
	} else if !p.isEnd() && !p.isLastToken() && (p.peek().Type == INCR || p.peek().Type == DECR) {
		return p.parseUpdateExpr()
	} else {
		return p.parseCallExpr()
	}
}

// callExpression ::= (identifier) ('(' arguments? ')')*;
func (p *Parser) parseCallExpr() (ast.Expr, error) {

	prev := p.current()
	var call ast.Expr
	call, err := p.parsePrimaryExpr()
	if err != nil {
		return nil, err
	}

	if p.isEnd() ||
		// maybe there's a nicer way to do this.
		// But we don't wanna parse it as call if it's a primitive type
		p.tokenTypeEqual(prev.Type, NUMBER, STRING, BOOLEAN) ||
		p.current().Type != LPAREN {
		return call, nil
	}

	for !p.isEnd() && p.current().Type == LPAREN {
		if err = p.consume(LPAREN); err != nil {
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

		call = &ast.CallExpr{Callee: call, Args: args}

	}

	return call, nil

}

// REVIEW
// // sliceExpression ::= identifier '[' expression ':' expression ' (':' expression)?]';
// func (p *Parser) parseSliceExpr() (ast.Expr, error) {
// 	var err error
// 	id, err := p.parseIdentifierExpr()
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err = p.consume(LBRACK); err != nil {
// 		return nil, err
// 	}

// 	low, err := p.parseExpr()
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err = p.consume(COLON); err != nil {
// 		return nil, err
// 	}
// 	high, err := p.parseExpr()
// 	if err != nil {
// 		return nil, err
// 	}
// 	var step ast.Expr
// 	if p.current().Type == COLON {
// 		p.next()
// 		step, err = p.parseExpr()
// 		if err != nil {
// 			return nil, err
// 		}

// 	}
// 	if err = p.consume(RBRACK); err != nil {
// 		return nil, err
// 	}

// 	return &ast.SliceExpr{Id: id, Low: low, High: high, Step: step}, nil

// }

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
// multiplicativeExpression ::= unaryExpression (multiplicativeOperator unaryExpression)*;
func (p *Parser) parseMultiplicativeExpr() (ast.Expr, error) {
	lhs, err := p.parseUnaryExpr()

	if err != nil {
		return nil, err
	}
	for p.pos < p.len {

		curr := p.current()
		if p.tokenTypeEqual(curr.Type, MUL, DIV, POW, MOD) {

			p.next()
			rhs, err := p.parseUnaryExpr()
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

// type ::= identifier | '(' identifier* ')' '=>' type;
func (p *Parser) parseTypeExpr() (*ast.TypeExpr, error) {

	if p.current().Type == LPAREN {
		p.next()

		params := []*ast.TypeExpr{}

		for !p.isEnd() && p.current().Type != RPAREN {

			paramType, err := p.parseIdentifierExpr()
			if err != nil {
				return nil, err
			}
			params = append(params, &ast.TypeExpr{Type: paramType})

			if p.current().Type == COMMA {
				p.next()
			}
		}

		if err := p.consume(RPAREN); err != nil {
			return nil, err
		}

		if err := p.consume(ARROW); err != nil {
			return nil, err
		}

		retType, err := p.parseIdentifierExpr()
		if err != nil {
			return nil, err
		}

		return &ast.TypeExpr{Type: &ast.FuncTypeExpr{Args: params, ReturnType: &ast.TypeExpr{Type: retType}}}, nil
	} else {

		t, err := p.parseIdentifierExpr()
		if err != nil {
			return nil, err
		}
		return &ast.TypeExpr{Type: t}, nil
	}

}
