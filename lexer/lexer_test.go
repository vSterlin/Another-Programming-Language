package lexer

import "testing"

func TestTokenTypes(t *testing.T) {
	for _, i := range tests {
		l := NewLexer(i.input)
		tok, err := l.getToken()
		if err != nil {
			t.Errorf("Did not expect error, got: %s", err)

		}
		if tok.Type != i.expected {
			t.Errorf("Expected %d, got: %d", i.expected, tok.Type)
		}

	}
}

func TestInvalidInput(t *testing.T) {
	l := NewLexer("$")
	_, err := l.getToken()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

var tests = []struct {
	input    string
	expected TokenType
}{
	{input: "let", expected: LET},
	{input: "if", expected: IF},
	{input: "else", expected: ELSE},
	{input: "while", expected: WHILE},
	{input: "return", expected: RETURN},

	{input: "+", expected: ADD},
	{input: "-", expected: SUB},
	{input: "*", expected: MUL},
	{input: "/", expected: DIV},
	{input: "=", expected: ASSIGN},

	{input: "(", expected: LPAREN},
	{input: ")", expected: RPAREN},

	{input: "[", expected: LBRACK},
	{input: "]", expected: RBRACK},

	{input: ",", expected: COMMA},

	{input: "1", expected: NUMBER},

	{input: "true", expected: BOOLEAN},
	{input: "false", expected: BOOLEAN},
}
