package interpreter

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
)

// Runtime types for the interpreter to print more nicely

type Number int
type Boolean bool
type String string

func (r Number) String() string  { return color.YellowString(strconv.Itoa(int(r))) }
func (r Boolean) String() string { return color.YellowString(strconv.FormatBool(bool(r))) }
func (r String) String() string  { return color.GreenString(fmt.Sprintf("\"%s\"", string(r))) }
