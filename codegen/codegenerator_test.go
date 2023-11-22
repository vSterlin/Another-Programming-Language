package codegen

import (
	"language/ast"
	"testing"
)

func TestLiteralsCodegen(t *testing.T) {
	tests := []struct {
		astNode  ast.Expr
		expected string
	}{
		{
			astNode:  &ast.NumberExpr{Val: 1},
			expected: "1",
		},
		{
			astNode:  &ast.StringExpr{Val: "hello"},
			expected: "\"hello\"",
		},
		{
			astNode:  &ast.BooleanExpr{Val: true},
			expected: "true",
		},
		{
			astNode:  &ast.BooleanExpr{Val: false},
			expected: "false",
		},
	}

	for _, test := range tests {
		cg := NewCodeGenerator()
		code, err := cg.genExpr(test.astNode)
		if err != nil {
			t.Errorf("Error generating code: %s", err)
		}

		if code != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, code)
		}
	}

}
