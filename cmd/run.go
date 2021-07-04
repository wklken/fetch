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
		path := args[0]

		run(path)
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
	fmt.Println("allKeys", allKeys)

	fmt.Printf("the case and data: %s, %+v\n", path, c)

	var resp *http.Response
	var err1 error
	if c.Request.Method == "get" {
		resp, err1 = http.Get(c.Request.URL)
	}
	//else if c.Request.Method == "post" {
	//	resp, err1 = http.Post(c.Request.URL)
	//}

	assert.NoError(err1)

	if c.Assert.StatusCode != 0 {
		assert.Equal(c.Assert.StatusCode, resp.StatusCode)
	}

	if allKeys.Has("assert.statuscode_in") {
		assert.In(resp.StatusCode, c.Assert.StatusCodeIn)
	}

	if c.Assert.Status != "" {
		assert.Equal(strings.ToLower(c.Assert.Status), strings.ToLower(http.StatusText(resp.StatusCode)))
	}

	// TODO: not set ? or is actually == 0, 不知道是==0, 还是没有配置unmarshall
	if allKeys.Has("assert.contentlength") {
		assert.Equal(c.Assert.ContentLength, resp.ContentLength)
	}
	//https://golang.org/src/net/http/status.go

	body, err := io.ReadAll(resp.Body)
	assert.NoError(err)
	assert.Equal(c.Assert.Body, string(body))

	if allKeys.Has("assert.body_contains") {
		assert.Contains(string(body), c.Assert.BodyContains)
	}
	if allKeys.Has("assert.body_not_contains") {
		assert.NotContains(string(body), c.Assert.BodyNotContains)
	}

	if allKeys.Has("assert.body_startswith") {
		assert.StartsWith(string(body), c.Assert.BodyStartsWith)
	}
	if allKeys.Has("assert.body_endswith") {
		assert.EndsWith(string(body), c.Assert.BodyEndsWith)
	}
}
