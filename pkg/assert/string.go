package assert

import (
	"fmt"
	"reflect"
	"strings"
)

type stringTestFunc func(s, i string) bool

func test(s, i interface{}, testFunc stringTestFunc) bool {
	sType := reflect.TypeOf(s)
	if sType == nil {
		return false
	}
	iType := reflect.TypeOf(i)
	if iType == nil {
		return false
	}

	if sType.Kind() != reflect.String || iType.Kind() != reflect.String {
		return false
	}

	return testFunc(s.(string), i.(string))
}

func StartsWith(s, prefix interface{}) (bool, string) {
	if test(s, prefix, strings.HasPrefix) {
		return true, "OK"
	}

	return false, fmt.Sprintf("startswith, string=`%v`, prefix=`%v`", prettyLine(s), prefix)
}

func EndsWith(s, suffix interface{}) (bool, string) {
	if test(s, suffix, strings.HasSuffix) {
		return true, "OK"
	}

	return false, fmt.Sprintf("endswith, string=`%v`, suffix=`%v`", prettyLine(s), suffix)
}

func NotStartsWith(s, prefix interface{}) (bool, string) {
	if test(s, prefix, strings.HasPrefix) {
		return false, fmt.Sprintf("not_startswith, string=`%v`, prefix=`%v`", prettyLine(s), prefix)
	}

	return true, "OK"
}

func NotEndsWith(s, suffix interface{}) (bool, string) {
	if test(s, suffix, strings.HasSuffix) {
		return false, fmt.Sprintf("not_endswith, string=`%v`, suffix=`%v`", prettyLine(s), suffix)
	}

	return true, "OK"
}
