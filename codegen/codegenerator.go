package codegen

import (
	"language/ast"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
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

	env *Env
}

type Env struct {
	vars    map[string]value.Value
	strings map[string]*ir.Global
	parent  *Env
}

func NewEnv(parent *Env) *Env {
	return &Env{
		vars:    make(map[string]value.Value),
		strings: make(map[string]*ir.Global),
		parent:  parent,
	}
}

func NewLLVMCodeGenerator() *LLVMCodeGenerator {
	module := ir.NewModule()
	return &LLVMCodeGenerator{
		module: module,
		env:    NewEnv(nil),
	}
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
