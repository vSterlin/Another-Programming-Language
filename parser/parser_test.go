package parser

import (
	"language/ast"
	"language/lexer"
	"testing"
)

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
