package ast

import (
	"encoding/json"
)

// Expressions
type Expr interface {
	exprNode()
	json.Marshaler
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

type BinaryExpr struct {
	Op  string `json:"operator"`
	Lhs Expr   `json:"left"`
	Rhs Expr   `json:"right"`
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
func (c *CallExpr) exprNode()       {}
func (a *ArrayExpr) exprNode()      {}
func (s *SliceExpr) exprNode()      {}

func (n *NumberExpr) MarshalJSON() ([]byte, error)     { return toJSON(*n, "numberExpression") }
func (b *BooleanExpr) MarshalJSON() ([]byte, error)    { return toJSON(*b, "booleanExpression") }
func (s *StringExpr) MarshalJSON() ([]byte, error)     { return toJSON(*s, "stringExpression") }
func (i *IdentifierExpr) MarshalJSON() ([]byte, error) { return toJSON(*i, "identifierExpression") }
func (b *BinaryExpr) MarshalJSON() ([]byte, error)     { return toJSON(*b, "binaryExpression") }
func (c *CallExpr) MarshalJSON() ([]byte, error)       { return toJSON(*c, "callExpression ") }
func (a *ArrayExpr) MarshalJSON() ([]byte, error)      { return toJSON(*a, "arrayExpression") }
func (s *SliceExpr) MarshalJSON() ([]byte, error)      { return toJSON(*s, "sliceExpression") }

// Statements
type Stmt interface {
	stmtNode()
	json.Marshaler
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

func (e *ExprStmt) MarshalJSON() ([]byte, error)     { return toJSON(e, "expressionStatement") }
func (v *VarDecStmt) MarshalJSON() ([]byte, error)   { return toJSON(*v, "variableDeclarationStatement") }
func (b *BlockStmt) MarshalJSON() ([]byte, error)    { return toJSON(*b, "blockStatement") }
func (w *WhileStmt) MarshalJSON() ([]byte, error)    { return toJSON(*w, "whileStatement") }
func (f *FuncDecStmt) MarshalJSON() ([]byte, error)  { return toJSON(*f, "funcDeclarationStatement") }
func (i *IfStmt) MarshalJSON() ([]byte, error)       { return toJSON(*i, "ifStatement") }
func (d *DeferStmt) MarshalJSON() ([]byte, error)    { return toJSON(*d, "deferStatement") }
func (r *RangeStmt) MarshalJSON() ([]byte, error)    { return toJSON(*r, "rangeStatement") }
func (i *IncrDecrStmt) MarshalJSON() ([]byte, error) { return toJSON(*i, "incrDecrStatement") }
func (v *VarAssignStmt) MarshalJSON() ([]byte, error) {
	return toJSON(*v, "variableAssignmentStatement")
}

type Program struct{ Stmts []Stmt }

// util
func toJSON(node any, typeName string) ([]byte, error) {

	m := make(map[string]any)

	m["type"] = typeName
	jsonStr, err := json.Marshal((node))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonStr, &m)
	if err != nil {
		return nil, err
	}

	jsonStr, err = json.Marshal(m)
	return (jsonStr), err
}
