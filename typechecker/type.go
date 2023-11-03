package typechecker

import "strings"

// type Type int

// const (
// 	Number Type = iota
// 	String
// 	Boolean
// 	Void

// 	INVALID
// )

// func (t Type) String() string {
// 	switch t {
// 	case Number:
// 		return "number"
// 	case String:
// 		return "string"
// 	case Boolean:
// 		return "boolean"
// 	case Void:
// 		return "void"
// 	default:
// 		return "invalid"
// 	}
// }

// func fromString(s string) Type {
// 	switch s {
// 	case "number", "int":
// 		return Number
// 	case "string":
// 		return String
// 	case "boolean":
// 		return Boolean
// 	case "void":
// 		return Void
// 	default:
// 		return INVALID
// 	}
// }

type Type interface {
	String() string
	Equals(Type) bool
}

type NumberType struct{}
type StringType struct{}
type BooleanType struct{}
type VoidType struct{}
type ArrowFuncType struct {
	Args       []Type
	ReturnType Type
}

type InvalidType struct{}

func (t NumberType) String() string  { return "number" }
func (t StringType) String() string  { return "string" }
func (t BooleanType) String() string { return "boolean" }
func (t VoidType) String() string    { return "void" }
func (t ArrowFuncType) String() string {
	args := []string{}
	for _, arg := range t.Args {
		args = append(args, arg.String())
	}

	return "func(" + strings.Join(args, ", ") + ") " + t.ReturnType.String()
}
func (t InvalidType) String() string { return "invalid" }

func (t NumberType) Equals(other Type) bool {
	_, ok := other.(NumberType)
	return ok
}

func (t StringType) Equals(other Type) bool {
	_, ok := other.(StringType)
	return ok
}

func (t BooleanType) Equals(other Type) bool {
	_, ok := other.(BooleanType)
	return ok
}

func (t VoidType) Equals(other Type) bool {
	_, ok := other.(VoidType)
	return ok
}

func (t ArrowFuncType) Equals(other Type) bool {
	otherFuncType, ok := other.(ArrowFuncType)
	if !ok {
		return false
	}
	if len(t.Args) != len(otherFuncType.Args) {
		return false
	}
	for i, arg := range t.Args {
		if !arg.Equals(otherFuncType.Args[i]) {
			return false
		}
	}
	return t.ReturnType.Equals(otherFuncType.ReturnType)
}

func (t InvalidType) Equals(other Type) bool {
	_, ok := other.(InvalidType)
	return ok
}

func fromString(s string) Type {
	switch s {
	case "number", "int":
		return NumberType{}
	case "string":
		return StringType{}
	case "boolean":
		return BooleanType{}
	case "void":
		return VoidType{}
	default:
		return InvalidType{}
	}
}

func areTypesEqual(expected Type, actual ...Type) bool {
	for _, a := range actual {
		if !expected.Equals(a) {
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

var Number = NumberType{}
var String = StringType{}
var Boolean = BooleanType{}
var Void = VoidType{}
var Invalid = InvalidType{}
