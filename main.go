package main

import (
	"encoding/json"
	"fmt"
	"language/codegen"
	"language/lexer"
	"language/parser"
	"os"
)

func main() {
	l := lexer.NewLexer(`
		func main() {
			defer print("lol")
			a := 1
			b := 2
		}
	`)
	tokens, _ := l.GetTokens()
	fmt.Println(tokens)
	p := parser.NewParser(tokens)

	prog := p.ParseProgram()
	jsonStr, _ := json.MarshalIndent(prog, "", "  ")

	fmt.Println(string(jsonStr))

	cg := codegen.NewJavascriptCodeGenerator()
	code := cg.Generate(prog)

	writeToFile(code)

}

func writeToFile(code string) {
	os.WriteFile("build/main.js", []byte(code), 0644)
}
