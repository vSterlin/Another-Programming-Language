package main

type Expr interface {
	exprNode()
}

type NumberExpr struct {
	Val int
}

func (n *NumberExpr) exprNode() {}

type VariableExpr struct {
	Name string
}

func (v *VariableExpr) exprNode() {}

type BinaryExpr struct {
	Op  string
	Lhs Expr
	Rhs Expr
}

func (b *BinaryExpr) exprNode() {}

type CallExpr struct {
	Callee string
	Args   []Expr
}

func (c *CallExpr) exprNode() {}

type Stmt interface {
	stmtNode()
}

type Program struct {
	Stmts []Stmt
}

func (p *Program) stmtNode() {}
