package assert

import (
	"reflect"
	"strings"
)

type AssertFunc func(expected interface{}, actual interface{}) (bool, string)

func prettyLine(s interface{}) interface{} {
	if reflect.TypeOf(s).Kind() == reflect.Array || reflect.TypeOf(s).Kind() == reflect.Slice {
		// TODO: pretty print it

		//x := make([]string, 0, len(a))
		//for _, i := range a {
		//	x = append(x, fmt.Sprintf("%#v", i))
		//}
		//fmt.Printf("[%s]\n", strings.Join(x, ", "))
	}

	if reflect.TypeOf(s).Kind() != reflect.String {
		return s
	}

	return strings.ReplaceAll(s.(string), "\n", "\\n")
}
