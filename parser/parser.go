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
	curr := p.current()
	if curr.Type != RPAREN {
		//	TODO: return error
		fmt.Println("BOOO!")
	}
	// eat the RPAREN
	p.next()

	return val
}

// arrayExpression ::= '[' (expression (',' expression)*)? ']';
func (p *Parser) parseArrayExpr() ast.Expr {
	p.next() // eat the LBRACK
	exprs := []ast.Expr{}
	for p.pos < p.len && p.current().Type != RBRACK {
		expr := p.parseExpr()
		exprs = append(exprs, expr)
		if p.current().Type == COMMA {
			p.next()
		}

	}
	// eat the RBRACK
	p.next()
	return &ast.ArrayExpr{Elements: exprs}

}

// primaryExpression ::= identifier | number | boolean | string | '(' expression ')' | callExpression | arrayExpression;
func (p *Parser) parsePrimaryExpr() ast.Expr {

	switch p.current().Type {

	case IDENTIFIER:
		// if followed by LPAREN, then it's a call expression
		if p.pos+1 < p.len && p.tokens[p.pos+1].Type == LPAREN {
			return p.parseCallExpr()
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
		fmt.Println("BOOO!")
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
	// eat the RPAREN
	p.next()

	return &ast.CallExpr{
		Callee: calleeId.(*ast.IdentifierExpr),
		Args:   args,
	}

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

// expression ::= primaryExpression | additiveExpression;
func (p *Parser) parseExpr() ast.Expr {
	return p.parseAdditiveExpr()
}

// TODO: I will get rid of this
// variableDeclarationStatement ::= 'let' identifier '=' expression;
func (p *Parser) parseVarDecStmt() ast.Stmt {
	p.next()
	id := p.parseIdentifierExpr()

	if p.current().Type != ASSIGN {
		return nil
	}
	p.next()
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
	p.next() // eat the func
	id := p.parseIdentifierExpr()
	if p.current().Type != LPAREN {
		return nil
	}
	p.next() // eat the LPAREN
	var args []*ast.IdentifierExpr
	for !p.isEnd() && p.current().Type != RPAREN {
		arg := p.parseIdentifierExpr()
		args = append(args, arg.(*ast.IdentifierExpr))
		if p.current().Type == COMMA {
			p.next()
		}
	}
	// TODO: check if the next token is RPAREN
	p.next() // eat the RPAREN
	body := p.parseBlockStmt()
	return &ast.FuncDecStmt{Id: id.(*ast.IdentifierExpr), Args: args, Body: body.(*ast.BlockStmt)}
}

// blockStatement ::= '{' statement* '}';
func (p *Parser) parseBlockStmt() ast.Stmt {
	p.next() // eat the LBRACE
	var stmts []ast.Stmt
	for p.pos < p.len && p.current().Type != RBRACE {
		stmts = append(stmts, p.parseStmt())
	}
	// eat the RBRACE
	p.next()
	return &ast.BlockStmt{Stmts: stmts}
}

// whileStatement ::= 'while' [expression] blockStatement;
func (p *Parser) parseWhileStmt() ast.Stmt {
	p.next() // eat the WHILE
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
	p.next() // eat the IF
	test := p.parseExpr()
	consequent := p.parseBlockStmt()
	var alternate ast.Stmt
	if !p.isEnd() && p.current().Type == ELSE {
		p.next() // eat the ELSE
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
// | whileStatement | functionDeclaration | ifStatement;
func (p *Parser) parseStmt() ast.Stmt {

	switch p.current().Type {
	case LET:
		return p.parseVarDecStmt()
	case LBRACE:
		return p.parseBlockStmt()
	case WHILE:
		return p.parseWhileStmt()
	case FUNC:
		return p.parseFuncDecStmt()
	case IF:
		return p.parseIfStmt()
	case IDENTIFIER:
		return p.parseVarAssignStmt()
	default:
		ex := p.parseExpr()
		return &ast.ExprStmt{Expr: ex}
	}
}

// program ::= statement*;
func (p *Parser) ParseProgram() *ast.Program {
	var stmts []ast.Stmt
	for p.pos < p.len {
		stmts = append(stmts, p.parseStmt())
	}
	return &ast.Program{Stmts: stmts}
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
