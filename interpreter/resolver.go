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
		fmt.Println("RUNTIME ERROR: variable already declared in this scope")
		return NewRuntimeError(fmt.Sprintf("variable %s already declared in this scope", name))
	}

	(*scope)[name] = false
	return nil
}

func (r *resolver) define(name string) {
	if r.scopes.isEmpty() {
		fmt.Printf("\"lol\": %v\n", "lol")
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
func (r *resolver) ResolveProgram(p *ast.Program) {
	for _, stmt := range p.Stmts {
		r.resolveStmt(stmt)
	}
}

// TODO: Might consider using a visitor pattern here
func (r *resolver) resolveStmt(stmt ast.Stmt) {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		r.resolveExpr(stmt.Expr)
	case *ast.BlockStmt:
		r.resolveBlockStmt(stmt)
	case *ast.VarAssignStmt:
		r.resolveVarAssignStmt(stmt)
	case *ast.FuncDecStmt:
		r.resolveFuncDecStmt(stmt)
	case *ast.ReturnStmt:
		r.resolveReturnStmt(stmt)
	case *ast.IfStmt:
		r.resolveIfStmt(stmt)
	case *ast.WhileStmt:
		r.resolveWhileStmt(stmt)
		// case *ast.ForStmt:
		// 	r.resolveForStmt(stmt)

	}
}

func (r *resolver) resolveBlockStmt(stmt *ast.BlockStmt) {
	r.beginScope()
	defer r.endScope()
	for _, stmt := range stmt.Stmts {
		r.resolveStmt(stmt)
	}
}

func (r *resolver) resolveVarAssignStmt(stmt *ast.VarAssignStmt) {
	if stmt.Op == ":=" {
		r.declare(stmt.Id.Name)
		r.resolveExpr(stmt.Init)
		r.define(stmt.Id.Name)
	} else {
		r.resolveExpr(stmt.Init)
		r.resolveLocal(stmt.Id, stmt.Id.Name)
	}
}

func (r *resolver) resolveFuncDecStmt(stmt *ast.FuncDecStmt) {
	r.declare(stmt.Id.Name)
	r.define(stmt.Id.Name)

	r.resolveFunction(stmt, functionFuncType)
}

func (r *resolver) resolveFunction(funcDec *ast.FuncDecStmt, funcType functionType) {

	enclosingFunc := r.currentFunc
	r.currentFunc = funcType

	r.beginScope()
	// TODO: review this
	defer func() {
		r.endScope()
		r.currentFunc = enclosingFunc
	}()

	for _, arg := range funcDec.Args {
		r.declare(arg.Name)
		r.define(arg.Name)
	}
	r.resolveStmt(funcDec.Body)
}

func (r *resolver) resolveReturnStmt(stmt *ast.ReturnStmt) error {
	if r.currentFunc == noneFuncType {
		fmt.Println("RUNTIME ERROR: can't return from top-level code")
		return NewRuntimeError("can't return from top-level code")
	}
	if stmt.Arg != nil {
		r.resolveExpr(stmt.Arg)
	}
	return nil
}

func (r *resolver) resolveIfStmt(stmt *ast.IfStmt) {
	r.resolveExpr(stmt.Test)
	r.resolveStmt(stmt.Consequent)
	if stmt.Alternate != nil {
		r.resolveStmt(stmt.Alternate)
	}
}

func (r *resolver) resolveWhileStmt(stmt *ast.WhileStmt) {
	r.resolveExpr(stmt.Test)
	r.resolveStmt(stmt.Body)
}

// Expressions
func (r *resolver) resolveExpr(expr ast.Expr) error {
	switch expr := expr.(type) {
	case *ast.IdentifierExpr:
		return r.resolveIdentifierExpr(expr)
	case *ast.BinaryExpr:
		r.resolveBinaryExpr(expr)
	case *ast.LogicalExpr:
		r.resolveLogicalExpr(expr)
	case *ast.CallExpr:
		r.resolveCallExpr(expr)
	}
	return nil
}

func (r *resolver) resolveIdentifierExpr(expr *ast.IdentifierExpr) error {
	if !r.scopes.isEmpty() && r.scopes.peek().isDefined(expr.Name) && !r.scopes.peek().isInitialized(expr.Name) {
		fmt.Println("RUNTIME ERROR: can't resolve local variable in own initializer")
		return NewRuntimeError(fmt.Sprintf("can't resolve local variable %s in own initializer", expr.Name))
	}

	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *resolver) resolveBinaryExpr(expr *ast.BinaryExpr) {
	r.resolveExpr(expr.Lhs)
	r.resolveExpr(expr.Rhs)
}
func (r *resolver) resolveLogicalExpr(expr *ast.LogicalExpr) {
	r.resolveExpr(expr.Lhs)
	r.resolveExpr(expr.Rhs)
}
func (r *resolver) resolveCallExpr(expr *ast.CallExpr) {
	r.resolveExpr(expr.Callee)

	for _, arg := range expr.Args {
		r.resolveExpr(arg)
	}
}
