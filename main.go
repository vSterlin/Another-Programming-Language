package main

import (
	"language/codegen"
	"language/lexer"
	"language/parser"
	"os"
)

func main() {
	l := lexer.NewLexer(`
	x := [1,2,3,4,5,6,7,8,9,10]
	y = x[1:8]
	z = y[1:5:2]
	`)
	tokens, _ := l.GetTokens()
	// fmt.Println(tokens)
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
