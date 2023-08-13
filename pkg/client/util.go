package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"

	"github.com/wklken/fetch/pkg/util"
)

const (
	maxResponseBodyLength = 2048
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

// GetContentType returns the Content-Type header of the request.
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

	return strings.Join(newLines, "\n")
}

func dumpRequest(req *http.Request) (logs []string) {
	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		logs = append(logs, fmt.Sprintf("dump request fail: %s", err))
	} else {
		logs = append(logs, fmt.Sprintf("%s", prettyFormatDump(dump, "> ")))
	}
	return
}

// from: https://github.com/henvic/httpretty/blob/master/printer.go

var binaryMediatypes = map[string]struct{}{
	"application/pdf":               {},
	"application/postscript":        {},
	"image":                         {}, // for practical reasons, any image (including SVG) is considered binary data
	"audio":                         {},
	"application/ogg":               {},
	"video":                         {},
	"application/vnd.ms-fontobject": {},
	"font":                          {},
	"application/x-gzip":            {},
	"application/zip":               {},
	"application/x-rar-compressed":  {},
	"application/wasm":              {},
}

func isBinaryMediatype(mediatype string) bool {
	if _, ok := binaryMediatypes[mediatype]; ok {
		return true
	}

	if parts := strings.SplitN(mediatype, "/", 2); len(parts) == 2 {
		if _, ok := binaryMediatypes[parts[0]]; ok {
			return true
		}
	}

	return false
}

func dumpResponse(resp *http.Response) (logs []string) {
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		logs = append(logs, fmt.Sprintf("dump response fail: %s", err))
	} else {

		if contentType := resp.Header.Get("Content-Type"); contentType != "" && isBinaryMediatype(contentType) {
			logs = append(logs, fmt.Sprintf("response: * body contains binary data"))
			return
		}

		respLines := prettyFormatDump(dump, "< ")

		// fmt.Println("the contentLength:", resp.ContentLength)
		if resp.ContentLength > maxResponseBodyLength || len(respLines) > maxResponseBodyLength {
			actualLength := resp.ContentLength
			if actualLength == -1 {
				actualLength = int64(len(respLines))
			}

			logs = append(logs, fmt.Sprintf("response: * body is too long (%d bytes) to print, skipping  (longer than %d bytes)\n", actualLength, maxResponseBodyLength))
			logs = append(logs, fmt.Sprintf("%s", util.TruncateString(respLines, maxResponseBodyLength)))
			return
		}

		logs = append(logs, fmt.Sprintf("%s", prettyFormatDump(dump, "< ")))
	}
	return logs
}

func parseBodyIfGotAFile(caseDir string, body string) (content string, err error) {
	content = body
	if body != "" && strings.HasPrefix(body, "@") {
		bodyFilePath := strings.TrimPrefix(body, "@")

		// NOTE: should be relative path to the `path`
		realBodyFilePath := filepath.Join(caseDir, bodyFilePath)

		if _, err = os.Stat(realBodyFilePath); os.IsNotExist(err) {
			return
		}

		var dat []byte
		dat, err = os.ReadFile(realBodyFilePath)
		if err != nil {
			return
		}

		content = string(dat)
	}

	return content, nil
}

func parseCookiesIfGotAFile(caseDir string, cookie string) (cookies []*http.Cookie, err error) {
	cookieFilePath := strings.TrimPrefix(cookie, "@")

	realCookieFilePath := filepath.Join(caseDir, cookieFilePath)

	if _, err = os.Stat(realCookieFilePath); os.IsNotExist(err) {
		return
	}

	var dat []byte
	dat, err = os.ReadFile(realCookieFilePath)
	if err != nil {
		return
	}

	err = json.Unmarshal(dat, &cookies)
	return
}

func saveCookies(caseDir string, savedPath string, cookieJsonBytes []byte, resp *http.Response) (err error) {
	// cookies := jar.Cookies(resp.Request.URL)
	// var bs []byte
	// bs, err = json.Marshal(cookies)
	// if err != nil {
	// 	return
	// }
	cookiePath := filepath.Join(caseDir, savedPath)
	// log.Info("saved cookie into %s", cookiePath)
	err = os.WriteFile(cookiePath, cookieJsonBytes, 0o644)
	if err != nil {
		return
	}
	return
}

func saveResponseBody(caseDir string, savedPath string, respBody []byte) (err error) {
	respBodyPath := filepath.Join(caseDir, savedPath)
	err = os.WriteFile(respBodyPath, respBody, 0o644)
	return
}
