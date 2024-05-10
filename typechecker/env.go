package typechecker

type Env struct {
	parent *Env
	vars   map[string]Type
	types  map[string]Type
}

func NewEnv(parent *Env) *Env {
	return &Env{
		parent: parent,
		vars:   make(map[string]Type),
		types:  defaultTypes(),
	}
}

func (e *Env) Define(name string, t Type) {
	e.vars[name] = t
}

func (e *Env) Assign(name string, t Type) error {
	_, foundEnv, err := e.Get(name)

	if err == nil {
		return err
	}

	foundEnv.vars[name] = t

	return nil
}

func (e *Env) Get(name string) (Type, *Env, error) {
	t, ok := e.vars[name]

	if ok {
		return t, e, nil
	}

	if e.parent != nil {
		return e.parent.Get(name)
	}

	return Invalid, nil, NewTypeError("undefined variable: " + name)

}

func defaultTypes() map[string]Type {
	return map[string]Type{
		"int":    Number,
		"number": Number,
		"string": String,
		"bool":   Boolean,
		"void":   Void,
	}
}

func (e *Env) DefineType(name string, t Type) {
	e.types[name] = t
}

func (e *Env) ResolveType(name string) (Type, error) {
	t, ok := e.types[name]

	if ok {
		return t, nil
	} else if e.parent != nil {
		return e.parent.ResolveType(name)
	}

	return Invalid, NewTypeError("undefined type: " + name)
}

func GetGlobalFuncReturnType(name string) (Type, bool) {
	switch name {
	case "print":
		return Void, true
	default:
		return Invalid, false
	}
}
