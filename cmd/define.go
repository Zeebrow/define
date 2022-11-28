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
		define.GetMW(CliArgs.Word, define.GlobalConfig.MWDictionaryApiKey)
		// mw := define.NewApi(define.GlobalConfig.MWDictionaryApiKey)
		// mw.Headword = CliArgs.Word
		// mw.Define()
	}

	os.Exit(0)
}
