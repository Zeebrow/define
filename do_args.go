package main

import (
	"flag"
	"fmt"
	"os"
)

type CLIArgs struct {
	// more = use Merriam-Webster
	stdin       bool
	dev         bool
	store       string
	debug       bool
	dictApiKey  string
	thesApiKey  string
	cfgFilepath string
	word        string
	version     bool
}

func (c *CLIArgs) printDebug() {
	fmt.Printf("-------------cli args-------------\n")
	fmt.Printf("q: %v\n", c.dev)
	fmt.Printf("stdin: %v\n", c.stdin)
	fmt.Printf("debug: %v\n", c.debug)
	fmt.Printf("word to define: '%v\n", flag.Arg(0))
	fmt.Printf("get config from file: %v\n", c.cfgFilepath)
	fmt.Printf("version: %v\n", c.version)
}

func (cliargs *CLIArgs) DoArgs() {

	const (
		stdinHelp       = "Read a json from stdin instead of calling api"
		devHelp         = "Use an alternative dictionary with more reliable api responses"
		cfgFilepathHelp = "override location to config file (default is .MW-api-keys in home directory)"
		versionHelp     = "print the version and exit."
	)

	flag.BoolVar(&cliargs.dev, "q", false, devHelp)
	flag.BoolVar(&cliargs.stdin, "stdin", false, stdinHelp)
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
		fmt.Printf("one word at a time, please! (You entered %v)\n", flag.Args())
		flag.Usage()
		os.Exit(1)
	} else {
		cliargs.word = flag.Arg(0)
	}
	return
}
