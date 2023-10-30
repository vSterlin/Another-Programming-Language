package codegen

import (
	"fmt"
)

func (cg *CodeGenerator) genImports() string {
	importStr := ""
	for _, imp := range cg.imports {
		importStr += fmt.Sprintf("#include <%s.h>\n", imp)
	}

	importStr += "\n"
	return importStr
}

func (cg *CodeGenerator) genTabs() string {
	tabs := ""
	for i := 0; i < cg.identLevel; i++ {
		tabs += "\t"
	}
	return tabs
}

func cType(t string) string {
	switch t {
	case "int":
		return "int"
	case "string":
		return "char*"
	default:
		return ""
	}
}
