package codegen

import (
	"fmt"
	"language/ast"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

var (
	I32 = types.I32
)

type CodeGenerator interface {
	Generate(prog *ast.Program) string
}

type LLVMCodeGenerator struct {
	module       *ir.Module
	currentBlock *ir.Block
}

func NewLLVMCodeGenerator() *LLVMCodeGenerator {

	module := ir.NewModule()
	return &LLVMCodeGenerator{
		module: module,
	}
}

// Statements
func (cg *LLVMCodeGenerator) generateStmt(stmt ast.Stmt) value.Value {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		return cg.generateExprStmt(stmt)
	default:
		return nil
	}

}

func (cg *LLVMCodeGenerator) generateExprStmt(stmt *ast.ExprStmt) value.Value {
	return cg.generateExpr(stmt.Expr)
}

func (cg *LLVMCodeGenerator) generateFuncDecStmt(stmt *ast.FuncDecStmt) {

}

func (cg *LLVMCodeGenerator) generateBlockStmt(stmt *ast.BlockStmt) {

}

// Expressions
func (cg *LLVMCodeGenerator) generateExpr(expr ast.Expr) value.Value {
	switch expr := expr.(type) {
	case *ast.BinaryExpr:
		return cg.generateBinaryExpr(expr)
	default:
		return nil
	}
}

func generateNumberExpr(expr *ast.NumberExpr) *constant.Int {
	return constant.NewInt(I32, int64(expr.Val))
}

func (cg *LLVMCodeGenerator) generateBinaryExpr(expr *ast.BinaryExpr) value.Value {
	lhs := generateNumberExpr(expr.Lhs.(*ast.NumberExpr))
	rhs := generateNumberExpr(expr.Rhs.(*ast.NumberExpr))

	var resExpr constant.Constant
	switch expr.Op {
	case "+":
		resExpr = constant.NewAdd(lhs, rhs)
	case "-":
		resExpr = constant.NewSub(lhs, rhs)
	}

	// fmt.Println(resExpr)
	return resExpr
}

func (cg *LLVMCodeGenerator) Generate(prog *ast.Program) string {
	m := cg.module
	// for _, stmt := range prog.Stmts {
	// 	cg.generateStmt(stmt)
	// }

	fn := m.NewFunc("main", types.I32)
	block := fn.NewBlock("entry")

	// for now
	res := cg.generateStmt(prog.Stmts[0])

	fmt.Println(res)

	// res  := cg.generateStmt(prog.Stmts[0])

	block.NewRet(res)

	return m.String()
}
