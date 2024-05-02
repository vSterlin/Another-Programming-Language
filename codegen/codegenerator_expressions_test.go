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

func TestBinaryExprCodegen(t *testing.T) {
	tests := []struct {
		astNode  ast.Expr
		expected string
	}{
		{
			astNode:  buildExpr("1 + 1"),
			expected: "1 + 1",
		},
		{
			astNode:  buildExpr("\"hello\" + \"world\""),
			expected: "\"hello\" + \"world\"",
		},
		{
			astNode:  buildExpr("1 == 1"),
			expected: "1 == 1",
		},
		{
			astNode:  buildExpr("1 != 1"),
			expected: "1 != 1",
		},
		{
			astNode:  buildExpr("1 > 1"),
			expected: "1 > 1",
		},
		{
			astNode:  buildExpr("1 < 1"),
			expected: "1 < 1",
		},
		{
			astNode:  buildExpr("1 >= 1"),
			expected: "1 >= 1",
		},
		{
			astNode:  buildExpr("1 <= 1"),
			expected: "1 <= 1",
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

func TestLogicalExprCodegen(t *testing.T) {
	tests := []struct {
		astNode  ast.Expr
		expected string
	}{
		{
			astNode:  buildExpr("true && true"),
			expected: "true && true",
		},
		{
			astNode:  buildExpr("true && false"),
			expected: "true || false",
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

func TestCallExprCodegen(t *testing.T) {
	tests := []struct {
		astNode  ast.Expr
		expected string
	}{
		{
			astNode:  buildExpr("foo()"),
			expected: "foo()",
		},
		{
			astNode:  buildExpr("foo(1)"),
			expected: "foo(1)",
		},
		{
			astNode:  buildExpr("foo(1, 2)"),
			expected: "foo(1, 2)",
		},

		{
			astNode:  buildExpr("foo(\"bar\")"),
			expected: "foo(\"bar\")",
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
		astNode  ast.Expr
		expected string
	}{
		{
			astNode:  buildExpr("() => {}"),
			expected: "[=]() mutable {}",
		},
		{
			astNode:  buildExpr("(a int) => {}"),
			expected: "[=](int a) mutable {}",
		},
		{
			astNode:  buildExpr("(a int, b string) => {}"),
			expected: "[=](int a, std::string b) mutable {}",
		},
	}

	for _, test := range tests {
		cg := NewCodeGenerator()
		code, err := cg.genExpr(test.astNode)
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
		astNode  ast.Expr
		expected string
	}{
		{
			astNode:  buildExpr("!true"),
			expected: "!true",
		},
		{
			astNode:  buildExpr("!false"),
			expected: "!false",
		},
		{
			astNode:  buildExpr("!!false"),
			expected: "!!false",
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

func TestGroupingExprCodegen(t *testing.T) {
	tests := []struct {
		astNode  ast.Expr
		expected string
	}{
		{
			astNode:  buildExpr("i++"),
			expected: "i++",
		},
		{
			astNode:  buildExpr("i--"),
			expected: "i--",
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

// Statement codegen tests

// helpers
func buildExpr(code string) ast.Expr {
	l := lexer.NewLexer(code)
	tokens, _ := l.GetTokens()
	p := parser.NewParser(tokens)
	prog, _ := p.ParseProgram()
	return prog.Stmts[0].(*ast.ExprStmt).Expr
}
