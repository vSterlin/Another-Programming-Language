package interpreter

type ReturnValue struct {
	value any
}

func NewReturnValue(value any) *ReturnValue {
	return &ReturnValue{value: value}
}

func (r *ReturnValue) Value() any {
	return r.value
}
