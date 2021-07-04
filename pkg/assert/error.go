package assert

import "fmt"

func NoError(err error) bool {
	if err != nil {
		fmt.Println("FAIL: got an error")
		return false
	}
	return true
}
