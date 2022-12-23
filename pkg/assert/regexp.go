package assert

import (
	"fmt"
	"regexp"
)

func regexpMatch(text string, expr string) bool {
	r := regexp.MustCompile(expr)
	return r.FindStringIndex(text) != nil
}

func Regexp(text, expr interface{}) (bool, string) {
	if test(text, expr, regexpMatch) {
		return true, "OK"
	}

	// if regexpMatch(expr, text) {
	// 	return true, "OK"
	// }
	return false, fmt.Sprintf("regexp | `%v` should match `%v`", prettyLine(text), expr)
}

func NotRegexp(text, expr interface{}) (bool, string) {
	if !test(text, expr, regexpMatch) {
		return true, "OK"
	}
	// if !regexpMatch(expr, text) {
	// 	return true, "OK"
	// }

	return false, fmt.Sprintf("not_regexp | `%v` should not match `%v`", prettyLine(text), expr)
}
