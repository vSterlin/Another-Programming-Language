package codegen

import (
	"language/ast"
)

type CodeGenerator struct {
	env     *Env
	program string
}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		env:     NewEnv(nil),
		program: "",
	}
}

func (cg *CodeGenerator) Gen(prog *ast.Program) string {
	p := ""
	p += "#include <stdio.h>\n"
	p += "#include <stdlib.h>\n"

	for _, stmt := range prog.Stmts {
		code, _ := cg.genStmt(stmt)
		p += code
	}

	return p
}
