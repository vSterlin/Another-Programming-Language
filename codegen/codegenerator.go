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
		code += fmt.Sprintf("%s\n", j.generateStmt(stmt))
	}

	code = strings.TrimSpace(code)
	return code
}

func (j *JavascriptCodeGenerator) generateStmt(stmt ast.Stmt) string {
	switch stmt := stmt.(type) {
	case *ast.VarDecStmt:
		return j.generateVarDecStmt(stmt)
	case *ast.ExprStmt:
		return j.generateExpr(stmt.Expr)
	case *ast.BlockStmt:
		return j.generateBlockStmt(stmt)
	default:
		return ""
	}
}

func (j *JavascriptCodeGenerator) generateBlockStmt(stmt *ast.BlockStmt) string {
	code := "{\n"
	for _, stmt := range stmt.Stmts {
		stmtCode := j.generateStmt(stmt)
		code += stmtCode + "\n"
	}

	code += "}"
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
	lhs := j.generateExpr(expr.Lhs)
	rhs := j.generateExpr(expr.Rhs)
	if expr.Op == "**" {
		return fmt.Sprintf("Math.pow(%s, %s)", lhs, rhs)
	}
	return fmt.Sprintf("(%s %s %s)", lhs, expr.Op, rhs)
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
