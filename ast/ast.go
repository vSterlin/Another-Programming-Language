package ast

// Expressions
type Expr interface {
	exprNode()
	// fmt.Stringer
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
	Callee *IdentifierExpr `json:"callee"`
	Args   []Expr          `json:"arguments"`
}

type ArrayExpr struct {
	Elements []Expr `json:"elements"`
}

type SliceExpr struct {
	Id   *IdentifierExpr `json:"identifier"`
	Low  Expr            `json:"low"`
	High Expr            `json:"high"`
	Step Expr            `json:"step"`
}

func (n *NumberExpr) exprNode()     {}
func (v *IdentifierExpr) exprNode() {}
func (b *BooleanExpr) exprNode()    {}
func (s *StringExpr) exprNode()     {}
func (b *BinaryExpr) exprNode()     {}
func (b *LogicalExpr) exprNode()    {}
func (c *CallExpr) exprNode()       {}
func (a *ArrayExpr) exprNode()      {}
func (s *SliceExpr) exprNode()      {}

// Statements
type Stmt interface {
	stmtNode()
	// fmt.Stringer
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
func (v *VarDecStmt) stmtNode()    {}
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

type Program struct{ Stmts []Stmt }
