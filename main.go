package main

import "fmt"

func main() {
	l := NewLexer("let x = 123 + (77 - 99 * 2) + \"test\"")

	tokens, err := l.GetTokens()

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, tok := range tokens {
		fmt.Println(tok)
	}

}
