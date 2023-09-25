package interpreter

import (
	"fmt"
	"language/ast"
)

type functionType string

const (
	noneFuncType     functionType = "NONE"
	functionFuncType functionType = "FUNCTION"
)

type resolver struct {
	interpreter *Interpreter
	scopes      scopeStack
	currentFunc functionType
}

func NewResolver(interpreter *Interpreter) *resolver {
	return &resolver{interpreter: interpreter, scopes: []*scope{}, currentFunc: noneFuncType}
}

func (r *resolver) beginScope() { r.scopes.push(&scope{}) }
func (r *resolver) endScope()   { r.scopes.pop() }

func (r *resolver) declare(name string) error {
	if r.scopes.isEmpty() {
		return nil
	}
	scope := r.scopes.peek()

	if _, ok := (*scope)[name]; ok {
		return NewRuntimeError(fmt.Sprintf("variable %s already declared in this scope", name))
	}

	(*scope)[name] = false
	return nil
}

func (r *resolver) define(name string) {
	if r.scopes.isEmpty() {
		return
	}
	scope := r.scopes.peek()
	(*scope)[name] = true
}

func (r *resolver) resolveLocal(expr ast.Expr, name string) {
	// traverse scopes from innermost
	for i := len(r.scopes) - 1; i >= 0; i-- {
		scope := r.scopes[i]
		if scope.isDefined(name) {
			r.interpreter.Resolve(expr, len(r.scopes)-1-i)
			return
		}

	}
}

// Statements
func (r *resolver) ResolveProgram(p *ast.Program) error {
	for _, stmt := range p.Stmts {
		err := r.resolveStmt(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO: Might consider using a visitor pattern here
func (r *resolver) resolveStmt(stmt ast.Stmt) error {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		return r.resolveExpr(stmt.Expr)
	case *ast.BlockStmt:
		return r.resolveBlockStmt(stmt)
	case *ast.VarAssignStmt:
		return r.resolveVarAssignStmt(stmt)
	case *ast.FuncDecStmt:
		return r.resolveFuncDecStmt(stmt)
	case *ast.ReturnStmt:
		return r.resolveReturnStmt(stmt)
	case *ast.IfStmt:
		return r.resolveIfStmt(stmt)
	case *ast.WhileStmt:
		return r.resolveWhileStmt(stmt)
		// case *ast.ForStmt:
		// 	r.resolveForStmt(stmt)
	}
	return nil
}

func (r *resolver) resolveBlockStmt(stmt *ast.BlockStmt) error {
	r.beginScope()
	defer r.endScope()
	for _, stmt := range stmt.Stmts {
		err := r.resolveStmt(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *resolver) resolveVarAssignStmt(stmt *ast.VarAssignStmt) error {
	if stmt.Op == ":=" {
		err := r.declare(stmt.Id.Name)
		if err != nil {
			return err
		}
		r.resolveExpr(stmt.Init)
		r.define(stmt.Id.Name)

	} else {
		r.resolveExpr(stmt.Init)
		r.resolveLocal(stmt.Id, stmt.Id.Name)
	}
	return nil
}

func (r *resolver) resolveFuncDecStmt(stmt *ast.FuncDecStmt) error {
	err := r.declare(stmt.Id.Name)
	if err != nil {
		return err
	}
	r.define(stmt.Id.Name)
	return r.resolveFunction(stmt, functionFuncType)
}

func (r *resolver) resolveFunction(funcDec *ast.FuncDecStmt, funcType functionType) error {

	enclosingFunc := r.currentFunc
	r.currentFunc = funcType

	r.beginScope()
	// TODO: review this
	defer func() {
		r.endScope()
		r.currentFunc = enclosingFunc
	}()

	for _, arg := range funcDec.Args {
		err := r.declare(arg.Name)
		if err != nil {
			return err
		}

		r.define(arg.Name)
	}
	return r.resolveStmt(funcDec.Body)
}

func (r *resolver) resolveReturnStmt(stmt *ast.ReturnStmt) error {
	if r.currentFunc == noneFuncType {
		return NewRuntimeError("can't return from top-level code")
	}
	if stmt.Arg != nil {
		r.resolveExpr(stmt.Arg)
	}
	return nil
}

func (r *resolver) resolveIfStmt(stmt *ast.IfStmt) error {
	err := r.resolveExpr(stmt.Test)
	if err != nil {
		return err
	}
	err = r.resolveStmt(stmt.Consequent)
	if err != nil {
		return err
	}

	if stmt.Alternate != nil {
		err = r.resolveStmt(stmt.Alternate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *resolver) resolveWhileStmt(stmt *ast.WhileStmt) error {
	r.resolveExpr(stmt.Test)
	err := r.resolveStmt(stmt.Body)
	if err != nil {
		return err
	}
	return nil
}

// Expressions
func (r *resolver) resolveExpr(expr ast.Expr) error {
	switch expr := expr.(type) {
	case *ast.IdentifierExpr:
		return r.resolveIdentifierExpr(expr)
	case *ast.BinaryExpr:
		return r.resolveBinaryExpr(expr)
	case *ast.LogicalExpr:
		return r.resolveLogicalExpr(expr)
	case *ast.CallExpr:
		return r.resolveCallExpr(expr)
	}
	return nil
}

func (r *resolver) resolveIdentifierExpr(expr *ast.IdentifierExpr) error {
	if !r.scopes.isEmpty() && r.scopes.peek().isDefined(expr.Name) && !r.scopes.peek().isInitialized(expr.Name) {
		return NewRuntimeError(fmt.Sprintf("can't resolve local variable %s in own initializer", expr.Name))
	}

	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *resolver) resolveBinaryExpr(expr *ast.BinaryExpr) error {
	err := r.resolveExpr(expr.Lhs)
	if err != nil {
		return err
	}
	return r.resolveExpr(expr.Rhs)
}
func (r *resolver) resolveLogicalExpr(expr *ast.LogicalExpr) error {
	err := r.resolveExpr(expr.Lhs)
	if err != nil {
		return err
	}
	return r.resolveExpr(expr.Rhs)
}
func (r *resolver) resolveCallExpr(expr *ast.CallExpr) error {
	err := r.resolveExpr(expr.Callee)
	if err != nil {
		return err
	}

	for _, arg := range expr.Args {
		err = r.resolveExpr(arg)
		if err != nil {
			return err
		}
	}
	return nil
}
