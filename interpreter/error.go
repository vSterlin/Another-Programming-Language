package interpreter

type RuntimeError struct {
	Msg string
}

func (r *RuntimeError) Error() string {
	return r.Msg
}

func NewRuntimeError(msg string) *RuntimeError {
	return &RuntimeError{Msg: msg}
}
