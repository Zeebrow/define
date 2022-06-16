package main

import (
	"fmt"
	"os"
)

func main() {
	var args = DoArgs()
	GlobalConfig.SetConfig(args)
	fmt.Printf("mw dictionary key: %s\nmw thesaurus key: %s\n", GlobalConfig.MWDictionaryApiKey, GlobalConfig.MWThesaurusApiKey)
	GlobalConfig.printDebug()

	if args.more {
		GetMW(args.word, args.cfgFilepath)
	} else {
		DevPrintMeanings(args.word)
	}

	os.Exit(0)
}
