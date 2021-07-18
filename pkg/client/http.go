package client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wklken/httptest/pkg/config"

	"github.com/wklken/httptest/pkg/log"
)

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
) (resp *http.Response, latency int64, err error) {
	// NOTE: if c.Request.Body begin with `@`, means it's a file
	requestBody, err := parseBodyIfGotAFile(caseDir, body)
	if err != nil {
		//logRunCaseFail(path, &c, "Read body file content fail: body=@%s err=%s", c.Request.Body, err)
		//stats.failCaseCount += 1
		return
	}

	httpMethod := strings.ToUpper(method)

	var req *http.Request

	switch httpMethod {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		req, err = http.NewRequest(httpMethod, url, nil)
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		if hasBody {
			req, err = http.NewRequest(httpMethod, url, strings.NewReader(requestBody))
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

	if !auth.Empty() {
		req.SetBasicAuth(auth.Username, auth.Password)
	}
	//req.AddCookie(&http.Cookie{
	//	Name:  "aaa",
	//	Value: "123",
	//})
	//fmt.Println("see cookies", req.Cookies())

	// dump request, for debug
	if debug {
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			log.Info("DEBUG request: dump err %s", err)
		} else {
			log.Info("DEBUG request: \n%s", prettyFormatDump(dump, "> "))
		}
	}

	// do send
	start := time.Now()

	jar, err := cookiejar.New(nil)
	if err != nil {
		return
	}
	client := &http.Client{
		Jar: jar,
	}
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	latency = time.Since(start).Milliseconds()

	//fmt.Println("cookie jar:", jar)
	//fmt.Println("response cookies:", resp.Cookies())
	//for _, ck := range resp.Cookies() {
	//	ck.String()
	//}

	if hook.SaveCookie != "" {
		cookies := jar.Cookies(resp.Request.URL)
		//fmt.Println("save cookies:", cookies)
		bs, err1 := json.Marshal(cookies)
		if err1 != nil {
			err = err1
			return
		}
		cookiePath := filepath.Join(caseDir, hook.SaveCookie)
		//log.Info("saved cookie into %s", cookiePath)
		err = ioutil.WriteFile(cookiePath, bs, 0644)
		if err != nil {
			return
		}
	}

	// dump request, for debug
	if debug {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.Info("DEBUG response: dump err %s", err)
		} else {
			log.Info("DEBUG request: \n%s", prettyFormatDump(dump, "< "))
		}
	}

	return
}
