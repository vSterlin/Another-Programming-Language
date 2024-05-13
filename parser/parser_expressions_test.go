package parser

import (
	"language/ast"
	"language/lexer"
	"testing"
)

func TestParseNumberExpr(t *testing.T) {
	p := NewParser(getTokens("1"))
	expr, err := p.parseNumberExpr()
	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	e, ok := expr.(*ast.NumberExpr)

	if !ok {
		t.Errorf("Expected NumberExpr, got: %T", expr)
	}

	if e.Val != 1 {
		t.Errorf("Expected 1, got: %d", e.Val)
	}
}

func TestParseBooleanExpr(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		p := NewParser(getTokens(tt.input))
		expr, err := p.parseBooleanExpr()
		if err != nil {
			t.Errorf("Expected no error, got: %s", err)
		}

		e, ok := expr.(*ast.BooleanExpr)

		if !ok {
			t.Errorf("Expected BooleanExpr, got: %T", expr)
		}

		if e.Val != tt.want {
			t.Errorf("Expected %t, got: %t", tt.want, e.Val)
		}

	}
}

func TestParseStringExpr(t *testing.T) {
	t.Error("Not implemented")
}

func TestParseIdentifierExpr(t *testing.T) {
	t.Error("Not implemented")
}

func TestParseParenExpr(t *testing.T) {
	t.Error("Not implemented")
}

func TestParseArrowFunc(t *testing.T) {
	t.Error("Not implemented")
}

func TestParsePrimaryExpr(t *testing.T) {
	t.Error("Not implemented")
}

func TestParseUpdateExpr(t *testing.T) {
	t.Error("Not implemented")
}

func TestParseUnaryExpr(t *testing.T) {
	t.Error("Not implemented")
}

func TestParseCallExpr(t *testing.T) {
	t.Error("Not implemented")
}

func TestParseEqualityExpr(t *testing.T) {
	t.Error("Not implemented")
}

func TestParseRelationalExpr(t *testing.T) {
	t.Error("Not implemented")
}

func TestParseAdditiveExpr(t *testing.T) {
	t.Error("Not implemented")
}

func TestParseMultiplicativeExpr(t *testing.T) {
	t.Error("Not implemented")
}

func TestParseAndExpr(t *testing.T) {
	t.Error("Not implemented")
}

func TestParseOrExpr(t *testing.T) {
	t.Error("Not implemented")
}

func TestParseTypeExpr(t *testing.T) {
	t.Error("Not implemented")
}

func getTokens(code string) []*lexer.Token {
	l := lexer.NewLexer(code)
	tokens, _ := l.GetTokens()
	return tokens
}
