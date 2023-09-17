package main

import (
	"language/ast"
	"language/lexer"
	"language/parser"
	"testing"
)

func TestTypecheck(t *testing.T) {

	tc := NewTypechecker(&ast.Program{})

	tests := []struct {
		expr     ast.Expr
		expected Type
	}{
		{expr: &ast.NumberExpr{Val: 1}, expected: NumberType},
		{expr: &ast.BooleanExpr{Val: true}, expected: BooleanType},
		{expr: &ast.BinaryExpr{Op: "+", Lhs: &ast.NumberExpr{Val: 1}, Rhs: &ast.NumberExpr{Val: 2}}, expected: NumberType},
	}
	for _, i := range tests {
		if tc.typeofExpr(i.expr, nil) != i.expected {
			t.Errorf("Expected %s, got: %s", i.expected, tc.typeofExpr(i.expr, nil))
		}
	}

}

func TestExpectTypeEqual(t *testing.T) {

	tc := NewTypechecker(&ast.Program{})
	tests := []struct {
		types    []Type
		expected bool
	}{
		{types: []Type{NumberType, NumberType}, expected: true},
		{types: []Type{BooleanType, BooleanType}, expected: true},
		{types: []Type{NumberType, BooleanType}, expected: false},
		{types: []Type{NumberType, NumberType, NumberType}, expected: true},
	}

	for _, i := range tests {
		if tc.expectTypeEqual(i.types[0], i.types[1:]...) != i.expected {
			t.Errorf("Expected %t, got: %t", i.expected, tc.expectTypeEqual(i.types[0], i.types[1:]...))
		}
	}

}

// func TestVarDec(t *testing.T) {

// 	tokens, _ := NewLexer("let x = 1").GetTokens()

// 	prog := NewParser(tokens).parseProgram()

// 	tc := NewTypechecker(prog)

// }

func TestGlobalVar(t *testing.T) {

	prog := buildProgram("let x = 1")

	tc := NewTypechecker(prog)

	expectedType := StringType

	tc.typeofProgram(prog)

	actualType := tc.typeofVar(&ast.IdentifierExpr{Name: "VERSION"}, tc.globalTypeEnv)

	if expectedType != actualType {
		t.Errorf("Expected %s, got: %s", expectedType, actualType)
	}
}

func TestVar(t *testing.T) {

	prog := buildProgram("let x = 1")

	tc := NewTypechecker(prog)

	expectedType := NumberType

	typeEnv := NewTypeEnv(nil, nil)

	typeEnv.DefineVar("x", NumberType)

	actualType := tc.typeofVar(&ast.IdentifierExpr{Name: "x"}, typeEnv)

	if expectedType != actualType {
		t.Errorf("Expected %s, got: %s", expectedType, actualType)
	}

	nonExistentType := tc.typeofVar(&ast.IdentifierExpr{Name: "y"}, typeEnv)

	if nonExistentType != UndefinedType {
		t.Errorf("Expected %s, got: %s", UndefinedType, nonExistentType)
	}
}

func TestVarDec(t *testing.T) {

	prog := buildProgram("let x = 1")

	tc := NewTypechecker(prog)

	expectedType := NumberType

	tc.typeofProgram(prog)

	actualType := tc.typeofVar(&ast.IdentifierExpr{Name: "x"}, tc.globalTypeEnv)

	if expectedType != actualType {
		t.Errorf("Expected %s, got: %s", expectedType, actualType)
	}
}

func buildProgram(code string) *ast.Program {
	tokens, _ := lexer.NewLexer(code).GetTokens()
	prog, _ := parser.NewParser(tokens).ParseProgram()
	return prog
}
