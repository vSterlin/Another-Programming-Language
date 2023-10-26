package codegen

import (
	"language/ast"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
)

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

	switch val := value.(type) {

	case *ir.InstAlloca:
		load := block.NewLoad((val.ElemType), value)
		return load
	case *ir.Global:
		load := block.NewLoad((val.Typ.ElemType), value)
		return load
	default:
		return value
	}

}
