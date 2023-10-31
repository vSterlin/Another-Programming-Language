package typechecker

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

func areTypesEqual(expected Type, actual ...Type) bool {
	for _, a := range actual {
		if expected != a {
			return false
		}
	}
	return true
}

type TypeError struct {
	text string
}

func (t TypeError) Error() string {
	return "type error: " + t.text
}

func NewTypeError(text string) error {
	return &TypeError{text}
}
