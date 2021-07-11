package assert

import (
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

	if sType.Kind() == reflect.String || iType.Kind() == reflect.String {
		return false
	}

	return testFunc(s.(string), i.(string))
}

func StartsWith(s, prefix interface{}) bool {
	if test(s, prefix, strings.HasPrefix) {
		OK()
		return true
	}

	Fail("FAIL: startswith, string=%v, prefix=%v\n", s, prefix)
	return false
}

func EndsWith(s, suffix interface{}) bool {
	if test(s, suffix, strings.HasPrefix) {
		OK()
		return true
	}

	Fail("FAIL: endswith, string=%v, suffix=%v\n", s, suffix)
	return false
}
