package typechecker

type Env struct {
	parent    *Env
	vars      map[string]Type
	functions map[string]FuncType
}

func NewEnv(parent *Env) *Env {
	return &Env{
		parent:    parent,
		vars:      make(map[string]Type),
		functions: make(map[string]FuncType),
	}
}

func (e *Env) Define(name string, t Type) {
	e.vars[name] = t
}

func (e *Env) Assign(name string, t Type) error {
	_, ok := e.vars[name]

	if !ok {
		return NewTypeError("undefined variable: " + name)
	}

	e.vars[name] = t

	return nil
}

func (e *Env) Get(name string) (Type, error) {
	t, ok := e.vars[name]

	if ok {
		return t, nil
	}

	if e.parent != nil {
		return e.parent.Get(name)
	}

	return Invalid, NewTypeError("undefined variable: " + name)

}

func (e *Env) DefineFunction(name string, stmt FuncType) {
	e.functions[name] = stmt
}

func (e *Env) GetFunction(name string) (FuncType, error) {
	f, ok := e.functions[name]

	if ok {
		return f, nil
	}

	if e.parent != nil {
		return e.parent.GetFunction(name)
	}

	return FuncType{}, NewTypeError("undefined function: " + name)
}
