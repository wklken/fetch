package client

import "net/http"

// from gin
func filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

// ContentType returns the Content-Type header of the request.
func GetContentType(header http.Header) string {
	return filterFlags(header.Get("Content-Type"))
}
