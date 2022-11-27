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
)

type MWRawAPI struct {
	Headword     string
	Resp         []byte
	isGood       bool
	goodresponse GoodResponse
	badresponse  []string
}

func (mwresp *MWRawAPI) lookup(headword string, stdin bool) error {
	mwresp.Headword = headword
	var t []byte
	var url string
	var err error

	if stdin {
		reader := bufio.NewReader(os.Stdin)
		t, err = reader.ReadBytes('\n')
		if err != io.EOF && err != nil {
			fmt.Println(err)
			return err
		}
	} else {
		url = fmt.Sprintf("https://www.dictionaryapi.com/api/v3/references/collegiate/json/%s?key=%s", headword, GlobalConfig.MWDictionaryApiKey)
		resp, err := http.Get(url)

		if err != nil {
			fmt.Printf("Failed to get data from '%s': '%v'", url, err)
			return err
		}

		t, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading raw response from MW server: %v\n", err)
			fmt.Printf("Server responded: %d\n", resp.StatusCode)
			return err
		}
	}
	if GlobalConfig.Debug {
		fmt.Printf("url: %s\n", url)
	}
	mwresp.Resp = t

	return nil
}

type GoodResponse struct {
	Headword string
	Entries  []Entry
}

type Prs struct {
	Mw string `json:"mw,omitempty"`
}
type HWI struct {
	Prs []Prs `json:"prs,omitempty"`
}

// top-level respsonse object
type Entry struct {
	Fl       string
	Meta     MWMetadata
	Hwi      HWI      `json:"hwi,omitempty"`
	Shortdef []string `json:"shortdef"`
}

type MWMetadata struct {
	Id      string
	Uuid    string
	Sort    string
	Src     string //ignore
	Section string // Indicates the section that the entry belongs to in print
	//  lists all of the entry's headwords, variants, inflections, undefined entry words, and defined run-on phrases.
	// Each stem string is a valid search term that should match this entry.
	Stem      map[string]string
	Offensive bool
}

func (meta *MWMetadata) hom() string {
	if strings.Contains(meta.Id, ":") {
		return strings.Split(meta.Id, ":")[0]
	}
	return meta.Id
}

func (meta *MWMetadata) homNum() string {
	if strings.Contains(meta.Id, ":") {
		return strings.Split(meta.Id, ":")[1]
	}
	return ""
}

// This exists because I can't work with the raw api response
// from Merriam-Webster. It's an array of X.
// response = 200, regardless of whether word exists in dictionary,
// so we cram response into structs and see where it doesn't fit.
func (sus *MWRawAPI) judge() { //aka unmarshal?
	sus.isGood = false
	err := json.Unmarshal(sus.Resp, &sus.badresponse)
	if err != nil {
		// potentially a good response
		err = json.Unmarshal(sus.Resp, &sus.goodresponse.Entries)
		if err != nil {
			fmt.Println(string(sus.Resp))
			ProgInfo.NewBugReport("Bug: Could not unmarshal json into empty array", err.Error())
			os.Exit(1)
		}
		sus.goodresponse.Headword = sus.Headword
		sus.isGood = true
	}
}

func (e *Entry) printShortdefs() {
	// fmt.Printf("\t%s (%s)\n\t%s\n", e.Fl, e.Meta.Id, strings.Repeat("-", len(e.Fl)))
	var printString string
	if e.Meta.homNum() != "" {
		printString = fmt.Sprintf("\t%s (%s - Homonym %s)\n", e.Fl, e.Meta.hom(), e.Meta.homNum())
	} else {
		printString = fmt.Sprintf("\t%s (%s)\n", e.Fl, e.Meta.Id)
	}
	underline := fmt.Sprintf("\t%s\n", strings.Repeat("-", len(printString)))
	fmt.Printf("%s", printString)
	fmt.Printf("%s", underline)
	for i, v := range e.Shortdef {
		fmt.Printf("\t(%d/%d)\t%s\n", i+1, len(e.Shortdef), v)
	}
	fmt.Println()
}

type outputJSON struct {
	Headword      string
	HomonymGroups []HomonymEntry
}
type HomonymEntry struct {
	HomonymSense string
	PartOfSpeech string
	Definitions  []string
}

func (gr *GoodResponse) sortByHomonym() outputJSON {
	oj := &outputJSON{Headword: gr.Headword}
	for _, e := range gr.Entries {
		hEntry := &HomonymEntry{HomonymSense: e.Meta.Id, PartOfSpeech: e.Fl, Definitions: e.Shortdef}
		oj.HomonymGroups = append(oj.HomonymGroups, *hEntry)
	}
	return *oj
}

func (gr *GoodResponse) doForEntries() {
	if GlobalConfig.Debug {
		fmt.Printf("DEBUG: Number of entries: %d\n\n", len(gr.Entries))
	}
	var prevWasHomonym bool
	var currentHasHomonym bool
	/*group definitions together by homonym, separate by crappy little line*/
	for n, v := range gr.Entries {
		if n == 0 {
			prevWasHomonym = true
		}
		if v.Meta.homNum() != "" {
			currentHasHomonym = true
		} else {
			currentHasHomonym = false
		}

		if currentHasHomonym && !prevWasHomonym {
			prevWasHomonym = true
			fmt.Println("_________")
			v.printShortdefs()
		} else if currentHasHomonym && prevWasHomonym {
			prevWasHomonym = true
			v.printShortdefs()
		} else if !currentHasHomonym && prevWasHomonym {
			prevWasHomonym = false
			fmt.Println("_________")
			v.printShortdefs()
		} else { //!currentHasHomonym && !prevWasHomonym
			prevWasHomonym = false
			fmt.Println("_________")
			v.printShortdefs()
		}

	}
}

func GetMW(headword string, stdin bool) {
	var r MWRawAPI
	err := r.lookup(headword, stdin)
	if err != nil {
		panic(err)
	}
	r.judge()
	if r.isGood {
		// r.goodresponse.doForEntries()
		homonyms := r.goodresponse.sortByHomonym()
		outpoutJson, err := json.MarshalIndent(homonyms, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(outpoutJson))
		// fmt.Printf("'%s':\n\n", homonyms.headword)
		// for i, h := range homonyms.HomonymGroups {
		// 	fmt.Println("-------------------")
		// 	fmt.Printf("sense %d - %s (%s)\n", i+1, h.HomonymSense, h.PartOfSpeech)
		// 	for j, d := range h.Definitions {
		// 		fmt.Printf("\t%d) %s\n", j+1, d)
		// 	}
		// }
	} else {
		fmt.Printf("Bad response: \n")
		fmt.Printf("%v\n", r.badresponse)
	}
}
