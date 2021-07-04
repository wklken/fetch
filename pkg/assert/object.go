package assert

import (
	"fmt"

	"github.com/stretchr/testify/assert"
)

func Equal(expected interface{}, actual interface{}) {
	equal := assert.ObjectsAreEqual(expected, actual)
	if !equal {
		fmt.Printf("FAIL: not equal, expected=%v, actual=%v\n", expected, actual)
		//fmt.Println()
	}
}
