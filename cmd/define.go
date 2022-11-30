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
		mwDictionary := define.NewDictionary(define.GlobalConfig.MWDictionaryApiKey)
		definitions := mwDictionary.Define(CliArgs.Word).GetSimpleHomonymJSON()
		definitions.Print()
	}

	os.Exit(0)
}
