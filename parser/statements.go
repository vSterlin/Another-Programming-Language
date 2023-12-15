package parser

import (
	"language/ast"
	. "language/lexer"
)

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

	return &ast.RangeStmt{Id: id, Expr: ex, Body: body}, nil
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

	return &ast.VarAssignStmt{Id: id.(*ast.IdentifierExpr), Init: ex, Op: assignOp}, nil
}

// functionDeclaration ::= 'func' identifier '(' (identifier (',' identifier)*)? ')' blockStatement;
func (p *Parser) parseFuncDecStmt(funcType string) (ast.Stmt, error) {

	// to handle methods
	if funcType == "func" {
		if err := p.consume(FUNC); err != nil {
			return nil, err
		}
	}

	id, err := p.parseIdentifierExpr()
	if err != nil {
		return nil, err
	}

	if err := p.consume(LPAREN); err != nil {
		return nil, err
	}

	var params []*ast.Param
	for !p.isEnd() && p.current().Type != RPAREN {
		paramName, err := p.parseIdentifierExpr()
		if err != nil {
			return nil, err
		}
		paramType, err := p.parseTypeExpr()

		param := &ast.Param{Id: paramName, Type: paramType}

		params = append(params, param)
		if p.current().Type == COMMA {
			p.next()
		}
	}

	if err := p.consume(RPAREN); err != nil {
		return nil, err
	}

	var retType *ast.TypeExpr = &ast.TypeExpr{Type: &ast.IdentifierExpr{Name: "void"}}

	if p.current().Type != LBRACE {
		retType, err = p.parseTypeExpr()
	}

	body, err := p.parseBlockStmt()
	if err != nil {
		return nil, err
	}
	return &ast.FuncDecStmt{Id: id, Args: params, Body: body, ReturnType: retType}, nil
}

// blockStatement ::= '{' statement* '}';
func (p *Parser) parseBlockStmt() (*ast.BlockStmt, error) {
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

// returnStatement ::= 'return' [expression];
func (p *Parser) parseReturnStmt() (ast.Stmt, error) {
	if err := p.consume(RETURN); err != nil {
		return nil, err
	}
	arg, err := p.parseExpr()
	if err != nil {
		return &ast.ReturnStmt{}, nil

	}
	return &ast.ReturnStmt{Arg: arg}, nil
}

// classDeclaration ::= 'class' identifier '{' (functionDeclaration)* '}';
func (p *Parser) parseClassDecStmt() (ast.Stmt, error) {
	if err := p.consume(CLASS); err != nil {
		return nil, err
	}
	id, err := p.parseIdentifierExpr()
	if err != nil {
		return nil, err
	}
	if err := p.consume(LBRACE); err != nil {
		return nil, err
	}
	var methods []*ast.FuncDecStmt = []*ast.FuncDecStmt{}
	for !p.isEnd() && p.current().Type != RBRACE {
		method, err := p.parseFuncDecStmt("method")
		if err != nil {
			return nil, err
		}
		methods = append(methods, method.(*ast.FuncDecStmt))
	}
	if err := p.consume(RBRACE); err != nil {
		return nil, err
	}
	return &ast.ClassDecStmt{Id: id, Methods: methods}, nil
}

func (p *Parser) parseTypeAliasStmt() (ast.Stmt, error) {
	if err := p.consume(TYPE); err != nil {
		return nil, err
	}
	id, err := p.parseIdentifierExpr()
	if err != nil {
		return nil, err
	}

	typ, err := p.parseTypeExpr()
	if err != nil {
		return nil, err
	}
	return &ast.TypeAliasStmt{Id: id, Type: typ}, nil
}

// statement ::= expression | variableDeclarationStatement
// | variableAssignmentStatement | blockStatement
// | whileStatement | functionDeclaration
// | ifStatement | deferStatement | rangeStatement | returnStatement;
func (p *Parser) parseStmt() (ast.Stmt, error) {

	switch p.current().Type {

	case LBRACE:
		return p.parseBlockStmt()
	case WHILE:
		return p.parseWhileStmt()
	case FUNC:
		return p.parseFuncDecStmt("func")
	case IF:
		return p.parseIfStmt()
	case DEFER:
		return p.parseDeferStmt()
	case FOR:
		return p.parseRangeStmt()
	case IDENTIFIER, THIS:
		return p.parseVarAssignStmt()
	case RETURN:
		return p.parseReturnStmt()
	case CLASS:
		return p.parseClassDecStmt()
	case TYPE:
		return p.parseTypeAliasStmt()
	default:
		ex, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &ast.ExprStmt{Expr: ex}, nil
	}
}
