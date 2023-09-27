package interpreter

import (
	"fmt"

	"github.com/fatih/color"
)

type Class struct {
	name string
}

func NewClass(name string) *Class {
	return &Class{name: name}
}

func (c *Class) String() string { return color.BlueString(fmt.Sprintf("<class %s>", c.name)) }

func (c *Class) Call(i *Interpreter, args []any) any {
	instance := NewInstance(c)
	return instance
}

type Instance struct {
	class  *Class
	fields map[string]any
}

func NewInstance(class *Class) *Instance {
	return &Instance{class: class, fields: map[string]any{}}
}

func (i *Instance) Get(name string) (any, error) {
	if val, ok := i.fields[name]; ok {
		return val, nil
	}
	return nil, NewRuntimeError(fmt.Sprintf("undefined property '%s'", name))
}

func (i *Instance) Set(name string, val any) {
	i.fields[name] = val
}

func (i *Instance) String() string {
	return color.BlueString(fmt.Sprintf("<instance of %s>", i.class.name))
}
