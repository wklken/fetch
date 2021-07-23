package client

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

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
	hook config.Hook,
	debug bool,
) (resp *http.Response, hasRedirect bool, latency int64, err error) {
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
			//fmt.Printf("has body, the headers: %+v\n", headers)

			bodyReader := strings.NewReader(requestBody)

			// TODO: support msgpack? and the gzip?
			//if ct, ok := headers["content-type"]; ok {
			//	switch strings.ToLower(ct) {
			//	case binding.MIMEMSGPACK, binding.MIMEMSGPACK2:
			//		bs, err1 := msgpack.Marshal(requestBody)
			//		if err1 != nil {
			//			return fmt.Errorf("do msgpack encode fail")
			//
			//		}
			//
			//	}
			//}

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

	dumpRequest(debug, req)

	// do send
	start := time.Now()
	jar, err := cookiejar.New(nil)
	if err != nil {
		return
	}
	client := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			hasRedirect = true
			return nil
		},
	}
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	latency = time.Since(start).Milliseconds()
	//fmt.Println("hasRedirect: ", hasRedirect)

	if hook.SaveCookie != "" {
		err = saveCookies(caseDir, hook.SaveCookie, jar, resp)
		if err != nil {
			return
		}
	}

	dumpResponse(debug, resp)

	return
}
