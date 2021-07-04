package assert

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/stretchr/testify/assert"
)

func Contains(s, contains interface{}) bool {
	ok, found := includeElement(s, contains)
	if !ok {
		fmt.Println("FAIL: contains error")
		return false
	}
	if !found {
		fmt.Printf("FAIL: contains, list=%v, contains=%v\n", s, contains)
		return false
	}

	return true
}

func NotContains(s, contains interface{}) bool {
	ok, found := includeElement(s, contains)
	if !ok {
		fmt.Println("FAIL: not contains error")
		return false
	}
	if found {
		fmt.Printf("FAIL: not contains, list=%v, not_contains=%v\n", s, contains)
		return false
	}

	return true
}

func In(element, s interface{}) bool {
	return Contains(s, element)
}

func NotIn(element, s interface{}) bool {
	return !Contains(s, element)
}

// NOTE: from testify
// containsElement try loop over the list check if the list includes the element.
// return (false, false) if impossible.
// return (true, false) if element was not found.
// return (true, true) if element was found.
func includeElement(list interface{}, element interface{}) (ok, found bool) {

	listValue := reflect.ValueOf(list)
	elementValue := reflect.ValueOf(element)
	defer func() {
		if e := recover(); e != nil {
			ok = false
			found = false
		}
	}()

	if reflect.TypeOf(list).Kind() == reflect.String {
		return true, strings.Contains(listValue.String(), elementValue.String())
	}

	if reflect.TypeOf(list).Kind() == reflect.Map {
		mapKeys := listValue.MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			if assert.ObjectsAreEqual(mapKeys[i].Interface(), element) {
				return true, true
			}
		}
		return true, false
	}

	for i := 0; i < listValue.Len(); i++ {
		if assert.ObjectsAreEqual(listValue.Index(i).Interface(), element) {
			return true, true
		}
	}
	return true, false

}
