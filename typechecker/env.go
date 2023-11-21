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
