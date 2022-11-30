package define

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type MWRawAPI struct {
	DictApiKey  string
	Headword    string
	Response    *Response
	Suggestions *[]string
}

func NewApi(key string) *MWRawAPI {
	var mw MWRawAPI
	mw.DictApiKey = key
	mw.Response = nil
	mw.Suggestions = nil
	return &mw
}

type Response struct {
	Headword string
	Entries  []Entry
}

type Entry struct {
	Fl       string
	Meta     MWMetadata
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

func (mw *MWRawAPI) Define(headword string) {
	mw.Headword = headword
	var t []byte
	url := fmt.Sprintf("https://www.dictionaryapi.com/api/v3/references/collegiate/json/%s?key=%s", headword, mw.DictApiKey)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to get data from '%s': '%v'", url, err)
		panic(err)
	}
	// fmt.Printf("%d\n", resp.StatusCode)

	t, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading raw response from MW server: %v\n", err)
		fmt.Printf("Server responded: %d\n", resp.StatusCode)
		panic(err)
	}
	err = mw.judge(t)
	if err != nil {
		fmt.Printf("error judging: %s\n", err)
	}
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
func (sus *MWRawAPI) judge(t []byte) error {
	var sugs []string
	var resp Response

	resp.Headword = sus.Headword
	err := json.Unmarshal(t, &sugs)
	if err != nil {
		err = json.Unmarshal(t, &resp.Entries)
		if err != nil {
			return err
		}
		sus.Response = &resp
		return nil
	} else {
		sus.Suggestions = &sugs
	}
	return errors.New("error unmarshalling response")
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

type HomonymJSON struct {
	Headword      string
	HomonymGroups []HomonymEntry
}

/*context for Definitions*/
type HomonymEntry struct {
	SenseIndex   int
	HomonymSense string
	PartOfSpeech string
	Definitions  []string
}

func (h *HomonymJSON) Print() {
	underline := func(word string) string {
		underline := ""
		for i := 0; i < len(word); i++ {
			underline += "-"
		}
		return underline
	}
	ul := underline("Definition of " + h.Headword + "'': ")

	hjTemplate := `
Definitions for '{{.Headword}}':
` + ul + `
{{range $HOMONYMGROUPS := .HomonymGroups}}
	┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	┃ {{.PartOfSpeech}} ({{.HomonymSense}}) 
	┠────────────────────────────────────────{{range $DEF := .Definitions}}
	┃ - {{.}}{{end}}
{{end}}` + "\n"
	ht := template.New("HomonymGrouping")
	t, err := ht.Parse(hjTemplate)
	if err != nil {
		fmt.Printf("template error: %v\n", err)
	}
	err = t.Execute(os.Stdout, h)
	if err != nil {
		fmt.Printf("template execute error: %v\n", err)
	}

}

func (r *Response) GroupByHomonym() HomonymJSON {
	var oj HomonymJSON
	oj.Headword = r.Headword
	for i, e := range r.Entries {
		hEntry := &HomonymEntry{SenseIndex: i + 1, HomonymSense: e.Meta.Id, PartOfSpeech: e.Fl, Definitions: e.Shortdef}
		oj.HomonymGroups = append(oj.HomonymGroups, *hEntry)
	}
	return oj
}

func (r *Response) doForEntries() {
	var prevWasHomonym bool
	var currentHasHomonym bool
	/*group definitions together by homonym, separate by crappy little line*/
	for n, v := range r.Entries {
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

func GetMW(headword string, key string) {
	mw := NewApi(key)
	mw.Define(headword)
	if mw.Response != nil {
		homonyms := mw.Response.GroupByHomonym()
		homonyms.Print()
		// outpoutJson, err := json.MarshalIndent(homonyms, "", "  ")
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Println(string(outpoutJson))
	} else {
		fmt.Printf("'%s' isn't a word. Did you mean one of these?\ni%v", headword, mw.Suggestions)
	}
}
