package codegen

import (
	"fmt"
	"language/ast"
	"strings"
)

// Statements
func (cg *CodeGenerator) genStmt(stmt ast.Stmt) error {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		return cg.genExprStmt(stmt)
	case *ast.FuncDecStmt:
		return cg.genFuncDecStmt(stmt)

	case *ast.BlockStmt:
		return cg.genBlockStmt(stmt)

	case *ast.VarAssignStmt:
		return cg.genVarAssignStmt(stmt)

	// case *ast.IfStmt:
	// 	return cg.genIfStmt(stmt)
	// case *ast.WhileStmt:
	// 	return cg.genWhileStmt(stmt)
	case *ast.ReturnStmt:
		return cg.genReturnStmt(stmt)

	default:
		return fmt.Errorf("unknown statement type: %T", stmt)
	}

}

func (cg *CodeGenerator) genExprStmt(stmt *ast.ExprStmt) error {
	_, err := cg.genExpr(stmt.Expr)
	return err
}

func (cg *CodeGenerator) genFuncDecStmt(stmt *ast.FuncDecStmt) error {

	cg.isInGlobal = false
	funcName, err := cg.genExpr(stmt.Id)
	if err != nil {
		return err
	}

	retType, err := cg.genExpr(stmt.ReturnType)
	if err != nil {
		return err
	}

	args := []string{}

	for _, arg := range stmt.Args {
		argStr := fmt.Sprintf("%s %s", cTypeFromAst(arg.Type), arg.Id.Name)
		args = append(args, argStr)
	}

	argsStr := strings.Join(args, ", ")

	cg.write(fmt.Sprintf("%s %s(%s)", retType, funcName, argsStr))

	err = cg.genBlockStmt(stmt.Body)
	if err != nil {
		return err
	}

	cg.isInGlobal = true

	return nil

}

func (cg *CodeGenerator) genBlockStmt(stmt *ast.BlockStmt) error {
	prevIsInGlobal := cg.isInGlobal
	cg.isInGlobal = false
	// stmts := ""
	cg.indent++
	tabs := cg.genTabs()
	cg.write("{\n")
	for _, stmt := range stmt.Stmts {
		err := cg.genStmt(stmt)
		if err != nil {
			return err
		}

	}

	cg.indent--

	tabs = tabs[1:]

	cg.isInGlobal = prevIsInGlobal

	cg.write("\n}\n")

	return nil
}

func (cg *CodeGenerator) genVarAssignStmt(stmt *ast.VarAssignStmt) error {

	id := stmt.Id.Name

	if stmt.Op == ":=" {
		varType := inferFromAstNode(stmt.Init)
		cg.write(fmt.Sprintf("%s %s = ", varType, id))
	} else {
		cg.write(fmt.Sprintf("%s = ", id))
	}

	_, err := cg.genExpr(stmt.Init)
	if err != nil {
		return err
	}
	cg.write(";")
	return nil

}

// func (cg *CodeGenerator) genIfStmt(stmt *ast.IfStmt) (string, error) {

// 	test, err := cg.genExpr(stmt.Test)

// 	if err != nil {
// 		return "", err
// 	}

// 	err = cg.genStmt(stmt.Consequent)

// 	if err != nil {
// 		return "", err
// 	}

// 	if stmt.Alternate != nil {
// 		alternate, err := cg.genStmt(stmt.Alternate)

// 		if err != nil {
// 			return "", err
// 		}

// 		return fmt.Sprintf("if (%s) %s else %s", test, body, alternate), nil
// 	}

// 	return fmt.Sprintf("if (%s) %s", test, body), nil

// }

// func (cg *CodeGenerator) genWhileStmt(stmt *ast.WhileStmt) (string, error) {

// 	test, err := cg.genExpr(stmt.Test)
// 	if err != nil {
// 		return "", err
// 	}
// 	body, err := cg.genStmt(stmt.Body)
// 	if err != nil {
// 		return "", err
// 	}

// 	return fmt.Sprintf("while (%s) %s", test, body), nil

// }

func (cg *CodeGenerator) genReturnStmt(stmt *ast.ReturnStmt) error {

	cg.write(("return"))

	if stmt.Arg != nil {
		cg.write(" ")
		_, err := cg.genExpr(stmt.Arg)
		if err != nil {
			return err
		}
	}

	cg.write(";")
	return nil
}
