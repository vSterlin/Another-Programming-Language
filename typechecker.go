package main

import (
	"language/ast"
)

type Typechecker struct {
	program       *ast.Program
	globalTypeEnv *TypeEnv
}

func NewTypechecker(program *ast.Program) *Typechecker {

	global := NewTypeEnv(nil, TypeMap{
		"VERSION": StringType,
	})
	return &Typechecker{
		program:       program,
		globalTypeEnv: global,
	}
}

type Type int

const (
	// Undefined has to be 0th cause maps assign 0 by default
	// if you look up non existent key
	// alternatively we can check if the key exists
	UndefinedType Type = iota

	BooleanType
	NumberType
	StringType
)

func (t Type) String() string {

	strMap := map[Type]string{
		UndefinedType: "undefined",
		BooleanType:   "boolean",
		NumberType:    "number",
	}
	str, ok := strMap[t]
	if !ok {
		return "undefined"
	}
	return str
}

func (tc *Typechecker) typeofBinaryExpr(ex *ast.BinaryExpr) Type {
	lhsType := tc.typeofExpr(ex.Lhs, nil)
	rhsType := tc.typeofExpr(ex.Rhs, nil)

	if !tc.expectTypeEqual(NumberType, lhsType, rhsType) {
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

func (tc *Typechecker) typeofExpr(ex ast.Expr, typeEnv *TypeEnv) Type {
	switch ex := ex.(type) {
	case *ast.NumberExpr:
		return NumberType
	case *ast.BooleanExpr:
		return BooleanType
	case *ast.BinaryExpr:
		return tc.typeofBinaryExpr(ex)
	case *ast.IdentifierExpr:
		return tc.typeofVar(ex, typeEnv)
	default:
		return UndefinedType
	}
}

// TODO: error handling
func (tc *Typechecker) typeofVar(id *ast.IdentifierExpr, typeEnv *TypeEnv) Type {
	varType := typeEnv.LookupVar(id.Name)
	return varType
}

func (tc *Typechecker) typeofVarDec(stmt *ast.VarDecStmt, typeEnv *TypeEnv) Type {
	valueType := tc.typeofExpr(stmt.Init, typeEnv)
	typeEnv.DefineVar(stmt.Id.Name, valueType)
	return valueType
}

func (tc *Typechecker) typeofStmt(stmt ast.Stmt, typeEnv *TypeEnv) Type {

	switch stmt := stmt.(type) {

	// case *ast.VarDecStmt:
	// 	return tc.typeofVarDec(stmt, typeEnv)
	case *ast.ExprStmt:
		return tc.typeofExpr(stmt.Expr, typeEnv)
	default:
		return UndefinedType
	}
}

func (tc *Typechecker) typeofProgram(program *ast.Program) {

	for _, stmt := range program.Stmts {
		tc.typeofStmt(stmt, tc.globalTypeEnv)
	}

}

// Helpers
func (tc *Typechecker) expectTypeEqual(expected Type, actual ...Type) bool {
	for _, a := range actual {
		if expected != a {
			return false
		}
	}
	return true
}

// Type Env

type TypeMap map[string]Type
type TypeEnv struct {
	env    TypeMap
	parent *TypeEnv
}

func NewTypeEnv(parent *TypeEnv, env TypeMap) *TypeEnv {
	if env == nil {
		env = make(TypeMap)
	}
	return &TypeEnv{
		env:    env,
		parent: parent,
	}
}

func (te *TypeEnv) DefineVar(name string, t Type) {
	(te).env[name] = t
}

func (te *TypeEnv) LookupVar(name string) Type {
	return (te).env[name]
}
