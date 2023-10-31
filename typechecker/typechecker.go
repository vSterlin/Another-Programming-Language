package typechecker

import "language/ast"

type TypeChecker struct {
	env                *Env
	currentFuncRetType Type
}

func NewTypeChecker() *TypeChecker {
	return &TypeChecker{
		env:                NewEnv(nil),
		currentFuncRetType: INVALID,
	}
}

func (t *TypeChecker) Check(prog *ast.Program) error {
	for _, stmt := range prog.Stmts {
		err := t.checkStmt(stmt)
		if err != nil {
			return err
		}
	}

	// Not sure what to return here.
	return nil
}
