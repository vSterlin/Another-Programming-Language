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

	default:
		panic("unknown expression type")

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
		return INVALID, err
	}
	rhs, err := t.checkExpr(expr.Rhs)
	if err != nil {
		return INVALID, err
	}

	switch expr.Op {
	case ast.ADD, ast.SUB, ast.MUL, ast.DIV, ast.LT, ast.GT, ast.LTE, ast.GTE:
		if expr.Op == ast.ADD && areTypesEqual(lhs, rhs, String) {
			return String, nil
		}
		if !areTypesEqual(lhs, rhs, Number) {
			return INVALID, NewTypeError(fmt.Sprintf("expected %s, got %s", Number, lhs))
		}

	case ast.EQ, ast.NEQ:
		if !areTypesEqual(lhs, rhs) {
			return INVALID, NewTypeError(fmt.Sprintf("expected %s, got %s", lhs, rhs))
		}

	}

	return lhs, nil
}

func (t *TypeChecker) checkLogicalExpr(expr *ast.LogicalExpr) (Type, error) {

	lhs, err := t.checkExpr(expr.Lhs)
	if err != nil {
		return INVALID, err
	}
	rhs, err := t.checkExpr(expr.Rhs)
	if err != nil {
		return INVALID, err
	}

	if !areTypesEqual(lhs, rhs, Boolean) {
		return INVALID, NewTypeError(
			fmt.Sprintf("expected both operands to be of type %s, got %s and %s",
				Boolean, lhs, rhs))
	}

	return Boolean, nil
}

func (t *TypeChecker) checkIdentifierExpr(expr *ast.IdentifierExpr) (Type, error) {
	return t.env.Get(expr.Name)
}
