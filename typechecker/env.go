package typechecker

type Env struct {
	parent *Env
	vars   map[string]string
}

func NewEnv(parent *Env) *Env {
	return &Env{
		parent: parent,
		vars:   make(map[string]string),
	}
}
