package interpreter

import (
	"fmt"
	"language/ast"

	"github.com/fatih/color"
)

type Caller interface {
	Call(i *Interpreter, args []any) any
}
type function struct {
	funcDef *ast.FuncDecStmt
	closure *Environment
}

func NewFunction(funcDef *ast.FuncDecStmt, closure *Environment) *function {
	return &function{funcDef: funcDef, closure: closure}
}

func (f *function) Call(i *Interpreter, args []any) any {
	env := NewEnvironment(f.closure)

	for i, arg := range f.funcDef.Args {
		env.Define(arg.Name, args[i])
	}

	retVal := i.evalBlockStmt(f.funcDef.Body, env)

	// unwrap return value
	if retObj, ok := retVal.(*ReturnValue); ok {
		return retObj.Value()
	} else {
		return nil
	}

}

func (f *function) String() string {
	return color.BlueString(fmt.Sprintf("<function %s>", f.funcDef.Id.Name))
}

type PrintFunction struct{ function }

func (p *PrintFunction) Call(i *Interpreter, args []any) any {
	fmt.Println(args...)
	return nil
}

func (p *PrintFunction) String() string {
	return color.BlueString("<function print>")
}

func NewGlobalFunctions() map[string]Caller {
	return map[string]Caller{
		"print": &PrintFunction{},
	}
}
