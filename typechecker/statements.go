package typechecker

import (
	"fmt"
	"language/ast"
)

func (t *TypeChecker) checkStmt(stmt ast.Stmt) error {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		t, err := t.checkExpr(stmt.Expr)
		fmt.Println(t)
		return err
	default:
		panic("unknown statement type")
	}
}
