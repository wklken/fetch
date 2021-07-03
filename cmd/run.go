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
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"httptest/pkg/config"
	"net/http"
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

	fmt.Println("the case and data", path, c)

	if c.Request.Method == "get" {
		resp, err := http.Get(c.Request.URL)

		assertNoError(err)

		assertEqual(c.Assert.StatusCode, resp.StatusCode)

		//t := new(testing.T)
		//assert := assert.New(t)
		//fmt.Println("c.Assert", c.Assert.Status)
		//assert.NoError(err, "should not error")
		//assert.Equal(c.Assert.Status, resp.StatusCode, "status should be")
	}
}

func assertNoError(err error) {
	if err != nil {
		fmt.Println("FAIL: got an error")
	}
}

func assertEqual(expected interface{}, actual interface{}) {
	equal := assert.ObjectsAreEqual(expected, actual)
	if !equal {
		fmt.Printf("FAIL: not equal, expected=%d, actual=%d\n", expected,actual)
		fmt.Println()
	}
}