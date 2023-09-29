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

func (i *Interpreter) Interpret() ([]any, error) {
	return i.evalProgram(i.program)
}

func (i *Interpreter) Resolve(expr ast.Expr, depth int) {
	i.locals[expr] = depth
}

// Expressions
func (i *Interpreter) evalExpr(expr ast.Expr) (any, error) {
	switch expr := expr.(type) {
	case *ast.NumberExpr:
		return i.evalNumberExpr(expr), nil
	case *ast.BinaryExpr:
		return i.evalBinaryExpr(expr), nil
	case *ast.LogicalExpr:
		return i.evalLogicalExpr(expr), nil
	case *ast.BooleanExpr:
		return i.evalBooleanExpr(expr), nil
	case *ast.StringExpr:
		return i.evalStringExpr(expr), nil
	case *ast.IdentifierExpr:
		return i.evalIdentifierExpr(expr)
	case *ast.CallExpr:
		return i.evalCallExpr(expr)
	case *ast.MemberExpr:
		return i.evalMemberExpr(expr)
	case *ast.ThisExpr:
		return i.evalThisExpr(expr), nil
	default:
		return nil, nil
	}
}

func (i *Interpreter) evalNumberExpr(expr *ast.NumberExpr) any   { return Number(expr.Val) }
func (i *Interpreter) evalBooleanExpr(expr *ast.BooleanExpr) any { return Boolean(expr.Val) }
func (i *Interpreter) evalStringExpr(expr *ast.StringExpr) any   { return String(expr.Val) }
func (i *Interpreter) evalIdentifierExpr(expr *ast.IdentifierExpr) (any, error) {
	val, err := i.lookUpVariable(expr.Name, expr)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (i *Interpreter) evalCallExpr(expr *ast.CallExpr) (any, error) {

	callExpr, _ := i.evalExpr(expr.Callee)
	callee, ok := callExpr.(Caller)

	if !ok {
		return nil, NewRuntimeError("the expression is not callable")
	}

	args := []any{}
	for _, arg := range expr.Args {
		evaluatedExpr, _ := i.evalExpr(arg)
		args = append(args, evaluatedExpr)
	}

	return callee.Call(i, args), nil
}

func (i *Interpreter) evalBinaryExpr(expr *ast.BinaryExpr) any {
	lhs, _ := i.evalExpr(expr.Lhs)
	rhs, _ := i.evalExpr(expr.Rhs)

	lhsNum, isLhsNum := lhs.(Number)
	rhsNum, isRhsNum := rhs.(Number)

	lhsStr, isLhsStr := lhs.(String)
	rhsStr, isRhsStr := rhs.(String)

	switch expr.Op {
	case "+":
		if isLhsNum && isRhsNum {
			return lhsNum + rhsNum
		} else if isLhsStr && isRhsStr {
			return lhsStr + rhsStr
		} else {
			return nil
		}
	case "-":
		return lhsNum - rhsNum
	case "*":
		return lhsNum * rhsNum
	case "/":
		return lhsNum / rhsNum
	case "%":
		return lhsNum % rhsNum
	case "**":
		return Number(math.Pow(float64(lhs.(Number)), float64(rhsNum)))
	case "<":
		return Boolean(lhsNum < rhsNum)
	case ">":
		return Boolean(lhsNum > rhsNum)
	case "<=":
		return Boolean(lhsNum <= rhsNum)
	case ">=":
		return Boolean(lhsNum >= rhsNum)
	case "==":
		return Boolean((lhs == rhs))
	case "!=":
		return Boolean((lhs != rhs))
	default:
		return nil
	}
}

func (i *Interpreter) evalLogicalExpr(expr *ast.LogicalExpr) any {
	lhs, _ := i.evalExpr(expr.Lhs)

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

	rhs, _ := i.evalExpr(expr.Rhs)
	rhsBool := rhs.(Boolean)

	return rhsBool
}

func (i *Interpreter) evalMemberExpr(expr *ast.MemberExpr) (any, error) {

	obj, _ := i.evalExpr(expr.Obj)
	instance, ok := obj.(*Instance)
	if !ok {
		return nil, NewRuntimeError("the object is not an instance")
	}

	// for now handle identifier exp only
	field, err := instance.Get(expr.Prop.(*ast.IdentifierExpr).Name)

	return field, err
}

func (i *Interpreter) evalThisExpr(expr *ast.ThisExpr) any {
	this, err := i.lookUpVariable("this", expr)
	if err != nil {
		fmt.Println("error with this expression")
	}
	return this
}

// Statements
func (i *Interpreter) evalStmt(stmt ast.Stmt) (any, error) {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		return i.evalExpr(stmt.Expr)
	case *ast.VarAssignStmt:
		return i.evalVarAssignStmt(stmt), nil
	case *ast.FuncDecStmt:
		return i.evalFuncDecStmt(stmt), nil
	case *ast.BlockStmt:
		// NewEnvironment(i.env) creates a new environment with the current environment as its parent
		return i.evalBlockStmt(stmt, NewEnvironment(i.env)), nil
	case *ast.ReturnStmt:
		return i.evalReturnStmt(stmt), nil
	case *ast.IfStmt:
		return i.evalIfStmt(stmt), nil
	case *ast.WhileStmt:
		return i.evalWhileStmt(stmt), nil
	case *ast.ClassDecStmt:
		return i.evalClassDecStmt(stmt), nil
	case *ast.SetStmt:
		return i.evalSetStmt(stmt)
	default:
		return nil, NewRuntimeError("unknown statement")
	}
}

func (i *Interpreter) evalIfStmt(stmt *ast.IfStmt) any {
	var retVal any
	test, _ := i.evalExpr(stmt.Test)
	if test.(Boolean) {
		retVal, _ = i.evalStmt(stmt.Consequent.(*ast.BlockStmt))
	} else {
		retVal, _ = i.evalStmt(stmt.Alternate)
	}
	if retVal != nil {
		return retVal
	} else {
		return nil
	}
}

func (i *Interpreter) evalWhileStmt(stmt *ast.WhileStmt) any {
	test, _ := i.evalExpr(stmt.Test)
	for test.(Boolean) {
		retVal, _ := i.evalStmt(stmt.Body.(*ast.BlockStmt))
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
		retVal, _ := i.evalStmt(stmt)
		if retVal != nil {
			return retVal
		}
	}
	return nil
}

func (i *Interpreter) evalVarAssignStmt(stmt *ast.VarAssignStmt) error {
	varName := stmt.Id.Name
	varValue, _ := i.evalExpr(stmt.Init)
	if stmt.Op == ":=" {
		i.env.Define(varName, varValue)
		return nil
	} else { // "="
		distance, ok := i.locals[stmt.Id]
		if ok {
			return i.env.AssignAt(distance, varName, varValue)
		} else {
			return i.env.Assign(varName, varValue)
		}
	}
}

func (i *Interpreter) evalFuncDecStmt(stmt *ast.FuncDecStmt) any {
	fn := NewFunction(stmt, i.env)
	i.env.Define(stmt.Id.Name, fn)
	return nil
}

func (i *Interpreter) evalReturnStmt(stmt *ast.ReturnStmt) any {
	evaluatedExpr, _ := i.evalExpr(stmt.Arg)
	return NewReturnValue(evaluatedExpr)
}

func (i *Interpreter) evalSetStmt(stmt *ast.SetStmt) (any, error) {
	lhs, _ := i.evalExpr(stmt.Lhs)
	instance, ok := lhs.(*Instance)
	if !ok {
		return nil, NewRuntimeError("the object is not an instance")
	}
	evaluatedExpr, _ := i.evalExpr(stmt.Val)
	instance.Set(stmt.Name, evaluatedExpr)
	return nil, nil
}

func (i *Interpreter) evalProgram(p *ast.Program) ([]any, error) {
	stmts := []any{}
	for _, stmt := range p.Stmts {
		evaluatedStmt, err := i.evalStmt(stmt)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, evaluatedStmt)
	}
	return stmts, nil
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
