package main

import "fmt"

type Interpreter struct {
}

type RuntimeValue interface {
	value() any
}

type NumberValue struct {
	Val int
}

func (n *NumberValue) value() any {
	return n.Val
}

func (n *NumberValue) String() string {
	return fmt.Sprintf("numberValue(%d)", n.Val)
}

func (i *Interpreter) evaluateBinaryExpr(ex *BinaryExpr) RuntimeValue {
	lhs := i.evaluate(ex.Lhs)
	rhs := i.evaluate(ex.Rhs)

	lhsNum, lhsOk := lhs.(*NumberValue)
	rhsNum, rhsOk := rhs.(*NumberValue)

	if !lhsOk || !rhsOk {
		return nil
	}

	res := 0
	switch ex.Op {
	case "+":
		res = lhsNum.Val + rhsNum.Val
	case "-":
		res = lhsNum.Val - rhsNum.Val
	case "*":
		res = lhsNum.Val * rhsNum.Val
	case "/":
		res = lhsNum.Val / rhsNum.Val

	default:
		return nil
	}

	return &NumberValue{Val: res}
}

func (i *Interpreter) evaluate(ex Expr) RuntimeValue {
	switch ex := ex.(type) {
	case *NumberExpr:
		return &NumberValue{Val: ex.Val}
	case *BinaryExpr:
		return i.evaluateBinaryExpr(ex)
	default:
		return nil
	}
}
