package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type GlobalConfiguration struct {
	Version            string
	MWDictionaryApiKey string
	MWThesaurusApiKey  string
	ConfigFilepath     string
}

var MWDictionaryApiKey string
var MWThesaurusApiKey string
var Version = "dev"

func GetDefaultFilepath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Could not get user home directory... searching current directory for config")
		return ".MW-api-keys"
	} else {
		return fmt.Sprintf("%s", filepath.Join(dirname, ".MW-api-keys"))
	}
}

func (cfg *GlobalConfiguration) getVersion() string {
	return cfg.Version
}

func (cfg *GlobalConfiguration) SetConfig(cliArgs CLIArgs) {
	cfg.Version = Version
	// 5th hard-coded variables
	cfg.ConfigFilepath = "/home/zeebrow/.local/etc/define.conf"
	// 4th injected build variables
	cfg.MWDictionaryApiKey = MWDictionaryApiKey
	cfg.MWThesaurusApiKey = MWThesaurusApiKey

	// 3rd config file default

	// @tonotdo

	// 2nd is Environment
	if os.Getenv("MW_DICTIONARY_API_KEY") != "" {
		cfg.MWDictionaryApiKey = os.Getenv("MW_DICTIONARY_API_KEY")
	}
	if os.Getenv("MW_THESAURUS_API_KEY") != "" {
		cfg.MWThesaurusApiKey = os.Getenv("MW_THESAURUS_API_KEY")
	}

	// 1st CLI args take precedence
	// might not be a thing, idk
	if cliArgs.cfgFilepath == "" {
		cfg.ConfigFilepath = GetDefaultFilepath()
	} else {
		//cfg.ConfigFilepath = filepath
		//fmt.Printf("Setting cfg filepath from cli arg: %s\n", cliArgs.cfgFilepath)
		cfg.ConfigFilepath = cliArgs.cfgFilepath
	}

	// @todo, n00b
	// How to iterate thru struct and check for nil string?
	// v := reflect.ValueOf(cfg)
	// for i := 0; i < v.NumField(); i++ {
	// 	if v.Field(i).IsNil() {
	// 		os.Exit(1)
	// 	}
	// }
}
