package codegen

import (
	"language/ast"
	"language/lexer"
	"language/parser"
	"strings"
	"testing"
)

func TestExprStmtCodegen(t *testing.T) {
	tests := []struct {
		astNode  ast.Stmt
		expected string
	}{
		{
			astNode:  buildStmt("1"),
			expected: "1;",
		},
		{
			astNode:  buildStmt("\"hello\""),
			expected: "\"hello\";",
		},
		{
			astNode:  buildStmt("true"),
			expected: "true;",
		},
		{
			astNode:  buildStmt("false"),
			expected: "false;",
		},
	}

	for _, test := range tests {
		cg := NewCodeGenerator()
		code, err := cg.genStmt(test.astNode)
		if err != nil {
			t.Errorf("Error generating code: %s", err)
		}

		if code != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, code)
		}
	}
}

func TestFuncDecStmtCodegen(t *testing.T) {

	tests := []struct {
		astNode  ast.Stmt
		expected string
	}{
		{
			astNode:  buildStmt("func foo() {}"),
			expected: "void foo() {}",
		},
		{
			astNode:  buildStmt("func foo() int {}"),
			expected: "int foo() {}",
		},
		{
			astNode:  buildStmt("func foo(a int) int {}"),
			expected: "int foo(int a) {}",
		},
		{
			astNode:  buildStmt("func foo(a int, b string) int {}"),
			expected: "int foo(int a, std::string b) {}",
		},
	}

	for _, test := range tests {
		cg := NewCodeGenerator()
		code, err := cg.genStmt(test.astNode)
		if err != nil {
			t.Errorf("Error generating code: %s", err)
		}

		// TODO: review
		code = strings.ReplaceAll(code, "\n", "")

		if code != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, code)
		}
	}

}

// helpers
func buildStmt(code string) ast.Stmt {
	l := lexer.NewLexer(code)
	tokens, _ := l.GetTokens()
	p := parser.NewParser(tokens)
	prog, _ := p.ParseProgram()
	return prog.Stmts[0]
}
