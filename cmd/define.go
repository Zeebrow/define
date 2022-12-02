package main

import (
	"os"

	"github.com/Zeebrow/define/define"
)

func main() {
	var CliArgs define.CliArgs
	CliArgs.GetCliArgs()
	err := define.GlobalConfig.LoadFromFile()
	if err != nil {
		panic(err)
	}

	mwDictionary := define.NewDictionary(define.GlobalConfig.MWDictionaryApiKey)
	//@@@ Where am I supposed to ask the user if they are sure what they want to define should be in the dictionary?
	definitions, err := mwDictionary.Lookup(CliArgs.Word)
	if err != nil {
		definitions.PrintSuggestions()
		os.Exit(1)
	}
	entries := definitions.GetSimpleHomonymJSON()
	entries.Print()

	os.Exit(0)
}
