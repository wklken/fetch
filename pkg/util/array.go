package util

import "strings"

func StringArrayToLower(ss []string) []string {
	l := make([]string, 0, len(ss))
	for _, s := range ss {
		l = append(l, strings.ToLower(s))
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
