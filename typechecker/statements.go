package typechecker

import "language/ast"

func (t *TypeChecker) checkStmt(stmt ast.Stmt) error {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		return t.checkExpr(stmt.Expr)
	default:
		panic("unknown statement type")
	}
}
