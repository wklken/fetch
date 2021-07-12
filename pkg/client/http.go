package client

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

func Send(
	method string,
	url string,
	hasBody bool, body string,
	headers map[string]string,
	debug bool,
) (resp *http.Response, latency int64, err error) {
	httpMethod := strings.ToUpper(method)

	var req *http.Request

	// TODO: get params => from url + params(append)
	//params := make(url.Values)
	//params.Add("key1", "value1")
	//params.Add("key2", "value2")
	//req.URL.RawQuery = params.Encode()

	switch httpMethod {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		req, err = http.NewRequest(httpMethod, url, nil)
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		if hasBody {
			req, err = http.NewRequest(httpMethod, url, strings.NewReader(body))
		} else {
			req, err = http.NewRequest(httpMethod, url, nil)
		}
	}
	if err != nil {
		fmt.Println("error: make request fail ", err)
		return
	}

	// set header
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	// dump request, for debug
	if debug {
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			fmt.Printf("DEBUG request: dump err %s\n", err)
		} else {
			fmt.Printf("DEBUG request: \n%s\n", prettyFormatDump(dump, "> "))
		}
	}
	//fmt.Printf("DEBUG request: \n%s\n", strings.ReplaceAll(string(dump), "\n", ">\r\n"))

	// do send
	start := time.Now()
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("error: do request fail ", err)
		return
	}

	latency = time.Since(start).Milliseconds()

	// dump request, for debug
	if debug {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Printf("DEBUG response: dump err %s\n", err)
		} else {
			fmt.Printf("DEBUG request: \n%s\n", prettyFormatDump(dump, "< "))
		}
	}

	return
}
