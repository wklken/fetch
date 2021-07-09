/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"httptest/pkg/assert"
	"httptest/pkg/config"
	"httptest/pkg/util"
	"io"
	"net/http"
	"strings"

	"github.com/fatih/color"

	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("args required")
			return
		}
		//path := args[0]
		for _, path := range args {

			run(path)
		}

		fmt.Print("done")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var (
	Info = color.New(color.FgWhite).PrintfFunc()
	Tip  = color.New(color.FgYellow).PrintfFunc()
)

func run(path string) {

	v, err := config.ReadFromFile(path)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	var c config.Case
	err = v.Unmarshal(&c)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	allKeys := util.NewStringSetWithValues(v.AllKeys())
	//fmt.Println("allKeys", allKeys)

	Tip("Run Case: %s | %s | [%s %s]\n", path, c.Title, strings.ToUpper(c.Request.Method), c.Request.URL)
	//fmt.Printf("the case and data: %s, %+v\n", path, c)

	var resp *http.Response
	var err1 error
	if c.Request.Method == "get" {
		resp, err1 = http.Get(c.Request.URL)
	}
	// TODO: post / head / put / delete / Patch
	//else if c.Request.Method == "post" {
	//	resp, err1 = http.Post(c.Request.URL)
	//}

	assert.NoError(err1)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(err)

	// TODO:
	//   1. POST
	//   2. json response assert
	//   3. content-type assert
	//   4. latency assert
	//   5. set timeout=x, each case?
	//   6. `-e env.toml` support envs => can render

	bodyStr := string(body)

	contentType := GetContentType(resp.Header)

	type Ctx struct {
		f        assert.AssertFunc
		element1 interface{}
		element2 interface{}
	}

	// TODO: how to keep the order!!!!!!
	keyAssertFuncs := map[string]Ctx{
		// statuscode
		"assert.statuscode": {
			f:        assert.Equal,
			element1: resp.StatusCode,
			element2: c.Assert.StatusCode,
		},
		"assert.statuscode_lt": {
			f:        assert.Less,
			element1: resp.StatusCode,
			element2: c.Assert.StatusCodeLt,
		},
		"assert.statuscode_lte": {
			f:        assert.LessOrEqual,
			element1: resp.StatusCode,
			element2: c.Assert.StatusCodeLte,
		},
		"assert.statuscode_gt": {
			f:        assert.Greater,
			element1: resp.StatusCode,
			element2: c.Assert.StatusCodeGt,
		},
		"assert.statuscode_gte": {
			f:        assert.GreaterOrEqual,
			element1: resp.StatusCode,
			element2: c.Assert.StatusCodeGte,
		},
		"assert.statuscode_in": {
			f:        assert.In,
			element1: resp.StatusCode,
			element2: c.Assert.StatusCodeIn,
		},
		// status
		"assert.status": {
			f:        assert.Equal,
			element1: strings.ToLower(http.StatusText(resp.StatusCode)),
			element2: strings.ToLower(c.Assert.Status),
		},
		// TODO: status_in
		"assert.contenttype": {
			f:        assert.Equal,
			element1: strings.ToLower(contentType),
			element2: strings.ToLower(c.Assert.ContentType),
		},
		// TODO: contentType_in

		// contentlength
		"assert.contentlength": {
			f:        assert.Equal,
			element1: resp.ContentLength,
			element2: c.Assert.ContentLength,
		},
		"assert.contentlength_lt": {
			f:        assert.Less,
			element1: resp.ContentLength,
			element2: c.Assert.ContentLengthLt,
		},
		"assert.contentlength_lte": {
			f:        assert.LessOrEqual,
			element1: resp.ContentLength,
			element2: c.Assert.ContentLengthLte,
		},
		"assert.contentlength_gt": {
			f:        assert.Greater,
			element1: resp.ContentLength,
			element2: c.Assert.ContentLengthGt,
		},
		"assert.contentlength_gte": {
			f:        assert.GreaterOrEqual,
			element1: resp.ContentLength,
			element2: c.Assert.ContentLengthGte,
		},
		// body
		"assert.body": {
			f:        assert.Equal,
			element1: bodyStr,
			element2: c.Assert.Body,
		},
		"assert.body_contains": {
			f:        assert.Contains,
			element1: bodyStr,
			element2: c.Assert.BodyContains,
		},
		"assert.body_not_contains": {
			f:        assert.NotContains,
			element1: bodyStr,
			element2: c.Assert.BodyNotContains,
		},
		"assert.body_startswith": {
			f:        assert.StartsWith,
			element1: bodyStr,
			element2: c.Assert.BodyStartsWith,
		},
		"assert.body_endswith": {
			f:        assert.EndsWith,
			element1: bodyStr,
			element2: c.Assert.BodyEndsWith,
		},
	}

	for key, ctx := range keyAssertFuncs {
		if allKeys.Has(key) {
			Info("%s: ", key)
			ctx.f(ctx.element1, ctx.element2)
		}
	}

	// TODO: =============================================

	// parse
	//fmt.Println("the response content type", resp.Header.Get("Content-Type"))
	// content-type: text/html; charset=utf-8
	// TODO: 如果是application/json, 直接转成json path? assert?
	b := Default("post", GetContentType(resp.Header))
	if b != nil {
		var i map[string]interface{}
		err = b.BindBody(body, &i)
		assert.NoError(err)

		//fmt.Println("the json body", i)
	}

	// headers
	//resp.Proto
	//resp.ProtoMajor
	//resp.ProtoMinor

	//contentType := resp.Header.Get("Content-type")
	//dump, err := httputil.DumpResponse(resp, true)
	// HTTP/1.1 200 OK
	//\r\nContent-Length: 76
	//\r\nContent-Type: text/plain; charset=utf-8
	//\r\nDate: Wed, 19 Jul 1972 19:00:00 GMT
	//\r\n\r\nGo is a general-purpose language designed with systems programming in mind."

	//fmt.Println("all headers:", resp.Header)
	//all headers: map[
	//Access-Control-Allow-Credentials:[true]
	//Access-Control-Allow-Origin:[*]
	//Content-Length:[18]
	//Content-Type:[text/html; charset=utf-8]
	//Date:[Mon, 05 Jul 2021 15:48:57 GMT] Server:[gunicorn/19.9.0]]

	//< Content-Type: text/css
	//< Content-Length: 7832

	//< x-proxy-by: SmartGate-IDC
	//< set-cookie: x-client-ssid=17a7755884a-545e86e6b94f752cc79eafa76f816a8c6886e4b0; path=/; domain=.oa.com; HttpOnly
	//< set-cookie: x-host-key-front=17a77558864-76daf86f1e36a8d8f09e66fe6873bd540953e430; path=/; domain=.oa.com; HttpOnly
	//< set-cookie: x_host_key=17a7755885f-b5ac4643c31f9b89e5d10625aee04b93bdebc0df; path=/; domain=.oa.com; HttpOnly
	//< set-cookie: x-host-key-ngn=17a7755884a-5ad277cc537abfe3a26b8ec683f670c8fa9ff0a0; path=/; domain=.oa.com; HttpOnly
	//< x-forwarded-for: 9.146.99.128,9.19.161.39,10.14.87.133,9.218.225.9
	//< Date: Mon, 05 Jul 2021 15:42:12 GMT
	//< Connection: keep-alive
	//< Vary: Accept-Encoding
	//< Last-Modified: Sun, 25 Apr 2021 09:27:01 GMT
	//< ETag: "608535e5-1e98"
	//< Accept-Ranges: bytes
	//< x-rio-seq: kqqskkfq-147822534

}

func Default(method, contentType string) binding.BindingBody {
	//if method == http.MethodGet {
	//	return binding.Form
	//}

	switch contentType {
	case binding.MIMEJSON:
		return binding.JSON
	case binding.MIMEXML, binding.MIMEXML2:
		return binding.XML
	case binding.MIMEPROTOBUF:
		return binding.ProtoBuf
	case binding.MIMEMSGPACK, binding.MIMEMSGPACK2:
		return binding.MsgPack
	case binding.MIMEYAML:
		return binding.YAML
		//case binding.MIMEMultipartPOSTForm:
		//	return binding.FormMultipart
		//default: // case MIMEPOSTForm:
		//	return binding.Form
	}
	return nil
}

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

//https://golang.org/src/net/http/status.go
