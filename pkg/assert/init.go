package assert

import (
	"fmt"
	"reflect"
	"strings"
)

type AssertFunc func(expected interface{}, actual interface{}) (bool, string)

func prettyLine(s interface{}) interface{} {
	if reflect.TypeOf(s).Kind() == reflect.Array || reflect.TypeOf(s).Kind() == reflect.Slice {
		s := reflect.ValueOf(s)

		x := make([]string, 0, s.Len())
		for i := 0; i < s.Len(); i++ {
			x = append(x, fmt.Sprintf("%#v", s.Index(i)))
		}
		return fmt.Sprintf("[%s]", strings.Join(x, ", "))
	}

	if reflect.TypeOf(s).Kind() != reflect.String {
		return s
	}

	return strings.ReplaceAll(s.(string), "\n", "\\n")
}
