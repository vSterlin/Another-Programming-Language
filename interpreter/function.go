package interpreter

import (
	"fmt"
	"language/ast"

	"github.com/fatih/color"
)

type Function interface {
	Call(i *Interpreter, args []any) any
}
type function struct {
	FuncDef *ast.FuncDecStmt
}

func NewFunction(funcDef *ast.FuncDecStmt) *function {
	return &function{FuncDef: funcDef}
}

func (f *function) Call(i *Interpreter, args []any) any {
	env := NewEnvironment(i.env)

	for i, arg := range f.FuncDef.Args {
		env.Define(arg.Name, args[i])
	}

	retVal := i.evalBlockStmt(f.FuncDef.Body, env)

	// unwrap return value
	if retObj, ok := retVal.(*ReturnValue); ok {
		return retObj.Value()
	} else {
		return nil
	}

}

func (f *function) String() string {
	return color.BlueString(fmt.Sprintf("<function %s>", f.FuncDef.Id.Name))
}

type PrintFunction struct{ function }

func (p *PrintFunction) Call(i *Interpreter, args []any) any {
	fmt.Println(args...)
	return nil
}

func (p *PrintFunction) String() string {
	return color.BlueString("<function print>")
}

func NewGlobalFunctions() map[string]Function {
	return map[string]Function{
		"print": &PrintFunction{},
	}
}
