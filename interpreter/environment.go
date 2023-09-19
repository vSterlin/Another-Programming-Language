package interpreter

type Environment struct {
	values map[string]any
	parent *Environment
}

func NewEnvironment(parent *Environment) *Environment {
	valueMap := map[string]any{}
	return &Environment{values: valueMap, parent: parent}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Assign(name string, value any) error {
	_, ok := e.values[name]
	if !ok && e.parent == nil {
		return NewRuntimeError("undefined variable: " + name)
	}

	if ok {
		e.values[name] = value
		return nil
	} else {
		return e.parent.Assign(name, value)
	}
}

func (e *Environment) Get(name string) (any, error) {
	_, ok := e.values[name]
	if ok {
		return e.values[name], nil
	} else {
		if e.parent != nil {
			return e.parent.Get(name)
		} else {
			return nil, NewRuntimeError("undefined variable: " + name)
		}
	}
}
