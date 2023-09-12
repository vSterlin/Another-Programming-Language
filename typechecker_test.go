package main

import (
	"testing"
)

func TestTypecheck(t *testing.T) {

	tc := NewTypechecker(&Program{})

	tests := []struct {
		expr     Expr
		expected Type
	}{
		{expr: &NumberExpr{Val: 1}, expected: NumberType},
		{expr: &BooleanExpr{Val: true}, expected: BooleanType},
	}
	for _, i := range tests {
		if tc.typecheckExpr(i.expr) != i.expected {
			t.Errorf("Expected %s, got: %s", i.expected, tc.typecheckExpr(i.expr))
		}
	}

}

func TestExpectTypeEqual(t *testing.T) {

	tc := NewTypechecker(&Program{})
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
