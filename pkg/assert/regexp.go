package assert

import (
	"fmt"
	"regexp"
)

func regexpMatch(text string, expr string) bool {
	r := regexp.MustCompile(expr)
	return r.FindStringIndex(text) != nil
}

func Matches(text, expr interface{}) (bool, string) {
	if test(text, expr, regexpMatch) {
		return true, "OK"
	}

	return false, fmt.Sprintf("matches | `%v` should match `%v`", prettyLine(text), expr)
}

func NotMatches(text, expr interface{}) (bool, string) {
	if !test(text, expr, regexpMatch) {
		return true, "OK"
	}

	return false, fmt.Sprintf("not_matches | `%v` should not match `%v`", prettyLine(text), expr)
}
