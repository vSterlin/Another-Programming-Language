package codegen

import (
	"fmt"
	"language/ast"
	"strings"
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
	case "boolean":
		return "bool"
	case "void":
		return "void"
	default:
		return ""
	}
}

func cTypeFromAst(typeNode *ast.TypeExpr) string {
	switch t := typeNode.Type.(type) {
	case *ast.IdentifierExpr:
		return cType(t.Name)
	case *ast.FuncTypeExpr:
		return cFuncTypeFromAst(t)
	default:
		return "auto"
	}

}

func cFuncTypeFromAst(typeNode *ast.FuncTypeExpr) string {
	args := []string{}
	for _, arg := range typeNode.Args {
		argStr := cTypeFromAst(arg)
		args = append(args, argStr)
	}

	retType := cTypeFromAst(typeNode.ReturnType)

	argsStr := fmt.Sprintf("std::function<%s(%s)>", retType, strings.Join(args, ", "))

	return argsStr

}

func inferFromAstNode(node ast.Expr) string {
	switch t := node.(type) {
	case *ast.NumberExpr:
		return Number
	case *ast.StringExpr:
		return String
	case *ast.BooleanExpr:
		return Bool
	case *ast.ArrowFunc:
		args := []*ast.TypeExpr{}
		for _, arg := range t.Args {
			args = append(args, arg.Type)
		}
		funcType := &ast.FuncTypeExpr{
			Args:       args,
			ReturnType: t.ReturnType,
		}
		return cFuncTypeFromAst(funcType)
	}

	return ""

}

const (
	Number = "int"
	String = "std::string"
	Bool   = "bool"
)
