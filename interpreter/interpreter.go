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

	globalEnv := NewEnvironment(nil)
	globalFuncs := NewGlobalFunctions()
	for name, fn := range globalFuncs {
		globalEnv.Define(name, fn)
	}

	globalEnv.Define("VERSION", "0.0.1")
	return &Interpreter{program: program, env: globalEnv}
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

func (i *Interpreter) evalCallExpr(expr *ast.CallExpr) any {

	callee := i.evalIdentifierExpr(expr.Callee).(Function)

	args := []any{}
	for _, arg := range expr.Args {
		args = append(args, i.evalExpr(arg))
	}

	return callee.Call(i, args)
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
	case "%":
		return lhs.(int) % rhs.(int)
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

	lhsBool := lhs.(bool)

	if expr.Op == ast.OR {
		if lhsBool {
			return true
		}
	}
	if expr.Op == ast.AND {
		if !lhsBool {
			return false
		}
	}

	rhs := i.evalExpr(expr.Rhs)
	return rhs
}

// Statements
func (i *Interpreter) evalStmt(stmt ast.Stmt) any {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		return i.evalExpr(stmt.Expr)
	case *ast.VarAssignStmt:
		return i.evalVarAssignStmt(stmt)
	case *ast.FuncDecStmt:
		return i.evalFuncDecStmt(stmt)
	case *ast.BlockStmt:
		// NewEnvironment(i.env) creates a new environment with the current environment as its parent
		return i.evalBlockStmt(stmt, NewEnvironment(i.env))
	case *ast.ReturnStmt:
		return i.evalReturnStmt(stmt)
	case *ast.IfStmt:
		return i.evalIfStmt(stmt)
	case *ast.WhileStmt:
		return i.evalWhileStmt(stmt)
	default:
		return nil
	}
}

func (i *Interpreter) evalIfStmt(stmt *ast.IfStmt) any {
	var retVal any
	if i.evalExpr(stmt.Test).(bool) {
		retVal = i.evalStmt(stmt.Consequent.(*ast.BlockStmt))
	} else {
		retVal = i.evalStmt(stmt.Alternate)
	}
	if retVal != nil {
		return retVal
	} else {
		return nil
	}
}

func (i *Interpreter) evalWhileStmt(stmt *ast.WhileStmt) any {
	for i.evalExpr(stmt.Test).(bool) {
		retVal := i.evalStmt(stmt.Body.(*ast.BlockStmt))
		if retVal != nil {
			return retVal
		}
	}
	return nil
}

func (i *Interpreter) evalBlockStmt(stmt *ast.BlockStmt, env *Environment) any {
	currEnv := i.env
	i.env = env
	defer (func() { i.env = currEnv })()

	// stmts := []any{}
	for _, stmt := range stmt.Stmts {
		retVal := i.evalStmt(stmt)
		if retVal != nil {
			return retVal
		}
	}
	return nil
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

	return nil
}

func (i *Interpreter) evalFuncDecStmt(stmt *ast.FuncDecStmt) any {
	fn := NewFunction(stmt)
	i.env.Define(stmt.Id.Name, fn)
	return nil
}

func (i *Interpreter) evalReturnStmt(stmt *ast.ReturnStmt) any {
	return i.evalExpr(stmt.Arg)
}

func (i *Interpreter) evaluateProgram(p *ast.Program) []any {
	stmts := []any{}
	for _, stmt := range p.Stmts {
		evaluatedStmt := i.evalStmt(stmt)
		stmts = append(stmts, evaluatedStmt)
	}
	return stmts
}
