package util

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
