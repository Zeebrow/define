package main

import (
	"flag"
	"fmt"
	"os"
)

type CLIArgs struct {
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
	version     bool
}

func (c *CLIArgs) printDebug() {
	fmt.Printf("-------------cli args-------------\n")
	fmt.Printf("m: %v\n", c.more)
	fmt.Printf("debug: %v\n", c.debug)
	fmt.Printf("word to define: %v\n", c.word)
	fmt.Printf("get config from file: %v\n", c.cfgFilepath)
	fmt.Printf("NArg: %d\n", flag.NArg())
	fmt.Printf("NFlag: %d\n", flag.NFlag())
	fmt.Printf("version: %v\n", c.version)
}

func DoArgs() (cliargs CLIArgs) {

	const (
		more_help       = "Print more detailed definitions"
		nsfw_help       = "Print potentially offensive definitions"
		cfgFilepathHelp = "override location to config file (default is .MW-api-keys in home directory)"
		versionHelp     = "print the version and exit."
	)

	flag.BoolVar(&cliargs.more, "m", false, more_help)
	flag.BoolVar(&cliargs.synonyms, "s", false, "Get Synonym")
	flag.BoolVar(&cliargs.antonyms, "a", false, "Get Antonym")
	flag.BoolVar(&cliargs.debug, "debug", false, "Print debug output")
	flag.StringVar(&cliargs.dictApiKey, "dict-api-key", "", "Overwrite any configuration of MW_DICTIONARY_API_KEY")
	flag.StringVar(&cliargs.thesApiKey, "thes-api-key", "", "Overwrite any configuration of MW_THESAURUS_API_KEY")
	flag.StringVar(&cliargs.cfgFilepath, "f", GetDefaultFilepath(), cfgFilepathHelp)
	flag.BoolVar(&cliargs.version, "version", false, versionHelp)
	flag.Parse()
	if cliargs.debug {
		GlobalConfig.Debug = true
		cliargs.printDebug()
		GlobalConfig.printDebug()
	}
	if cliargs.version {
		fmt.Println(ProgInfo.GetInfo())
		os.Exit(0)
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
