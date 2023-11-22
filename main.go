package main

import (
	"fmt"

	"language/ast"
	"language/codegen"
	"language/lexer"
	"language/parser"
	"language/typechecker"
	"os"
)

func main() {

	compile()

}

var PRINT_AST = false

func buildAST(code string) *ast.Program {
	l := lexer.NewLexer(code)
	tokens, _ := l.GetTokens()
	// fmt.Println(tokens)
	p := parser.NewParser(tokens)
	prog, err := p.ParseProgram()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if PRINT_AST {
		fmt.Println(prog)
		// jsonStr, _ := json.MarshalIndent(prog, "", "  ")
		// fmt.Println(string(jsonStr))
	}
	return prog
}

func compile() {

	// read source.vs file and pass it to the lexer
	code, err := os.ReadFile("./source.vs")
	if err != nil {
		fmt.Println(err)
		return
	}

	prog := buildAST(string(code))

	tc := typechecker.NewTypeChecker()

	err = tc.Check(prog)

	if err != nil {
		fmt.Println(err)
		return
	}

	cg := codegen.NewCodeGenerator()
	output := cg.Gen(prog)
	writeToFile(output)
}

func writeToFile(code string) {
	os.WriteFile("build/out.cpp", []byte(code), 0644)
}
