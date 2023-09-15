package main

import (
	"fmt"
	"language/codegen"
	"language/lexer"
	"language/parser"
	"os"
)

func main() {
	l := lexer.NewLexer(`
	func add(a, b) {
		a + b
	}

	add(1, 2)
	`)
	tokens, _ := l.GetTokens()
	fmt.Println(tokens)
	p := parser.NewParser(tokens)

	prog := p.ParseProgram()
	// fmt.Println(prog)

	cg := codegen.NewJavascriptCodeGenerator()
	code := cg.Generate(prog)

	writeToFile(code)
}

func writeToFile(code string) {
	os.WriteFile("build/main.js", []byte(code), 0644)
}
