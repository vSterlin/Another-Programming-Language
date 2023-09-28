package ast

import "fmt"

// Expressions
type Expr interface {
	exprNode()
	fmt.Stringer
}

type NumberExpr struct {
	Val int ``
}
type BooleanExpr struct {
	Val bool `json:"value"`
}
type StringExpr struct {
	Val string `json:"value"`
}
type IdentifierExpr struct {
	Name string `json:"name"`
}

type MemberExpr struct {
	Obj  Expr
	Prop Expr
}

type BinOp string

const (
	ADD BinOp = "+"
	SUB BinOp = "-"
	MUL BinOp = "*"
	DIV BinOp = "/"
	MOD BinOp = "%"
	POW BinOp = "**"
)

type BinaryExpr struct {
	Op  BinOp `json:"operator"`
	Lhs Expr  `json:"left"`
	Rhs Expr  `json:"right"`
}

type LogicalOp string

const (
	AND LogicalOp = "&&"
	OR  LogicalOp = "||"
)

type LogicalExpr struct {
	Op  LogicalOp `json:"operator"`
	Lhs Expr      `json:"left"`
	Rhs Expr      `json:"right"`
}

type CallExpr struct {
	Callee Expr   `json:"callee"`
	Args   []Expr `json:"arguments"`
}

type ArrayExpr struct {
	Elements []Expr `json:"elements"`
}

// TODO: maybe this should be a member expr
type SliceExpr struct {
	Id   *IdentifierExpr `json:"identifier"`
	Low  Expr            `json:"low"`
	High Expr            `json:"high"`
	Step Expr            `json:"step"`
}

type ThisExpr struct{}

func (n *NumberExpr) exprNode()     {}
func (v *IdentifierExpr) exprNode() {}
func (b *BooleanExpr) exprNode()    {}
func (s *StringExpr) exprNode()     {}
func (b *BinaryExpr) exprNode()     {}
func (b *LogicalExpr) exprNode()    {}
func (c *CallExpr) exprNode()       {}
func (a *ArrayExpr) exprNode()      {}
func (s *SliceExpr) exprNode()      {}
func (m *MemberExpr) exprNode()     {}
func (t *ThisExpr) exprNode()       {}

func (n *NumberExpr) String() string     { return fmt.Sprintf("number(%d)", n.Val) }
func (v *IdentifierExpr) String() string { return fmt.Sprintf("identifier(%s)", v.Name) }
func (b *BooleanExpr) String() string    { return fmt.Sprintf("boolean(%t)", b.Val) }
func (s *StringExpr) String() string     { return fmt.Sprintf("string(%s)", s.Val) }
func (b *BinaryExpr) String() string     { return fmt.Sprintf("binary(%s, %s, %s)", b.Lhs, b.Op, b.Rhs) }
func (b *LogicalExpr) String() string    { return fmt.Sprintf("logical(%s, %s, %s)", b.Lhs, b.Op, b.Rhs) }
func (c *CallExpr) String() string       { return fmt.Sprintf("call(%s)", c.Callee) }
func (a *ArrayExpr) String() string      { return fmt.Sprintf("array(%s)", a.Elements) }
func (s *SliceExpr) String() string      { return fmt.Sprintf("slice(%s)", s.Id) }
func (m *MemberExpr) String() string     { return fmt.Sprintf("member(%s, %s)", m.Obj, m.Prop) }
func (t *ThisExpr) String() string       { return ("this") }

// Statements
type Stmt interface {
	stmtNode()
	fmt.Stringer
}

type ExprStmt struct{ Expr Expr }

type VarDecStmt struct {
	Id   *IdentifierExpr `json:"identifier"`
	Init Expr            `json:"init"`
}

type VarAssignStmt struct {
	Id   *IdentifierExpr `json:"identifier"`
	Op   string          `json:"operator"`
	Init Expr            `json:"init"`
}

type SetStmt struct {
	Lhs  Expr   `json:"object"`
	Name string `json:"name"`
	Val  Expr   `json:"value"`
}

type BlockStmt struct {
	Stmts []Stmt `json:"statements"`
}

type WhileStmt struct {
	Test Expr `json:"test"`
	Body Stmt `json:"body"`
}

type FuncDecStmt struct {
	Id   *IdentifierExpr   `json:"identifier"`
	Args []*IdentifierExpr `json:"arguments"`
	Body *BlockStmt        `json:"body"`
}

type IfStmt struct {
	Test       Expr `json:"test"`
	Consequent Stmt `json:"consequent"`
	Alternate  Stmt `json:"alternate"`
}

type DeferStmt struct {
	Call *CallExpr `json:"call"`
}

type RangeStmt struct {
	Id   *IdentifierExpr `json:"identifier"`
	Expr Expr            `json:"expression"`
	Body *BlockStmt      `json:"body"`
}

type IncrDecrStmt struct {
	Expr Expr   `json:"expression"`
	Op   string `json:"operator"`
}

type ReturnStmt struct {
	Arg Expr `json:"argument"`
}

type ClassDecStmt struct {
	Id      *IdentifierExpr `json:"identifier"`
	Methods []*FuncDecStmt  `json:"methods"`
}

func (e *ExprStmt) stmtNode()      {}
func (v *VarAssignStmt) stmtNode() {}
func (b *BlockStmt) stmtNode()     {}
func (w *WhileStmt) stmtNode()     {}
func (f *FuncDecStmt) stmtNode()   {}
func (i *IfStmt) stmtNode()        {}
func (d *DeferStmt) stmtNode()     {}
func (r *RangeStmt) stmtNode()     {}
func (i *IncrDecrStmt) stmtNode()  {}
func (r *ReturnStmt) stmtNode()    {}
func (c *ClassDecStmt) stmtNode()  {}
func (v *SetStmt) stmtNode()       {}
func (p *Program) stmtNode()       {}

func (e *ExprStmt) String() string      { return fmt.Sprintf("expr(%s)", e.Expr) }
func (v *VarAssignStmt) String() string { return fmt.Sprintf("var(%s)", v.Id) }
func (b *BlockStmt) String() string     { return fmt.Sprintf("block(%s)", b.Stmts) }
func (w *WhileStmt) String() string     { return fmt.Sprintf("while(%s)", w.Test) }
func (f *FuncDecStmt) String() string   { return fmt.Sprintf("func(%s)", f.Id) }
func (i *IfStmt) String() string        { return fmt.Sprintf("if(%s)", i.Test) }
func (d *DeferStmt) String() string     { return fmt.Sprintf("defer(%s)", d.Call) }
func (r *RangeStmt) String() string     { return fmt.Sprintf("range(%s)", r.Id) }
func (i *IncrDecrStmt) String() string  { return fmt.Sprintf("update(%s, %s)", i.Expr, i.Op) }
func (r *ReturnStmt) String() string    { return fmt.Sprintf("return(%s)", r.Arg) }
func (c *ClassDecStmt) String() string  { return fmt.Sprintf("class(%s, methods(%s))", c.Id, c.Methods) }
func (v *SetStmt) String() string       { return fmt.Sprintf("set(%s, %s, %s)", v.Lhs, v.Name, v.Val) }
func (p *Program) String() string       { return fmt.Sprintf("program(%s)", p.Stmts) }

type Program struct{ Stmts []Stmt }
