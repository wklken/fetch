package util

import "strings"

// TruncateBytes truncate []byte to specific length
func TruncateBytes(content []byte, length int) []byte {
	if len(content) > length {
		return content[:length]
	}
	return content
}

// TruncateBytesToString ...
func TruncateBytesToString(content []byte, length int) string {
	s := TruncateBytes(content, length)
	return string(s)
}

// TruncateString truncate string to specific length
func TruncateString(s string, n int) string {
	if n > len(s) {
		return s
	}
	return s[:n]
}

func OmitMiddle(s string, head int, tail int) string {
	if len(s) <= head+tail {
		return s
	}

	return s[:head] + "..." + s[len(s)-tail:]
}

func PrettyStringSlice(s []string) string {
	return "[" + strings.Join(s, ", ") + "]"
}

func StringArrayMapFunc(elements []string, f func(string) string) (result []string) {
	for _, element := range elements {
		result = append(result, f(element))
	}
	return
}
