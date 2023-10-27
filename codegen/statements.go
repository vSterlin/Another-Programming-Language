package codegen

import (
	"language/ast"
)

// Statements
func (cg *CodeGenerator) genStmt(stmt ast.Stmt) (string, error) {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		return cg.genExprStmt(stmt)
	case *ast.FuncDecStmt:
		return "", cg.genFuncDecStmt(stmt)
	case *ast.BlockStmt:
		return "", cg.genBlockStmt(stmt)
	case *ast.VarAssignStmt:
		return "", cg.genVarAssignStmt(stmt)
	case *ast.IfStmt:
		return "", cg.genIfStmt(stmt)
	case *ast.WhileStmt:
		return "", cg.genWhileStmt(stmt)
	case *ast.ReturnStmt:
		return "", cg.genReturnStmt(stmt)
	default:
		return "", nil
	}

}

func (cg *CodeGenerator) genExprStmt(stmt *ast.ExprStmt) (string, error) {
	return cg.genExpr(stmt.Expr)

}

func (cg *CodeGenerator) genFuncDecStmt(stmt *ast.FuncDecStmt) error {
	return nil

}

func (cg *CodeGenerator) genBlockStmt(stmt *ast.BlockStmt) error {
	return nil
}

func (cg *CodeGenerator) genVarAssignStmt(stmt *ast.VarAssignStmt) error {
	return nil

}

func (cg *CodeGenerator) genIfStmt(stmt *ast.IfStmt) error {
	return nil

}

func (cg *CodeGenerator) genWhileStmt(stmt *ast.WhileStmt) error {
	return nil

}

func (cg *CodeGenerator) genReturnStmt(stmt *ast.ReturnStmt) error {
	return nil
}
