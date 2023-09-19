package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType int

const (
	keyword_beg TokenType = iota
	LET
	IF
	ELSE
	FOR
	WHILE
	FUNC
	RETURN
	DEFER
	RANGE

	keyword_end

	operator_beg

	DECLARE
	ASSIGN

	ADD
	SUB
	MUL
	DIV
	POW
	MOD

	INCR
	DECR

	AND
	OR

	GT
	LT
	GTE
	LTE
	EQ
	NEQ

	LPAREN
	RPAREN

	LBRACK
	RBRACK

	LBRACE
	RBRACE

	COMMA
	COLON

	operator_end

	NUMBER
	BOOLEAN
	IDENTIFIER
	STRING

	EOF
	UNKNOWN
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input string
	pos   int
	len   int
}

var keywords map[string]TokenType = map[string]TokenType{
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"for":    FOR,
	"while":  WHILE,
	"func":   FUNC,
	"return": RETURN,
	"defer":  DEFER,
	"range":  RANGE,

	"true":  BOOLEAN,
	"false": BOOLEAN,
}

var operators map[string]TokenType = map[string]TokenType{

	":=": DECLARE,
	"=":  ASSIGN,

	"+":  ADD,
	"-":  SUB,
	"*":  MUL,
	"/":  DIV,
	"**": POW,
	"%":  MOD,

	"++": INCR,
	"--": DECR,

	"&&": AND,
	"||": OR,

	">":  GT,
	"<":  LT,
	">=": GTE,
	"<=": LTE,
	"==": EQ,
	"!=": NEQ,

	"(": LPAREN,
	")": RPAREN,

	"[": LBRACK,
	"]": RBRACK,

	"{": LBRACE,
	"}": RBRACE,

	",": COMMA,
	":": COLON,
}

func NewLexer(input string) *Lexer {
	input = strings.TrimSpace(input)
	return &Lexer{input: input, len: len(input), pos: 0}
}

func (l *Lexer) getToken() (*Token, error) {
	if isWhitespace(l.current()) {
		l.skipWhitespace()
	}

	if l.pos >= l.len {
		return &Token{Type: EOF}, nil
	}

	if tok := l.tryTokenizeIdentifier(); tok != nil {
		return tok, nil
	} else if tok = l.tryTokenizeNumber(); tok != nil {
		return tok, nil
	} else if tok = l.tryTokenizeString(); tok != nil {
		return tok, nil
	} else if tok = l.tryTokenizeOperator(); tok != nil {
		return tok, nil
	} else {
		return nil, (fmt.Errorf("invalid token %c at position %d", l.current(), l.pos))
	}

}

func (l *Lexer) GetTokens() ([]*Token, error) {
	var tokens []*Token
	for l.pos < l.len {
		tok, err := l.getToken()
		if err != nil {
			return tokens, err
		}
		tokens = append(tokens, tok)
	}
	return tokens, nil

}

func (l *Lexer) tryTokenizeIdentifier() *Token {
	if isAlpha(l.current()) {
		val := ""
		for l.pos < (l.len) && isAlpha(l.current()) {
			val += string(l.current())
			l.next()
		}
		if tokType, ok := keywords[val]; ok {
			return &Token{Type: tokType, Value: val}
		}

		return &Token{Type: IDENTIFIER, Value: val}
	}

	return nil
}

func (l *Lexer) tryTokenizeNumber() *Token {
	if isNumber(l.current()) {
		val := ""
		for l.pos < (l.len) && isNumber(l.current()) {
			val += string(l.current())
			l.next()
		}
		return &Token{Type: NUMBER, Value: val}
	}

	return nil
}

func (l *Lexer) tryTokenizeString() *Token {
	if l.current() == '"' {
		val := ""
		l.next()
		for l.pos < (l.len) && l.current() != '"' {
			val += string(l.current())
			l.next()
		}
		l.next()
		return &Token{Type: STRING, Value: val}
	}
	return nil
}

func (l *Lexer) tryTokenizeOperator() *Token {

	if l.current() == '&' && l.peek() == '&' {
		l.next()
		l.next()
		return &Token{Type: AND, Value: "&&"}
	} else if l.current() == '|' && l.peek() == '|' {
		l.next()
		l.next()
		return &Token{Type: OR, Value: "||"}
	}

	if tokType, ok := operators[string(l.current())]; ok {
		val := string(l.current())

		if l.pos < (l.len - 1) {

			switch tokType {
			case COLON:
				if l.peek() == '=' {
					val += string(l.peek())
					l.next()
					tokType = DECLARE
				}
			case MUL:
				if l.peek() == '*' {
					val += string(l.peek())
					l.next()
					tokType = POW
				}
			case ASSIGN:
				if l.peek() == '=' {
					val += string(l.peek())
					l.next()
					tokType = EQ
				}
			case GT:
				if l.peek() == '=' {
					val += string(l.peek())
					l.next()
					tokType = GTE
				}
			case LT:
				if l.peek() == '=' {
					val += string(l.peek())
					l.next()
					tokType = LTE
				}
			case ADD:
				if l.peek() == '+' {
					val += string(l.peek())
					l.next()
					tokType = INCR
				}
			case SUB:
				if l.peek() == '-' {
					val += string(l.peek())
					l.next()
					tokType = DECR
				}
			}
		}

		l.next()
		return &Token{Type: tokType, Value: val}
	}
	return nil

}

func (l *Lexer) current() rune {
	return rune(l.input[l.pos])
}

func (l *Lexer) next() {
	l.pos++
}

func (l *Lexer) peek() rune {
	return rune(l.input[l.pos+1])
}

func (l *Lexer) skipWhitespace() {
	for l.pos < l.len && isWhitespace(l.current()) {
		l.next()
	}
}

func isWhitespace(char rune) bool {
	return unicode.IsSpace((char))
}

func isAlpha(char rune) bool {
	return unicode.IsLetter(char)
}

func isNumber(char rune) bool {
	return unicode.IsNumber(char)
}

func (t *Token) IsKeyword() bool {
	return t.Type > keyword_beg && t.Type < keyword_end
}

func (t *Token) IsOperator() bool {
	return t.Type > operator_beg && t.Type < operator_end
}

func (t *Token) String() string {

	if t.IsKeyword() {
		return "keyword(" + t.Value + ")"
	} else if t.IsOperator() {
		return "operator(" + t.Value + ")"
	}

	switch t.Type {
	case IDENTIFIER:
		return "identifier(" + t.Value + ")"
	case NUMBER:
		return "number(" + t.Value + ")"
	case STRING:
		return "string(\"" + t.Value + "\")"
	case BOOLEAN:
		return "boolean(" + t.Value + ")"
	default:
		return "unknown"
	}

}
