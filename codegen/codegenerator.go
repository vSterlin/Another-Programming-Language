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
	Str = types.I8Ptr
)

type CodeGenerator interface {
	Gen(prog *ast.Program) string
}

type LLVMCodeGenerator struct {
	module       *ir.Module
	currentBlock *ir.Block
	mainFunc     *ir.Func
}

func NewLLVMCodeGenerator() *LLVMCodeGenerator {

	module := ir.NewModule()
	return &LLVMCodeGenerator{
		module: module,
	}
}

// Statements
func (cg *LLVMCodeGenerator) genStmt(stmt ast.Stmt) value.Value {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		return cg.genExprStmt(stmt)
	case *ast.FuncDecStmt:
		return cg.genFuncDecStmt(stmt)
	case *ast.BlockStmt:
		return cg.genBlockStmt(stmt)
	case *ast.VarAssignStmt:
		return cg.genVarAssignStmt(stmt)
	default:
		return nil
	}

}

func (cg *LLVMCodeGenerator) genExprStmt(stmt *ast.ExprStmt) value.Value {
	return cg.genExpr(stmt.Expr)
}

func (cg *LLVMCodeGenerator) genFuncDecStmt(stmt *ast.FuncDecStmt) *ir.Func {
	fn := cg.module.NewFunc(stmt.Id.Name, I32)
	block := fn.NewBlock("entry")

	// to keep track of the current block to add stuff to
	prevBlock := cg.currentBlock
	cg.currentBlock = block

	cg.genBlockStmt(stmt.Body)

	block.NewRet(constant.NewInt(I32, 0))

	cg.currentBlock = prevBlock
	return fn
}

func (cg *LLVMCodeGenerator) genBlockStmt(stmt *ast.BlockStmt) value.Value {
	for _, stmt := range stmt.Stmts {
		cg.genStmt(stmt)
	}
	return nil
}

func (cg *LLVMCodeGenerator) genVarAssignStmt(stmt *ast.VarAssignStmt) value.Value {
	varName := stmt.Id.Name
	init := cg.genExpr(stmt.Init)
	fmt.Println(varName, init)

	block := cg.getCurrentBlock()
	initType := init.Type()
	alloc := block.NewAlloca(initType)
	block.NewStore(init, alloc)

	return nil
}

// Expressions
func (cg *LLVMCodeGenerator) genExpr(expr ast.Expr) value.Value {
	switch expr := expr.(type) {
	case *ast.BinaryExpr:
		return cg.genBinaryExpr(expr)
	case *ast.NumberExpr:
		return genNumberExpr(expr)
	case *ast.StringExpr:
		return genStringExpr(expr)
	default:
		return nil
	}
}

// Literals start
func genNumberExpr(expr *ast.NumberExpr) *constant.Int {
	return constant.NewInt(I32, int64(expr.Val))
}

func genStringExpr(expr *ast.StringExpr) *constant.CharArray {
	return constant.NewCharArrayFromString(expr.Val)
}

// Literals end

func (cg *LLVMCodeGenerator) genBinaryExpr(expr *ast.BinaryExpr) value.Value {

	lhs := genNumberExpr(expr.Lhs.(*ast.NumberExpr))
	rhs := genNumberExpr(expr.Rhs.(*ast.NumberExpr))

	block := cg.getCurrentBlock()

	switch expr.Op {
	case "+":
		return block.NewAdd(lhs, rhs)
	case "-":
		return block.NewSub(lhs, rhs)
	case "*":
		return block.NewMul(lhs, rhs)
	case "/":
		return block.NewSDiv(lhs, rhs)
	default:
		return nil
	}

}

func (cg *LLVMCodeGenerator) Gen(prog *ast.Program) string {
	m := cg.module

	fn := m.NewFunc("main", I32)
	block := fn.NewBlock("entry")
	cg.mainFunc = fn

	setupExternal(m)

	for _, stmt := range prog.Stmts {
		cg.genStmt(stmt)
	}

	block.NewRet(constant.NewInt(I32, 0))

	return m.String()
}

// Helpers
func setupExternal(m *ir.Module) {
	m.NewFunc("printf", I32, ir.NewParam("", Str))
}

func (cg *LLVMCodeGenerator) getCurrentBlock() *ir.Block {
	currentBlock := cg.currentBlock
	if currentBlock == nil {
		currentBlock = cg.mainFunc.Blocks[0]
	}
	return currentBlock
}
