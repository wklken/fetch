package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch is a command line http test tool maintain the api test cases via git and pure text",
	Long:  `A command lin http test tool. Complete documentation is available at https://github.com/wklken/fetch`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
