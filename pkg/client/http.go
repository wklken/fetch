package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	net_url "net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/goccy/go-json"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/wklken/fetch/pkg/config"
	"github.com/wklken/fetch/pkg/version"
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
	maxRedirects int,
	hook config.Hook,
	timeout int,
	proxy string,
	debug bool,
) (resp *http.Response, respBody []byte, redirectCount int64, latency int64, debugLogs []string, err error) {
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
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", fmt.Sprintf("fetch/%s", version.Version))
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
			redirectCount += 1
			return http.ErrUseLastResponse
		}
	} else if maxRedirects > 0 {
		checkRedirect = func(req *http.Request, via []*http.Request) error {
			redirectCount += 1
			if redirectCount > int64(maxRedirects) {
				return http.ErrUseLastResponse
			}
			return nil
		}
	} else {
		checkRedirect = func(req *http.Request, via []*http.Request) error {
			redirectCount += 1
			return nil
		}
	}

	transport := &http.Transport{}
	if proxy != "" {
		parsedProxyUrl, err1 := net_url.Parse(proxy)
		if err1 != nil {
			err = err1
			return
		}
		transport.Proxy = http.ProxyURL(parsedProxyUrl)
	}
	client := &http.Client{
		Jar:           jar,
		CheckRedirect: checkRedirect,
		Timeout:       time.Duration(timeout) * time.Millisecond,
		Transport:     transport,
	}
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	latency = time.Since(start).Milliseconds()
	// fmt.Println("hasRedirect: ", hasRedirect)

	// get the cookies
	cookies := jar.Cookies(resp.Request.URL)
	var cookieJsonBytes []byte
	cookieJsonBytes, err = json.Marshal(cookies)
	if err != nil {
		// FIXME: the error should not make the request fail
		return
	}

	if hook.SaveCookie != "" {
		err = saveCookies(caseDir, hook.SaveCookie, cookieJsonBytes, resp)
		if err != nil {
			return
		}
	}

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	// reset the resp body
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(respBody))

	if hook.SaveResponse != "" {
		err = saveResponseBody(caseDir, hook.SaveResponse, respBody)
		if err != nil {
			return
		}
	}

	var headerJsonBytes []byte
	headerJsonBytes, err = json.Marshal(resp.Header)
	if err != nil {
		// FIXME: the error should not make the request fail
		return
	}

	// TODO: execute support get parsed result as context, return and merge with the context, for next case
	if hook.Exec != "" {
		processDir, _ := os.Getwd()
		execDir := filepath.Join(processDir, caseDir)
		err = executeCmd(hook.Exec, execDir, resp.StatusCode, headerJsonBytes, respBody, cookieJsonBytes, debug)
		if err != nil {
			return
		}
	}

	// FIXME: put it here? or put it after the assertions?
	if hook.Sleep > 0 {
		time.Sleep(time.Duration(hook.Sleep) * time.Millisecond)
	}

	if debug {
		debugLogs = append(debugLogs, dumpResponse(resp)...)
	}

	return
}

func executeCmd(
	command string,
	execDir string,
	statusCode int,
	headerJsonBytes []byte,
	respBody []byte,
	cookieJsonBytes []byte,
	debug bool,
) (err error) {
	finalCmd := fmt.Sprintf("%s $statusCode $headerJson $respBody $cookieJson", command)
	fmt.Println("finalCmd: ", finalCmd)

	cmd := exec.Command(
		"bash",
		"-c",
		finalCmd,
	)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("statusCode=%d", statusCode))
	cmd.Env = append(cmd.Env, fmt.Sprintf("headerJson=%s", string(headerJsonBytes)))
	cmd.Env = append(cmd.Env, fmt.Sprintf("respBody=%s", string(respBody)))
	cmd.Env = append(cmd.Env, fmt.Sprintf("cookieJson=%s", string(cookieJsonBytes)))
	cmd.Dir = execDir
	var o []byte
	o, err = cmd.Output()
	if err != nil {
		fmt.Println("\nrun shell fail, err=", err)
		return
	}

	if debug {
		fmt.Println("\nrun shell success, output=\n", string(o))
	}

	return nil
}
