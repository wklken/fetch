package assertion

import (
	"reflect"

	"github.com/jmespath/go-jmespath"
	"github.com/wklken/httptest/pkg/assert"
	"github.com/wklken/httptest/pkg/config"
	"github.com/wklken/httptest/pkg/util"
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

			// fmt.Printf("%T, %T", actualValue, expectedValue)
			// make float64 compare with int64
			if reflect.TypeOf(actualValue).Kind() == reflect.Float64 && reflect.TypeOf(expectedValue).Kind() == reflect.Int64 {
				actualValue = int64(actualValue.(float64))
			}

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
