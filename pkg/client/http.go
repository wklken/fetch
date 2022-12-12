package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/goccy/go-json"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/wklken/httptest/pkg/config"
)

func Send(
	caseDir string,
	method string,
	url string,
	hasBody bool, body string,
	headers map[string]string,
	cookie string,
	auth config.BasicAuth,
	disableRedirect bool,
	hook config.Hook,
	timeout int,
	debug bool,
) (resp *http.Response, hasRedirect bool, latency int64, debugLogs []string, err error) {
	// NOTE: if c.Request.Body begin with `@`, means it's a file
	requestBody, err := parseBodyIfGotAFile(caseDir, body)
	if err != nil {
		return
	}

	var req *http.Request
	httpMethod := strings.ToUpper(method)
	switch httpMethod {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		req, err = http.NewRequest(httpMethod, url, nil)
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		if hasBody {
			// fmt.Printf("has body, the headers: %+v\n", headers)

			var bodyReader io.Reader

			// TODO: support msgpack? and the gzip?
			if ct, ok := headers["content-type"]; ok {
				switch strings.ToLower(ct) {
				case binding.MIMEMSGPACK, binding.MIMEMSGPACK2:
					var obj map[string]interface{}
					err0 := json.Unmarshal([]byte(requestBody), &obj)
					if err0 != nil {
						err = fmt.Errorf("try to validate the request body valid json fail, %w", err0)
						return
					}

					bs, err1 := msgpack.Marshal(obj)
					if err1 != nil {
						err = fmt.Errorf("do msgpack encode fail, err=%w", err)
						return
					}
					bodyReader = bytes.NewReader(bs)
				default:
					bodyReader = strings.NewReader(requestBody)
				}
			} else {
				bodyReader = strings.NewReader(requestBody)
			}

			req, err = http.NewRequest(httpMethod, url, bodyReader)
		} else {
			req, err = http.NewRequest(httpMethod, url, nil)
		}
	default:
		err = errors.New("http method not supported yet")
		return
	}
	if err != nil {
		return
	}

	// set header
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	// set header: basic_auth
	if !auth.Empty() {
		req.SetBasicAuth(auth.Username, auth.Password)
	}

	// set cookie
	if cookie != "" {
		if !strings.HasPrefix(cookie, "@") {
			req.Header.Set("Cookie", cookie)
		} else {
			cookies, err1 := parseCookiesIfGotAFile(caseDir, cookie)
			if err1 != nil {
				err = err1
				return
			}

			for _, c := range cookies {
				req.AddCookie(c)
			}
		}
	}

	if debug {
		debugLogs = append(debugLogs, dumpRequest(req)...)
	}

	// do send
	start := time.Now()
	jar, err := cookiejar.New(nil)
	if err != nil {
		return
	}

	var checkRedirect func(req *http.Request, via []*http.Request) error
	if disableRedirect {
		checkRedirect = func(req *http.Request, via []*http.Request) error {
			hasRedirect = true
			return http.ErrUseLastResponse
		}
	} else {
		checkRedirect = func(req *http.Request, via []*http.Request) error {
			hasRedirect = true
			return nil
		}
	}

	client := &http.Client{
		Jar:           jar,
		CheckRedirect: checkRedirect,
		Timeout:       time.Duration(timeout) * time.Millisecond,
	}
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	latency = time.Since(start).Milliseconds()
	// fmt.Println("hasRedirect: ", hasRedirect)

	if hook.SaveCookie != "" {
		err = saveCookies(caseDir, hook.SaveCookie, jar, resp)
		if err != nil {
			return
		}
	}

	if debug {
		debugLogs = append(debugLogs, dumpResponse(resp)...)
	}

	return
}
