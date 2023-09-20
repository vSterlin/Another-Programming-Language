package interpreter

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
)

type RuntimeNumber int
type RuntimeBoolean bool
type RuntimeString string

func (r RuntimeNumber) String() string  { return color.YellowString(strconv.Itoa(int(r))) }
func (r RuntimeBoolean) String() string { return color.YellowString(strconv.FormatBool(bool(r))) }
func (r RuntimeString) String() string  { return color.GreenString(fmt.Sprintf("\"%s\"", string(r))) }
