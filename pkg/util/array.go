package util

func ToLower(ss []string) []string {
	l := make([]string, 0, len(ss))
	for _, s := range ss {
		l = append(l, s)
	}
	return l
}

func ItemInIntArray(item int, array []int) bool {
	for _, a := range array {
		if a == item {
			return true
		}
	}
	return false
}
