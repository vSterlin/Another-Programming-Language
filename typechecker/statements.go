package typechecker

import (
	"language/ast"
)

func (t *TypeChecker) checkStmt(stmt ast.Stmt) error {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		_, err := t.checkExpr(stmt.Expr)

		return err
	default:
		panic("unknown statement type")
	}
}
