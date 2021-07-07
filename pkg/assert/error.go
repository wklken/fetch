package assert

func NoError(err error) bool {
	if err != nil {
		Fail("FAIL: got an error")
		return false
	}
	return true
}
