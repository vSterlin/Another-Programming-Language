package codegen

import (
	"fmt"
	"language/ast"
	"strings"
)

type CodeGenerator struct {
	imports []string
	indent  int
}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{

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
			// TODO: do it more efficiently
			lines := strings.Split(code, "\n")
			code = strings.Join(lines, "\n\t")
			code = "\t" + code
			main += code + "\n"
		}

	}

	main = fmt.Sprintf("int main() {\n%s\nreturn 0;\n}", main)

	res := cg.genImports() + funcs + main

	return res
}
