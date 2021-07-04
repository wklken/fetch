package assert

import (
	"fmt"
	"strings"
)

func StartsWith(s, prefix string) bool {
	if !strings.HasPrefix(s, prefix) {
		fmt.Printf("FAIL: startswith, string=%v, prefix=%v\n", s, prefix)
		return false
	}
	return true

}

func EndsWith(s, suffix string) bool {
	if !strings.HasSuffix(s, suffix) {
		fmt.Printf("FAIL: endswith, string=%v, suffix=%v\n", s, suffix)
		return false
	}
	return true

}
