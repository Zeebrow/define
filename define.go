package main

import "os"

// Globals
// do_args.go
var Debug bool

// config.go
var GlobalConfig GlobalConfiguration

func main() {
	GlobalConfig.SetConfig()
	var args = DoArgs()

	if args.more {
		GetMW(args.word, args.nsfw)
	} else {
		DevPrintMeanings(args.word, args.nsfw)
	}

	os.Exit(0)
}
