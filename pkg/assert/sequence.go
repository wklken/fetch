/* MIT License
 * Copyright (c) 2012-2020 Mat Ryer, Tyler Bunnell and contributors.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/* NOTE: copied from https://github.com/stretchr/testify/assert/assertion_compare.go and modified
 *  The original versions of the files are MIT licensed
 */

package assert

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/stretchr/testify/assert"
)

func Contains(s, contains interface{}) (bool, string) {
	ok, found := includeElement(s, contains)
	if !ok {
		return false, "contains error"
	}
	if !found {
		return false, fmt.Sprintf("contains, sequence=`%v`, contains=`%v`", prettyLine(s), contains)
	}

	return true, "OK"
}

func NotContains(s, contains interface{}) (bool, string) {
	ok, found := includeElement(s, contains)
	if !ok {
		return false, "not contains error"
	}
	if found {
		return false, fmt.Sprintf("not contains, sequence=`%v`, not_contains=`%v`", prettyLine(s), contains)
	}

	return true, "OK"
}

func In(element, s interface{}) (bool, string) {
	ok, _ := Contains(s, element)
	if !ok {
		// if string
		return false, fmt.Sprintf("in, element=`%v`, sequence=`%v`", element, prettyLine(s))
	}

	return true, "OK"
}

func NotIn(element, s interface{}) (bool, string) {
	ok, _ := NotContains(s, element)
	if !ok {
		return false, fmt.Sprintf("not_in, element=`%v`, sequence=`%v`", element, prettyLine(s))
	}

	return true, "OK"
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
