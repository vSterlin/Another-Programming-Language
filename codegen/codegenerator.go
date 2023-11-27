package codegen

import (
	"fmt"
	"language/ast"
)

type CodeGenerator struct {
	env        *Env
	imports    []string
	indent     int
	mainBuf    string
	codeBuf    string
	isInGlobal bool
}

func NewCodeGenerator() *CodeGenerator {

	return &CodeGenerator{
		env: NewEnv(nil),

		imports: []string{
			"iostream",
			"string"},

		indent:     0,
		mainBuf:    "",
		codeBuf:    "",
		isInGlobal: true,
	}
}

func (cg *CodeGenerator) Gen(prog *ast.Program) string {
	p := ""

	for _, stmt := range prog.Stmts {
		err := cg.genStmt(stmt)
		if err != nil {
			panic(err)
		}
	}

	p = cg.genImports() + p

	fmt.Printf("cg.codeBuf: %v\n", cg.codeBuf)

	fmt.Printf("cg.mainBuf: %v\n", cg.mainBuf)

	return p
}

func (cg *CodeGenerator) write(code string) {

	if cg.isInGlobal {
		cg.mainBuf += code
	} else {
		cg.codeBuf += code
	}
}
