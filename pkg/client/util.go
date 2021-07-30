package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"

	"github.com/wklken/httptest/pkg/util"

	"github.com/wklken/httptest/pkg/log"
)

const (
	maxResponseBodyLength = 1024
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
func dumpRequest(debug bool, req *http.Request) {
	// dump request, for debug
	if debug {
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			log.Info("DEBUG request: dump err %s", err)
		} else {
			log.Info("DEBUG request: \n%s", prettyFormatDump(dump, "> "))
		}
	}
}

// from: https://github.com/henvic/httpretty/blob/master/printer.go

var binaryMediatypes = map[string]struct{}{
	"application/pdf":               struct{}{},
	"application/postscript":        struct{}{},
	"image":                         struct{}{}, // for practical reasons, any image (including SVG) is considered binary data
	"audio":                         struct{}{},
	"application/ogg":               struct{}{},
	"video":                         struct{}{},
	"application/vnd.ms-fontobject": struct{}{},
	"font":                          struct{}{},
	"application/x-gzip":            struct{}{},
	"application/zip":               struct{}{},
	"application/x-rar-compressed":  struct{}{},
	"application/wasm":              struct{}{},
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

func dumpResponse(debug bool, resp *http.Response) {
	// dump request, for debug
	if debug {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.Info("DEBUG response: dump err %s", err)
		} else {

			if contentType := resp.Header.Get("Content-Type"); contentType != "" && isBinaryMediatype(contentType) {
				log.Info("DEBUG response: * body contains binary data")
				return
			}

			respLines := prettyFormatDump(dump, "< ")

			//fmt.Println("the contentLength:", resp.ContentLength)
			if resp.ContentLength > maxResponseBodyLength || len(respLines) > maxResponseBodyLength {
				actualLength := resp.ContentLength
				if actualLength == -1 {
					actualLength = int64(len(respLines))
				}

				log.Info("DEBUG response* body is too long (%d bytes) to print, skipping  (longer than %d bytes)\n", actualLength, maxResponseBodyLength)
				log.Info("%s", util.TruncateString(respLines, maxResponseBodyLength))
				return
			}

			log.Info("DEBUG response: \n%s", prettyFormatDump(dump, "< "))
		}
	}

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
		dat, err = ioutil.ReadFile(realBodyFilePath)
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
	dat, err = ioutil.ReadFile(realCookieFilePath)
	if err != nil {
		return
	}

	err = json.Unmarshal(dat, &cookies)
	return
}

func saveCookies(caseDir string, savedPath string, jar *cookiejar.Jar, resp *http.Response) (err error) {
	cookies := jar.Cookies(resp.Request.URL)
	var bs []byte
	bs, err = json.Marshal(cookies)
	if err != nil {
		return
	}
	cookiePath := filepath.Join(caseDir, savedPath)
	//log.Info("saved cookie into %s", cookiePath)
	err = ioutil.WriteFile(cookiePath, bs, 0644)
	if err != nil {
		return
	}
	return
}
