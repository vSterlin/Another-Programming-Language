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
	class *Class
}

func NewInstance(class *Class) *Instance {
	return &Instance{class: class}
}

func (i *Instance) String() string {
	return color.BlueString(fmt.Sprintf("<instance of %s>", i.class.name))
}
