package typechecker

import (
	"fmt"
	"language/ast"
	"language/lexer"
	"language/parser"
	"testing"
)

func TestCheckExpr(t *testing.T) {

	tests := []struct {
		expr     ast.Expr
		expected Type
	}{
		{expr: buildExpr("1"), expected: Number},
		{expr: buildExpr("\"hello\""), expected: String},
		{expr: buildExpr("true"), expected: Boolean},
		{expr: buildExpr("1 + 1"), expected: Number},
		{expr: buildExpr("true && false"), expected: Boolean},
		{expr: buildExpr("() => {}"), expected: FuncType{Args: []Type{}, ReturnType: Void}},
		{expr: buildExpr("() number => { return 1 }"), expected: FuncType{Args: []Type{}, ReturnType: Number}},
		{
			expr:     buildExpr("(a number, b number) number => { return 1 }"),
			expected: FuncType{Args: []Type{Number, Number}, ReturnType: Number},
		},
	}

	tc := NewTypeChecker()

	for _, i := range tests {
		typ, err := tc.checkExpr(i.expr)

		if err != nil {
			t.Errorf("Expected no error, got: %s", err)
		}

		if !typ.Equals(i.expected) {
			t.Errorf("Expected %s, got: %s", i.expected, typ)
		}

	}
}

func TestVarDecStmt(t *testing.T) {

	definedVarProg := buildProgram(`
		a := 1
		a
	`)

	tc := NewTypeChecker()

	err := tc.Check(definedVarProg)

	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	undefinedVarProg := buildProgram("iDontExist")

	err = tc.Check(undefinedVarProg)

	if err == nil {
		t.Errorf("Expected error, got none")
	}

}

func TestArrowFuncScope(t *testing.T) {

	prog := buildProgram(`
		a := 1
		() => { a }
	`)

	tc := NewTypeChecker()

	err := tc.Check(prog)

	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

}

func TestArrowFuncScope2(t *testing.T) {

	prog := buildProgram(`
	func makeCounter() () => number {
		count := 0
		return () => {
			count = count + 1
			return count
		}
	}
	`)

	tc := NewTypeChecker()

	err := tc.Check(prog)

	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

}

// helpers
func buildProgram(code string) *ast.Program {
	tokens, _ := lexer.NewLexer(code).GetTokens()
	p := parser.NewParser(tokens)
	prog, err := p.ParseProgram()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return prog
}

func buildExpr(code string) ast.Expr {
	prog := buildProgram(code)
	return prog.Stmts[0].(*ast.ExprStmt).Expr
}
