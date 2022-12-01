package main

import (
	"os"

	"github.com/Zeebrow/define/define"
)

func main() {
	// var args = DoArgs()
	var CliArgs define.CliArgs
	CliArgs.GetCliArgs()
	err := define.GlobalConfig.LoadFromFile()
	if err != nil {
		panic(err)
	}

	mwDictionary := define.NewDictionary(define.GlobalConfig.MWDictionaryApiKey)
	definitions := mwDictionary.Define(CliArgs.Word).GetSimpleHomonymJSON()
	definitions.Print()

	os.Exit(0)
}
