package parser

import (
	"language/ast"
	"language/lexer"
	"reflect"
	"testing"
)

func TestParseNumberExpr(t *testing.T) {
	p := NewParser(getTokens("1"))
	expr, err := p.parseNumberExpr()
	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	e, ok := expr.(*ast.NumberExpr)

	if !ok {
		t.Errorf("Expected NumberExpr, got: %T", expr)
	}

	if e.Val != 1 {
		t.Errorf("Expected 1, got: %d", e.Val)
	}
}

func TestParseBooleanExpr(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		p := NewParser(getTokens(tt.input))
		expr, err := p.parseBooleanExpr()
		if err != nil {
			t.Errorf("Expected no error, got: %s", err)
		}

		e, ok := expr.(*ast.BooleanExpr)

		if !ok {
			t.Errorf("Expected BooleanExpr, got: %T", expr)
		}

		if e.Val != tt.want {
			t.Errorf("Expected %t, got: %t", tt.want, e.Val)
		}

	}
}

func TestParseStringExpr(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`""`, ""},
		{`"hello"`, "hello"},
		{`"hello world"`, "hello world"},
		{`"1"`, "1"},
		{`"true"`, "true"},
		{`"false"`, "false"},
	}

	for _, tt := range tests {
		p := NewParser(getTokens(tt.input))
		expr, err := p.parseStringExpr()
		if err != nil {
			t.Errorf("Expected no error, got: %s", err)
		}

		e, ok := expr.(*ast.StringExpr)

		if !ok {
			t.Errorf("Expected StringExpr, got: %T", expr)
		}

		if e.Val != tt.want {
			t.Errorf("Expected %s, got: %s", tt.want, e.Val)
		}

	}

}

func TestParseIdentifierExpr(t *testing.T) {

	id := "x"
	p := NewParser(getTokens(id))
	expr, err := p.parseIdentifierExpr()
	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	if expr.Name != id {
		t.Errorf("Expected %s, got: %s", id, expr.Name)
	}

}

func TestParseParenExpr(t *testing.T) {
	tests := []struct {
		srcCode string
		want    ast.Expr
	}{
		{srcCode: "(1)", want: &ast.NumberExpr{}},
		{srcCode: "(true)", want: &ast.BooleanExpr{}},
		{srcCode: `("hello")`, want: &ast.StringExpr{}},
	}

	for _, tt := range tests {
		p := NewParser(getTokens(tt.srcCode))
		expr, err := p.parseParenExpr()
		if err != nil {
			t.Errorf("Expected no error, got: %s", err)
		}

		if reflect.TypeOf(tt.want) != reflect.TypeOf(expr) {
			t.Errorf("Expected %T, got: %T", tt.want, expr)
		}

	}

}

func TestParseArrowFunc(t *testing.T) {
	code := `() => {}`
	p := NewParser(getTokens(code))
	expr, err := p.parseArrowFunc()
	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	_, ok := expr.(*ast.ArrowFunc)
	if !ok {
		t.Errorf("Expected ArrowFunc, got: %T", expr)
	}

	code = `(x int) => {}`
	p = NewParser(getTokens(code))
	expr, err = p.parseArrowFunc()
	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	arrFunc, ok := expr.(*ast.ArrowFunc)
	if !ok {
		t.Errorf("Expected ArrowFunc, got: %T", expr)
	}

	if len(arrFunc.Args) != 1 {
		t.Errorf("Expected 1, got: %d", len(arrFunc.Args))
	}

	if arrFunc.Args[0].Id.Name != "x" {
		t.Errorf("Expected x, got: %s", arrFunc.Args[0].Id.Name)
	}

}

func TestParsePrimaryExpr(t *testing.T) {
	tests := []struct {
		input string
		want  ast.Expr
	}{
		{"1", &ast.NumberExpr{}},
		{"true", &ast.BooleanExpr{}},
		{"false", &ast.BooleanExpr{}},
		{`""`, &ast.StringExpr{}},
	}
	for _, tt := range tests {
		p := NewParser(getTokens(tt.input))
		expr, err := p.parsePrimaryExpr()
		if err != nil {
			t.Errorf("Expected no error, got: %s", err)
		}

		if reflect.TypeOf(expr) != reflect.TypeOf(tt.want) {
			t.Errorf("Expected %T, got: %T", tt.want, expr)
		}
	}
}

func TestParseUpdateExpr(t *testing.T) {

	tests := []string{"x++", "x--"}

	want := &ast.UpdateExpr{}

	for _, tt := range tests {
		p := NewParser(getTokens(tt))
		expr, err := p.parseUpdateExpr()
		if err != nil {
			t.Errorf("Expected no error, got: %s", err)
		}

		if reflect.TypeOf(expr) != reflect.TypeOf(want) {
			t.Errorf("Expected %T, got: %T", want, expr)
		}
	}
}

func TestParseCallExpr(t *testing.T) {
	test := "foo()"
	p := NewParser(getTokens(test))
	expr, err := p.parseCallExpr()
	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	_, ok := expr.(*ast.CallExpr)
	if !ok {
		t.Errorf("Expected CallExpr, got: %T", expr)
	}

}

func TestParseUnaryExpr(t *testing.T) {

	tests := []struct {
		input string
		want  ast.Expr
	}{
		{"!x", &ast.UnaryExpr{}},
		{"x++", &ast.UpdateExpr{}},
		{"x--", &ast.UpdateExpr{}},
		{"foo()", &ast.CallExpr{}},
	}

	for _, tt := range tests {
		p := NewParser(getTokens(tt.input))
		expr, err := p.parseUnaryExpr()
		if err != nil {
			t.Errorf("Expected no error, got: %s", err)
		}

		if reflect.TypeOf(expr) != reflect.TypeOf(tt.want) {
			t.Errorf("Expected %T, got: %T", tt.want, expr)
		}
	}
}

func TestParseEqualityExpr(t *testing.T) {
	test := "x == y"
	p := NewParser(getTokens(test))
	expr, err := p.parseEqualityExpr()
	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	_, ok := expr.(*ast.BinaryExpr)
	if !ok {
		t.Errorf("Expected BinaryExpr, got: %T", expr)
	}
}

func TestParseRelationalExpr(t *testing.T) {
	tests := []string{
		"x < y",
		"x > y",
		"x <= y",
		"x >= y",
	}

	want := &ast.BinaryExpr{}

	for _, tt := range tests {
		p := NewParser(getTokens(tt))
		expr, err := p.parseRelationalExpr()
		if err != nil {
			t.Errorf("Expected no error, got: %s", err)
		}

		if reflect.TypeOf(expr) != reflect.TypeOf(want) {
			t.Errorf("Expected %T, got: %T", want, expr)
		}
	}
}

func TestParseAdditiveExpr(t *testing.T) {
	tests := []string{
		"x + y",
		"x - y",
	}

	want := &ast.BinaryExpr{}

	for _, tt := range tests {
		p := NewParser(getTokens(tt))
		expr, err := p.parseAdditiveExpr()
		if err != nil {
			t.Errorf("Expected no error, got: %s", err)
		}

		if reflect.TypeOf(expr) != reflect.TypeOf(want) {
			t.Errorf("Expected %T, got: %T", want, expr)
		}

	}
}

func TestParseMultiplicativeExpr(t *testing.T) {
	tests := []string{
		"x * y",
		"x / y",
		"x % y",
		"x ** y",
	}

	want := &ast.BinaryExpr{}

	for _, tt := range tests {
		p := NewParser(getTokens(tt))
		expr, err := p.parseMultiplicativeExpr()
		if err != nil {
			t.Errorf("Expected no error, got: %s", err)
		}

		if reflect.TypeOf(expr) != reflect.TypeOf(want) {
			t.Errorf("Expected %T, got: %T", want, expr)
		}
	}
}

func TestParseAndExpr(t *testing.T) {

	test := "x && y"
	p := NewParser(getTokens(test))
	expr, err := p.parseAndExpr()
	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	_, ok := expr.(*ast.LogicalExpr)
	if !ok {
		t.Errorf("Expected LogicalExpr, got: %T", expr)
	}

}

func TestParseOrExpr(t *testing.T) {
	test := "x || y"
	p := NewParser(getTokens(test))
	expr, err := p.parseOrExpr()
	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	_, ok := expr.(*ast.LogicalExpr)
	if !ok {
		t.Errorf("Expected LogicalExpr, got: %T", expr)
	}
}

func TestParseTypeExpr(t *testing.T) {
	tests := []string{
		"int",
		"bool",
		"string",
		"void",
	}

	want := &ast.TypeExpr{}

	for _, tt := range tests {
		p := NewParser(getTokens(tt))
		expr, err := p.parseTypeExpr()
		if err != nil {
			t.Errorf("Expected no error, got: %s", err)
		}

		if reflect.TypeOf(expr) != reflect.TypeOf(want) {
			t.Errorf("Expected %T, got: %T", want, expr)
		}
	}
}

func getTokens(code string) []*lexer.Token {
	l := lexer.NewLexer(code)
	tokens, _ := l.GetTokens()
	return tokens
}
