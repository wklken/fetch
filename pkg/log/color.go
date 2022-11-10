package log

import (
	"github.com/fatih/color"
)

var quiet bool

func BeQuiet(q bool) {
	quiet = q
}

func Tip(format string, a ...interface{}) {
	if !quiet {
		color.New(color.FgYellow).PrintfFunc()(format+"\n", a...)
	}
}

func Warning(format string, a ...interface{}) {
	if !quiet {
		color.New(color.FgYellow).PrintfFunc()(format+"\n", a...)
	}
}

func Infof(format string, a ...interface{}) {
	if !quiet {
		color.New(color.FgWhite).PrintfFunc()(format, a...)
	}
}

func Info(format string, a ...interface{}) {
	if !quiet {
		color.New(color.FgWhite).PrintfFunc()(format+"\n", a...)
	}
}

func Pass() {
	if !quiet {
		color.New(color.FgGreen).PrintfFunc()("Pass\n")
	}
}

func Error(format string, a ...interface{}) {
	if !quiet {
		color.New(color.FgHiRed).PrintfFunc()(format+"\n", a...)
	}
}

func Fail(format string, a ...interface{}) {
	if !quiet {
		color.New(color.FgRed).PrintfFunc()("FAIL: "+format+"\n", a...)
	}
}
