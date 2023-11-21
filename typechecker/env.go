package typechecker

type Env struct {
	parent    *Env
	vars      map[string]Type
	functions map[string]FuncType
}

func NewEnv(parent *Env) *Env {
	return &Env{
		parent: parent,
		vars:   make(map[string]Type),
	}
}

func (e *Env) Define(name string, t Type) {
	e.vars[name] = t
}

func (e *Env) Assign(name string, t Type) error {
	_, err, foundEnv := e.Get(name)

	if err == nil {
		return err
	}

	foundEnv.vars[name] = t

	return nil
}

func (e *Env) Get(name string) (Type, error, *Env) {
	t, ok := e.vars[name]

	if ok {
		return t, nil, e
	}

	if e.parent != nil {
		return e.parent.Get(name)
	}

	return Invalid, NewTypeError("undefined variable: " + name), nil

}
