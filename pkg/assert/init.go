package assert

import "github.com/fatih/color"

var Fail = color.New(color.FgRed).PrintfFunc()

func OK() {
	color.New(color.FgGreen).PrintfFunc()("OK\n")
}

type AssertFunc func(expected interface{}, actual interface{}) bool
