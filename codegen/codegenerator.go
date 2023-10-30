package codegen

import (
	"language/ast"
)

type CodeGenerator struct {
	env        *Env
	imports    []string
	identLevel int
}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		env: NewEnv(nil),

		imports: []string{
			"sdio",
			"stdlib"},

		identLevel: 0,
	}
}

func (cg *CodeGenerator) Gen(prog *ast.Program) string {
	p := ""

	for _, stmt := range prog.Stmts {
		code, _ := cg.genStmt(stmt)
		p += code
	}

	p = cg.genImports() + p

	return p
}
