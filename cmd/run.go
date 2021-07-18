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
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/jmespath/go-jmespath"
	"github.com/spf13/cobra"

	"github.com/wklken/httptest/pkg/assert"
	"github.com/wklken/httptest/pkg/client"
	"github.com/wklken/httptest/pkg/config"
	"github.com/wklken/httptest/pkg/log"
	"github.com/wklken/httptest/pkg/util"
)

const tableTPL = `
┌─────────────────────────┬─────────────────┬─────────────────┬─────────────────┐
│                         │           total │              ok │            fail │
├─────────────────────────┼─────────────────┼─────────────────┼─────────────────┤
│                   cases │          %6d │          %6d │          %6d │
├─────────────────────────┼─────────────────┼─────────────────┼─────────────────┤
│              assertions │          %6d │          %6d │          %6d │
├─────────────────────────┴─────────────────┴─────────────────┴─────────────────┤
│ total run duration: %6d ms                                                 │
└───────────────────────────────────────────────────────────────────────────────┘`
const (
	DebugEnvName = "HTTPTEST_DEBUG"
)

var (
	verbose = false
	quiet   = false
	cfgFile string
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run cases",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Error("args required")
			os.Exit(1)
			return
		}

		var runConfig config.RunConfig
		if cfgFile != "" {
			if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
				log.Error("config file not exists: %s", err)
				os.Exit(1)
				return
			}
			cv, err := config.ReadFromFile(cfgFile)
			if err != nil {
				log.Error("read config file fail: path=%s, err=%s", cfgFile, err)
				return
			}
			err = cv.Unmarshal(&runConfig)
			if err != nil {
				log.Error("parse config file fail: path=%s, err=%s", cfgFile, err)
				return
			}

			log.Info("runConfig: %v", runConfig)
		}

		totalStats := Stats{}

		log.BeQuiet(quiet)

		start := time.Now()
		for _, path := range args {
			s := run(path, &runConfig)
			totalStats.Add(s)

			// if got fail assert, the case is fail
			if s.failAssertCount > 0 {
				totalStats.failCaseCount += 1
			} else {
				totalStats.okCaseCount += 1
			}
		}
		latency := time.Since(start).Milliseconds()

		log.Info(tableTPL,
			len(args), totalStats.okCaseCount, totalStats.failCaseCount,
			totalStats.okAssertCount+totalStats.failAssertCount, totalStats.okAssertCount, totalStats.failAssertCount,
			latency)
		if totalStats.failCaseCount > 0 {
			log.Info("the execute result: 1")
			os.Exit(1)
		} else {
			log.Info("the execute result: 0")
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// -v verbose
	runCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode")
	// -q quiet
	runCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "be quiet")

	// -e dev.toml
	runCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file(like dev.toml/prod.toml")
}

type Stats struct {
	okCaseCount     int64
	failCaseCount   int64
	okAssertCount   int64
	failAssertCount int64
}

func (s *Stats) Add(s1 Stats) {
	s.okAssertCount += s1.okAssertCount
	s.failAssertCount += s1.failAssertCount
}

func logRunCaseFail(path string, c *config.Case, format string, a ...interface{}) {
	log.Tip("Run Case: %s | %s | [%s %s]", path, c.Title, strings.ToUpper(c.Request.Method), c.Request.URL)
	log.Error(format, a...)
}

func run(path string, runConfig *config.RunConfig) (stats Stats) {
	v, err := config.ReadFromFile(path)
	if err != nil {
		log.Tip("Run Case: %s", path)
		log.Error("read fail: %s", err)
		stats.failCaseCount += 1
		return
	}
	var c config.Case
	err = v.Unmarshal(&c)
	if err != nil {
		log.Tip("Run Case: %s", path)
		log.Error("parse fail: %s", err)
		return
	}
	allKeys := util.NewStringSetWithValues(v.AllKeys())
	//fmt.Println("allKeys", allKeys)
	//fmt.Printf("the case and data: %s, %+v", path, c)

	// do render
	if runConfig.Render && len(runConfig.Env) > 0 {
		c.Render(runConfig.Env)
	}

	debug := (verbose || strings.ToLower(os.Getenv(DebugEnvName)) == "true" || runConfig.Debug) && !quiet

	resp, latency, err := client.Send(
		filepath.Dir(path),
		c.Request.Method, c.Request.URL, allKeys.Has("request.body"), c.Request.Body, c.Request.Header, c.Request.Cookie, c.Request.BasicAuth, c.Hook, debug)
	if err != nil {
		logRunCaseFail(path, &c, "Send HTTP Request fail: %s", err)
		stats.failCaseCount += 1
		return
	}

	log.Tip("Run Case: %s | %s | [%s %s] | %dms", path, c.Title, strings.ToUpper(c.Request.Method), c.Request.URL, latency)

	stats = doAssertions(allKeys, resp, c, latency)
	return
}

func doAssertions(allKeys *util.StringSet, resp *http.Response, c config.Case, latency int64) (stats Stats) {

	body, err := io.ReadAll(resp.Body)
	// TODO: handle err
	assert.NoError(err)

	bodyStr := strings.TrimSuffix(string(body), "\n")
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
		"assert.statuscode_not_in": {
			f:        assert.NotIn,
			element1: resp.StatusCode,
			element2: c.Assert.StatusCodeNotIn,
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
		"assert.body_not_startswith": {
			f:        assert.NotStartsWith,
			element1: bodyStr,
			element2: c.Assert.BodyNotStartsWith,
		},
		"assert.body_not_endswith": {
			f:        assert.NotEndsWith,
			element1: bodyStr,
			element2: c.Assert.BodyNotEndsWith,
		},
	}

	for key, ctx := range keyAssertFuncs {
		if allKeys.Has(key) {
			log.Infof("%s: ", key)
			ok, message := ctx.f(ctx.element1, ctx.element2)
			if ok {
				log.OK()
				stats.okAssertCount += 1
			} else {
				log.Fail(message)
				stats.failAssertCount += 1
			}
		}
	}

	if len(c.Assert.Header) > 0 {
		for key, value := range c.Assert.Header {
			log.Infof("assert.header.%s: ", key)
			ok, message := assert.Equal(resp.Header.Get(key), value)
			if ok {
				log.OK()
				stats.okAssertCount += 1
			} else {
				log.Fail(message)
				stats.failAssertCount += 1
			}

		}
	}

	var jsonData interface{}
	if contentType == binding.MIMEJSON {
		err = binding.JSON.BindBody(body, &jsonData)
		if err != nil {
			log.Fail("binding.json fail: %s", err)
			stats.failAssertCount += int64(len(c.Assert.Json))
			return
		}

		if allKeys.Has("assert.json") && len(c.Assert.Json) > 0 {
			s1 := doJsonAssertions(jsonData, c.Assert.Json)
			stats.Add(s1)
		}
	}

	//   5. set timeout=x, each case?

	return
}

func doJsonAssertions(jsonData interface{}, jsons []config.AssertJson) (stats Stats) {
	for _, dj := range jsons {
		path := dj.Path
		expectedValue := dj.Value
		log.Infof("assert.json.%s: ", path)

		if jsonData == nil {
			ok, message := assert.Equal(nil, expectedValue)
			if ok {
				stats.okAssertCount += 1
			} else {
				log.Fail(message)
				stats.failAssertCount += 1
			}
			continue
		}

		actualValue, err := jmespath.Search(path, jsonData)
		if err != nil {
			log.Fail("search json data fail, err=%s, path=%s, expected=%s", err, path, expectedValue)
		} else {
			// missing
			if actualValue == nil {
				_, message := assert.Equal(nil, expectedValue)
				log.Fail(message)
				stats.failAssertCount += 1
				continue
			}

			//fmt.Printf("%T, %T", actualValue, expectedValue)
			// make float64 compare with int64
			if reflect.TypeOf(actualValue).Kind() == reflect.Float64 && reflect.TypeOf(expectedValue).Kind() == reflect.Int64 {
				actualValue = int64(actualValue.(float64))
			}

			// not working there
			//#[[assert.json]]
			//#path = 'json.array[0:3]'
			//#value =  [1, 2, 3]

			ok, message := assert.Equal(actualValue, expectedValue)
			if ok {
				log.OK()
				stats.okAssertCount += 1
			} else {
				log.Fail(message)
				stats.failAssertCount += 1
			}
		}
	}

	return
}
