package main

import (
	"fmt"
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
	definitions, err := mwDictionary.Lookup(CliArgs.Word)
	if err != nil {
		fmt.Printf("\n%s\n%s\n\n", err, "Consider trying one of these instead!")
		definitions.PrintSuggestions()
		return
		// os.Exit(1)
	}
	entries := definitions.GetSimpleHomonymJSON()
	entries.Print()
	os.Exit(0)
}
