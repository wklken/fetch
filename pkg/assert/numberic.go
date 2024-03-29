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
)

// Greater asserts that the first element is greater than the second
func Greater(e1 interface{}, e2 interface{}) (bool, string) {
	e1Kind := reflect.ValueOf(e1).Kind()
	e2Kind := reflect.ValueOf(e2).Kind()
	if e1Kind != e2Kind {
		return false, "greater error, elemtnes should be the same type"
	}

	res, isComparable := compare(e1, e2, e1Kind)
	if !isComparable {
		return false, fmt.Sprintf("Can not compare type `%s`", reflect.TypeOf(e1))
	}

	if res != -1 {
		return false, fmt.Sprintf("greater | `%v` is not greater than `%v`", e1, e2)
	}

	return true, "OK"
}

func GreaterOrEqual(e1 interface{}, e2 interface{}) (bool, string) {
	e1Kind := reflect.ValueOf(e1).Kind()
	e2Kind := reflect.ValueOf(e2).Kind()
	if e1Kind != e2Kind {
		return false, "greater or equal error, elements should be the same type"
	}

	res, isComparable := compare(e1, e2, e1Kind)
	if !isComparable {
		return false, fmt.Sprintf("Can not compare type `%s`", reflect.TypeOf(e1))
	}

	if res != -1 && res != 0 {
		return false, fmt.Sprintf("greater_or_equal | `%v` is not greater than or equal to `%v`", e1, e2)
	}

	return true, "OK"
}

func Less(e1 interface{}, e2 interface{}) (bool, string) {
	e1Kind := reflect.ValueOf(e1).Kind()
	e2Kind := reflect.ValueOf(e2).Kind()
	if e1Kind != e2Kind {
		return false, "less error, elements should be the same type"
	}

	res, isComparable := compare(e1, e2, e1Kind)
	if !isComparable {
		return false, fmt.Sprintf("Can not compare type `%s`", reflect.TypeOf(e1))
	}

	if res != 1 {
		return false, fmt.Sprintf("less | `%v` is not less than `%v`", e1, e2)
	}

	return true, "OK"
}

func LessOrEqual(e1 interface{}, e2 interface{}) (bool, string) {
	e1Kind := reflect.ValueOf(e1).Kind()
	e2Kind := reflect.ValueOf(e2).Kind()
	if e1Kind != e2Kind {
		return false, "less or equal error, elements should be the same type"
	}

	res, isComparable := compare(e1, e2, e1Kind)
	if !isComparable {
		return false, fmt.Sprintf("Can not compare type `%s`", reflect.TypeOf(e1))
	}

	if res != 1 && res != 0 {
		return false, fmt.Sprintf("less_or_equal | `%v` is not less than or equal to `%v`", e1, e2)
	}

	return true, "OK"
}

func compare(obj1, obj2 interface{}, kind reflect.Kind) (int, bool) {
	switch kind {
	case reflect.Int:
		{
			intobj1 := obj1.(int)
			intobj2 := obj2.(int)
			if intobj1 > intobj2 {
				return -1, true
			}
			if intobj1 == intobj2 {
				return 0, true
			}
			if intobj1 < intobj2 {
				return 1, true
			}
		}
	case reflect.Int8:
		{
			int8obj1 := obj1.(int8)
			int8obj2 := obj2.(int8)
			if int8obj1 > int8obj2 {
				return -1, true
			}
			if int8obj1 == int8obj2 {
				return 0, true
			}
			if int8obj1 < int8obj2 {
				return 1, true
			}
		}
	case reflect.Int16:
		{
			int16obj1 := obj1.(int16)
			int16obj2 := obj2.(int16)
			if int16obj1 > int16obj2 {
				return -1, true
			}
			if int16obj1 == int16obj2 {
				return 0, true
			}
			if int16obj1 < int16obj2 {
				return 1, true
			}
		}
	case reflect.Int32:
		{
			int32obj1 := obj1.(int32)
			int32obj2 := obj2.(int32)
			if int32obj1 > int32obj2 {
				return -1, true
			}
			if int32obj1 == int32obj2 {
				return 0, true
			}
			if int32obj1 < int32obj2 {
				return 1, true
			}
		}
	case reflect.Int64:
		{
			int64obj1 := obj1.(int64)
			int64obj2 := obj2.(int64)
			if int64obj1 > int64obj2 {
				return -1, true
			}
			if int64obj1 == int64obj2 {
				return 0, true
			}
			if int64obj1 < int64obj2 {
				return 1, true
			}
		}
	case reflect.Uint:
		{
			uintobj1 := obj1.(uint)
			uintobj2 := obj2.(uint)
			if uintobj1 > uintobj2 {
				return -1, true
			}
			if uintobj1 == uintobj2 {
				return 0, true
			}
			if uintobj1 < uintobj2 {
				return 1, true
			}
		}
	case reflect.Uint8:
		{
			uint8obj1 := obj1.(uint8)
			uint8obj2 := obj2.(uint8)
			if uint8obj1 > uint8obj2 {
				return -1, true
			}
			if uint8obj1 == uint8obj2 {
				return 0, true
			}
			if uint8obj1 < uint8obj2 {
				return 1, true
			}
		}
	case reflect.Uint16:
		{
			uint16obj1 := obj1.(uint16)
			uint16obj2 := obj2.(uint16)
			if uint16obj1 > uint16obj2 {
				return -1, true
			}
			if uint16obj1 == uint16obj2 {
				return 0, true
			}
			if uint16obj1 < uint16obj2 {
				return 1, true
			}
		}
	case reflect.Uint32:
		{
			uint32obj1 := obj1.(uint32)
			uint32obj2 := obj2.(uint32)
			if uint32obj1 > uint32obj2 {
				return -1, true
			}
			if uint32obj1 == uint32obj2 {
				return 0, true
			}
			if uint32obj1 < uint32obj2 {
				return 1, true
			}
		}
	case reflect.Uint64:
		{
			uint64obj1 := obj1.(uint64)
			uint64obj2 := obj2.(uint64)
			if uint64obj1 > uint64obj2 {
				return -1, true
			}
			if uint64obj1 == uint64obj2 {
				return 0, true
			}
			if uint64obj1 < uint64obj2 {
				return 1, true
			}
		}
	case reflect.Float32:
		{
			float32obj1 := obj1.(float32)
			float32obj2 := obj2.(float32)
			if float32obj1 > float32obj2 {
				return -1, true
			}
			if float32obj1 == float32obj2 {
				return 0, true
			}
			if float32obj1 < float32obj2 {
				return 1, true
			}
		}
	case reflect.Float64:
		{
			float64obj1 := obj1.(float64)
			float64obj2 := obj2.(float64)
			if float64obj1 > float64obj2 {
				return -1, true
			}
			if float64obj1 == float64obj2 {
				return 0, true
			}
			if float64obj1 < float64obj2 {
				return 1, true
			}
		}
	case reflect.String:
		{
			stringobj1 := obj1.(string)
			stringobj2 := obj2.(string)
			if stringobj1 > stringobj2 {
				return -1, true
			}
			if stringobj1 == stringobj2 {
				return 0, true
			}
			if stringobj1 < stringobj2 {
				return 1, true
			}
		}
	}

	return 0, false
}
