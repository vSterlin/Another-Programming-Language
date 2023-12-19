package typechecker

import (
	"fmt"
	"language/ast"
	"strings"
)

type Type interface {
	String() string
	Equals(Type) bool
	IsFunc() bool
}

type NumberType struct{}
type StringType struct{}
type BooleanType struct{}
type VoidType struct{}
type FuncType struct {
	Args       []Type
	ReturnType Type
}

type InvalidType struct{}

func (t NumberType) String() string  { return "number" }
func (t StringType) String() string  { return "string" }
func (t BooleanType) String() string { return "boolean" }
func (t VoidType) String() string    { return "void" }
func (t FuncType) String() string {
	args := []string{}
	for _, arg := range t.Args {
		args = append(args, arg.String())
	}

	return fmt.Sprintf("func(%s) => %s", strings.Join(args, ", "), t.ReturnType.String())
}
func (t InvalidType) String() string { return "invalid" }

func (t NumberType) IsFunc() bool  { return false }
func (t StringType) IsFunc() bool  { return false }
func (t BooleanType) IsFunc() bool { return false }
func (t VoidType) IsFunc() bool    { return false }
func (t FuncType) IsFunc() bool    { return true }
func (t InvalidType) IsFunc() bool { return false }

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

func (t FuncType) Equals(other Type) bool {
	otherFuncType, ok := other.(FuncType)
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

// to set return type for func calls
func toAstNode(t Type) *ast.TypeExpr {

	switch t := t.(type) {
	case NumberType:
		return &ast.TypeExpr{Type: &ast.IdentifierExpr{Name: "number"}}
	case StringType:
		return &ast.TypeExpr{Type: &ast.IdentifierExpr{Name: "string"}}
	case BooleanType:
		return &ast.TypeExpr{Type: &ast.IdentifierExpr{Name: "boolean"}}
	case VoidType:
		return &ast.TypeExpr{Type: &ast.IdentifierExpr{Name: "void"}}
	case FuncType:
		args := []*ast.TypeExpr{}
		for _, arg := range t.Args {
			args = append(args, toAstNode(arg))
		}
		return &ast.TypeExpr{
			Type: &ast.FuncTypeExpr{
				Args:       args,
				ReturnType: toAstNode(t.ReturnType),
			},
		}
	default:
		panic("invalid type")

	}

}

// it also modifies astNode
func resolveType(astNode *ast.TypeExpr, env *Env) (Type, error) {
	switch nodeType := astNode.Type.(type) {
	case *ast.IdentifierExpr:
		t, err := env.ResolveType(nodeType.Name)
		if err != nil {
			return Invalid, err
		}

		// TODO: maybe do it somewhere else
		astNode.Type = toAstNode(t).Type

		return t, nil
	case *ast.FuncTypeExpr:
		args := []Type{}
		for _, arg := range nodeType.Args {
			t, err := resolveType(arg, env)
			if err != nil {
				return Invalid, err
			}
			args = append(args, t)
			// TODO: maybe do it somewhere else
			arg.Type = toAstNode(t).Type
		}
		retType, err := resolveType(nodeType.ReturnType, env)
		if err != nil {
			return Invalid, err
		}
		// TODO: maybe do it somewhere else
		nodeType.ReturnType.Type = toAstNode(retType).Type

		fmt.Printf("%#v", FuncType{
			Args:       args,
			ReturnType: retType,
		})

		return FuncType{
			Args:       args,
			ReturnType: retType,
		}, nil
	}

	panic("invalid type")
	return Invalid, nil
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
