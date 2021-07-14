package assert

import (
	"reflect"
	"strings"

	"github.com/fatih/color"
)

var Fail = color.New(color.FgRed).PrintfFunc()

func OK() {
	color.New(color.FgGreen).PrintfFunc()("OK\n")
}

type AssertFunc func(expected interface{}, actual interface{}) bool

func prettyLine(s interface{}) interface{} {
	if reflect.TypeOf(s).Kind() != reflect.String {
		return s
	}

	return strings.ReplaceAll(s.(string), "\n", "\\n")
}
