package main

import (
	"os"
)

type GlobalConfiguration struct {
	MWDictionaryApiKey string
	MWThesaurusApiKey  string
	ConfigFilepath     string
}

var MWDictionaryApiKey string
var MWThesaurusApiKey string

func (cfg *GlobalConfiguration) SetConfig() {

	// 4th injected build variables
	cfg.MWDictionaryApiKey = MWDictionaryApiKey
	cfg.MWThesaurusApiKey = MWThesaurusApiKey

	// 3rd config file default

	// 2nd is Environment
	cfg.MWDictionaryApiKey = os.Getenv("MW_DICTIONARY_API_KEY")
	cfg.MWThesaurusApiKey = os.Getenv("MW_THESAURUS_API_KEY")
	cfg.ConfigFilepath = "asdf"

	// 1st CLI args take precedence
	// might not be a thing, idk

	// How to iterate thru struct and check for nil string?
	// v := reflect.ValueOf(cfg)
	// for i := 0; i < v.NumField(); i++ {
	// 	if v.Field(i).IsNil() {
	// 		os.Exit(1)
	// 	}
	// }
}
