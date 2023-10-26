package codegen

import (
	"language/ast"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/value"
)

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
	case *ast.WhileStmt:
		return cg.genWhileStmt(stmt)
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

	if stmt.Op == ":=" {

		initType := init.Type()
		alloc := block.NewAlloca(initType)
		block.NewStore(init, alloc)

		cg.env.vars[varName] = alloc
	} else {
		alloc := cg.env.vars[varName]
		block.NewStore(init, alloc)
	}

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

func (cg *LLVMCodeGenerator) genWhileStmt(stmt *ast.WhileStmt) value.Value {

	block := cg.getCurrentBlock()

	fn := block.Parent

	whileBlock := fn.NewBlock("")
	bodyBlock := fn.NewBlock("")
	exitBlock := fn.NewBlock("")

	prevExitBlock := cg.getExitBlock()
	cg.exitBlock = whileBlock

	block.NewBr(whileBlock)

	cg.currentBlock = whileBlock
	test := cg.genExpr(stmt.Test)
	whileBlock.NewCondBr(test, bodyBlock, exitBlock)

	cg.currentBlock = bodyBlock
	cg.genStmt(stmt.Body)

	if bodyBlock.Term == nil {
		bodyBlock.NewBr(whileBlock)
	}

	cg.currentBlock = exitBlock

	if exitBlock.Term == nil {
		if prevExitBlock != nil {
			exitBlock.NewBr(prevExitBlock)
		}
		// else {
		// 	exitBlock.NewBr(whileBlock)
		// }
	}
	return nil
}

func (cg *LLVMCodeGenerator) genReturnStmt(stmt *ast.ReturnStmt) value.Value {
	val := cg.genExpr(stmt.Arg)
	block := cg.getCurrentBlock()
	block.NewRet(val)
	return val
}
