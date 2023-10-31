package main

import (
	"bufio"
	"fmt"

	"language/ast"
	"language/codegen"
	"language/interpreter"
	"language/lexer"
	"language/parser"
	"language/typechecker"
	"os"
)

func main() {

	compile(`
	func sum(a int, b int) int {
			return a + b
	}

	x := sum(1, 2)

	x()


	`)

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

func interpret(code string) {

	prog := buildAST(code)
	i := interpreter.NewInterpreter(prog)
	resolver := interpreter.NewResolver(i)
	err := resolver.ResolveProgram(prog)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = i.Interpret()
	if err != nil {
		fmt.Println(err)
		return
	}

	// for _, evaluatedStmt := range evaluatedProgram {
	// 	fmt.Println(evaluatedStmt)
	// }

}

func compile(code string) {
	prog := buildAST(code)

	tc := typechecker.NewTypeChecker()

	err := tc.Check(prog)

	if err != nil {
		fmt.Println(err)
		return
	}

	cg := codegen.NewCodeGenerator()
	output := cg.Gen(prog)
	writeToFile(output)
}

func repl() {
	for {
		fmt.Print(">> ")
		bufioReader := bufio.NewReader(os.Stdin)
		code, err := bufioReader.ReadString('\n')
		if err != nil {
			return
		}

		interpret(code)
	}
}

func writeToFile(code string) {
	os.WriteFile("build/out.cpp", []byte(code), 0644)
}
