package typechecker

import "errors"

type Type int

const (
	Number Type = iota
	String
	Boolean

	INVALID
)

func (t Type) String() string {
	switch t {
	case Number:
		return "number"
	case String:
		return "string"
	case Boolean:
		return "boolean"
	default:
		return "invalid"
	}
}

// variadic function
func expectTypesEqual(expected Type, actual ...Type) bool {
	for _, a := range actual {
		if expected != a {
			return false
		}
	}
	return true
}

func NewTypeError(text string) error {
	return errors.New("type error: " + text)
}
