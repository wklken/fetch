package assert

import (
	"fmt"
	"httptest/pkg/util"

	"github.com/stretchr/testify/assert"
)

func Equal(actual interface{}, expected interface{}) bool {
	equal := assert.ObjectsAreEqual(actual, expected)
	if !equal {
		actualStr := util.TruncateString(fmt.Sprintf("%v", actual), 100)
		Fail("FAIL: not equal, expected=%v, actual=%v\n", expected, actualStr)
		return false
	} else {
		OK()
		return true
	}

}
