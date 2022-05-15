package main

import (
	"flag"
	"fmt"
	"os"
)

type CLIArgs struct {
	// nsfw = include MW offensive: true
	nsfw bool
	// more = use Merriam-Webster
	more        bool
	store       string
	debug       bool
	synonyms    bool
	antonyms    bool
	dictApiKey  string
	thesApiKey  string
	cfgFilepath string
	word        string
}

// var Usage = func() {
// 	// const usage_template = ``
// 	ul := chalk.Underline.TextStyle("word")
// 	fmt.Printf("usage: define %v\n", ul)
// }

func printDebug(c *CLIArgs) {

	fmt.Printf("m: %v\n", c.more)
	fmt.Printf("NSFW?: %v\n", c.nsfw)
	fmt.Printf("debug: %v\n", c.debug)
	fmt.Printf("store: %v\n", c.store)
	fmt.Printf("word to define: %v\n", c.word)
	fmt.Printf("get config from file: %v\n", c.cfgFilepath)
	fmt.Printf("NArg: %d\n", flag.NArg())
	fmt.Printf("NFlag: %d\n", flag.NFlag())
}

func DoArgs() (cliargs CLIArgs) {

	const (
		more_help       = "Print more detailed definitions"
		nsfw_help       = "Print potentially offensive definitions"
		cfgFilepathHelp = "override location to config file (default is .MW-api-keys in home directory)"
	)

	flag.BoolVar(&cliargs.more, "m", false, more_help)
	flag.BoolVar(&cliargs.nsfw, "x", false, nsfw_help)
	flag.BoolVar(&cliargs.synonyms, "s", false, "Save output to file")
	flag.BoolVar(&cliargs.antonyms, "a", false, "Save output to file")
	flag.BoolVar(&cliargs.debug, "debug", false, "Print debug output")
	flag.StringVar(&cliargs.store, "w", "", "Save output to file")
	flag.StringVar(&cliargs.dictApiKey, "dict-api-key", "", "Overwrite any configuration of MW_DICTIONARY_API_KEY")
	flag.StringVar(&cliargs.thesApiKey, "thes-api-key", "", "Overwrite any configuration of MW_THESAURUS_API_KEY")
	flag.StringVar(&cliargs.cfgFilepath, "f", "", cfgFilepathHelp)

	flag.Parse()
	if cliargs.debug {
		Debug = true
		printDebug(&cliargs)
	}
	if flag.NArg() < 1 {
		fmt.Println("Specify a word to look up!")
		flag.Usage()
		os.Exit(1)
	} else if flag.NArg() > 1 {
		fmt.Printf("Woah there, only one word at a time. (You entered %v)\n", flag.Args())
		flag.Usage()
		os.Exit(1)
	} else {
		cliargs.word = flag.Arg(0)
	}
	return

}
