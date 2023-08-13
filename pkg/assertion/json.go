package assertion

import (
	"fmt"
	"reflect"

	"github.com/jmespath/go-jmespath"
	"github.com/wklken/fetch/pkg/assert"
	"github.com/wklken/fetch/pkg/config"
	"github.com/wklken/fetch/pkg/util"
)

func DoJSONAssertions(jsonData interface{}, jsons []config.AssertJSON) (stats util.Stats) {
	for _, dj := range jsons {
		path := dj.Path
		expectedValue := dj.Value
		stats.AddInfofMessage("assert.json.%s: ", path)

		if jsonData == nil {
			ok, message := assert.Equal(nil, expectedValue)
			if ok {
				stats.IncrOkAssertCount()
			} else {
				stats.AddFailMessage(message)
				stats.IncrFailAssertCount()
			}
			continue
		}

		actualValue, err := jmespath.Search(path, jsonData)
		if err != nil {
			// log.Fail("search json data fail, err=%s, path=%s, expected=%s", err, path, expectedValue)
			stats.AddFailMessage("search json data fail, err=%s, path=%s, expected=%s", err, path, expectedValue)
			stats.IncrFailAssertCount()
		} else {
			// missing
			if actualValue == nil {
				_, message := assert.Equal(nil, expectedValue)
				stats.AddFailMessage(message)
				stats.IncrFailAssertCount()
				continue
			}

			actualValueKind := reflect.TypeOf(actualValue).Kind()
			expectedValueKind := reflect.TypeOf(expectedValue).Kind()

			// cast to same type, float64 or int64
			if actualValueKind != expectedValueKind && isNumberKind(actualValueKind) && isNumberKind(expectedValueKind) {
				if isFloatKind(actualValueKind) || isFloatKind(expectedValueKind) {
					newActualValue, err := toFloat64(actualValue)
					if err == nil {
						actualValue = newActualValue
					}
					newExpectedValue, err := toFloat64(expectedValue)
					if err == nil {
						expectedValue = newExpectedValue
					}
				} else {
					newActualValue, err := toInt64(actualValue)
					if err == nil {
						actualValue = newActualValue
					}
					newExpectedValue, err := toInt64(expectedValue)
					if err == nil {
						expectedValue = newExpectedValue
					}
				}
			}

			// fmt.Printf("%T, %T", actualValue, expectedValue)
			// make float64 compare with int64
			// if reflect.TypeOf(actualValue).Kind() == reflect.Float64 && reflect.TypeOf(expectedValue).Kind() == reflect.Int64 {
			// 	actualValue = int64(actualValue.(float64))
			// }

			// not working there
			//#[[assert.json]]
			//#path = 'json.array[0:3]'
			//#value =  [1, 2, 3]

			ok, message := assert.Equal(actualValue, expectedValue)
			if ok {
				stats.AddPassMessage()
				stats.IncrOkAssertCount()
			} else {
				stats.AddFailMessage(message)
				stats.IncrFailAssertCount()
			}
		}
	}

	return
}

func isNumberKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int64,
		reflect.Float64,
		reflect.Int,
		reflect.Float32,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32:
		return true
	default:
		return false
	}
}

func isFloatKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Float64, reflect.Float32:
		return true
	default:
		return false
	}
}

func toInt64(i interface{}) (int64, error) {
	switch s := i.(type) {
	case int:
		return int64(s), nil
	case int64:
		return s, nil
	case int32:
		return int64(s), nil
	case int16:
		return int64(s), nil
	case int8:
		return int64(s), nil
	case uint:
		return int64(s), nil
	case uint64:
		// NOTE: precision lost
		return int64(s), nil
	case uint32:
		return int64(s), nil
	case uint16:
		return int64(s), nil
	case uint8:
		return int64(s), nil
		// NOTE: only cast between int*, no float32/float64
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int64", i, i)
	}
}

func toFloat64(i interface{}) (float64, error) {
	switch s := i.(type) {
	case float64:
		return s, nil
	case float32:
		return float64(s), nil
	case int:
		return float64(s), nil
	case int64:
		// NOTE: precision lost
		return float64(s), nil
	case int32:
		return float64(s), nil
	case int16:
		return float64(s), nil
	case int8:
		return float64(s), nil
	case uint:
		return float64(s), nil
	case uint64:
		// NOTE: precision lost
		return float64(s), nil
	case uint32:
		return float64(s), nil
	case uint16:
		return float64(s), nil
	case uint8:
		return float64(s), nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to float64", i, i)
	}
}
