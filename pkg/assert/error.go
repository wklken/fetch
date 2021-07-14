package assert

import "fmt"

func NoError(err error) (bool, string) {
	if err != nil {
		return false, fmt.Sprintf("got an error, err=`%s`", err)
	}
	return true, ""
}
