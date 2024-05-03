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
		srcCode  string
		expected string
	}{
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

	tests := []struct {
		srcCode  string
		expected string
	}{
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

		// TODO: review
		code = strings.ReplaceAll(code, " ", "")
		code = strings.ReplaceAll(code, "\n", "")
		code = strings.ReplaceAll(code, "\t", "")

		if code != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, code)
		}
	}
}

func TestVarAssignStmtCodegen(t *testing.T) {

	tests := []struct {
		srcCode  string
		expected string
	}{
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

/*


func (cg *CodeGenerator) genVarAssignStmt(stmt *ast.VarAssignStmt) (string, error) {

	id := stmt.Id.Name

	init, err := cg.genExpr(stmt.Init)
	if err != nil {
		return "", err
	}
	if stmt.Op == ":=" {

		varType := inferFromAstNode(stmt.Init)

		fmt.Printf("varType: %#v\n", varType)

		return fmt.Sprintf("%s %s = %s;", varType, id, init), nil
	} else {
		return fmt.Sprintf("%s = %s;", id, init), nil
	}

}

func (cg *CodeGenerator) genIfStmt(stmt *ast.IfStmt) (string, error) {

	test, err := cg.genExpr(stmt.Test)

	if err != nil {
		return "", err
	}

	body, err := cg.genStmt(stmt.Consequent)

	if err != nil {
		return "", err
	}

	if stmt.Alternate != nil {
		alternate, err := cg.genStmt(stmt.Alternate)

		if err != nil {
			return "", err
		}

		return fmt.Sprintf("if (%s) %s else %s", test, body, alternate), nil
	}

	return fmt.Sprintf("if (%s) %s", test, body), nil

}

func (cg *CodeGenerator) genWhileStmt(stmt *ast.WhileStmt) (string, error) {

	test, err := cg.genExpr(stmt.Test)
	if err != nil {
		return "", err
	}
	body, err := cg.genStmt(stmt.Body)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("while (%s) %s", test, body), nil

}

func (cg *CodeGenerator) genReturnStmt(stmt *ast.ReturnStmt) (string, error) {

	returnedVal, err := cg.genExpr(stmt.Arg)

	if err != nil {
		return "return;", nil
	}
	return fmt.Sprintf("return %s;", returnedVal), nil
}

*/

// helpers
func buildStmt(code string) ast.Stmt {
	l := lexer.NewLexer(code)
	tokens, _ := l.GetTokens()
	p := parser.NewParser(tokens)
	prog, _ := p.ParseProgram()
	return prog.Stmts[0]
}
