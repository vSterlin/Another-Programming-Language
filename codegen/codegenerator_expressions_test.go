package codegen

import (
	"language/ast"
	"language/lexer"
	"language/parser"
	"strings"
	"testing"
)

// Expression codegen tests

func TestLiteralsCodegen(t *testing.T) {
	tests := []struct {
		srcCode  string
		expected string
	}{
		{
			srcCode:  "1",
			expected: "1",
		},
		{
			srcCode:  "\"hello\"",
			expected: "\"hello\"",
		},
		{
			srcCode:  "true",
			expected: "true",
		},
		{
			srcCode:  "false",
			expected: "false",
		},
	}

	for _, test := range tests {
		cg := NewCodeGenerator()

		astNode := buildExpr(test.srcCode)
		code, err := cg.genExpr(astNode)
		if err != nil {
			t.Errorf("Error generating code: %s", err)
		}

		if code != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, code)
		}
	}

}

func TestBinaryExprCodegen(t *testing.T) {
	tests := []struct {
		srcCode  string
		expected string
	}{
		{
			srcCode:  "1 + 1",
			expected: "1 + 1",
		},
		{
			srcCode:  "1 - 1",
			expected: "1 - 1",
		},
		{
			srcCode:  "1 * 1",
			expected: "1 * 1",
		},
		{
			srcCode:  "1 / 1",
			expected: "1 / 1",
		},
		{
			srcCode:  "\"hello\" + \"world\"",
			expected: "\"hello\" + \"world\"",
		},
		{
			srcCode:  "1 == 1",
			expected: "1 == 1",
		},
		{
			srcCode:  "1 != 1",
			expected: "1 != 1",
		},
		{
			srcCode:  "1 > 1",
			expected: "1 > 1",
		},
		{
			srcCode:  "1 < 1",
			expected: "1 < 1",
		},
		{
			srcCode:  "1 >= 1",
			expected: "1 >= 1",
		},
		{
			srcCode:  "1 <= 1",
			expected: "1 <= 1",
		},
	}

	for _, test := range tests {
		cg := NewCodeGenerator()

		astNode := buildExpr(test.srcCode)
		code, err := cg.genExpr(astNode)
		if err != nil {
			t.Errorf("Error generating code: %s", err)
		}

		if code != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, code)
		}
	}

}

func TestLogicalExprCodegen(t *testing.T) {
	tests := []struct {
		srcCode  string
		expected string
	}{
		{
			srcCode:  "true && true",
			expected: "true && true",
		},
		{
			srcCode:  "true || false",
			expected: "true || false",
		},
	}

	for _, test := range tests {
		cg := NewCodeGenerator()

		astNode := buildExpr(test.srcCode)
		code, err := cg.genExpr(astNode)
		if err != nil {
			t.Errorf("Error generating code: %s", err)
		}

		if code != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, code)
		}
	}

}

func TestCallExprCodegen(t *testing.T) {
	tests := []struct {
		srcCode  string
		expected string
	}{
		{
			srcCode:  "foo()",
			expected: "foo()",
		},
		{
			srcCode:  "foo(1)",
			expected: "foo(1)",
		},
		{
			srcCode:  "foo(1, 2)",
			expected: "foo(1, 2)",
		},
		{
			srcCode:  "foo(\"bar\")",
			expected: "foo(\"bar\")",
		},
	}

	for _, test := range tests {
		cg := NewCodeGenerator()

		astNode := buildExpr(test.srcCode)
		code, err := cg.genExpr(astNode)
		if err != nil {
			t.Errorf("Error generating code: %s", err)
		}

		if code != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, code)
		}
	}
}

func TestIdentifierExprCodegen(t *testing.T) {
	tests := []struct {
		astNode  ast.Expr
		expected string
	}{
		{
			astNode:  &ast.IdentifierExpr{Name: "foo"},
			expected: "foo",
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

func TestTypeExprCodegen(t *testing.T) {
	tests := []struct {
		astNode  ast.Expr
		expected string
	}{
		{
			astNode:  &ast.TypeExpr{Type: &ast.IdentifierExpr{Name: "int"}},
			expected: "int",
		},
		{
			astNode:  &ast.TypeExpr{Type: &ast.IdentifierExpr{Name: "string"}},
			expected: "std::string",
		},
		{
			astNode:  &ast.TypeExpr{Type: &ast.IdentifierExpr{Name: "boolean"}},
			expected: "bool",
		},
		{
			astNode:  &ast.TypeExpr{Type: &ast.IdentifierExpr{Name: "void"}},
			expected: "void",
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

func TestArrowFuncCodegen(t *testing.T) {
	tests := []struct {
		srcCode  string
		expected string
	}{
		{
			srcCode:  "() => {}",
			expected: "[=]() mutable {}",
		},
		{
			srcCode:  "(a int) => {}",
			expected: "[=](int a) mutable {}",
		},
		{
			srcCode:  "(a int, b string) => {}",
			expected: "[=](int a, std::string b) mutable {}",
		},
	}

	for _, test := range tests {
		cg := NewCodeGenerator()

		astNode := buildExpr(test.srcCode)
		code, err := cg.genExpr(astNode)
		if err != nil {
			t.Errorf("Error generating code: %s", err)
		}

		// TODO: review this
		code = strings.ReplaceAll(code, "\n", "")
		if code != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, code)
		}
	}
}

func TestUnaryExprCodegen(t *testing.T) {
	tests := []struct {
		srcCode  string
		expected string
	}{
		{
			srcCode:  "!true",
			expected: "!true",
		},
		{
			srcCode:  "!false",
			expected: "!false",
		},
		{
			srcCode:  "!!false",
			expected: "!!false",
		},
	}

	for _, test := range tests {
		cg := NewCodeGenerator()

		astNode := buildExpr(test.srcCode)
		code, err := cg.genExpr(astNode)
		if err != nil {
			t.Errorf("Error generating code: %s", err)
		}

		if code != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, code)
		}

	}
}

func TestGroupingExprCodegen(t *testing.T) {
	tests := []struct {
		srcCode  string
		expected string
	}{
		{
			srcCode:  "i++",
			expected: "i++",
		},
		{
			srcCode:  "i--",
			expected: "i--",
		},
	}

	for _, test := range tests {
		cg := NewCodeGenerator()

		astNode := buildExpr(test.srcCode)
		code, err := cg.genExpr(astNode)
		if err != nil {
			t.Errorf("Error generating code: %s", err)
		}

		if code != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, code)
		}

	}
}

// helpers
func buildExpr(code string) ast.Expr {
	l := lexer.NewLexer(code)
	tokens, _ := l.GetTokens()
	p := parser.NewParser(tokens)
	prog, _ := p.ParseProgram()
	return prog.Stmts[0].(*ast.ExprStmt).Expr
}
