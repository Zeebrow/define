package define

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

type ProgramInfo struct {
	ProgramName string
	Version     string
	BuildDate   string
	CommitHash  string
	OS          string
	GoVer       string
}

func (p *ProgramInfo) GetInfo() string {
	return fmt.Sprintf("%s version '%s' (%s) build date: '%s'", p.ProgramName, p.Version, p.CommitHash, p.BuildDate)
}

/*
App globals
*/
var GlobalConfig GlobalConfiguration
var MWDictionaryApiKey string
var MWThesaurusApiKey string
var ConfigFileName = ".define.conf"
var ProgramName = "define"
var Version = "dev"
var BuildDate = ""
var CommitHash = ""
var ProgInfo = ProgramInfo{
	ProgramName: ProgramName,
	Version:     Version,
	BuildDate:   BuildDate,
	CommitHash:  CommitHash,
	OS:          runtime.GOOS,
	GoVer:       runtime.Version(),
}

type GlobalConfiguration struct {
	MWDictionaryApiKey string
	// MWThesaurusApiKey  string
}

type ConfigFile struct {
	Dictionary string `json:"dictionary"`
	// Thesaurus  string `json:"thesaurus,omitempty"`
}

func GetDefaultFilepath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Could not get user home directory... searching current directory for config")
		return ConfigFileName
	} else {
		return filepath.Join(dirname, ConfigFileName)
	}
}

func (cfg *GlobalConfiguration) LoadFromFile() error {
	var ks ConfigFile
	c, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	cfgJsonFilepath := path.Join(c, "define", "define.conf")
	r, err := os.ReadFile(cfgJsonFilepath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(r, &ks)
	if err != nil {
		return err
	}
	cfg.MWDictionaryApiKey = ks.Dictionary
	// cfg.MWThesaurusApiKey = ks.Thesaurus
	return nil
}
