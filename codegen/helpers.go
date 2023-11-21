package codegen

import (
	"fmt"
	"language/ast"
)

func (cg *CodeGenerator) genImports() string {
	importStr := ""
	for _, imp := range cg.imports {
		importStr += fmt.Sprintf("#include <%s>\n", imp)
	}

	importStr += "\n"
	return importStr
}

func (cg *CodeGenerator) genTabs() string {
	tabs := ""
	for i := 0; i < cg.indent; i++ {
		tabs += "\t"
	}
	return tabs
}

func cType(t string) string {
	switch t {
	case "int", "number":
		return "int"
	case "string":
		return "std::string"
	default:
		return ""
	}
}

// TODO: review
func cTypeFromAst(typeNode *ast.TypeExpr) string {
	switch t := typeNode.Type.(type) {
	case *ast.IdentifierExpr:
		return cType(t.Name)
	default:
		return "auto"
	}

}
