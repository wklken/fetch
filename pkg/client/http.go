package client

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/wklken/httptest/pkg/config"

	"github.com/wklken/httptest/pkg/log"
)

func Send(
	method string,
	url string,
	hasBody bool, body string,
	headers map[string]string,
	cookie string,
	auth config.BasicAuth,
	debug bool,
) (resp *http.Response, latency int64, err error) {
	httpMethod := strings.ToUpper(method)

	var req *http.Request

	switch httpMethod {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		req, err = http.NewRequest(httpMethod, url, nil)
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		if hasBody {
			req, err = http.NewRequest(httpMethod, url, strings.NewReader(body))
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
		req.Header.Set("Cookie", cookie)
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
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	latency = time.Since(start).Milliseconds()

	//fmt.Println("response cookies:", resp.Cookies())

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
