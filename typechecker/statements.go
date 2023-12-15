package typechecker

import (
	"fmt"
	"language/ast"
)

func (t *TypeChecker) checkStmt(stmt ast.Stmt) error {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		_, err := t.checkExpr(stmt.Expr)
		return err
	case *ast.BlockStmt:
		return t.checkBlockStmt(stmt, NewEnv(t.env))
	case *ast.VarAssignStmt:
		return t.checkVarAssignStmt(stmt)
	case *ast.FuncDecStmt:
		return t.checkFuncDecStmt(stmt)
	case *ast.IfStmt:
		return t.checkIfStmt(stmt)
	case *ast.WhileStmt:
		return t.checkWhileStmt(stmt)
	case *ast.ReturnStmt:
		return t.checkReturnStmt(stmt)
	case *ast.TypeAliasStmt:
		return t.checkTypeAliasStmt(stmt)

	default:
		return NewTypeError(fmt.Sprintf("unknown statement type: %T", stmt))
	}
}

func (t *TypeChecker) checkBlockStmt(stmt *ast.BlockStmt, env *Env) error {

	prevEnv := t.env
	t.env = env
	defer func() { t.env = prevEnv }()

	needsToReturn := !t.isInLoop && !t.currentFuncRetType.Equals(Invalid) && !t.currentFuncRetType.Equals(Void)

	if len(stmt.Stmts) == 0 {
		if needsToReturn {
			return NewTypeError("missing return statement")
		}
		return nil
	}

	for i, s := range stmt.Stmts {
		err := t.checkStmt(s)
		if err != nil {
			return err
		}

		if needsToReturn {
			if i == len(stmt.Stmts)-1 && !hasReturned(s) {
				return NewTypeError("missing return statement")
			}
		}
	}

	return nil
}

func hasReturned(stmt ast.Stmt) bool {
	switch stmt := stmt.(type) {
	case *ast.ReturnStmt:
		return true
	case *ast.IfStmt:
		return hasReturned(stmt.Consequent) && hasReturned(stmt.Alternate)
	case *ast.WhileStmt:
		return hasReturned(stmt.Body)
	case *ast.BlockStmt:
		return len(stmt.Stmts) > 0 && hasReturned(stmt.Stmts[len(stmt.Stmts)-1])
	default:
		return false
	}
}

func (t *TypeChecker) checkVarAssignStmt(stmt *ast.VarAssignStmt) error {

	initType, err := t.checkExpr(stmt.Init)
	if err != nil {
		return err
	}

	if stmt.Op == ":=" {
		t.env.Define(stmt.Id.Name, initType)
		return nil
	} else {

		foundVar, foundEnv, err := t.env.Get(stmt.Id.Name)
		if err != nil {
			return err
		}

		if !areTypesEqual(foundVar, initType) {
			return NewTypeError(fmt.Sprintf("cannot assign value of type %s to variable of type %s", initType, foundVar))
		}

		return foundEnv.Assign(stmt.Id.Name, initType)
	}

}

func (t *TypeChecker) checkFuncDecStmt(stmt *ast.FuncDecStmt) error {

	retType, err := resolveType(stmt.ReturnType, t.env)
	if err != nil {
		return err
	}

	funcType := FuncType{
		Args:       []Type{},
		ReturnType: retType,
	}
	for _, param := range stmt.Args {
		paramType, err := resolveType(param.Type, t.env)

		if err != nil {
			return err
		}

		t.env.Define(param.Id.Name, paramType)

		funcType.Args = append(funcType.Args, paramType)
	}

	t.env.Define(stmt.Id.Name, funcType)

	prevFuncRetType := t.currentFuncRetType
	t.currentFuncRetType = retType

	err = t.checkBlockStmt(stmt.Body, NewEnv(t.env))
	if err != nil {
		return err
	}

	t.currentFuncRetType = prevFuncRetType

	return nil
}

func (t *TypeChecker) checkIfStmt(stmt *ast.IfStmt) error {

	testType, err := t.checkExpr(stmt.Test)
	if err != nil {
		return err
	}

	if !areTypesEqual(testType, Boolean) {
		return NewTypeError(fmt.Sprintf("expected %s, got %s", Boolean, testType))
	}

	err = t.checkStmt(stmt.Consequent)
	if err != nil {
		return err
	}

	if stmt.Alternate != nil {
		return t.checkStmt(stmt.Alternate)
	}

	return nil
}

func (t *TypeChecker) checkWhileStmt(stmt *ast.WhileStmt) error {

	testType, err := t.checkExpr(stmt.Test)
	if err != nil {
		return err
	}

	if !areTypesEqual(testType, BooleanType{}) {
		return NewTypeError(fmt.Sprintf("expected %s, got %s", BooleanType{}, testType))
	}

	t.isInLoop = true
	err = t.checkStmt(stmt.Body)
	t.isInLoop = false

	return err

}

// TODO: review
func (t *TypeChecker) checkReturnStmt(stmt *ast.ReturnStmt) error {

	expectedType := t.currentFuncRetType

	fmt.Printf("expectedType: %v\n", expectedType)

	if expectedType == Invalid && t.currentArrowFuncType == nil {
		return NewTypeError("return statement outside of function")
	}

	var actualType Type = Void

	if !expectedType.Equals(Void) {
		t, err := t.checkExpr(stmt.Arg)
		if err != nil {
			return err
		}
		actualType = t
	} else {
		if stmt.Arg != nil {
			return NewTypeError("expected no return value")
		}
	}

	if t.currentArrowFuncType != nil {
		fmt.Printf("actualType: %v\n", actualType)
		t.currentArrowFuncType.ReturnType = actualType
		return nil
	}

	if !areTypesEqual(expectedType, actualType) {
		return NewTypeError(fmt.Sprintf("expected return type %s, got %s", expectedType, actualType))
	}

	return nil
}

func (t *TypeChecker) checkTypeAliasStmt(stmt *ast.TypeAliasStmt) error {

	aliasType, err := resolveType(stmt.Type, t.env)

	if err != nil {
		return err
	}

	t.env.DefineType(stmt.Id.Name, aliasType)

	return nil
}
