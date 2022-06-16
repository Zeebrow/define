package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	// TODO: but is it really?
	// Do you REALLY wanna maintain a log for this stupid thing?
	// "log"
)

type MWRawAPIResp struct {
	// This needs to handle both good AND bad responses
	// All responses will have status 200.
	// not sure how to implement UnmarshalJSON interface, if that's best.
	Resp []byte
}

type BothResps struct {
	isGood       bool
	goodresponse GoodResponse
	badresponse  BadResponse
}

type BadResponse struct {
	Suggestions []string
}

type GoodResponse struct {
	Entries []Entry
}

type Entry struct {
	Fl   string
	Meta MWMetadata
	// Shortdef is probably all we need
	Shortdef []string
	Def      []Sense `json:"def"`
}

// General datastructure based on Merriam-Webster's API documentation
// here: https://dictionaryapi.com/products/json#term-headword
type MWDefn struct {
	// The word to be defined
	Headword string

	// grammatical context for the headword.
	// "formal", "slang", etc
	Label string

	// I think this is the glyphy representation of the headword
	Pronunciation []byte
}

type SSEQ struct {
	Senses []Sense
}

type GeneralSense struct {
	Type string
	// Challenge: Unmarshall to SenseObject,
	// based on GeneralSense.Type
	// I THINK this has merit...
	SenseObject map[string]interface{}
}

type Etymology struct {
}

// Collection (list) of headwords and their definitions
type Sense struct {
	// index of a particular sense in a list of senses.
	SenseNumber string   `json:"sn"`
	SLS         []string `json:"sls,omitempty"`

	// Defining text. Meat and potatoes.
	DT DefiningText `json:"dt"`
}

// Basically a janky map, in array form
type DefiningText struct {
	DTArray []DTItem
}

type DTItem struct {
	// https://dictionaryapi.com/products/json#sec-2.dt
	Item map[string]string
}

type MWMetadata struct {
	Id   string
	Uuid string
	Sort string
	Src  string //ignore

	// Indicates the section that the entry belongs to in print
	Section string

	//  lists all of the entry's headwords, variants, inflections, undefined entry words, and defined run-on phrases.
	// Each stem string is a valid search term that should match this entry.
	Stem      map[string]string
	Offensive bool
}

// This exists because I can't work with the raw api response
// from Merriam-Webster. It's an array of X.
// response = 200, regardless of whether word exists in dictionary,
// so we cram response into structs and see where it doesn't fit.
func (sus *MWRawAPIResp) judge() (resps *BothResps) {
	resps = new(BothResps)
	resps.isGood = false
	err := json.Unmarshal(sus.Resp, &resps.badresponse)
	if err != nil {
		// fmt.Printf("Could not unmarshal raw api response into badresponse (Good thing). %s\n", err)
		err = json.Unmarshal(sus.Resp, &resps.badresponse)
		if err != nil {
			// fmt.Printf("Could not unmarshal raw api response into badresponse (Good thing). %s\n", err)
			err = json.Unmarshal(sus.Resp, &resps.goodresponse.Entries)
			if err != nil {
				fmt.Printf("WTF error - Could not unmarshal into 'empty array'. %v\n", err)
				fmt.Printf("You definitely didn't see this coming. uhhh.... exiting!\n")
				os.Exit(1)
			}
			resps.isGood = true
		}
	}
	return
}

func (e *Entry) printShortdefs() {
	fmt.Printf("\t%s (%s)\n\t%s\n", e.Fl, e.Meta.Id, strings.Repeat("-", len(e.Fl)))
	for i, v := range e.Shortdef {
		fmt.Printf("\t(%d/%d)\t%s\n", i+1, len(e.Shortdef), v)
	}
	fmt.Println()
}

func (gr *GoodResponse) doForEntries() {
	if GlobalConfig.Debug {
		ne := len(gr.Entries)
		fmt.Printf("DEBUG: Number of entries: %d\n\n", ne)
	}
	for _, v := range gr.Entries {
		// fmt.Printf("%v", v)
		v.printShortdefs()
	}
	fmt.Println()
}

func (gr *GoodResponse) PrintRawMWResponse() {
	fmt.Printf("%v", gr.Entries)
}

func GetMW(headword string, stdin bool) {
	var t []byte
	var err error
	var url string
	if stdin {
		reader := bufio.NewReader(os.Stdin)
		t, err = reader.ReadBytes('\n')
		if err != io.EOF && err != nil {
			fmt.Println(err)
		}
	} else {
		url = fmt.Sprintf("https://www.dictionaryapi.com/api/v3/references/collegiate/json/%s?key=%s", headword, GlobalConfig.MWDictionaryApiKey)
		resp, err := http.Get(url)

		if err != nil {
			fmt.Printf("Failed to get data from '%s': '%v'", url, err)
			os.Exit(1)
		}

		t, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading raw response from MW server: %v\n", err)
			fmt.Printf("Server responded: %d\n", resp.StatusCode)
			os.Exit(4)
		}
	}

	var r MWRawAPIResp
	r.Resp = t
	br := r.judge()
	if br.isGood {
		fmt.Println(headword + ":")
		br.goodresponse.doForEntries()
	} else {
		fmt.Printf("Bad response: \n")
		fmt.Printf("%v\n", br.badresponse.Suggestions)
	}

	if GlobalConfig.Debug {
		fmt.Printf("url: %s\n", url)
		// fmt.Println(rawResp)
	}
	return
}
