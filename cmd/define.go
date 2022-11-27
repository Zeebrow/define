package main

import (
	"os"

	"github.com/Zeebrow/define/define"
)

var CliArgs define.CliArgs

func main() {
	// var args = DoArgs()
	CliArgs.DoArgs()
	define.GlobalConfig.SetConfig(CliArgs)

	if CliArgs.Dev {
		define.DevPrintMeanings(CliArgs.Word)
	} else {
		define.GetMW(CliArgs.Word, CliArgs.Stdin)
	}

	os.Exit(0)
}
