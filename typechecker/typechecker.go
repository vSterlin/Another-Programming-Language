package typechecker

import "language/ast"

type TypeChecker struct {
	env *Env
}

func NewTypeChecker() *TypeChecker {
	return &TypeChecker{
		env: NewEnv(nil),
	}
}

func (t *TypeChecker) Check(prog *ast.Program) error {
	for _, stmt := range prog.Stmts {
		err := t.checkStmt(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}
