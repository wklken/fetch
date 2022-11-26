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
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/jmespath/go-jmespath"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/cobra"

	"github.com/wklken/httptest/pkg/assert"
	"github.com/wklken/httptest/pkg/client"
	"github.com/wklken/httptest/pkg/config"
	"github.com/wklken/httptest/pkg/log"
	"github.com/wklken/httptest/pkg/util"
)

const (
	DebugEnvName = "HTTPTEST_DEBUG"
)

var (
	verbose   = false
	quiet     = false
	cfgFile   string
	parallels = 1
)

type RunInParallelArgs struct {
	Path      string
	RunConfig *config.RunConfig
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run cases",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
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

		if len(args) == 0 && len(runConfig.Order) == 0 {
			log.Error("args required")
			os.Exit(1)
			return
		}

		// parse files in order to run
		orderedCases, err := util.GetRunningOrderedFiles(args, runConfig.Order)
		if err != nil {
			log.Error("parse config file `Order` fail, err=%s", err)
			os.Exit(1)
			return
		}

		totalStats := util.Stats{}
		log.BeQuiet(quiet)

		start := time.Now()

		if parallels <= 1 {
			for _, path := range orderedCases {
				s := run(path, &runConfig)

				// 2. collect the result
				totalStats.MergeAssertCount(s)

				// FIXME: log one by one, not at the last

				if runConfig.FailFast && !(s.AllPassed()) {
					log.Info("failFast=True, quit, the execute result: 1")
					os.Exit(1)
				}
			}
		} else {
			var wg sync.WaitGroup
			sc := util.StatsCollection{}
			p, _ := ants.NewPoolWithFunc(parallels, func(i interface{}) {
				defer wg.Done()
				args := i.(RunInParallelArgs)

				s := run(args.Path, args.RunConfig)
				if runConfig.FailFast && !(s.AllPassed()) {
					log.Info("failFast=True, quit, the execute result: 1")
					// FIXME: should stop all, not only the goroutine
					os.Exit(1)
				}

				sc.Add(s)
				// FIXME: log one by one, not at the last
			})
			defer p.Release()

			for _, path := range orderedCases {
				// TODO: -p 10 to run in parallel with 10 goroutines
				wg.Add(1)

				args := RunInParallelArgs{
					Path:      path,
					RunConfig: &runConfig,
				}
				p.Invoke(args)
			}
			wg.Wait()
			totalStats = sc.GetStats()
		}

		totalStats.PrintMessages()

		latency := time.Since(start).Milliseconds()
		totalStats.Report(latency)
		if totalStats.AllPassed() {
			log.Info("the execute result: 0")
		} else {
			log.Info("the execute result: 1")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// -v verbose
	runCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode")
	// -q quiet
	runCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "be quiet")
	// -p parallel
	runCmd.PersistentFlags().IntVarP(&parallels, "parallel", "p", 1, "run in parallel")

	// -e dev.toml
	runCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file(like dev.toml/prod.toml")
}

// func logRunCaseFail(path string, c *config.Case, format string, a ...interface{}) {
// 	log.Tip("Run Case: %s | %s | [%s %s]", path, c.Title, strings.ToUpper(c.Request.Method), c.Request.URL)
// 	log.Error(format, a...)
// }

func run(path string, runConfig *config.RunConfig) (stats util.Stats) {
	// FIXME: should collect the log instead of print it(in parallel will be a problem)
	v, err := config.ReadFromFile(path)
	if err != nil {
		stats.AddTipMessage("Run Case: %s", path)
		stats.AddErrorMessage("read fail: %s", err)
		stats.IncrFailCaseCount()
		return
	}
	// read lines, for display the failed asset line number
	fileLines, err := config.ReadLines(path)
	if err != nil {
		stats.AddTipMessage("Run Case: %s", path)
		stats.AddErrorMessage("read fail: %s", err)
		stats.IncrFailCaseCount()
		return
	}

	var c config.Case
	err = v.Unmarshal(&c)
	if err != nil {
		stats.AddTipMessage("Run Case: %s", path)
		stats.AddErrorMessage("parse fail: %s", err)
		return
	}
	// set the content
	c.FileLines = fileLines

	allKeys := util.NewStringSetWithValues(v.AllKeys())
	// fmt.Println("allKeys", allKeys)
	// fmt.Printf("the case and data: %s, %+v", path, c)

	// do render
	if runConfig.Render && len(runConfig.Env) > 0 {
		finalEnv := runConfig.Env
		// the priority of env in case is higher than env in config
		if len(c.Env) > 0 {
			for k, v := range c.Env {
				finalEnv[k] = v
			}
		}
		c.Render(runConfig.Env)
	}

	debug := (verbose || strings.ToLower(os.Getenv(DebugEnvName)) == "true" || runConfig.Debug) && !quiet
	timeout := runConfig.Timeout
	if c.Config.Timeout > 0 {
		timeout = c.Config.Timeout
	}

	// support repeat, if got repeat, on case will be repeat N times, as N cases
	repeat := 1
	if c.Config.Repeat > 0 {
		repeat = c.Config.Repeat
	}

	for i := 0; i < repeat; i++ {
		var (
			resp        *http.Response
			hasRedirect bool
			latency     int64
			err2        error
			count       int
		)
		for {
			resp, hasRedirect, latency, err2 = client.Send(
				filepath.Dir(path),
				c.Request.Method,
				c.Request.URL,
				allKeys.Has("request.body"),
				c.Request.Body,
				c.Request.Header,
				c.Request.Cookie,
				c.Request.BasicAuth,
				c.Hook,
				timeout,
				debug,
			)

			if c.Config.Retry.Enable && count < c.Config.Retry.Count &&
				(err2 != nil || util.ItemInIntArray(resp.StatusCode, c.Config.Retry.StatusCodes)) {
				time.Sleep(time.Duration(c.Config.Retry.Interval) * time.Millisecond)
				count++
				continue
			} else {
				break
			}

		}
		title := c.Title
		if repeat > 1 {
			title = fmt.Sprintf("%s (%d/%d)", c.Title, i+1, repeat)
		}

		if err2 != nil {
			stats.AddTipMessage(
				"Run Case: %s | %s | [%s %s] | %dms",
				path,
				title,
				strings.ToUpper(c.Request.Method),
				c.Request.URL,
				latency,
			)

			if !allKeys.Has("assert.error_contains") {
				stats.AddErrorMessage("Send HTTP Request fail: %s", err2)
				stats.IncrFailCaseCount()
			} else {
				// do assert with error_contains
				s1 := doErrorAssertions(c, err2)
				stats.MergeAssertCount(s1)
			}

			if repeat > 1 && i < repeat-1 {
				continue
			}
			return
		}

		stats.AddTipMessage(
			"Run Case: %s | %s | [%s %s] | %dms",
			path,
			title,
			strings.ToUpper(c.Request.Method),
			c.Request.URL,
			latency,
		)

		s := doAssertions(allKeys, resp, c, hasRedirect, latency)
		stats.MergeAssertCount(s)
	}

	return
}

func doAssertions(
	allKeys *util.StringSet,
	resp *http.Response,
	c config.Case,
	hasRedirect bool,
	latency int64,
) (stats util.Stats) {
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
			stats.AddInfofMessage("%s: ", ka.key)
			ok, message := ka.ctx.f(ka.ctx.element1, ka.ctx.element2)
			if ok {
				stats.AddPassMessage()
				stats.IncrOkAssertCount()
			} else {
				// the ka.key is like assert.latency_lt
				lineNumber := c.GuessAssertLineNumber(ka.key)
				if lineNumber > 0 {
					message = fmt.Sprintf("line:%d | %s", lineNumber, message)
				}
				stats.AddFailMessage(message)
				stats.IncrFailAssertCount()
			}
		}
	}

	if len(c.Assert.Header) > 0 {
		for key, value := range c.Assert.Header {
			stats.AddInfoMessage("assert.header.%s: ", key)
			ok, message := assert.Equal(resp.Header.Get(key), value)
			if ok {
				stats.AddPassMessage()
				stats.IncrOkAssertCount()
			} else {
				// the ka.key is like assert.latency_lt
				lineNumber := c.GuessAssertLineNumber(key)
				if lineNumber > 0 {
					message = fmt.Sprintf("line:%d | %s", lineNumber, message)
				}
				stats.AddFailMessage(message)
				stats.IncrFailAssertCount()
			}

		}
	}

	var jsonData interface{}
	if contentType == binding.MIMEJSON || contentType == binding.MIMEMSGPACK || contentType == binding.MIMEMSGPACK2 {
		var f binding.BindingBody
		switch contentType {
		case binding.MIMEJSON:
			f = binding.JSON
		case binding.MIMEMSGPACK, binding.MIMEMSGPACK2:
			f = binding.MsgPack
		default:
			f = nil
		}

		if f != nil {
			// err = binding.JSON.BindBody(body, &jsonData)
			err = f.BindBody(body, &jsonData)
			if err != nil {
				stats.AddFailMessage("binding.json fail: %s", err)
				stats.IncrFailAssertCountByN(int64(len(c.Assert.Json)))
				return
			}

			if allKeys.Has("assert.json") && len(c.Assert.Json) > 0 {
				s1 := doJsonAssertions(jsonData, c.Assert.Json)
				stats.MergeAssertCount(s1)
			}
		}
	}

	return
}

func doJsonAssertions(jsonData interface{}, jsons []config.AssertJson) (stats util.Stats) {
	for _, dj := range jsons {
		path := dj.Path
		expectedValue := dj.Value
		stats.AddInfofMessage("assert.json.%s: ", path)

		if jsonData == nil {
			ok, message := assert.Equal(nil, expectedValue)
			if ok {
				stats.IncrOkAssertCount()
			} else {
				stats.AddFailMessage(message)
				stats.IncrFailAssertCount()
			}
			continue
		}

		actualValue, err := jmespath.Search(path, jsonData)
		if err != nil {
			// log.Fail("search json data fail, err=%s, path=%s, expected=%s", err, path, expectedValue)
			stats.AddFailMessage("search json data fail, err=%s, path=%s, expected=%s", err, path, expectedValue)
		} else {
			// missing
			if actualValue == nil {
				_, message := assert.Equal(nil, expectedValue)
				stats.AddFailMessage(message)
				stats.IncrFailAssertCount()
				continue
			}

			// fmt.Printf("%T, %T", actualValue, expectedValue)
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
				stats.AddPassMessage()
				stats.IncrOkAssertCount()
			} else {
				stats.AddFailMessage(message)
				stats.IncrFailAssertCount()
			}
		}
	}

	return
}

func doErrorAssertions(c config.Case, err error) (stats util.Stats) {
	stats.AddInfofMessage("assert.error_contains: ")
	ok, message := assert.Contains(err.Error(), c.Assert.ErrorContains)
	if ok {
		stats.AddPassMessage()
		stats.IncrOkAssertCount()
	} else {
		// the ka.key is like assert.latency_lt
		lineNumber := c.GuessAssertLineNumber("assert.error_contains")
		if lineNumber > 0 {
			message = fmt.Sprintf("line:%d | %s", lineNumber, message)
		}
		stats.AddFailMessage(message)
		stats.IncrFailAssertCount()
	}

	return
}
