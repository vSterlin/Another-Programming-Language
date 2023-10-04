package codegen

import (
	"fmt"
	"language/ast"
	"strings"

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
	case *ast.CallExpr:
		return cg.genCallExpr(expr)
	default:
		return nil
	}
}

// Literals start
func genNumberExpr(expr *ast.NumberExpr) *constant.Int {
	return constant.NewInt(I32, int64(expr.Val))
}

func genStringExpr(expr *ast.StringExpr) *constant.CharArray {
	text := strings.Replace(expr.Val, "\\n", "\n", -1)
	str := constant.NewCharArrayFromString(text)
	return str
}

// Literals end

func (cg *LLVMCodeGenerator) genBinaryExpr(expr *ast.BinaryExpr) value.Value {

	lhs := cg.genExpr(expr.Lhs)
	rhs := cg.genExpr(expr.Rhs)

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

func (cg *LLVMCodeGenerator) genCallExpr(expr *ast.CallExpr) value.Value {
	fn := cg.getFunction(expr.Callee.(*ast.IdentifierExpr).Name)
	if fn == nil {
		panic("Function not found")
	}

	arg := cg.genExpr(expr.Args[0]).(*constant.CharArray)

	block := cg.getCurrentBlock()

	m := cg.module
	argPtr := m.NewGlobalDef("argPtr", arg)
	gep := constant.NewGetElementPtr(arg.Typ, argPtr, constant.NewInt(I32, 0), constant.NewInt(I32, 0))

	// argPtr := block.NewAlloca(arg.Type())
	// block.NewStore(arg, argPtr)
	// gep := block.NewGetElementPtr(arg.Type(), argPtr, constant.NewInt(I32, 0), constant.NewInt(I32, 0))

	block.NewCall(fn, gep)

	fmt.Println("arg", arg, arg.Type())
	return nil
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

func (cg *LLVMCodeGenerator) getFunction(name string) *ir.Func {

	for _, f := range cg.module.Funcs {
		if name == "print" && f.Name() == "printf" {
			return f
		}
		if f.Name() == name {
			return f
		}
	}
	return nil
}
