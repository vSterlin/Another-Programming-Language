package typechecker

import (
	"fmt"
	"language/ast"
)

func (t *TypeChecker) checkExpr(expr ast.Expr) (Type, error) {
	switch expr := expr.(type) {

	case *ast.NumberExpr:
		return t.checkNumberExpr(expr)
	case *ast.StringExpr:
		return t.checkStringExpr(expr)
	case *ast.BooleanExpr:
		return t.checkBooleanExpr(expr)

	case *ast.BinaryExpr:
		return t.checkBinaryExpr(expr)
	case *ast.LogicalExpr:
		return t.checkLogicalExpr(expr)
	case *ast.IdentifierExpr:
		return t.checkIdentifierExpr(expr)

	case *ast.ArrowFunc:
		return t.checkArrowFunc(expr)
	case *ast.CallExpr:
		return t.checkCallExpr(expr)

	case *ast.UnaryExpr:
		return t.checkUnaryExpr(expr)
	case *ast.UpdateExpr:
		return t.checkUpdateExpr(expr)

	default:
		return Invalid, NewTypeError(fmt.Sprintf("unknown expression type: %T", expr))

	}
}

func (t *TypeChecker) checkNumberExpr(expr *ast.NumberExpr) (Type, error) {
	return Number, nil
}

func (t *TypeChecker) checkStringExpr(expr *ast.StringExpr) (Type, error) {
	return String, nil
}

func (t *TypeChecker) checkBooleanExpr(expr *ast.BooleanExpr) (Type, error) {
	return Boolean, nil
}

func (t *TypeChecker) checkBinaryExpr(expr *ast.BinaryExpr) (Type, error) {

	lhs, err := t.checkExpr(expr.Lhs)
	if err != nil {
		return Invalid, err
	}
	rhs, err := t.checkExpr(expr.Rhs)
	if err != nil {
		return Invalid, err
	}

	switch expr.Op {
	case ast.ADD, ast.SUB, ast.MUL, ast.DIV, ast.LT, ast.GT, ast.LTE, ast.GTE:
		if expr.Op == ast.ADD && areTypesEqual(lhs, rhs, String) {
			return String, nil
		}
		if !areTypesEqual(lhs, rhs, Number) {
			return Invalid, NewTypeError(fmt.Sprintf("expected %s, got %s", Number, lhs))
		}

	case ast.EQ, ast.NEQ:
		if !areTypesEqual(lhs, rhs) {
			return Invalid, NewTypeError(fmt.Sprintf("expected %s, got %s", lhs, rhs))
		}

	}

	typ := lhs
	switch expr.Op {
	case ast.ADD, ast.SUB, ast.MUL, ast.DIV:
		typ = Number
	case ast.LT, ast.GT, ast.LTE, ast.GTE, ast.EQ, ast.NEQ:
		typ = Boolean
	}

	expr.Type = toAstNode(typ)

	return typ, nil

}

func (t *TypeChecker) checkLogicalExpr(expr *ast.LogicalExpr) (Type, error) {

	lhs, err := t.checkExpr(expr.Lhs)
	if err != nil {
		return Invalid, err
	}
	rhs, err := t.checkExpr(expr.Rhs)
	if err != nil {
		return Invalid, err
	}

	if !areTypesEqual(lhs, rhs, Boolean) {
		return Invalid, NewTypeError(
			fmt.Sprintf("expected both operands to be of type %s, got %s and %s",
				Boolean, lhs, rhs))
	}

	return Boolean, nil
}

func (t *TypeChecker) checkIdentifierExpr(expr *ast.IdentifierExpr) (Type, error) {
	typ, _, err := t.env.Get(expr.Name)
	return typ, err
}

func (t *TypeChecker) checkArrowFunc(expr *ast.ArrowFunc) (Type, error) {

	retType, err := resolveType(expr.ReturnType, t.env)

	if err != nil {
		return retType, err
	}

	funcType := FuncType{
		Args:       []Type{},
		ReturnType: retType,
	}
	prevFuncRetType := t.currentFuncRetType
	t.currentFuncRetType = retType
	prevArrowFuncType := t.currentArrowFuncType
	t.currentArrowFuncType = &funcType

	defer func() {
		t.currentFuncRetType = prevFuncRetType
		t.currentArrowFuncType = prevArrowFuncType
	}()

	for _, param := range expr.Args {
		paramType, err := resolveType(param.Type, t.env)
		if err != nil {
			return Invalid, err
		}
		t.env.Define(param.Id.Name, paramType)
		funcType.Args = append(funcType.Args, paramType)
	}

	err = t.checkBlockStmt(expr.Body, NewEnv(t.env))
	if err != nil {
		return Invalid, err
	}

	return funcType, nil

}

func (t *TypeChecker) checkCallExpr(expr *ast.CallExpr) (Type, error) {

	// TODO:
	funcName := (expr.Callee.(*ast.IdentifierExpr)).Name

	if globalVar, exists := GetGlobalFuncReturnType(funcName); exists {
		return globalVar, nil
	}

	funcVar, _, err := t.env.Get(funcName)

	// REVIEW:
	funcDef, ok := funcVar.(FuncType)
	if !ok {
		return Invalid, NewTypeError(fmt.Sprintf("expected %s to be a function", funcName))
	}

	if err != nil {
		return Invalid, NewTypeError(fmt.Sprintf("undefined function: %s", funcName))
	}

	if len(funcDef.Args) != len(expr.Args) {
		return Invalid, NewTypeError(
			fmt.Sprintf("expected %d arguments, got %d",
				len(funcDef.Args), len(expr.Args)))
	}

	for i, arg := range expr.Args {
		argType, err := t.checkExpr(arg)
		if err != nil {
			return Invalid, err
		}
		expectedType := (funcDef.Args[i])
		if !areTypesEqual(argType, expectedType) {
			return Invalid, NewTypeError(
				fmt.Sprintf("expected argument %d to be of type %s, got %s",
					i+1, expectedType, argType))
		}
	}

	retType := funcDef.ReturnType

	expr.ReturnType = toAstNode(retType)

	return retType, nil
}

func (t *TypeChecker) checkUnaryExpr(expr *ast.UnaryExpr) (Type, error) {
	argType, err := t.checkExpr(expr.Arg)
	if err != nil {
		return Invalid, err
	}

	if !areTypesEqual(argType, Boolean) {
		return Invalid, NewTypeError(fmt.Sprintf("expected %s, got %s", Boolean, argType))
	}

	return Boolean, nil
}

func (t *TypeChecker) checkUpdateExpr(expr *ast.UpdateExpr) (Type, error) {
	argType, err := t.checkExpr(expr.Arg)
	if err != nil {
		return Invalid, err
	}

	if !areTypesEqual(argType, Number) {
		return Invalid, NewTypeError(fmt.Sprintf("expected %s, got %s", Number, argType))
	}

	return Number, nil
}
