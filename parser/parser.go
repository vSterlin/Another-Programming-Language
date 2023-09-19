package parser

import (
	"fmt"
	"language/ast"
	. "language/lexer"
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

// deferStatement ::= 'defer' callExpression;
func (p *Parser) parseDeferStmt() (ast.Stmt, error) {
	if err := p.consume(DEFER); err != nil {
		return nil, err
	}
	ex, err := p.parseCallExpr()
	if err != nil {
		return nil, err
	}
	return &ast.DeferStmt{Call: ex.(*ast.CallExpr)}, nil
}

// rangeStatement ::= 'for' identifierExpression ':=' 'range' expression blockStatement;
func (p *Parser) parseRangeStmt() (ast.Stmt, error) {

	if err := p.consume(FOR); err != nil {
		return nil, err
	}

	id, err := p.parseIdentifierExpr()

	if err != nil {
		return nil, err
	}

	if err := p.consume(DECLARE); err != nil {
		return nil, err
	}

	if err := p.consume(RANGE); err != nil {
		return nil, err
	}
	ex, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	body, err := p.parseBlockStmt()

	if err != nil {
		return nil, err
	}

	return &ast.RangeStmt{Id: id.(*ast.IdentifierExpr), Expr: ex, Body: body.(*ast.BlockStmt)}, nil
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
				lhs = &ast.BinaryExpr{Op: "*", Lhs: lhs, Rhs: rhs}
			case DIV:
				lhs = &ast.BinaryExpr{Op: "/", Lhs: lhs, Rhs: rhs}
			case POW:
				lhs = &ast.BinaryExpr{Op: "**", Lhs: lhs, Rhs: rhs}
			case MOD:
				lhs = &ast.BinaryExpr{Op: "%", Lhs: lhs, Rhs: rhs}
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
		lhs = &ast.LogicalExpr{Op: "&&", Lhs: lhs, Rhs: rhs}

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
		lhs = &ast.LogicalExpr{Op: "||", Lhs: lhs, Rhs: rhs}

	}

	return lhs, nil
}

// expression ::= logicOrExpression;
func (p *Parser) parseExpr() (ast.Expr, error) {
	return p.parseOrExpr()
}

// variableAssignmentStatement ::= identifier ('=' | ':=') expression;
func (p *Parser) parseVarAssignStmt() (ast.Stmt, error) {
	id, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.isEnd() || (p.current().Type != ASSIGN && p.current().Type != DECLARE) {
		return &ast.ExprStmt{Expr: id}, nil
	}
	assignOp := p.current().Value
	p.next()

	ex, _ := p.parseExpr()

	// if err != nil {
	// 	return nil, err
	// }
	return &ast.VarAssignStmt{Id: id.(*ast.IdentifierExpr), Init: ex, Op: assignOp}, nil
}

// functionDeclaration ::= 'func' identifier '(' (identifier (',' identifier)*)? ')' blockStatement;
func (p *Parser) parseFuncDecStmt() (ast.Stmt, error) {

	if err := p.consume(FUNC); err != nil {
		return nil, err
	}

	id, err := p.parseIdentifierExpr()
	if err != nil {
		return nil, err
	}

	if err := p.consume(LPAREN); err != nil {
		return nil, err
	}

	var args []*ast.IdentifierExpr
	for !p.isEnd() && p.current().Type != RPAREN {
		arg, err := p.parseIdentifierExpr()
		if err != nil {
			return nil, err
		}
		args = append(args, arg.(*ast.IdentifierExpr))
		if p.current().Type == COMMA {
			p.next()
		}
	}

	if err := p.consume(RPAREN); err != nil {
		return nil, err
	}

	body, err := p.parseBlockStmt()
	if err != nil {
		return nil, err
	}
	return &ast.FuncDecStmt{Id: id.(*ast.IdentifierExpr), Args: args, Body: body.(*ast.BlockStmt)}, nil
}

// blockStatement ::= '{' statement* '}';
func (p *Parser) parseBlockStmt() (ast.Stmt, error) {
	if err := p.consume(LBRACE); err != nil {
		return nil, err
	}
	var stmts []ast.Stmt
	for p.pos < p.len && p.current().Type != RBRACE {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	if err := p.consume(RBRACE); err != nil {
		return nil, err
	}
	return &ast.BlockStmt{Stmts: stmts}, nil
}

// whileStatement ::= 'while' [expression] blockStatement;
func (p *Parser) parseWhileStmt() (ast.Stmt, error) {
	var err error
	if err = p.consume(WHILE); err != nil {
		return nil, err
	}
	var test ast.Expr
	if p.current().Type == LBRACE {
		test = &ast.BooleanExpr{Val: true}
	} else {
		test, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}
	body, err := p.parseBlockStmt()
	if err != nil {
		return nil, err
	}
	return &ast.WhileStmt{Test: test, Body: body}, nil
}

// ifStatement ::= 'if' expression blockStatement ('else if' expression blockStatement)* ('else' blockStatement)?;
func (p *Parser) parseIfStmt() (ast.Stmt, error) {
	var err error
	if err = p.consume(IF); err != nil {
		return nil, err
	}

	test, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	consequent, err := p.parseBlockStmt()
	if err != nil {
		return nil, err
	}
	var alternate ast.Stmt

	if !p.isEnd() && p.current().Type == ELSE {

		if err := p.consume(ELSE); err != nil {
			return nil, err
		}

		if p.current().Type == IF {
			alternate, err = p.parseIfStmt()
			if err != nil {
				return nil, err
			}

		} else {
			alternate, err = p.parseBlockStmt()
			if err != nil {
				return nil, err
			}

		}
	}
	return &ast.IfStmt{Test: test, Consequent: consequent, Alternate: alternate}, nil
}

// statement ::= expression | variableDeclarationStatement
// | variableAssignmentStatement | blockStatement
// | whileStatement | functionDeclaration | ifStatement | deferStatement | rangeStatement;
func (p *Parser) parseStmt() (ast.Stmt, error) {

	switch p.current().Type {

	case LBRACE:
		return p.parseBlockStmt()
	case WHILE:
		return p.parseWhileStmt()
	case FUNC:
		return p.parseFuncDecStmt()
	case IF:
		return p.parseIfStmt()
	case DEFER:
		return p.parseDeferStmt()
	case FOR:
		return p.parseRangeStmt()
	case IDENTIFIER:
		return p.parseVarAssignStmt()
	default:
		ex, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &ast.ExprStmt{Expr: ex}, nil
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
