package codegen

import (
	"fmt"
	"language/ast"
	"strings"
)

// Statements
func (cg *CodeGenerator) genStmt(stmt ast.Stmt) (string, error) {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		return cg.genExprStmt(stmt)
	case *ast.FuncDecStmt:
		return cg.genFuncDecStmt(stmt)
	case *ast.BlockStmt:
		return cg.genBlockStmt(stmt)
	case *ast.VarAssignStmt:
		return cg.genVarAssignStmt(stmt)
	case *ast.IfStmt:
		return cg.genIfStmt(stmt)
	case *ast.WhileStmt:
		return cg.genWhileStmt(stmt)
	case *ast.ReturnStmt:
		return cg.genReturnStmt(stmt)
	default:
		return "", fmt.Errorf("unknown statement type: %T", stmt)
	}

}

func (cg *CodeGenerator) genExprStmt(stmt *ast.ExprStmt) (string, error) {
	expr, err := cg.genExpr(stmt.Expr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s;", expr), nil
}

// TODO: fix
func (cg *CodeGenerator) genFuncDecStmt(stmt *ast.FuncDecStmt) (string, error) {

	funcName, err := cg.genExpr(stmt.Id)
	if err != nil {
		return "", err
	}

	retType, err := cg.genExpr(stmt.ReturnType)
	if err != nil {
		return "", err
	}

	body, err := cg.genBlockStmt(stmt.Body)
	if err != nil {
		return "", err
	}

	args := []string{}

	for _, arg := range stmt.Args {
		argStr := fmt.Sprintf("%s %s", cTypeFromAst(arg.Type), arg.Id.Name)
		args = append(args, argStr)
	}

	argsStr := strings.Join(args, ", ")

	return fmt.Sprintf("%s %s(%s) %s", retType, funcName, argsStr, body), nil

}

func (cg *CodeGenerator) genBlockStmt(stmt *ast.BlockStmt) (string, error) {

	stmts := ""
	cg.indent++
	tabs := cg.genTabs()
	for _, stmt := range stmt.Stmts {
		code, err := cg.genStmt(stmt)
		if err != nil {
			return "", err
		}

		stmts += fmt.Sprintf("%s%s\n", tabs, code)
	}

	cg.indent--

	tabs = tabs[1:]

	return fmt.Sprintf("{\n%s%s}\n", stmts, tabs), nil
}

func (cg *CodeGenerator) genVarAssignStmt(stmt *ast.VarAssignStmt) (string, error) {

	id := stmt.Id.Name

	init, err := cg.genExpr(stmt.Init)
	if err != nil {
		return "", err
	}
	if stmt.Op == ":=" {

		varType := ""

		switch stmt.Init.(type) {

		case *ast.NumberExpr:
			varType = "int"
		case *ast.StringExpr:
			varType = "std::string"
		case *ast.BooleanExpr:
			varType = "bool"
		default:

			// TODO: fix
			varType = "auto"
		}

		return fmt.Sprintf("%s %s = %s;", varType, id, init), nil
	} else {
		return fmt.Sprintf("%s = %s;", id, init), nil
	}

}

func (cg *CodeGenerator) genIfStmt(stmt *ast.IfStmt) (string, error) {

	test, err := cg.genExpr(stmt.Test)

	if err != nil {
		return "", err
	}

	body, err := cg.genStmt(stmt.Consequent)

	if err != nil {
		return "", err
	}

	if stmt.Alternate != nil {
		alternate, err := cg.genStmt(stmt.Alternate)

		if err != nil {
			return "", err
		}

		return fmt.Sprintf("if (%s) %s else %s", test, body, alternate), nil
	}

	return fmt.Sprintf("if (%s) %s", test, body), nil

}

func (cg *CodeGenerator) genWhileStmt(stmt *ast.WhileStmt) (string, error) {

	test, err := cg.genExpr(stmt.Test)
	if err != nil {
		return "", err
	}
	body, err := cg.genStmt(stmt.Body)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("while (%s) %s", test, body), nil

}

func (cg *CodeGenerator) genReturnStmt(stmt *ast.ReturnStmt) (string, error) {

	returnedVal, err := cg.genExpr(stmt.Arg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("return %s;", returnedVal), nil
}
