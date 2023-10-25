package main

import (
	"bufio"
	"fmt"

	"language/ast"
	"language/codegen"
	"language/interpreter"
	"language/lexer"
	"language/parser"
	"os"
)

func main() {

	compileToLLVM(`


	func a() int {
		if true {
			1 + 1
			if true {
				10 * 10
			} else {
				10 / 10
			}
			return 99999999 - 10
		} else {
			return 2 - 2
		}
		return 1001200 * 999
	}

	func b() int {
		return 10
	}

	x := a()
	z := b()
`)

	// compileToLLVM(`
	// 	if true {
	// 		if true {
	// 			10 * 10
	// 		} else {
	// 			10 / 10
	// 		}
	// 		99999999 - 10
	// 	} else {
	// 		2 - 2
	// 	}

	// 	10000 * 999
	// `)
}

var PRINT_AST = true

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

func compileToLLVM(code string) {
	prog := buildAST(code)
	cg := codegen.NewLLVMCodeGenerator()
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
	os.WriteFile("build/out.ll", []byte(code), 0644)
}
