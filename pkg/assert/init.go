package assert

import (
	"reflect"
	"strings"
)

type AssertFunc func(expected interface{}, actual interface{}) (bool, string)

func prettyLine(s interface{}) interface{} {
	if reflect.TypeOf(s).Kind() != reflect.String {
		return s
	}

	return strings.ReplaceAll(s.(string), "\n", "\\n")
}
