package assert

import (
	"reflect"
)

// Greater asserts that the first element is greater than the second
func Greater(e1 interface{}, e2 interface{}) bool {
	e1Kind := reflect.ValueOf(e1).Kind()
	e2Kind := reflect.ValueOf(e2).Kind()
	if e1Kind != e2Kind {
		Fail("FAIL: greater error, elemtnes should be the same type")
		return false
	}

	res, isComparable := compare(e1, e2, e1Kind)
	if !isComparable {
		Fail("FAIL: Can not compare type \"%s\"\n", reflect.TypeOf(e1))
		return false
	}

	if res != -1 {
		Fail("FAIL: greater, \"%v\" is not greater than \"%v\"\n", e1, e2)
		return false
	}

	OK()
	return true
}

func GreaterOrEqual(e1 interface{}, e2 interface{}) bool {
	e1Kind := reflect.ValueOf(e1).Kind()
	e2Kind := reflect.ValueOf(e2).Kind()
	if e1Kind != e2Kind {
		Fail("FAIL: greater or equal error, elements should be the same type")
		return false
	}

	res, isComparable := compare(e1, e2, e1Kind)
	if !isComparable {
		Fail("FAIL: Can not compare type \"%s\"\n", reflect.TypeOf(e1))
		return false
	}

	if res != -1 && res != 0 {
		Fail("FAIL: greater or equal, \"%v\" is not greater than or equal to \"%v\"\n", e1, e2)
		return false
	}

	OK()
	return true
}

func Less(e1 interface{}, e2 interface{}) bool {
	e1Kind := reflect.ValueOf(e1).Kind()
	e2Kind := reflect.ValueOf(e2).Kind()
	if e1Kind != e2Kind {
		Fail("FAIL: less error, elements should be the same type")
		return false
	}

	res, isComparable := compare(e1, e2, e1Kind)
	if !isComparable {
		Fail("FAIL: Can not compare type \"%s\"\n", reflect.TypeOf(e1))
		return false
	}

	if res != 1 {
		Fail("FAIL: less, \"%v\" is not less than \"%v\"\n", e1, e2)
		return false
	}

	OK()
	return true
}

func LessOrEqual(e1 interface{}, e2 interface{}) bool {
	e1Kind := reflect.ValueOf(e1).Kind()
	e2Kind := reflect.ValueOf(e2).Kind()
	if e1Kind != e2Kind {
		Fail("FAIL: less or equal error, elements should be the same type")
		return false
	}

	res, isComparable := compare(e1, e2, e1Kind)
	if !isComparable {
		Fail("FAIL: Can not compare type \"%s\"\n", reflect.TypeOf(e1))
		return false
	}

	if res != 1 && res != 0 {
		Fail("FAIL: less or equal, \"%v\" is not less than or equal to \"%v\"\n", e1, e2)
		return false
	}

	OK()
	return true
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
