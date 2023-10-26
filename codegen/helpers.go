package codegen

import (
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

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
