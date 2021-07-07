package assert

import (
	"strings"
)

func StartsWith(s, prefix string) bool {
	if !strings.HasPrefix(s, prefix) {
		Fail("FAIL: startswith, string=%v, prefix=%v\n", s, prefix)
		return false
	}
	return true

}

func EndsWith(s, suffix string) bool {
	if !strings.HasSuffix(s, suffix) {
		Fail("FAIL: endswith, string=%v, suffix=%v\n", s, suffix)
		return false
	}
	return true
}
