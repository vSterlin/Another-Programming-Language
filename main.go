package main

import (
	"fmt"
	"os"
)

func main() {
	l := NewLexer(`
	print(11 + 2 - 3)
	`)

	tokens, _ := l.GetTokens()
	p := NewParser(tokens)
	prog := p.parseProgram()
	fmt.Println(prog)

	cg := NewJavascriptCodeGenerator()
	code := cg.Generate(prog)

	writeToFile(code)
}

func writeToFile(code string) {
	os.WriteFile("build/main.js", []byte(code), 0644)
}
