package main

import (
	"os"
)

var CliArgs CLIArgs

func main() {
	// var args = DoArgs()
	CliArgs.DoArgs()
	GlobalConfig.SetConfig(CliArgs)

	if CliArgs.dev {
		DevPrintMeanings(CliArgs.word)
	} else {
		GetMW(CliArgs.word, CliArgs.stdin)
	}

	os.Exit(0)
}
