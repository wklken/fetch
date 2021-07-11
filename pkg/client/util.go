package client

import (
	"fmt"
	"net/http"
	"strings"
)

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

func prettyFormatDump(dump []byte, linePrefix string) string {
	s := string(dump)

	parts := strings.Split(s, "\n")
	newLines := make([]string, 0, len(parts))
	for _, p := range parts {
		newLines = append(newLines, fmt.Sprintf("%s%s", linePrefix, p))
	}

	//fmt.Printf("DEBUG request: \n%s\n", dump)
	//fmt.Printf("DEBUG request: \n%s", strings.Join(newLines, "\n"))
	return strings.Join(newLines, "\n")
}
