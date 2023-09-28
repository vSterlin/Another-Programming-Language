package interpreter

import (
	"fmt"
	"language/ast"
	"math"
)

type Interpreter struct {
	program *ast.Program
	env     *Environment
	locals  map[ast.Expr]int
}

func NewInterpreter(program *ast.Program) *Interpreter {

	globalEnv := NewEnvironment(nil)
	globalFuncs := NewGlobalFunctions()
	for name, fn := range globalFuncs {
		globalEnv.Define(name, fn)
	}

	globalEnv.Define("VERSION", "0.0.1")

	locals := map[ast.Expr]int{}
	return &Interpreter{
		program: program,
		env:     globalEnv,
		locals:  locals,
	}
}

func (i *Interpreter) Interpret() []any {
	return i.evalProgram(i.program)
}

func (i *Interpreter) Resolve(expr ast.Expr, depth int) {
	i.locals[expr] = depth
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
	case *ast.MemberExpr:
		return i.evalMemberExpr(expr)
	case *ast.ThisExpr:
		return i.evalThisExpr(expr)
	default:
		return nil
	}
}

func (i *Interpreter) evalNumberExpr(expr *ast.NumberExpr) any   { return Number(expr.Val) }
func (i *Interpreter) evalBooleanExpr(expr *ast.BooleanExpr) any { return Boolean(expr.Val) }
func (i *Interpreter) evalStringExpr(expr *ast.StringExpr) any   { return String(expr.Val) }
func (i *Interpreter) evalIdentifierExpr(expr *ast.IdentifierExpr) any {
	val, err := i.lookUpVariable(expr.Name, expr)
	if err != nil {
		fmt.Println(err)
	}
	return val
}

func (i *Interpreter) evalCallExpr(expr *ast.CallExpr) any {

	callee, ok := i.evalExpr(expr.Callee).(Caller)

	if !ok {
		fmt.Println("the expressions is not callable")
	}

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
		switch lhs := lhs.(type) {

		case String:
			return (lhs) + rhs.(String)
		case Number:
			return (lhs) + rhs.(Number)
			// TODO: handle errors
		default:
			return nil
		}

	case "-":
		return (lhs).(Number) - rhs.(Number)
	case "*":
		return (lhs).(Number) * rhs.(Number)
	case "/":
		return (lhs).(Number) / rhs.(Number)
	case "%":
		return (lhs).(Number) % rhs.(Number)
	case "**":
		return Number(math.Pow(float64(lhs.(Number)), float64(rhs.(Number))))
	case "<":
		return Boolean((lhs).(Number) < rhs.(Number))
	case ">":
		return Boolean((lhs).(Number) > rhs.(Number))
	case "<=":
		return Boolean((lhs).(Number) <= rhs.(Number))
	case ">=":
		return Boolean((lhs).(Number) >= rhs.(Number))
	case "==":
		return Boolean((lhs == rhs))
	case "!=":
		return Boolean((lhs != rhs))
	default:
		return nil
	}
}

func (i *Interpreter) evalLogicalExpr(expr *ast.LogicalExpr) any {
	lhs := i.evalExpr(expr.Lhs)

	lhsBool := lhs.(Boolean)

	if expr.Op == ast.OR {
		if lhsBool {
			return Boolean(true)
		}
	}
	if expr.Op == ast.AND {
		if !lhsBool {
			return Boolean(false)
		}
	}

	rhs := i.evalExpr(expr.Rhs)
	rhsBool := rhs.(Boolean)

	return rhsBool
}

func (i *Interpreter) evalMemberExpr(expr *ast.MemberExpr) any {

	obj := i.evalExpr(expr.Obj)
	instance, ok := obj.(*Instance)
	if !ok {
		fmt.Println("the object is not an instance")
		return nil
	}

	// for now handle identifier exp only
	field, err := instance.Get(expr.Prop.(*ast.IdentifierExpr).Name)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return field
}

func (i *Interpreter) evalThisExpr(expr *ast.ThisExpr) any {
	this, err := i.lookUpVariable("this", expr)
	if err != nil {
		fmt.Println("error with this expression")
	}
	return this
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
	case *ast.ClassDecStmt:
		return i.evalClassDecStmt(stmt)
	case *ast.SetStmt:
		return i.evalSetStmt(stmt)
	default:
		return nil
	}
}

func (i *Interpreter) evalIfStmt(stmt *ast.IfStmt) any {
	var retVal any
	if i.evalExpr(stmt.Test).(Boolean) {
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
	for i.evalExpr(stmt.Test).(Boolean) {
		retVal := i.evalStmt(stmt.Body.(*ast.BlockStmt))
		if retVal != nil {
			return retVal
		}
	}
	return nil
}

func (i *Interpreter) evalClassDecStmt(stmt *ast.ClassDecStmt) any {
	i.env.Define(stmt.Id.Name, nil)

	methods := map[string]*function{}
	for _, method := range stmt.Methods {
		fn := NewFunction(method, i.env)
		methods[method.Id.Name] = fn
	}

	class := NewClass(stmt.Id.Name, methods)
	i.env.Assign(stmt.Id.Name, class)
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
		distance, ok := i.locals[stmt.Id]
		if ok {
			err := i.env.AssignAt(distance, varName, varValue)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			err := i.env.Assign(varName, varValue)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return nil
}

func (i *Interpreter) evalFuncDecStmt(stmt *ast.FuncDecStmt) any {
	fn := NewFunction(stmt, i.env)
	i.env.Define(stmt.Id.Name, fn)
	return nil
}

func (i *Interpreter) evalReturnStmt(stmt *ast.ReturnStmt) any {
	return NewReturnValue(i.evalExpr(stmt.Arg))
}

func (i *Interpreter) evalSetStmt(stmt *ast.SetStmt) any {
	lhs := i.evalExpr(stmt.Lhs)
	instance, ok := lhs.(*Instance)
	if !ok {
		fmt.Println("the object is not an instance")
	}
	instance.Set(stmt.Name, i.evalExpr(stmt.Val))
	return nil
}

func (i *Interpreter) evalProgram(p *ast.Program) []any {
	stmts := []any{}
	for _, stmt := range p.Stmts {
		evaluatedStmt := i.evalStmt(stmt)
		stmts = append(stmts, evaluatedStmt)
	}
	return stmts
}

func (i *Interpreter) lookUpVariable(name string, expr ast.Expr) (any, error) {
	distance, ok := i.locals[expr]
	if ok {
		return i.env.GetAt(distance, name)
	} else {
		// get from global env
		return i.env.Get(name)
	}
}
