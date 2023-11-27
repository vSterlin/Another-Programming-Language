package codegen

import (
	"fmt"
	"language/ast"
)

type CodeGenerator struct {
	env     *Env
	imports []string
	indent  int
}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		env: NewEnv(nil),

		imports: []string{
			"iostream",
			"string"},

		indent: 0,
	}
}

func (cg *CodeGenerator) Gen(prog *ast.Program) string {
	funcs := ""
	main := ""

	for _, stmt := range prog.Stmts {
		code, _ := cg.genStmt(stmt)
		if _, ok := stmt.(*ast.FuncDecStmt); ok {
			funcs += code + "\n"
		} else {
			main += code + "\n"
		}

	}

	main = fmt.Sprintf("int main() {\n%s\nreturn 0;\n}", main)

	res := cg.genImports() + funcs + main

	return res
}
