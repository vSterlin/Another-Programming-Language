package parser

import (
	"language/ast"
	"language/lexer"
	"testing"
)

func TestParseArrayExpr(t *testing.T) {
	l := lexer.NewLexer("[1, 2]")

	tokens, _ := l.GetTokens()

	p := NewParser(tokens)

	expr := p.parseArrayExpr()

	arrExpr, ok := expr.(*ast.ArrayExpr)

	if !ok {
		t.Errorf("Expected ArrayExpr, got: %T", expr)
	}

	if arrExpr.Elements[0].(*ast.NumberExpr).Val != 1 {
		t.Errorf("Expected 1, got: %d", arrExpr.Elements[0].(*ast.NumberExpr).Val)
	}

	if arrExpr.Elements[1].(*ast.NumberExpr).Val != 2 {
		t.Errorf("Expected 2, got: %d", arrExpr.Elements[1].(*ast.NumberExpr).Val)
	}
}

func TestParseArrayExprEmpty(t *testing.T) {
	l := lexer.NewLexer("[]")

	tokens, _ := l.GetTokens()

	p := NewParser(tokens)

	expr := p.parseArrayExpr()

	arrExpr, ok := expr.(*ast.ArrayExpr)

	if !ok {
		t.Errorf("Expected ArrayExpr, got: %T", expr)
	}

	arrExprLen := len(arrExpr.Elements)

	if arrExprLen != 0 {
		t.Errorf("Expected 0, got: %d", arrExprLen)
	}

}

func TestParseWhileStmt(t *testing.T) {
	l := lexer.NewLexer(`
		while (true) {}
		while {}
	`)

	tokens, _ := l.GetTokens()

	p := NewParser(tokens)

	stmt, err := p.parseWhileStmt()
	stmt2, err2 := p.parseWhileStmt()

	if err != nil || err2 != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	whileStmt, ok := stmt.(*ast.WhileStmt)
	whileStmt2, ok2 := stmt2.(*ast.WhileStmt)

	if !ok || !ok2 {
		t.Errorf("Expected WhileStmt, got: %T", stmt)
	}

	if whileStmt.Test.(*ast.BooleanExpr).Val != true {
		t.Errorf("Expected true, got: %t", whileStmt.Test.(*ast.BooleanExpr).Val)
	}

	if whileStmt2.Test.(*ast.BooleanExpr).Val != whileStmt.Test.(*ast.BooleanExpr).Val {
		t.Errorf("Expected true, got: %t", whileStmt2.Test.(*ast.BooleanExpr).Val)
	}

}

func TestConsume(t *testing.T) {
	l := lexer.NewLexer(`while`)
	tokens, _ := l.GetTokens()
	p := NewParser(tokens)
	err := p.consume(lexer.WHILE)

	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}
}

func TestConsumeError(t *testing.T) {
	l := lexer.NewLexer(`while`)
	tokens, _ := l.GetTokens()
	p := NewParser(tokens)
	err := p.consume(lexer.FOR)

	if err == nil {
		t.Error("Expected error but didn't get one")
	}
}
