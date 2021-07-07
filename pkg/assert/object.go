package assert

import (
	"github.com/stretchr/testify/assert"
)

func Equal(expected interface{}, actual interface{}) {
	equal := assert.ObjectsAreEqual(expected, actual)
	if !equal {
		Fail("FAIL: not equal, expected=%v, actual=%v\n", expected, actual)
	}
}
