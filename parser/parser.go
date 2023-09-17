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

func (p *Parser) parseNumberExpr() ast.Expr {
	val, err := strconv.Atoi(p.current().Value)
	if err != nil {
		return nil
	}
	p.next()

	return &ast.NumberExpr{Val: val}
}

func (p *Parser) parseBooleanExpr() ast.Expr {
	val, err := strconv.ParseBool(p.current().Value)
	if err != nil {
		return nil
	}
	p.next()

	return &ast.BooleanExpr{Val: val}
}

func (p *Parser) parseStringExpr() ast.Expr {
	val := p.current().Value
	p.next()

	return &ast.StringExpr{Val: val}
}

func (p *Parser) parseIdentifierExpr() ast.Expr {
	name := p.current().Value
	p.next()
	return &ast.IdentifierExpr{Name: name}
}

func (p *Parser) parseParenExpr() ast.Expr {
	p.next()
	val := p.parseExpr()
	p.consume(RPAREN)
	return val
}

// arrayExpression ::= '[' (expression (',' expression)*)? ']';
func (p *Parser) parseArrayExpr() ast.Expr {
	err := p.consume(LBRACK)
	if err != nil {
		return nil
	}
	exprs := []ast.Expr{}
	for p.pos < p.len && p.current().Type != RBRACK {
		expr := p.parseExpr()
		exprs = append(exprs, expr)
		if p.current().Type == COMMA {
			p.next()
		}

	}
	err = p.consume(RBRACK)
	if err != nil {
		return nil
	}
	return &ast.ArrayExpr{Elements: exprs}

}

// primaryExpression ::= identifier | number | boolean | string | '(' expression ')' | callExpression | arrayExpression;
func (p *Parser) parsePrimaryExpr() ast.Expr {

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

	return nil
}

// callExpression ::= identifier '(' (identifier (',' identifier)*)? ')';
func (p *Parser) parseCallExpr() ast.Expr {

	calleeId := p.parseIdentifierExpr()
	if p.current().Type != LPAREN {
		fmt.Println("parseCallExprError!")
	}
	p.next()

	args := []ast.Expr{}
	for !p.isEnd() && p.current().Type != RPAREN {
		arg := p.parseExpr()
		args = append(args, arg)
		if p.current().Type == COMMA {
			p.next()
		}
	}
	p.consume(RPAREN)
	return &ast.CallExpr{
		Callee: calleeId.(*ast.IdentifierExpr),
		Args:   args,
	}

}

// deferStatement ::= 'defer' callExpression;
func (p *Parser) parseDeferStmt() ast.Stmt {
	p.consume(DEFER)
	ex := p.parseCallExpr()
	return &ast.DeferStmt{Call: ex.(*ast.CallExpr)}
}

// rangeStatement ::= 'for' variableAssignmentStatement 'range' expression blockStatement;
func (p *Parser) parseRangeStmt() ast.Stmt {

	p.consume(FOR)
	id := p.parseVarAssignStmt()
	if p.current().Type != RANGE {
		fmt.Println("parseRangeStmtError!")
	}
	p.consume(RANGE)
	ex := p.parseExpr()
	body := p.parseBlockStmt()

	return &ast.RangeStmt{Id: id.(*ast.VarAssignStmt).Id, Expr: ex, Body: body.(*ast.BlockStmt)}
}

// sliceExpression ::= identifier '[' expression ':' expression ' (':' expression)?]';
func (p *Parser) parseSliceExpr() ast.Expr {
	id := p.parseIdentifierExpr()
	if p.current().Type != LBRACK {
		fmt.Println("parseDeferStmtError 1!")
	}
	p.next()
	low := p.parseExpr()
	if p.current().Type != COLON {
		fmt.Println("parseDeferStmtError 2!")
	}
	p.next()
	high := p.parseExpr()

	var step ast.Expr
	if p.current().Type == COLON {
		p.next()
		step = p.parseExpr()
	}

	if p.current().Type != RBRACK {
		fmt.Println("parseDeferStmtError 3!")
	}
	p.next()
	return &ast.SliceExpr{Id: id.(*ast.IdentifierExpr), Low: low, High: high, Step: step}

}

// equalityExpression ::= relationalExpression (equalityOperator relationalExpression)*;
func (p *Parser) parseEqualityExpr() ast.Expr {
	lhs := p.parseRelationalExpr()

	for !p.isEnd() {
		switch p.current().Type {
		case EQ, NEQ:
			curr := p.current()
			p.next()
			rhs := p.parseRelationalExpr()
			switch curr.Type {
			case EQ:
				lhs = &ast.BinaryExpr{Op: "==", Lhs: lhs, Rhs: rhs}
			case NEQ:
				lhs = &ast.BinaryExpr{Op: "!=", Lhs: lhs, Rhs: rhs}
			}
		default:
			return lhs

		}
	}
	return lhs
}

// relationalExpression ::= additiveExpression (relationalOperator additiveExpression)*;
func (p *Parser) parseRelationalExpr() ast.Expr {
	lhs := p.parseAdditiveExpr()

	for !p.isEnd() {
		switch p.current().Type {
		case LT, GT, LTE, GTE:
			curr := p.current()
			p.next()
			rhs := p.parseAdditiveExpr()
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
		default:
			return lhs

		}
	}
	return lhs

}

// additiveOperator ::= '+' | '-';
// additiveExpression ::= multiplicativeExpression (additiveOperator multiplicativeExpression)*;
func (p *Parser) parseAdditiveExpr() ast.Expr {
	lhs := p.parseMultiplicativeExpr()

	for p.pos < p.len && (p.current().Type == ADD || p.current().Type == SUB) {
		curr := p.current()
		p.next()
		rhs := p.parseMultiplicativeExpr()
		switch curr.Type {
		case ADD:
			lhs = &ast.BinaryExpr{Op: "+", Lhs: lhs, Rhs: rhs}
		case SUB:
			lhs = &ast.BinaryExpr{Op: "-", Lhs: lhs, Rhs: rhs}
		default:
			return lhs
		}
	}
	return lhs
}

// multiplicativeOperator ::= '*' | '/' | '**' | '%';
// multiplicativeExpression ::= primaryExpression (multiplicativeOperator primaryExpression)*;
func (p *Parser) parseMultiplicativeExpr() ast.Expr {
	lhs := p.parsePrimaryExpr()

	for p.pos < p.len {

		curr := p.current()
		switch curr.Type {
		case MUL, DIV, POW, MOD:
			p.next()
			rhs := p.parsePrimaryExpr()
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
				return lhs
			}

		default:
			return lhs
		}
	}
	return lhs
}

// expression ::= additiveExpression;
func (p *Parser) parseExpr() ast.Expr {
	return p.parseEqualityExpr()
}

// TODO: I will get rid of this
// variableDeclarationStatement ::= 'let' identifier '=' expression;
func (p *Parser) parseVarDecStmt() ast.Stmt {
	p.next()
	id := p.parseIdentifierExpr()

	err := p.consume(ASSIGN)
	if err != nil {
		return nil
	}

	ex := p.parseExpr()

	// TODO: do some type checking here
	return &ast.VarDecStmt{Id: id.(*ast.IdentifierExpr), Init: ex}
}

// variableAssignmentStatement ::= identifier ('=' | ':=') expression;
func (p *Parser) parseVarAssignStmt() ast.Stmt {
	id := p.parseExpr()
	if p.isEnd() || (p.current().Type != ASSIGN && p.current().Type != DECLARE) {
		return &ast.ExprStmt{Expr: id}
	}
	assignOp := p.current().Value
	p.next()
	ex := p.parseExpr()
	return &ast.VarAssignStmt{Id: id.(*ast.IdentifierExpr), Init: ex, Op: assignOp}
}

// functionDeclaration ::= 'func' identifier '(' (identifier (',' identifier)*)? ')' blockStatement;
func (p *Parser) parseFuncDecStmt() ast.Stmt {
	p.consume(FUNC)
	id := p.parseIdentifierExpr()

	p.consume(LPAREN)

	var args []*ast.IdentifierExpr
	for !p.isEnd() && p.current().Type != RPAREN {
		arg := p.parseIdentifierExpr()
		args = append(args, arg.(*ast.IdentifierExpr))
		if p.current().Type == COMMA {
			p.next()
		}
	}
	p.consume(RPAREN)
	body := p.parseBlockStmt()
	return &ast.FuncDecStmt{Id: id.(*ast.IdentifierExpr), Args: args, Body: body.(*ast.BlockStmt)}
}

// blockStatement ::= '{' statement* '}';
func (p *Parser) parseBlockStmt() ast.Stmt {
	p.consume(LBRACE)
	var stmts []ast.Stmt
	for p.pos < p.len && p.current().Type != RBRACE {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil
		}
		stmts = append(stmts, stmt)
	}
	p.consume(RBRACE)
	return &ast.BlockStmt{Stmts: stmts}
}

// whileStatement ::= 'while' [expression] blockStatement;
func (p *Parser) parseWhileStmt() ast.Stmt {
	p.consume(WHILE)
	var test ast.Expr
	if p.current().Type == LBRACE {
		test = &ast.BooleanExpr{Val: true}
	} else {
		test = p.parseExpr()
	}
	body := p.parseBlockStmt()
	return &ast.WhileStmt{Test: test, Body: body}
}

// ifStatement ::= 'if' expression blockStatement ('else if' expression blockStatement)* ('else' blockStatement)?;
func (p *Parser) parseIfStmt() ast.Stmt {
	p.consume(IF)
	test := p.parseExpr()
	consequent := p.parseBlockStmt()
	var alternate ast.Stmt
	if !p.isEnd() && p.current().Type == ELSE {
		p.consume(ELSE)
		if p.current().Type == IF {
			alternate = p.parseIfStmt()
		} else {
			alternate = p.parseBlockStmt()
		}
	}

	return &ast.IfStmt{Test: test, Consequent: consequent, Alternate: alternate}

}

// statement ::= expression | variableDeclarationStatement
// | variableAssignmentStatement | blockStatement
// | whileStatement | functionDeclaration | ifStatement | deferStatement | rangeStatement;
func (p *Parser) parseStmt() (ast.Stmt, error) {

	switch p.current().Type {
	case LET:
		return p.parseVarDecStmt(), nil
	case LBRACE:
		return p.parseBlockStmt(), nil
	case WHILE:
		return p.parseWhileStmt(), nil
	case FUNC:
		return p.parseFuncDecStmt(), nil
	case IF:
		return p.parseIfStmt(), nil
	case DEFER:
		return p.parseDeferStmt(), nil
	case FOR:
		return p.parseRangeStmt(), nil
	case IDENTIFIER:
		return p.parseVarAssignStmt(), nil
	default:
		ex := p.parseExpr()
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
