package parser

import "fmt"

type ParserError struct {
	Pos int
	Msg string
}

func (e *ParserError) Error() string {
	return fmt.Sprintf("%s, line: %d", e.Msg, e.Pos)
}

func NewParserError(pos int, msg string) *ParserError {
	return &ParserError{Pos: pos, Msg: msg}
}
