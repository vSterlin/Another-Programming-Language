package ast

import (
	"fmt"
	"strings"
)

// Expressions
type Expr interface{ exprNode() }

type NumberExpr struct{ Val int }
type BooleanExpr struct{ Val bool }
type StringExpr struct{ Val string }
type IdentifierExpr struct{ Name string }

type BinaryExpr struct {
	Op  string
	Lhs Expr
	Rhs Expr
}

type CallExpr struct {
	Callee *IdentifierExpr
	Args   []Expr
}

type ArrayExpr struct {
	Elements []Expr
}

func (n *NumberExpr) exprNode()     {}
func (v *IdentifierExpr) exprNode() {}
func (b *BooleanExpr) exprNode()    {}
func (s *StringExpr) exprNode()     {}
func (b *BinaryExpr) exprNode()     {}
func (c *CallExpr) exprNode()       {}
func (a *ArrayExpr) exprNode()      {}

func (n *NumberExpr) String() string     { return fmt.Sprintf("numberExpression(%d)", n.Val) }
func (b *BooleanExpr) String() string    { return fmt.Sprintf("booleanExpression(%t)", b.Val) }
func (s *StringExpr) String() string     { return fmt.Sprintf("stringExpression(%s)", s.Val) }
func (v *IdentifierExpr) String() string { return fmt.Sprintf("identifierExpression(%s)", v.Name) }
func (b *BinaryExpr) String() string {
	return fmt.Sprintf("binaryExpression(%s, %s, %s)", b.Lhs, b.Op, b.Rhs)
}
func (c *CallExpr) String() string {
	return fmt.Sprintf("callExpression(%s, %s)", c.Callee, c.Args)
}
func (a *ArrayExpr) String() string {
	return fmt.Sprintf("arrayExpression(%s)", a.Elements)
}

// Statements
type Stmt interface{ stmtNode() }

type ExprStmt struct{ Expr Expr }

type VarDecStmt struct {
	Id   *IdentifierExpr
	Init Expr
}

type BlockStmt struct {
	Stmts []Stmt
}

type WhileStmt struct {
	Test Expr
	Body Stmt
}

type IfStmt struct {
	Test       Expr
	Consequent Stmt
	Alternate  Stmt
}

func (e *ExprStmt) stmtNode()   {}
func (v *VarDecStmt) stmtNode() {}
func (b *BlockStmt) stmtNode()  {}
func (w *WhileStmt) stmtNode()  {}
func (i *IfStmt) stmtNode()     {}

func (e *ExprStmt) String() string { return fmt.Sprintf("expressionStatement(%s)", e.Expr) }
func (v *VarDecStmt) String() string {
	return fmt.Sprintf("variableDeclarationStatement(%s, %s)", v.Id, v.Init)
}
func (b *BlockStmt) String() string {
	stmts := ""
	for _, stmt := range b.Stmts {
		stmts += fmt.Sprintf("%s\n", stmt)
	}
	stmts = strings.TrimSpace(stmts)
	return fmt.Sprintf("blockStatement(%s)", stmts)
}
func (w *WhileStmt) String() string {
	return fmt.Sprintf("whileStatement(%s, %s)", w.Test, w.Body)
}
func (i *IfStmt) String() string {
	alternate := "nil"
	if i.Alternate != nil {
		alternate = fmt.Sprintf("%s", i.Alternate)
	}
	return fmt.Sprintf("ifStatement(%s, %s, %s)", i.Test, i.Consequent, alternate)
}

type Program struct{ Stmts []Stmt }

func (p *Program) String() string {

	stmts := ""
	for _, stmt := range p.Stmts {
		stmts += fmt.Sprintf("%s\n", stmt)
	}

	stmts = strings.TrimSpace(stmts)
	return fmt.Sprintf("program(%s)", stmts)
}
