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
	deferredStmts := []string{}
	for _, stmt := range program.Stmts {
		if stmt, ok := stmt.(*ast.DeferStmt); ok {
			deferredStmts = append([]string{j.generateStmt(stmt)}, deferredStmts...)
			continue
		}
		code += fmt.Sprintf("%s\n", j.generateStmt(stmt))
	}

	if len(deferredStmts) > 0 {
		code += "\n"
		code += strings.Join(deferredStmts, "\n")
		code += "\n"
	}
	code = strings.TrimSpace(code)
	return code
}

func (j *JavascriptCodeGenerator) generateStmt(stmt ast.Stmt) string {
	switch stmt := stmt.(type) {
	case *ast.VarAssignStmt:
		return j.generateVarAssignStmt(stmt)
	case *ast.ExprStmt:
		return j.generateExpr(stmt.Expr)
	case *ast.BlockStmt:
		return j.generateBlockStmt(stmt)
	case *ast.WhileStmt:
		return j.generateWhileStmt(stmt)
	case *ast.FuncDecStmt:
		return j.generateFuncDecStmt(stmt)
	case *ast.IfStmt:
		return j.generateIfStmt(stmt)
	case *ast.DeferStmt:
		return j.generateDeferStmt(stmt)
	case *ast.RangeStmt:
		return j.generateRangeStmt(stmt)
	default:
		return ""
	}
}

func (j *JavascriptCodeGenerator) generateRangeStmt(stmt *ast.RangeStmt) string {
	id := stmt.Id.Name
	expr := j.generateExpr(stmt.Expr)
	body := j.generateStmt(stmt.Body)
	return fmt.Sprintf("%s.forEach(%s => %s)", expr, id, body)
}

func (j *JavascriptCodeGenerator) generateDeferStmt(stmt *ast.DeferStmt) string {
	call := j.generateExpr(stmt.Call)
	// need to execute this call at the end of the statement it is in
	return call
}
func (j *JavascriptCodeGenerator) generateIfStmt(stmt *ast.IfStmt) string {
	test := j.generateExpr(stmt.Test)
	consequent := j.generateStmt(stmt.Consequent)
	alternate := j.generateStmt(stmt.Alternate)

	if alternate != "" {
		return fmt.Sprintf("if (%s) %s else %s", test, consequent, alternate)
	} else {
		return fmt.Sprintf("if (%s) %s", test, consequent)
	}
}

func (j *JavascriptCodeGenerator) generateWhileStmt(stmt *ast.WhileStmt) string {
	test := j.generateExpr(stmt.Test)
	body := j.generateStmt(stmt.Body)
	return fmt.Sprintf("while (%s) %s", test, body)
}

func (j *JavascriptCodeGenerator) generateFuncDecStmt(stmt *ast.FuncDecStmt) string {
	id := stmt.Id.Name
	args := []string{}
	for _, arg := range stmt.Args {
		args = append(args, arg.Name)
	}
	body := j.generateBlockStmt(stmt.Body)

	return fmt.Sprintf("function %s(%s)%s", id, strings.Join(args, ", "), body)
}

func (j *JavascriptCodeGenerator) generateBlockStmt(stmt *ast.BlockStmt) string {
	code := "{\n"
	deferredStmts := []string{}
	for _, stmt := range stmt.Stmts {
		if stmt, ok := stmt.(*ast.DeferStmt); ok {
			deferredStmts = append([]string{j.generateStmt(stmt)}, deferredStmts...)
			continue
		}
		stmtCode := j.generateStmt(stmt)
		code += stmtCode + "\n"
	}

	if len(deferredStmts) > 0 {
		code += "\n"
		code += strings.Join(deferredStmts, "\n")
		code += "\n"
	}
	code += "}"
	return code
}

func (j *JavascriptCodeGenerator) generateVarDecStmt(stmt *ast.VarDecStmt) string {
	return fmt.Sprintf("var %s = %s", stmt.Id.Name, j.generateExpr(stmt.Init))
}

func (j *JavascriptCodeGenerator) generateVarAssignStmt(stmt *ast.VarAssignStmt) string {
	init := j.generateExpr(stmt.Init)
	if stmt.Op == "=" {
		return fmt.Sprintf("%s = %s", stmt.Id.Name, init)
	} else { // ":="
		return fmt.Sprintf("var %s = %s", stmt.Id.Name, init)
	}
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
		return ""
		// return j.generateCallExpr(expr)
	case *ast.SliceExpr:
		return j.generateSliceExpr(expr)
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

// func (j *JavascriptCodeGenerator) generateCallExpr(expr *ast.CallExpr) string {

// 	args := []string{}
// 	for _, arg := range expr.Args {
// 		arg := j.generateExpr(arg)
// 		args = append(args, arg)
// 	}
// 	funcName := expr.Callee.Name

// 	if val, ok := globalFuncs[funcName]; ok {
// 		funcName = val
// 	}

// 	return fmt.Sprintf("%s(%s)", funcName, strings.Join(args, ", "))
// }

func (j *JavascriptCodeGenerator) generateSliceExpr(expr *ast.SliceExpr) string {
	id := expr.Id.Name
	low := j.generateExpr(expr.Low)
	high := j.generateExpr(expr.High)

	slice := fmt.Sprintf("%s.slice(%s, %s)", id, low, high)

	if expr.Step != nil && expr.Step.(*ast.NumberExpr).Val != 1 {
		step := j.generateExpr(expr.Step)
		slice = fmt.Sprintf("%s.filter((_, i) => i %% %s === 0)", slice, step)
	}
	return slice
}

func (j *JavascriptCodeGenerator) generateArrayExpr(expr *ast.ArrayExpr) string {
	elements := []string{}
	for _, e := range expr.Elements {
		elements = append(elements, j.generateExpr(e))
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}
