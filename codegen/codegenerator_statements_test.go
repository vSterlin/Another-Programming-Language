package codegen

import (
	"language/ast"
	"language/lexer"
	"language/parser"
	"strings"
	"testing"
)

type tests []struct {
	srcCode  string
	expected string
}

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

	tests := tests{
		{
			srcCode:  "func foo() {}",
			expected: "void foo() {}",
		},
		{
			srcCode:  "func foo() int {}",
			expected: "int foo() {}",
		},
		{
			srcCode:  "func foo(a int) int {}",
			expected: "int foo(int a) {}",
		},
		{
			srcCode:  "func foo(a int, b string) int {}",
			expected: "int foo(int a, std::string b) {}",
		},
	}

	for _, test := range tests {
		cg := NewCodeGenerator()

		astNode := buildStmt(test.srcCode)
		code, err := cg.genStmt(astNode)
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

func TestBlockStmtCodegen(t *testing.T) {

	tests := tests{
		{
			srcCode:  "{}",
			expected: "{}",
		},
		{
			srcCode:  "{1}",
			expected: "{1;}",
		},
	}

	for _, test := range tests {
		cg := NewCodeGenerator()

		astNode := buildStmt(test.srcCode)
		code, err := cg.genStmt(astNode)
		if err != nil {
			t.Errorf("Error generating code: %s", err)
		}

		code = removeWhitespace(code)

		if code != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, code)
		}
	}
}

func TestVarDecStmtCodegen(t *testing.T) {

	tests := tests{
		{
			srcCode:  "a := 1",
			expected: "int a = 1;",
		},
		{
			srcCode:  "a := \"hello\"",
			expected: "std::string a = \"hello\";",
		},
		{
			srcCode:  "a := true",
			expected: "bool a = true;",
		},
		{
			srcCode:  "a := false",
			expected: "bool a = false;",
		},
	}

	for _, test := range tests {
		cg := NewCodeGenerator()

		astNode := buildStmt(test.srcCode)

		code, err := cg.genStmt(astNode)
		if err != nil {
			t.Errorf("Error generating code: %s", err)
		}

		if code != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, code)
		}

	}
}

func TestVarAssignStmtCodegen(t *testing.T) {

	tests := tests{
		{
			srcCode:  `a = 1`,
			expected: `a = 1;`,
		},
		{
			srcCode:  `a = "hello"`,
			expected: `a = "hello";`,
		},
		{
			srcCode:  `a = true`,
			expected: `a = true;`,
		},
		{
			srcCode:  `a = b`,
			expected: `a = b;`,
		},
	}

	for _, test := range tests {
		cg := NewCodeGenerator()

		astNode := buildStmt(test.srcCode)

		code, err := cg.genStmt(astNode)
		if err != nil {
			t.Errorf("Error generating code: %s", err)
		}

		if code != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, code)
		}

	}
}

func TestIfStmtCodegen(t *testing.T) {
	tests := tests{
		{
			srcCode:  "if (true) {}",
			expected: "if (true) {}",
		},
		{
			srcCode:  "if (true) {} else {}",
			expected: "if (true) {} else {}",
		},
		{
			srcCode:  "if (true) {1} else {2}",
			expected: "if (true) {1;} else {2;}",
		},
	}

	for _, test := range tests {
		cg := NewCodeGenerator()

		astNode := buildStmt(test.srcCode)
		code, err := cg.genStmt(astNode)
		if err != nil {
			t.Errorf("Error generating code: %s", err)
		}

		code = removeWhitespace(code)

		if code != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, code)
		}
	}
}

func TestWhileStmtCodegen(t *testing.T) {
	tests := tests{
		{
			srcCode:  "while (true) {}",
			expected: "while (true) {}",
		},
		{
			srcCode:  "while (true) {1}",
			expected: "while (true) {1;}",
		},
	}

	for _, test := range tests {
		cg := NewCodeGenerator()

		astNode := buildStmt(test.srcCode)
		code, err := cg.genStmt(astNode)
		if err != nil {
			t.Errorf("Error generating code: %s", err)
		}

		code = removeWhitespace(code)

		if code != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, code)
		}
	}

}

func TestReturnStmtCodegen(t *testing.T) {

	tests := tests{
		// FIXME: this test is failing
		{
			srcCode:  "return",
			expected: "return;",
		},
		{
			srcCode:  "return 1",
			expected: "return 1;",
		},
	}

	for _, test := range tests {
		cg := NewCodeGenerator()

		astNode := buildStmt(test.srcCode)
		code, err := cg.genStmt(astNode)
		if err != nil {
			t.Errorf("Error generating code: %s", err)
		}

		code = removeWhitespace(code)

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

// TODO: review usage of this
func removeWhitespace(s string) string {
	newS := strings.ReplaceAll(s, "\n", "")
	newS = strings.ReplaceAll(newS, "\t", "")
	return newS
}
