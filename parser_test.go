package main

import (
	"language/ast"
	"testing"
)

func TestParseVarDecStmt(t *testing.T) {
	l := NewLexer("let x = 1")

	tokens, _ := l.GetTokens()

	p := NewParser(tokens)

	stmt := p.parseVarDecStmt()

	varDecStmt, ok := stmt.(*ast.VarDecStmt)

	if !ok {
		t.Errorf("Expected VarDecStmt, got: %T", stmt)
	}

	if varDecStmt.Id.Name != "x" {
		t.Errorf("Expected x, got: %s", varDecStmt.Id.Name)
	}

	if varDecStmt.Init.(*ast.NumberExpr).Val != 1 {
		t.Errorf("Expected 1, got: %d", varDecStmt.Init.(*ast.NumberExpr).Val)
	}

}
