package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Env struct {
	vars    map[string]value.Value
	strings map[string]*ir.Global
	types   map[string]types.Type
	parent  *Env
}

func NewEnv(parent *Env) *Env {
	return &Env{
		vars:    make(map[string]value.Value),
		strings: make(map[string]*ir.Global),
		types:   make(map[string]types.Type),
		parent:  parent,
	}
}

func (env *Env) Set(name string, val value.Value, typ types.Type) {

	env.vars[name] = val
	env.types[name] = typ
}
