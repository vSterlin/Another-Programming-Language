package main

import (
	"bufio"
	"encoding/json"
	"fmt"

	"language/ast"
	"language/codegen"
	"language/interpreter"
	"language/lexer"
	"language/parser"
	"os"
)

func main() {
	interpret(`
		func count(n){
			if n > 1 {
				count(n - 1)
			}
			print(n)
		}
		count(10)
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
		jsonStr, _ := json.MarshalIndent(prog, "", "  ")
		fmt.Println(string(jsonStr))
	}
	return prog
}

func compileToJS(code string) {
	prog := buildAST(code)
	fmt.Println(prog)
	cg := codegen.NewJavascriptCodeGenerator()
	output := cg.Generate(prog)
	writeToFile(output)
}

func interpret(code string) {

	prog := buildAST(code)
	i := interpreter.NewInterpreter(prog)
	// evaluatedProgram :=
	i.Interpret()

	// for _, evaluatedStmt := range evaluatedProgram {
	// 	fmt.Println(evaluatedStmt)
	// }

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
	os.WriteFile("build/main.js", []byte(code), 0644)
}
