package typechecker

import "language/ast"

func (t *TypeChecker) checkExpr(expr ast.Expr) error {
	switch expr := expr.(type) {

	case *ast.NumberExpr:
		return t.checkNumberExpr(expr)
	// case *ast.StringExpr:
	// 	t.checkStringExpr(expr)
	// case *ast.BooleanExpr:
	// 	t.checkBooleanExpr(expr)

	// case *ast.BinaryExpr:
	// 	t.checkBinaryExpr(expr)

	default:
		panic("unknown expression type")
	}
}

func (t *TypeChecker) checkNumberExpr(expr *ast.NumberExpr) error {
	return nil
}
