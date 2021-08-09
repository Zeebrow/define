package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	// TODO: but is it really?
	// Do you REALLY wanna maintain a log for this stupid thing?
	// "log"
)

type MWRawAPIResp struct {
	// Known issue:
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
	Def      []Sense
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
	// But again, fuck this golang BS for something so easy in Python
	SenseObject map[string]interface{}
}

type Etymology struct {
}

// Collection (list) of headwords and their definitions
type Sense struct {
	// index of a particular sense in a list of senses.
	SenseNumber string `json:"sn"`
	SLS         string `json:"sls,omitempty"`

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
			// fmt.Printf("Good response.\n")
			resps.isGood = true
		}
	}
	return
}

func get_dictionary_key() (key string) {
	// fi,err := os.UserHomeDir()
	// if err != nil {
	// 	fmt.Printf("could not get user home directory %s, exiting!\n", )
	// }
	// f := fi + ""
	f := ".MW-api-keys"
	r, err := ioutil.ReadFile(f)
	if err != nil {
		fmt.Printf("Could open API keys file: %s. Error: %v", f, err)
		// TODO: learn better exit strategy
		os.Exit(2)
	}

	// Do I really have to unmarshall here?
	// Gut says this is not the best way to do this
	type keys struct {
		Dictionary string `json:"dictionary"`
		Thesaurus  string `json:"thesaurus,omitempty"`
	}
	var ks keys
	err = json.Unmarshal(r, &ks)
	if err != nil {
		fmt.Printf("Could not get dictionary API keys from file: %s. Error: %v", f, err)
		// TODO: learn better exit strategy
		os.Exit(3)
	}
	key = ks.Dictionary
	return
}

// type byPos struct {
// 	pos         string
// 	definitions []string
// }

func (e *Entry) printShortdefs() {
	//fmt.Printf("\tFound %d shortdefs for %s (%s).\n", len(e.Shortdef), e.Meta.Id, e.Fl)
	fmt.Printf("\t%s\n\t%s\n", e.Fl, strings.Repeat("-", len(e.Fl)))
	for i, v := range e.Shortdef {
		fmt.Printf("\t(%d/%d)\t%s\n", i+1, len(e.Shortdef), v)
	}
	fmt.Println()
}

func (gr *GoodResponse) doForEntries() {
	if Debug {
		ne := len(gr.Entries)
		fmt.Printf("DEBUG: Number of entries: %d\n\n", ne)
	}
	for _, v := range gr.Entries {
		// fmt.Printf("%v", v)
		v.printShortdefs()
	}
	fmt.Println()
}

func (s *Sense) printSenseStuff() {
	fmt.Printf("%v\n", s)
}

func (gr *GoodResponse) PrintRawMWResponse() {
	fmt.Printf("%v", gr.Entries)
}

func GetMW(headword string, nsfw bool) {
	url := "https://www.dictionaryapi.com/api/v3/references/collegiate/json/"
	key_dict := get_dictionary_key()
	fetch := fmt.Sprintf("%s/%s?key=%s", url, headword, key_dict)
	resp, err := http.Get(fetch)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	raw_resp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading raw response from MW server: %v\n", err)
		fmt.Printf("Server responded: %d\n", resp.StatusCode)
		os.Exit(4)
	}

	// fmt.Printf("%s", raw_resp)
	var r MWRawAPIResp
	r.Resp = raw_resp
	br := r.judge()
	// fmt.Println("__________________________________________________________________________________")
	if br.isGood {
		fmt.Println(headword + ":")
		// br.printShortdefs()
		br.goodresponse.doForEntries()
		return
	} else {
		fmt.Printf("Bad response: \n")
		fmt.Printf("%v\n", br.badresponse.Suggestions)
	}
	// fmt.Printf("%s\n", raw_resp)
	// fmt.Println("__________________________________________________________________________________")
}
