package codegen

import (
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

	funcs := strings.Builder{}
	main := strings.Builder{}

	for _, stmt := range prog.Stmts {
		code, _ := cg.genStmt(stmt)
		if _, ok := stmt.(*ast.FuncDecStmt); ok {
			funcs.WriteString(code + "\n")
		} else {
			// TODO: do it more efficiently
			lines := strings.Split(code, "\n")
			code = strings.Join(lines, "\n\t")
			code = "\t" + code
			main.WriteString(code + "\n")
		}

	}

	res := strings.Builder{}
	res.WriteString(cg.genImports())
	res.WriteString(funcs.String())
	res.WriteString("int main() {\n")
	res.WriteString(main.String())
	res.WriteString("\nreturn 0;\n}")

	return res.String()

}
