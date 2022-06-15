package main

import (
	"os"
)

// Globals
// do_args.go
var Debug bool

// config.go
var GlobalConfig GlobalConfiguration

func main() {
	var args = DoArgs()
	GlobalConfig.SetConfig(args)
	// if args.version {
	// 	fmt.Println("version " + Version)
	// 	os.Exit(0)
	// }

	GlobalConfig.getVersion()

	if args.more {
		GetMW(args.word, args.nsfw, args.cfgFilepath)
	} else {
		DevPrintMeanings(args.word, args.nsfw)
	}

	os.Exit(0)
}
