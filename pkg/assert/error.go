package assert

func NoError(err error) bool {
	if err != nil {
		Fail("FAIL: got an error", err)
		return false
	}
	//OK()
	return true
}
