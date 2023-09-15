package codegen

import (
	"fmt"
	"language/ast"
	"strings"
)

type CodeGenerator interface {
	Generate(*ast.Program) string
}

type JavascriptCodeGenerator struct{}

func NewJavascriptCodeGenerator() *JavascriptCodeGenerator {
	return &JavascriptCodeGenerator{}
}

func (j *JavascriptCodeGenerator) Generate(program *ast.Program) string {
	code := ""
	for _, stmt := range program.Stmts {
		switch stmt := stmt.(type) {
		case *ast.VarDecStmt:
			code += j.generateVarDecStmt(stmt) + "\n"
		case *ast.ExprStmt:
			code += j.generateExpr(stmt.Expr) + "\n"
		}
	}

	code = strings.TrimSpace(code)
	return code
}

func (j *JavascriptCodeGenerator) generateVarDecStmt(stmt *ast.VarDecStmt) string {
	return fmt.Sprintf("var %s = %s", stmt.Id.Name, j.generateExpr(stmt.Init))
}

func (j *JavascriptCodeGenerator) generateExpr(expr ast.Expr) string {

	switch expr := expr.(type) {
	case *ast.NumberExpr:
		return fmt.Sprint(expr.Val)
	case *ast.BooleanExpr:
		return fmt.Sprint(expr.Val)
	case *ast.StringExpr:
		return fmt.Sprintf("\"%s\"", expr.Val)
	case *ast.IdentifierExpr:
		return expr.Name
	case *ast.BinaryExpr:
		return j.generateBinaryExpr(expr)
	case *ast.CallExpr:
		return j.generateCallExpr(expr)
	case *ast.ArrayExpr:
		return j.generateArrayExpr(expr)
	default:
		return ""
	}
}

func (j *JavascriptCodeGenerator) generateBinaryExpr(expr *ast.BinaryExpr) string {
	return fmt.Sprintf("(%s %s %s)", j.generateExpr(expr.Lhs), expr.Op, j.generateExpr(expr.Rhs))
}

var globalFuncs map[string]string = map[string]string{
	"print": "console.log",
}

func (j *JavascriptCodeGenerator) generateCallExpr(expr *ast.CallExpr) string {
	// TODO: handle multiple args
	arg := j.generateExpr(expr.Args[0])
	funcName := expr.Callee.Name

	if val, ok := globalFuncs[funcName]; ok {
		funcName = val
	}

	return fmt.Sprintf("%s(%s)", funcName, arg)
}

func (j *JavascriptCodeGenerator) generateArrayExpr(expr *ast.ArrayExpr) string {
	elements := []string{}
	for _, e := range expr.Elements {
		elements = append(elements, j.generateExpr(e))
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}
