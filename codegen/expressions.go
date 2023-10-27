package codegen

import (
	"language/ast"
	"strconv"
)

func (cg *CodeGenerator) genExpr(expr ast.Expr) (string, error) {
	switch expr := expr.(type) {

	case *ast.BinaryExpr:
		return cg.genBinaryExpr(expr)
	case *ast.NumberExpr:
		return cg.genNumberExpr(expr)
	case *ast.StringExpr:
		return "", cg.genStringExpr(expr)
	case *ast.BooleanExpr:
		return "", genBooleanExpr(expr)
	case *ast.LogicalExpr:
		return "", cg.genLogicalExpr(expr)
	case *ast.CallExpr:
		return "", cg.genCallExpr(expr)
	case *ast.IdentifierExpr:
		return "", cg.genIdentifierExpr(expr)
	default:
		return "", nil
	}
}

// Literals start
func (cg *CodeGenerator) genNumberExpr(expr *ast.NumberExpr) (string, error) {
	return strconv.Itoa(expr.Val), nil
}
func (cg *CodeGenerator) genStringExpr(expr *ast.StringExpr) error {

	return nil
}

func genBooleanExpr(expr *ast.BooleanExpr) error {
	return nil
}

// Literals end

func (cg *CodeGenerator) genLogicalExpr(expr *ast.LogicalExpr) error {
	return nil
}

func (cg *CodeGenerator) genBinaryExpr(expr *ast.BinaryExpr) (string, error) {
	lhs, err := cg.genExpr(expr.Lhs)
	if err != nil {
		return "", err
	}
	rhs, err := cg.genExpr(expr.Rhs)
	if err != nil {
		return "", err
	}

	res := ""
	switch expr.Op {
	case ast.ADD:
		res = lhs + " + " + rhs
	case ast.SUB:
		res = lhs + " - " + rhs
	case ast.MUL:
		res = lhs + " * " + rhs
	case ast.DIV:
		res = lhs + " / " + rhs
	}

	return res, nil
}

func (cg *CodeGenerator) genCallExpr(expr *ast.CallExpr) error {
	return nil
}

func (cg *CodeGenerator) genIdentifierExpr(expr *ast.IdentifierExpr) error {
	return nil
}
