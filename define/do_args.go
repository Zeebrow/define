package define

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type LambdaConfig struct {
	DictApiKey string
	Word       string
}

type CliArgs struct {
	DictApiKey string
	// ThesApiKey  string
	Word    string
	Version bool
}

func (cliargs *CliArgs) GetCliArgs() {

	const (
		cfgFilepathHelp = "override location to config file (default is .MW-api-keys in home directory)"
		versionHelp     = "print the version and exit."
	)

	flag.StringVar(&cliargs.DictApiKey, "dict-api-key", "", "Overwrite any configuration of MW_DICTIONARY_API_KEY")
	// flag.StringVar(&cliargs.ThesApiKey, "thes-api-key", "", "Overwrite any configuration of MW_THESAURUS_API_KEY")
	flag.BoolVar(&cliargs.Version, "version", false, versionHelp)
	flag.Parse()
	if cliargs.Version {
		fmt.Println(ProgInfo.GetInfo())
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		fmt.Println("Specify a word to look up!")
		flag.Usage()
		os.Exit(1)
	} else if flag.NArg() > 1 {
		cliargs.Word = strings.Join(flag.Args(), " ")
	} else {
		cliargs.Word = flag.Arg(0)
	}
}
