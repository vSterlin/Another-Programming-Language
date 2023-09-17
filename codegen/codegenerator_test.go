package codegen

import (
	"language/ast"
	"language/lexer"
	"language/parser"
	"testing"
)

var tests = []struct {
	srcCode  string
	expected string
}{
	{srcCode: "let a = 1", expected: "var a = 1"},
	{srcCode: "let a = 1 + 2", expected: "var a = (1 + 2)"},
	{srcCode: "let a = 1 + 2 * 3", expected: "var a = (1 + (2 * 3))"},
	{srcCode: "(1 + 2) * 3", expected: "((1 + 2) * 3)"},
	{srcCode: "\"test\"", expected: "\"test\""},
	{srcCode: "true", expected: "true"},
}

func TestVarDec(t *testing.T) {

	for _, i := range tests {

		prog := codeToAst(i.srcCode)

		cg := NewJavascriptCodeGenerator()
		generatedCode := cg.Generate(prog)

		if generatedCode != i.expected {
			t.Errorf("Expected %s, got %s", i.expected, generatedCode)
		}
	}
}

func codeToAst(code string) *ast.Program {
	l := lexer.NewLexer(code)
	tokens, _ := l.GetTokens()
	p := parser.NewParser(tokens)
	prog, _ := p.ParseProgram()
	return prog
}
