package assert

import (
	"fmt"
	"reflect"
	"regexp"
)

func RegexpMatch(text string, expr string) bool {
	r := regexp.MustCompile(expr)
	return r.FindStringIndex(text) != nil
}

func Matches(text, expr interface{}) (bool, string) {
	if test(text, expr, RegexpMatch) {
		return true, "OK"
	}

	return false, fmt.Sprintf("matches | `%v` should match `%v`", prettyLine(text), expr)
}

func NotMatches(text, expr interface{}) (bool, string) {
	if !test(text, expr, RegexpMatch) {
		return true, "OK"
	}

	return false, fmt.Sprintf("not_matches | `%v` should not match `%v`", prettyLine(text), expr)
}

func StringMatchesAll(text, exprs interface{}) (bool, string) {
	listValue := reflect.ValueOf(exprs)
	if reflect.TypeOf(exprs).Kind() != reflect.Slice {
		return false, fmt.Sprintf("matches_all| `%v` should be slice", prettyLine(exprs))
	}

	for i := 0; i < listValue.Len(); i++ {
		expr := listValue.Index(i).Interface()
		if !test(text, expr, RegexpMatch) {
			// FIXME: add index in error info
			return false, fmt.Sprintf("matches | `%v` should match `%v`", prettyLine(text), expr)
		}
	}

	return true, "OK"
}

func StringNotMatchesAll(text, exprs interface{}) (bool, string) {
	listValue := reflect.ValueOf(exprs)
	if reflect.TypeOf(exprs).Kind() != reflect.Slice {
		return false, fmt.Sprintf("not_matches_all| `%v` should be slice", prettyLine(exprs))
	}

	for i := 0; i < listValue.Len(); i++ {
		expr := listValue.Index(i).Interface()
		if test(text, expr, RegexpMatch) {
			// FIXME: add index in error info
			return false, fmt.Sprintf("not_matches | `%v` should not match `%v`", prettyLine(text), expr)
		}
	}

	return true, "OK"
}
