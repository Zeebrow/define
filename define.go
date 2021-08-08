package main

import "os"

// Global so other functions can use it
// This is set in do_args.go, maybe it should be there instead?
var Debug bool

func main() {
	var args = DoArgs()

	if args.more {
		GetMW(args.word, args.nsfw)
	} else {
		DevPrintMeanings(args.word, args.nsfw)
	}

	os.Exit(0)
}
