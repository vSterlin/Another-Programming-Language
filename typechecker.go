package main

import "fmt"

type Typechecker struct {
	program *Program
}

func NewTypechecker(program *Program) *Typechecker {
	return &Typechecker{
		program: program,
	}
}

type Type int

const (
	NumberType Type = iota
	BooleanType

	UndefinedType
)

func (t Type) String() string {
	switch t {
	case NumberType:
		return "number"
	case BooleanType:
		return "boolean"
	default:
		return "undefined"
	}
}

func (t *Typechecker) typecheckBinaryExpr(ex *BinaryExpr) Type {
	lhsType := t.typecheckExpr(ex.Lhs)
	rhsType := t.typecheckExpr(ex.Rhs)

	if !t.expectTypeEqual(NumberType, lhsType, rhsType) {
		return UndefinedType
	}

	switch ex.Op {
	case "+", "-", "*", "/":
		return NumberType
	// case "==", "!=", "<", "<=", ">", ">=":
	// 	return &BooleanType{}
	default:
		return UndefinedType
	}
}

func (t *Typechecker) typecheckExpr(ex Expr) Type {
	switch ex := ex.(type) {
	case *NumberExpr:
		return NumberType
	case *BooleanExpr:
		return BooleanType
	case *BinaryExpr:
		return t.typecheckBinaryExpr(ex)
	default:
		return UndefinedType
	}
}

func (t *Typechecker) typecheckStmt(stmt Stmt) Type {

	switch stmt := stmt.(type) {
	case *ExprStmt:
		return t.typecheckExpr(stmt.Expr)
	default:
		return UndefinedType
	}
}

func (t *Typechecker) typecheckProgram(program *Program) {

	for _, stmt := range program.Stmts {
		fmt.Println(t.typecheckStmt(stmt))
	}

}

// Helpers
func (t *Typechecker) expectTypeEqual(expected Type, actual ...Type) bool {
	for _, a := range actual {
		if expected != a {
			return false
		}
	}
	return true
}
