package define

import (
	"flag"
	"fmt"
	"os"
)

type CliArgs struct {
	// more = use Merriam-Webster
	Stdin       bool
	Dev         bool
	Store       string
	Debug       bool
	DictApiKey  string
	ThesApiKey  string
	CfgFilepath string
	Word        string
	Version     bool
}

func (c *CliArgs) printDebug() {
	fmt.Printf("-------------cli args-------------\n")
	fmt.Printf("q: %v\n", c.Dev)
	fmt.Printf("stdin: %v\n", c.Stdin)
	fmt.Printf("debug: %v\n", c.Debug)
	fmt.Printf("word to define: '%v\n", flag.Arg(0))
	fmt.Printf("get config from file: %v\n", c.CfgFilepath)
	fmt.Printf("version: %v\n", c.Version)
}

func (cliargs *CliArgs) DoArgs() {

	const (
		stdinHelp       = "Read a json from stdin instead of calling api"
		devHelp         = "Use an alternative dictionary with more reliable api responses"
		cfgFilepathHelp = "override location to config file (default is .MW-api-keys in home directory)"
		versionHelp     = "print the version and exit."
	)

	flag.BoolVar(&cliargs.Dev, "q", false, devHelp)
	flag.BoolVar(&cliargs.Stdin, "stdin", false, stdinHelp)
	flag.BoolVar(&cliargs.Debug, "debug", false, "Print debug output")
	flag.StringVar(&cliargs.DictApiKey, "dict-api-key", "", "Overwrite any configuration of MW_DICTIONARY_API_KEY")
	flag.StringVar(&cliargs.ThesApiKey, "thes-api-key", "", "Overwrite any configuration of MW_THESAURUS_API_KEY")
	flag.StringVar(&cliargs.CfgFilepath, "f", GetDefaultFilepath(), cfgFilepathHelp)
	flag.BoolVar(&cliargs.Version, "version", false, versionHelp)
	flag.Parse()
	if cliargs.Debug {
		GlobalConfig.Debug = true
		cliargs.printDebug()
		GlobalConfig.printDebug()
	}
	if cliargs.Version {
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
		cliargs.Word = flag.Arg(0)
	}
	return
}
