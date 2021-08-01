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

			if runConfig.FailFast && s.failAssertCount > 0 {
				log.Info("failFast=True, quit, the execute result: 1")
				os.Exit(1)
			}

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

	resp, hasRedirect, latency, err := client.Send(
		filepath.Dir(path),
		c.Request.Method, c.Request.URL, allKeys.Has("request.body"), c.Request.Body, c.Request.Header, c.Request.Cookie, c.Request.BasicAuth, c.Hook, debug)
	if err != nil {
		logRunCaseFail(path, &c, "Send HTTP Request fail: %s", err)
		stats.failCaseCount += 1
		return
	}

	log.Tip("Run Case: %s | %s | [%s %s] | %dms", path, c.Title, strings.ToUpper(c.Request.Method), c.Request.URL, latency)

	stats = doAssertions(allKeys, resp, c, hasRedirect, latency)
	return
}

func doAssertions(allKeys *util.StringSet, resp *http.Response, c config.Case, hasRedirect bool, latency int64) (stats Stats) {

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

	type keyAssert struct {
		key string
		ctx Ctx
	}

	// NOTE: the order
	keyAsserts := []keyAssert{
		// statuscode
		{
			key: "assert.statuscode",
			ctx: Ctx{
				f:        assert.Equal,
				element1: resp.StatusCode,
				element2: c.Assert.StatusCode,
			},
		},
		{
			key: "assert.statuscode_lt",
			ctx: Ctx{
				f:        assert.Less,
				element1: resp.StatusCode,
				element2: c.Assert.StatusCodeLt,
			},
		},
		{
			key: "assert.statuscode_lte",
			ctx: Ctx{
				f:        assert.LessOrEqual,
				element1: resp.StatusCode,
				element2: c.Assert.StatusCodeLte,
			},
		},
		{
			key: "assert.statuscode_gt",
			ctx: Ctx{
				f:        assert.Greater,
				element1: resp.StatusCode,
				element2: c.Assert.StatusCodeGt,
			},
		},
		{
			key: "assert.statuscode_gte",
			ctx: Ctx{
				f:        assert.GreaterOrEqual,
				element1: resp.StatusCode,
				element2: c.Assert.StatusCodeGte,
			},
		},
		{
			key: "assert.statuscode_in",
			ctx: Ctx{
				f:        assert.In,
				element1: resp.StatusCode,
				element2: c.Assert.StatusCodeIn,
			},
		},
		{
			key: "assert.statuscode_not_in",
			ctx: Ctx{
				f:        assert.NotIn,
				element1: resp.StatusCode,
				element2: c.Assert.StatusCodeNotIn,
			},
		},
		// status
		{
			key: "assert.status",
			ctx: Ctx{
				f:        assert.Equal,
				element1: strings.ToLower(http.StatusText(resp.StatusCode)),
				element2: strings.ToLower(c.Assert.Status),
			},
		},
		{
			key: "assert.status_in",
			ctx: Ctx{
				f:        assert.In,
				element1: strings.ToLower(http.StatusText(resp.StatusCode)),
				element2: util.ToLower(c.Assert.StatusIn),
			},
		},
		{
			key: "assert.status_not_in",
			ctx: Ctx{
				f:        assert.NotIn,
				element1: strings.ToLower(http.StatusText(resp.StatusCode)),
				element2: util.ToLower(c.Assert.StatusNotIn),
			},
		},
		{
			key: "assert.contenttype",
			ctx: Ctx{
				f:        assert.Equal,
				element1: strings.ToLower(contentType),
				element2: strings.ToLower(c.Assert.ContentType),
			},
		},
		{
			key: "assert.contenttype_in",
			ctx: Ctx{
				f:        assert.In,
				element1: strings.ToLower(contentType),
				element2: util.ToLower(c.Assert.ContentTypeIn),
			},
		},
		{
			key: "assert.contenttype_not_in",
			ctx: Ctx{
				f:        assert.NotIn,
				element1: strings.ToLower(contentType),
				element2: util.ToLower(c.Assert.ContentTypeNotIn),
			},
		},
		// contentlength
		{
			key: "assert.contentlength",
			ctx: Ctx{
				f:        assert.Equal,
				element1: resp.ContentLength,
				element2: c.Assert.ContentLength,
			},
		},
		{
			key: "assert.contentlength_lt",
			ctx: Ctx{
				f:        assert.Less,
				element1: resp.ContentLength,
				element2: c.Assert.ContentLengthLt,
			},
		},
		{
			key: "assert.contentlength_lte",
			ctx: Ctx{
				f:        assert.LessOrEqual,
				element1: resp.ContentLength,
				element2: c.Assert.ContentLengthLte,
			},
		},
		{
			key: "assert.contentlength_gt",
			ctx: Ctx{
				f:        assert.Greater,
				element1: resp.ContentLength,
				element2: c.Assert.ContentLengthGt,
			},
		},
		{
			key: "assert.contentlength_gte",
			ctx: Ctx{
				f:        assert.GreaterOrEqual,
				element1: resp.ContentLength,
				element2: c.Assert.ContentLengthGte,
			},
		},
		// latency
		{
			key: "assert.latency_lt",
			ctx: Ctx{
				f:        assert.Less,
				element1: latency,
				element2: c.Assert.LatencyLt,
			},
		},
		{
			key: "assert.latency_lte",
			ctx: Ctx{
				f:        assert.LessOrEqual,
				element1: latency,
				element2: c.Assert.LatencyLte,
			},
		},
		{
			key: "assert.latency_gt",
			ctx: Ctx{
				f:        assert.Greater,
				element1: latency,
				element2: c.Assert.LatencyGt,
			},
		},
		{
			key: "assert.latency_gte",
			ctx: Ctx{
				f:        assert.GreaterOrEqual,
				element1: latency,
				element2: c.Assert.LatencyGte,
			},
		},
		// body
		{
			key: "assert.body",
			ctx: Ctx{
				f:        assert.Equal,
				element1: bodyStr,
				element2: c.Assert.Body,
			},
		},
		{
			key: "assert.body_contains",
			ctx: Ctx{
				f:        assert.Contains,
				element1: bodyStr,
				element2: c.Assert.BodyContains,
			},
		},
		{
			key: "assert.body_not_contains",
			ctx: Ctx{
				f:        assert.NotContains,
				element1: bodyStr,
				element2: c.Assert.BodyNotContains,
			},
		},
		{
			key: "assert.body_startswith",
			ctx: Ctx{
				f:        assert.StartsWith,
				element1: bodyStr,
				element2: c.Assert.BodyStartsWith,
			},
		},
		{
			key: "assert.body_endswith",
			ctx: Ctx{
				f:        assert.EndsWith,
				element1: bodyStr,
				element2: c.Assert.BodyEndsWith,
			},
		},
		{
			key: "assert.body_not_startswith",
			ctx: Ctx{
				f:        assert.NotStartsWith,
				element1: bodyStr,
				element2: c.Assert.BodyNotStartsWith,
			},
		},
		{
			key: "assert.body_not_endswith",
			ctx: Ctx{
				f:        assert.NotEndsWith,
				element1: bodyStr,
				element2: c.Assert.BodyNotEndsWith,
			},
		},
		{
			key: "assert.hasredirect",
			ctx: Ctx{
				f:        assert.Equal,
				element1: hasRedirect,
				element2: c.Assert.HasRedirect,
			},
		},
		{
			key: "assert.proto",
			ctx: Ctx{
				f:        assert.Equal,
				element1: resp.Proto,
				element2: c.Assert.Proto,
			},
		},
		{
			key: "assert.protomajor",
			ctx: Ctx{
				f:        assert.Equal,
				element1: resp.ProtoMajor,
				element2: c.Assert.ProtoMajor,
			},
		},
		{
			key: "assert.protominor",
			ctx: Ctx{
				f:        assert.Equal,
				element1: resp.ProtoMinor,
				element2: c.Assert.ProtoMinor,
			},
		},
	}

	for _, ka := range keyAsserts {
		if allKeys.Has(ka.key) {
			log.Infof("%s: ", ka.key)
			ok, message := ka.ctx.f(ka.ctx.element1, ka.ctx.element2)
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
