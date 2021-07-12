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
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/jmespath/go-jmespath"

	"github.com/gin-gonic/gin/binding"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"httptest/pkg/assert"
	"httptest/pkg/client"
	"httptest/pkg/config"
	"httptest/pkg/util"
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

		//fmt.Println("done")
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

const (
	DebugEnvName = "HTTPTEST_DEBUG"
)

func run(path string) {

	//fmt.Println(os.Getenv(DebugEnvName), strings.ToUpper(os.Getenv(DebugEnvName)))
	debug := strings.ToLower(os.Getenv(DebugEnvName)) == "true"

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
	//fmt.Printf("the case and data: %s, %+v\n", path, c)

	resp, latency, err := client.Send(
		c.Request.Method, c.Request.URL, allKeys.Has("request.body"), c.Request.Body, c.Request.Header, debug)
	if err != nil {
		Tip("Run Case: %s | %s | [%s %s]\n", path, c.Title, strings.ToUpper(c.Request.Method), c.Request.URL)
		fmt.Println(err)
	}

	Tip("Run Case: %s | %s | [%s %s] | %dms\n", path, c.Title, strings.ToUpper(c.Request.Method), c.Request.URL, latency)

	doAssertions(allKeys, resp, c, latency)
}

func doAssertions(allKeys *util.StringSet, resp *http.Response, c config.Case, latency int64) {
	body, err := io.ReadAll(resp.Body)
	// TODO: handle err
	assert.NoError(err)

	bodyStr := string(body)
	contentType := client.GetContentType(resp.Header)

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

		// latency
		"assert.latency_lt": {
			f:        assert.Less,
			element1: latency,
			element2: c.Assert.LatencyLt,
		},
		"assert.latency_lte": {
			f:        assert.LessOrEqual,
			element1: latency,
			element2: c.Assert.LatencyLte,
		},
		"assert.latency_gt": {
			f:        assert.Greater,
			element1: latency,
			element2: c.Assert.LatencyGt,
		},
		"assert.latency_gte": {
			f:        assert.GreaterOrEqual,
			element1: latency,
			element2: c.Assert.LatencyGte,
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
			// TODO: break or not?
			ctx.f(ctx.element1, ctx.element2)
		}
	}

	var jsonData interface{}
	if contentType == binding.MIMEJSON {
		err = binding.JSON.BindBody(body, &jsonData)
		if err != nil {
			// TODO
			fmt.Println("binding.json fail", err)
			// ?
			return
		}
	}

	if allKeys.Has("assert.json") && len(c.Assert.Json) > 0 {
		doJsonAssertions(jsonData, c.Assert.Json)
	}

	//   6. `-e env.toml` support envs => can render
	//   5. set timeout=x, each case?

	// TODO: =============================================
}

func doJsonAssertions(jsonData interface{}, jsons []config.AssertJson) {
	for _, dj := range jsons {
		path := dj.Path
		expectedValue := dj.Value
		Info("assert.json.%s: ", path)

		if jsonData == nil {
			assert.Equal(nil, expectedValue)
			continue
		}

		actualValue, err := jmespath.Search(path, jsonData)
		if err != nil {
			assert.Fail("search json data fail, path=%s, expected=%s\n", err, path, expectedValue)
		} else {

			//fmt.Printf("%T, %T", actualValue, expectedValue)
			// make float64 compare with int64
			if reflect.TypeOf(actualValue).Kind() == reflect.Float64 && reflect.TypeOf(expectedValue).Kind() == reflect.Int64 {
				actualValue = int64(actualValue.(float64))
			}

			// not working there
			//#[[assert.json]]
			//#path = 'json.array[0:3]'
			//#value =  [1, 2, 3]

			assert.Equal(actualValue, expectedValue)
		}
	}

}
