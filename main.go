package main

import (
	"bufio"
	"fmt"

	"os"
	"strings"
)

func main() {
	l := NewLexer("true false")

	tokens, err := l.GetTokens()

	if err != nil {
		fmt.Println(err)
		return
	}

	// for _, tok := range tokens {
	// 	fmt.Println(tok)
	// }

	p := NewParser(tokens)

	prog := p.parseProgram()

	t := Typechecker{}

	t.typecheckProgram(prog)

	fmt.Println(prog)

	// fmt.Println(ex)
	// repl()

}

func repl() {
	// repl
	for {
		fmt.Print(">> ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n")
		l := NewLexer(text)
		tokens, err := l.GetTokens()

		if err != nil {
			fmt.Println(err)
			continue
		}
		p := NewParser(tokens)
		ex := p.parseExpr()
		fmt.Println(ex)
	}
}
