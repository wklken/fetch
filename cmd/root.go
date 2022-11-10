package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "httptest",
	Short: "httptest is a command line http test tool Maintain the api test cases via git and pure text",
	Long:  `A command lin http test tool. Complete documentation is available at https://github.com/wklken/httptest`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		// fmt.Println("hello there")
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
