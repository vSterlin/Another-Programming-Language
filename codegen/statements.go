package codegen

import (
	"fmt"
	"language/ast"
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
	return cg.genExpr(stmt.Expr)

}

func (cg *CodeGenerator) genFuncDecStmt(stmt *ast.FuncDecStmt) (string, error) {

	retType := cType(stmt.ReturnType.Name)

	id := stmt.Id.Name

	body, err := cg.genBlockStmt(stmt.Body)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s() %s", retType, id, body), nil

}

func (cg *CodeGenerator) genBlockStmt(stmt *ast.BlockStmt) (string, error) {

	stmts := ""
	tabs := cg.genTabs()
	for _, stmt := range stmt.Stmts {
		code, err := cg.genStmt(stmt)
		if err != nil {
			return "", err
		}

		stmts += fmt.Sprintf("%s%s\n", tabs, code)
	}

	return fmt.Sprintf("{\n%s%s}", stmts, tabs), nil
}

func (cg *CodeGenerator) genVarAssignStmt(stmt *ast.VarAssignStmt) (string, error) {

	varType := ""

	switch stmt.Init.(type) {

	case *ast.NumberExpr:
		varType = "int"
	case *ast.StringExpr:
		varType = "std::string"
	case *ast.BooleanExpr:
		varType = "bool"
	}

	id := stmt.Id.Name

	init, err := cg.genExpr(stmt.Init)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s = %s;", varType, id, init), nil

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
