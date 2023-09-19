package interpreter

import (
	"fmt"
	"language/ast"
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

	i.evalBlockStmt(f.FuncDef.Body, env)
	return nil

}

type PrintFunction struct{}

func (p *PrintFunction) Call(i *Interpreter, args []any) any {
	fmt.Println(args...)
	return nil
}

func NewGlobalFunctions() map[string]Function {
	return map[string]Function{
		"print": &PrintFunction{},
	}
}
