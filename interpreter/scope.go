package interpreter

type scope map[string]bool

func (s *scope) isDefined(name string) bool {
	_, ok := (*s)[name]
	return ok
}

func (s *scope) isInitialized(name string) bool {
	initialized, ok := (*s)[name]
	return ok && initialized
}

type scopeStack []*scope

func (s *scopeStack) peek() *scope {
	return (*s)[len(*s)-1]
}
func (s *scopeStack) push(scope *scope) {
	*s = append(*s, scope)
}

func (s *scopeStack) pop() *scope {
	scope, scopes := (*s)[len(*s)-1], (*s)[:len(*s)-1]
	*s = scopes
	return scope
}

func (s *scopeStack) isEmpty() bool {
	return len(*s) == 0
}
