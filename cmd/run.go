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
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/jmespath/go-jmespath"
	"github.com/panjf2000/ants/v2"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"

	"github.com/wklken/httptest/pkg/assertion"
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
	proxy     string
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

			// log.Info("runConfig: %v", runConfig)
		}

		if len(args) == 0 && len(runConfig.Order) == 0 {
			log.Error("args required, please input the case file path")
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
		// the log
		log.BeQuiet(quiet)

		// the progress bar
		totalCases := int64(len(orderedCases))
		var bar *progressbar.ProgressBar
		if quiet {
			bar = progressbar.DefaultSilent(totalCases)
		} else {
			bar = progressbar.Default(totalCases)
		}

		start := time.Now()
		if parallels <= 1 {
			for _, path := range orderedCases {
				s := run(path, &runConfig)

				bar.Add(1)

				// 2. collect the result
				totalStats.MergeAssertAndCaseCount(s)

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
				defer bar.Add(1)
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

	// --proxy http://myproxy:9999
	runCmd.PersistentFlags().StringVarP(&proxy, "proxy", "", "", "proxy for http request")

	// -e dev.toml
	runCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file(like dev.toml/prod.toml")
}

// func logRunCaseFail(path string, c *config.Case, format string, a ...interface{}) {
// 	log.Tip("Run Case: %s | %s | [%s %s]", path, c.Title, strings.ToUpper(c.Request.Method), c.Request.URL)
// 	log.Error(format, a...)
// }

type CaseContext struct {
	Env map[string]interface{}
}

func run(path string, runConfig *config.RunConfig) (stats util.Stats) {
	// TODO: the path is one single file, but may got more than one case!
	cases, err := config.ReadCasesFromFile(path)
	if err != nil {
		stats.AddTipMessage("Run Case: %s", path)
		stats.AddErrorMessage("read fail: %s", err)
		// FIXME: case it a file? or each section in file?
		stats.IncrFailFileCount()
		return
	}

	// to keep the envs between cases, parse from first case, and use the vars in the next case
	caseContext := CaseContext{
		Env: map[string]interface{}{},
	}

	for _, c := range cases {
		allKeys := util.NewStringSetWithValues(c.AllKeys)
		// do render, priority: caseContext.env > case.Env > runConfig.Env
		finalEnv := map[string]interface{}{}
		if runConfig.Env != nil {
			finalEnv = runConfig.Env
		}
		if len(c.Env) > 0 {
			for k, v := range c.Env {
				finalEnv[k] = v
			}
		}
		if len(caseContext.Env) > 0 {
			for k, v := range caseContext.Env {
				finalEnv[k] = v
			}
		}
		if len(finalEnv) > 0 {
			c.Render(caseContext.Env)
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
				resp          *http.Response
				redirectCount int64
				latency       int64
				debugLogs     []string
				err2          error
				count         int
			)
			for {
				resp, redirectCount, latency, debugLogs, err2 = client.Send(
					filepath.Dir(path),
					c.Request.Method,
					c.Request.URL,
					allKeys.Has("request.body"),
					c.Request.Body,
					c.Request.Header,
					c.Request.Cookie,
					c.Request.BasicAuth,
					c.Request.DisableRedirect,
					c.Hook,
					timeout,
					proxy,
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

			title := c.ID()
			if repeat > 1 {
				title = fmt.Sprintf("%s (%d/%d)", c.ID(), i+1, repeat)
			}
			stats.AddTipMessage(
				"Run Case: %s | [%s %s] | %dms",
				title,
				strings.ToUpper(c.Request.Method),
				c.Request.URL,
				latency,
			)

			if len(debugLogs) > 0 {
				for _, l := range debugLogs {
					stats.AddInfoMessage(l)
				}
			}

			if err2 != nil {
				if !allKeys.Has("assert.error_contains") {
					stats.AddErrorMessage("Send HTTP Request fail: %s", err2)
					stats.IncrFailCaseCount()
				} else {
					// do assert with error_contains
					s1 := assertion.DoErrorAssertions(c, err2)
					stats.MergeAssertCount(s1)
					if stats.GetFailAssertCount() > 0 {
						stats.IncrFailCaseCount()
					} else {
						stats.IncrOkCaseCount()
					}
				}

				if repeat > 1 && i < repeat-1 {
					continue
				}
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				stats.IncrFailCaseCount()
				continue
			}
			// assert.NoError(err)

			s := doAssertions(allKeys, resp, body, c, redirectCount, latency)
			// fmt.Printf("s: %+v\n", stats)
			stats.MergeAssertCount(s)
			if stats.GetFailAssertCount() > 0 {
				stats.IncrFailCaseCount()
			} else {
				stats.IncrOkCaseCount()
			}

			// do parse
			if len(c.Parse) > 0 {
				envs := doParse(c.Parse, body, resp.Header)
				if len(envs) > 0 {
					for k, v := range envs {
						caseContext.Env[k] = v
					}
				}
			}

		}

	}

	return
}

func doParse(parses []config.Parse, body []byte, header http.Header) map[string]interface{} {
	envs := make(map[string]interface{})
	for _, p := range parses {
		// parse header
		if p.Source == "header" {
			envs[p.Key] = header.Get(p.Header)
		}
		// parse body
		if p.Source == "body" {
			if p.Jmespath != "" {
				// FIXME: support msgpack here too
				var jsonData interface{}
				err := binding.JSON.BindBody(body, &jsonData)
				if err != nil {
					// log warning and do nothing
					log.Warning("parse body fail, parse.key=`%s`, err=`%v`", p.Key, err)
				}
				actualValue, err := jmespath.Search(p.Jmespath, jsonData)
				if err != nil {
					log.Warning("parse body and jmespath.Search fail, parse.key=`%s`, err=`%v`", p.Key, err)
					// log warning and do nothing
				}
				envs[p.Key] = actualValue
			}
			// FIXME: support parse html and xml
		}
	}
	return envs
}

func doAssertions(
	allKeys *util.StringSet,
	resp *http.Response,
	body []byte,
	c *config.Case,
	redirectCount int64,
	latency int64,
) (stats util.Stats) {
	// sometimes the content-length is -1(means unknown), we recalculate it
	// NOTE: I DON'T KNOW if the logical is ok or not
	if resp.ContentLength == -1 {
		resp.ContentLength = int64(len(body))
	}

	contentType := client.GetContentType(resp.Header)

	// normal key-value assert
	s := assertion.DoKeysAssertion(allKeys, resp, c, redirectCount, latency, contentType, body)
	stats.MergeAssertCount(s)

	// header assert
	if len(c.Assert.Header) > 0 || len(c.Assert.HeaderExists) > 0 {
		s1 := assertion.DoHeaderAssertions(c, resp.Header)
		stats.MergeAssertCount(s1)
	}

	// xml assert
	if allKeys.Has("assert.xml") && len(c.Assert.XML) > 0 {
		s1 := assertion.DoXMLAssertions(body, c.Assert.XML)
		stats.MergeAssertCount(s1)
	}

	// html assert
	if allKeys.Has("assert.html") && len(c.Assert.HTML) > 0 {
		s1 := assertion.DoHTMLAssertions(body, c.Assert.HTML)
		stats.MergeAssertCount(s1)
	}

	// yaml assert
	if allKeys.Has("assert.yaml") && len(c.Assert.YAML) > 0 {
		s1 := assertion.DoYAMLAssertions(body, c.Assert.YAML)
		stats.MergeAssertCount(s1)
	}

	// toml assert
	if allKeys.Has("assert.toml") && len(c.Assert.TOML) > 0 {
		s1 := assertion.DoTOMLAssertions(body, c.Assert.TOML)
		stats.MergeAssertCount(s1)
	}

	// cookie assert
	if (allKeys.Has("assert.cookie") && len(c.Assert.Cookie) > 0) || len(c.Assert.CookieExists) > 0 {
		s1 := assertion.DoCookieAssertions(c, resp.Cookies())
		stats.MergeAssertCount(s1)
	}

	// json/msgpack assert
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
			if len(body) != 0 {
				err := f.BindBody(body, &jsonData)
				if err != nil {
					stats.AddFailMessage("binding.json fail: %s", err)
					stats.IncrFailAssertCountByN(int64(len(c.Assert.JSON)))
					return
				}
			} else {
				// the http method: head got no response body, but header is application/json
				jsonData = nil
			}

			if allKeys.Has("assert.json") && len(c.Assert.JSON) > 0 {
				s1 := assertion.DoJSONAssertions(jsonData, c.Assert.JSON)
				stats.MergeAssertCount(s1)
			}
		}
	}

	return
}
