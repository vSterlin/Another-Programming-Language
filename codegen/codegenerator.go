package codegen

import (
	"language/ast"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

var (
	Bool = types.I1
	Char = types.I8
	I32  = types.I32
	Str  = types.I8Ptr
	Void = types.Void
)

type CodeGenerator interface {
	Gen(prog *ast.Program) string
}

type LLVMCodeGenerator struct {
	module       *ir.Module
	currentBlock *ir.Block
	exitBlock    *ir.Block
	mainFunc     *ir.Func
	env          *Env
}

type Env struct {
	vars   map[string]*ir.InstAlloca
	parent *Env
}

func NewEnv(parent *Env) *Env {
	return &Env{
		vars:   make(map[string]*ir.InstAlloca),
		parent: parent,
	}
}

func NewLLVMCodeGenerator() *LLVMCodeGenerator {
	module := ir.NewModule()
	return &LLVMCodeGenerator{
		module: module,
		env:    NewEnv(nil),
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
	case *ast.IfStmt:
		return cg.genIfStmt(stmt)
	case *ast.ReturnStmt:
		return cg.genReturnStmt(stmt)
	default:
		return nil
	}

}

func (cg *LLVMCodeGenerator) genExprStmt(stmt *ast.ExprStmt) value.Value {
	return cg.genExpr(stmt.Expr)
}

func (cg *LLVMCodeGenerator) genFuncDecStmt(stmt *ast.FuncDecStmt) *ir.Func {
	fnParams := make([]*ir.Param, len(stmt.Args))
	for i, arg := range stmt.Args {
		fnParams[i] = ir.NewParam(arg.Id.Name, llvmType(arg.Type.Name))
	}

	fn := cg.module.NewFunc(stmt.Id.Name, llvmType(stmt.ReturnType.Name), fnParams...)
	block := fn.NewBlock("entry")

	// to keep track of the current block to add stuff to
	prevBlock := cg.getCurrentBlock()
	cg.currentBlock = block

	cg.genBlockStmt(stmt.Body)

	block = cg.getCurrentBlock()

	if block.Term == nil {
		fn := block.Parent
		if fn.Sig.RetType == Void {
			block.NewRet(nil)
		} else {
			block.NewRet(constant.NewInt(I32, 0))
		}
	}

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

	block := cg.getCurrentBlock()
	initType := init.Type()
	alloc := block.NewAlloca(initType)
	block.NewStore(init, alloc)

	cg.env.vars[varName] = alloc

	return nil
}

func (cg *LLVMCodeGenerator) genIfStmt(stmt *ast.IfStmt) value.Value {

	// TODO: if statement will generate an extra branch
	// but it doesn't really affect anything
	test := cg.genExpr(stmt.Test)
	block := cg.getCurrentBlock()

	fn := block.Parent
	thenBlock := fn.NewBlock("")
	elseBlock := fn.NewBlock("")
	exitBlock := fn.NewBlock("")

	prevExitBlock := cg.getExitBlock()
	cg.exitBlock = exitBlock

	block.NewCondBr(test, thenBlock, elseBlock)

	cg.currentBlock = thenBlock
	cg.genStmt(stmt.Consequent)
	if thenBlock.Term == nil {
		thenBlock.NewBr(exitBlock)
	}

	cg.currentBlock = elseBlock
	cg.genStmt(stmt.Alternate)
	if elseBlock.Term == nil {
		elseBlock.NewBr(exitBlock)
	}

	cg.currentBlock = exitBlock
	if exitBlock.Term == nil {
		if prevExitBlock != nil {
			exitBlock.NewBr(prevExitBlock)
		}
		// else {
		// exitBlock.NewRet(constant.NewInt(I32, 11))
		// }
	}

	cg.exitBlock = prevExitBlock
	return nil

}

func (cg *LLVMCodeGenerator) genReturnStmt(stmt *ast.ReturnStmt) value.Value {
	val := cg.genExpr(stmt.Arg)
	block := cg.getCurrentBlock()
	block.NewRet(val)
	return val
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
	case *ast.BooleanExpr:
		return genBooleanExpr(expr)
	case *ast.LogicalExpr:
		return cg.genLogicalExpr(expr)
	case *ast.CallExpr:
		return cg.genCallExpr(expr)
	case *ast.IdentifierExpr:
		return cg.genIdentifierExpr(expr)
	default:
		return nil
	}
}

// Literals start
func genNumberExpr(expr *ast.NumberExpr) *constant.Int {
	return constant.NewInt(I32, int64(expr.Val))
}

func genStringExpr(expr *ast.StringExpr) *constant.CharArray {
	text := strings.Replace(expr.Val, "\\n", "\n", -1) + "\x00"
	str := llvmStr(text)
	return str
}

func genBooleanExpr(expr *ast.BooleanExpr) *constant.Int {
	return constant.NewBool(expr.Val)
}

// Literals end

func (cg *LLVMCodeGenerator) genLogicalExpr(expr *ast.LogicalExpr) value.Value {

	lhs := cg.genExpr(expr.Lhs)
	rhs := cg.genExpr(expr.Rhs)

	block := cg.getCurrentBlock()

	switch expr.Op {
	case ast.AND:
		return block.NewAnd(lhs, rhs)
	case ast.OR:
		return block.NewOr(lhs, rhs)
	default:
		return nil
	}
}

func (cg *LLVMCodeGenerator) genBinaryExpr(expr *ast.BinaryExpr) value.Value {

	lhs := cg.genExpr(expr.Lhs)
	rhs := cg.genExpr(expr.Rhs)

	block := cg.getCurrentBlock()

	switch expr.Op {
	case ast.ADD:
		return block.NewAdd(lhs, rhs)
	case ast.SUB:
		return block.NewSub(lhs, rhs)
	case ast.MUL:
		return block.NewMul(lhs, rhs)
	case ast.DIV:
		return block.NewSDiv(lhs, rhs)

	// relational
	case "<":
		return block.NewICmp(enum.IPredSLT, lhs, rhs)
	case "<=":
		return block.NewICmp(enum.IPredSLE, lhs, rhs)
	case ">":
		return block.NewICmp(enum.IPredSGT, lhs, rhs)
	case ">=":
		return block.NewICmp(enum.IPredSGE, lhs, rhs)
	case "==":
		return block.NewICmp(enum.IPredEQ, lhs, rhs)
	case "!=":
		return block.NewICmp(enum.IPredNE, lhs, rhs)

	default:
		return nil
	}

}

func (cg *LLVMCodeGenerator) genCallExpr(expr *ast.CallExpr) value.Value {
	fn := cg.getFunction(expr.Callee.(*ast.IdentifierExpr).Name)
	if fn == nil {
		panic("Function not found")
	}

	args := []value.Value{}
	for _, arg := range expr.Args {
		val := cg.genExpr(arg)

		args = append(args, val)
	}
	block := cg.getCurrentBlock()

	if expr.Callee.(*ast.IdentifierExpr).Name == "print" && fn.Name() == "printf" {
		return cg.genPrintCall(fn, args...)
	}
	if fn.Name() == "printf" {
		return cg.genPrintfCall(fn, args...)
	}

	return block.NewCall(fn, args...)
}

func (cg *LLVMCodeGenerator) genPrintfCall(fn *ir.Func, args ...value.Value) *ir.InstCall {

	block := cg.getCurrentBlock()
	argList := []value.Value{}
	for _, arg := range args {

		switch a := arg.(type) {
		case *constant.CharArray:

			gep := getElementPtrFromString(block, a)

			argList = append(argList, gep)
		default:
			argList = append(argList, a)
		}
	}

	return block.NewCall(fn, argList...)

}

func (cg *LLVMCodeGenerator) genPrintCall(fn *ir.Func, args ...value.Value) *ir.InstCall {

	formatStr := ""

	for _, arg := range args {
		switch arg.(type) {
		case *constant.CharArray:
			formatStr = formatStr + "%s\n"

		default:
			formatStr = formatStr + "%d\n"

		}
	}

	formatPtr := llvmStr(formatStr)
	args = append([]value.Value{formatPtr}, args...)

	return cg.genPrintfCall(fn, args...)

}

func (cg *LLVMCodeGenerator) genIdentifierExpr(expr *ast.IdentifierExpr) value.Value {

	currentFunc := cg.getCurrentBlock().Parent
	if currentFunc != nil {
		var value *ir.Param
		for _, param := range currentFunc.Params {
			if param.Name() == expr.Name {
				value = param
				break
			}
		}
		if value != nil {
			return value
		}
	}

	value := cg.env.vars[expr.Name]
	block := cg.getCurrentBlock()
	load := block.NewLoad((value.ElemType), value)
	return load
}

func (cg *LLVMCodeGenerator) Gen(prog *ast.Program) string {
	m := cg.module

	fn := m.NewFunc("main", I32)
	fn.NewBlock("entry")
	cg.mainFunc = fn

	setupExternal(m)

	for _, stmt := range prog.Stmts {
		cg.genStmt(stmt)
	}

	block := cg.getCurrentBlock()
	if block.Term == nil {
		block.NewRet(constant.NewInt(I32, 0))
	}

	// terrible hack but for now will sort main func to be placed at the bottom
	// because stuff that I put in the main is stuff from global scope
	m.Funcs = m.Funcs[1:]
	m.Funcs = append(m.Funcs, fn)

	return m.String()
}

// Helpers
func setupExternal(m *ir.Module) {
	f := m.NewFunc("printf", I32, ir.NewParam("", Str))
	f.Sig.Variadic = true

}

func (cg *LLVMCodeGenerator) getCurrentBlock() *ir.Block {
	currentBlock := cg.currentBlock
	if currentBlock == nil {
		currentBlock = cg.mainFunc.Blocks[0]
	}
	return currentBlock
}

func (cg *LLVMCodeGenerator) getExitBlock() *ir.Block {
	return cg.exitBlock
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

func llvmType(t string) types.Type {
	switch t {
	case "int":
		return I32
	case "string":
		return Str
	case "bool":
		return Bool
	default:
		return Void
	}
}

func getElementPtrFromString(block *ir.Block, str *constant.CharArray) *ir.InstGetElementPtr {
	zero := constant.NewInt(I32, 0)
	strLen := str.Typ.Len
	strPtr := block.NewAlloca((types.NewArray(strLen, Char)))
	block.NewStore(str, strPtr)
	gep := block.NewGetElementPtr(str.Typ, strPtr, zero, zero)
	return gep
}

func llvmStr(s string) *constant.CharArray {
	s = strings.Replace(s, "\\n", "\n", -1) + "\x00"
	return constant.NewCharArrayFromString(s)
}
