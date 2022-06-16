package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

/*
App globals
*/
var GlobalConfig GlobalConfiguration
var MWDictionaryApiKey string
var MWThesaurusApiKey string
var ConfigFileName = ".define.conf" // idk what to name the file lol.
var ProgramName = "define"
var Version = "dev"
var BuildDate = ""
var CommitHash = ""
var ProgInfo = ProgramInfo{
	ProgramName: ProgramName,
	Version:     Version,
	BuildDate:   BuildDate,
	CommitHash:  CommitHash,
}

type ProgramInfo struct {
	ProgramName string
	Version     string
	BuildDate   string
	CommitHash  string
}

func (p *ProgramInfo) GetInfo() string {
	return fmt.Sprintf("%s version %s (%s) built %s", p.ProgramName, p.Version, p.CommitHash, p.BuildDate)
}

type GlobalConfiguration struct {
	MWDictionaryApiKey string
	MWThesaurusApiKey  string
	ConfigFilepath     string
	Debug              bool
}

type ConfigFile struct {
	Dictionary string `json:"dictionary"`
	Thesaurus  string `json:"thesaurus,omitempty"`
}

func (gc *GlobalConfiguration) printDebug() {
	fmt.Printf("-------------global config-------------\n")
	fmt.Printf("MWDictionaryApiKey : %s\n", gc.MWDictionaryApiKey)
	fmt.Printf("MWThesaurusApiKey  : %s\n", gc.MWThesaurusApiKey)
	fmt.Printf("ConfigFilepath     : %s\n", gc.ConfigFilepath)
	fmt.Printf("Debug              : %v\n", gc.Debug)
}

func GetDefaultFilepath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Could not get user home directory... searching current directory for config")
		return ConfigFileName
	} else {
		return fmt.Sprintf("%s", filepath.Join(dirname, ConfigFileName))
	}
}

func (cfg *GlobalConfiguration) SetConfig(cliArgs CLIArgs) {
	// 5th hard-coded variables
	// 4th injected build variables
	cfg.MWDictionaryApiKey = MWDictionaryApiKey
	cfg.MWThesaurusApiKey = MWThesaurusApiKey

	// 3rd ???
	cfg.ConfigFilepath = GetDefaultFilepath()
	func(configfile string) {
		var ks ConfigFile
		r, err := ioutil.ReadFile(configfile)
		if err != nil {
			fmt.Printf("Could not open API keys file: %s. Error: %v", configfile, err)
			os.Exit(1)
		}

		err = json.Unmarshal(r, &ks)
		if err != nil {
			fmt.Printf("Could not get dictionary API keys from file: %s. Error: %v", configfile, err)
			os.Exit(1)
		}
		fmt.Printf("=========>%s\n", ks)
		cfg.MWDictionaryApiKey = ks.Dictionary
		cfg.MWThesaurusApiKey = ks.Thesaurus
	}(cfg.ConfigFilepath)
	// 2nd is Environment
	if os.Getenv("MW_CONFIG_FILEPATH") != "" {
		cfg.ConfigFilepath = os.Getenv("MW_CONFIG_FILEPATH")
	}
	if os.Getenv("MW_DICTIONARY_API_KEY") != "" {
		cfg.MWDictionaryApiKey = os.Getenv("MW_DICTIONARY_API_KEY")
	}
	if os.Getenv("MW_THESAURUS_API_KEY") != "" {
		cfg.MWThesaurusApiKey = os.Getenv("MW_THESAURUS_API_KEY")
	}
	// 1st CLI args
	if cliArgs.dictApiKey != "" {
		cfg.MWDictionaryApiKey = cliArgs.dictApiKey
	}
	if cliArgs.thesApiKey != "" {
		cfg.MWThesaurusApiKey = cliArgs.thesApiKey
	}
}
