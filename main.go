package main

import (
	"bufio"
	"fmt"

	"os"
	"strings"
)

func main() {
	l := NewLexer("let x = 1")

	tokens, err := l.GetTokens()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(tokens)
	// for _, tok := range tokens {
	// 	fmt.Println(tok)
	// }

	p := NewParser(tokens)

	stmt := p.parseVarDecStmt()

	fmt.Println(stmt)

	// prog := p.parseProgram()

	// t := Typechecker{}

	// t.typecheckProgram(prog)

	// fmt.Println(prog)

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
