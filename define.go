package main

import (
	"os"
)

func main() {
	var args = DoArgs()
	GlobalConfig.SetConfig(args)

	if args.dev {
		DevPrintMeanings(args.word)
	} else {
		GetMW(args.word, args.stdin)
	}

	os.Exit(0)
}
