package codegen

import (
	"fmt"
	"language/ast"
	"strconv"
	"strings"
)

func (cg *CodeGenerator) genExpr(expr ast.Expr) (string, error) {
	switch expr := expr.(type) {

	case *ast.BinaryExpr:
		return cg.genBinaryExpr(expr)
	case *ast.NumberExpr:
		return cg.genNumberExpr(expr)
	case *ast.StringExpr:
		return cg.genStringExpr(expr)
	case *ast.BooleanExpr:
		return genBooleanExpr(expr)
	case *ast.LogicalExpr:
		return cg.genLogicalExpr(expr)
	case *ast.CallExpr:
		return cg.genCallExpr(expr)
	case *ast.IdentifierExpr:
		return cg.genIdentifierExpr(expr)
	default:
		return "", nil
	}
}

// Literals start
func (cg *CodeGenerator) genNumberExpr(expr *ast.NumberExpr) (string, error) {
	return strconv.Itoa(expr.Val), nil
}
func (cg *CodeGenerator) genStringExpr(expr *ast.StringExpr) (string, error) {
	str := expr.Val
	return fmt.Sprintf("\"%s\"", str), nil
}

func genBooleanExpr(expr *ast.BooleanExpr) (string, error) {

	if expr.Val {
		return "true", nil
	} else {
		return "false", nil
	}

}

// Literals end

func (cg *CodeGenerator) genLogicalExpr(expr *ast.LogicalExpr) (string, error) {

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
	case ast.AND:
		res = lhs + " && " + rhs
	case ast.OR:
		res = lhs + " || " + rhs

	}

	return res, nil
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

	case "==":
		res = lhs + " == " + rhs
	case "!=":
		res = lhs + " != " + rhs

	case "<":
		res = lhs + " < " + rhs
	case "<=":
		res = lhs + " <= " + rhs
	case ">":
		res = lhs + " > " + rhs
	case ">=":
		res = lhs + " >= " + rhs
	}

	return res, nil
}

func (cg *CodeGenerator) genCallExpr(expr *ast.CallExpr) (string, error) {

	callee, err := cg.genExpr(expr.Callee)

	if err != nil {
		return "", err
	}

	args := []string{}
	for _, arg := range expr.Args {
		argStr, err := cg.genExpr(arg)
		if err != nil {
			return "", err
		}

		args = append(args, argStr)
	}

	argsStr := strings.Join(args, ", ")

	return fmt.Sprintf("%s(%s)", callee, argsStr), nil
}

func (cg *CodeGenerator) genIdentifierExpr(expr *ast.IdentifierExpr) (string, error) {
	// TODO: need to check env here
	return expr.Name, nil
}
