package interpreter

import (
	"fmt"
	"language/ast"
	"math"
)

type Interpreter struct {
	program *ast.Program
	env     *Environment
}

func NewInterpreter(program *ast.Program) *Interpreter {
	return &Interpreter{program: program, env: NewEnvironment(nil)}
}

func (i *Interpreter) Interpret() []any {
	return i.evaluateProgram(i.program)
}

// Expressions
func (i *Interpreter) evalExpr(expr ast.Expr) any {
	switch expr := expr.(type) {
	case *ast.NumberExpr:
		return i.evalNumberExpr(expr)
	case *ast.BinaryExpr:
		return i.evalBinaryExpr(expr)
	case *ast.LogicalExpr:
		return i.evalLogicalExpr(expr)
	case *ast.BooleanExpr:
		return i.evalBooleanExpr(expr)
	case *ast.StringExpr:
		return i.evalStringExpr(expr)
	case *ast.IdentifierExpr:
		return i.evalIdentifierExpr(expr)
	case *ast.CallExpr:
		return i.evalCallExpr(expr)
	default:
		return nil
	}
}

func (i *Interpreter) evalNumberExpr(expr *ast.NumberExpr) any   { return expr.Val }
func (i *Interpreter) evalBooleanExpr(expr *ast.BooleanExpr) any { return expr.Val }
func (i *Interpreter) evalStringExpr(expr *ast.StringExpr) any   { return expr.Val }
func (i *Interpreter) evalIdentifierExpr(expr *ast.IdentifierExpr) any {
	val, err := i.env.Get(expr.Name)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return val

}

// for now only handle print expressions
// maybe print should be a statement
func (i *Interpreter) evalCallExpr(expr *ast.CallExpr) any {

	switch expr.Callee.Name {
	case "print":
		fmt.Println(i.evalExpr(expr.Args[0]))
		return nil
	default:
		return nil
	}
}

func (i *Interpreter) evalBinaryExpr(expr *ast.BinaryExpr) any {
	lhs := i.evalExpr(expr.Lhs)
	rhs := i.evalExpr(expr.Rhs)

	switch expr.Op {
	case "+":
		return lhs.(int) + rhs.(int)
	case "-":
		return lhs.(int) - rhs.(int)
	case "*":
		return lhs.(int) * rhs.(int)
	case "/":
		return lhs.(int) / rhs.(int)
	case "**":
		return int(math.Pow(float64(lhs.(int)), float64(rhs.(int))))
	case "<":
		return lhs.(int) < rhs.(int)
	case ">":
		return lhs.(int) > rhs.(int)
	case "<=":
		return lhs.(int) <= rhs.(int)
	case ">=":
		return lhs.(int) >= rhs.(int)
	case "==":
		return lhs == rhs
	case "!=":
		return lhs != rhs
	default:
		return nil
	}
}

func (i *Interpreter) evalLogicalExpr(expr *ast.LogicalExpr) any {
	lhs := i.evalExpr(expr.Lhs)
	rhs := i.evalExpr(expr.Rhs)

	switch expr.Op {
	case "&&":
		return lhs.(bool) && rhs.(bool)
	case "||":
		return lhs.(bool) || rhs.(bool)
	default:
		return nil
	}
}

// Statements
func (i *Interpreter) evalStmt(stmt ast.Stmt) any {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		return i.evalExpr(stmt.Expr)
	case *ast.VarAssignStmt:
		return i.evalVarAssignStmt(stmt)
	case *ast.BlockStmt:
		// NewEnvironment(i.env) creates a new environment with the current environment as its parent
		return i.evalBlockStmt(stmt, NewEnvironment(i.env))
	default:
		return nil
	}
}

func (i *Interpreter) evalBlockStmt(stmt *ast.BlockStmt, env *Environment) any {
	currEnv := i.env
	i.env = env
	defer (func() { i.env = currEnv })()

	stmts := []any{}
	for _, stmt := range stmt.Stmts {

		evaluatedStmt := i.evalStmt(stmt)
		stmts = append(stmts, evaluatedStmt)
	}

	return stmts
}

func (i *Interpreter) evalVarAssignStmt(stmt *ast.VarAssignStmt) any {
	varName := stmt.Id.Name
	varValue := i.evalExpr(stmt.Init)
	if stmt.Op == ":=" {
		i.env.Define(varName, varValue)
	} else { // "="
		err := i.env.Assign(varName, varValue)
		if err != nil {
			fmt.Println(err)
		}
	}

	return varValue
}

func (i *Interpreter) evaluateProgram(p *ast.Program) []any {
	stmts := []any{}
	for _, stmt := range p.Stmts {
		evaluatedStmt := i.evalStmt(stmt)
		stmts = append(stmts, evaluatedStmt)
	}
	return stmts
}

// Environment
type Environment struct {
	values map[string]any
	parent *Environment
}

func NewEnvironment(parent *Environment) *Environment {
	globalValues := map[string]any{"VERSION": "0.0.1"}
	return &Environment{values: globalValues, parent: parent}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Assign(name string, value any) error {
	_, ok := e.values[name]
	if ok {
		e.values[name] = value
		return nil
	} else {
		return NewRuntimeError("undefined variable: " + name)
	}
}

func (e *Environment) Get(name string) (any, error) {
	_, ok := e.values[name]
	if ok {
		return e.values[name], nil
	} else {
		if e.parent != nil {
			return e.parent.Get(name)
		} else {
			return nil, NewRuntimeError("undefined variable: " + name)
		}
	}
}
