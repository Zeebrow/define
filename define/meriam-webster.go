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
	DictApiKey string
	Headword   string
	// DefinitionSet *DefinitionSet
	// Suggestions   *[]string
}

func NewDictionary(apiKey string) *MWRawAPI {
	var mw MWRawAPI
	mw.DictApiKey = apiKey
	// mw.DefinitionSet = nil
	// mw.Suggestions = nil
	return &mw
}

/*The return value for Define*/
type DefinitionSet struct {
	Headword    string
	Entries     []Entry
	Suggestions *[]string
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

// @@@ maybe the suggestions should be returned as an error
func (mw *MWRawAPI) Lookup(headword string) (*DefinitionSet, error) {
	var ds DefinitionSet = DefinitionSet{}
	mw.Headword = headword
	ds.Headword = headword
	var t []byte
	url := fmt.Sprintf("https://www.dictionaryapi.com/api/v3/references/collegiate/json/%s?key=%s", headword, mw.DictApiKey)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to get data from '%s': '%v'", url, err)
		panic(err)
	}
	// fmt.Printf("%d\n", resp.StatusCode) // good feild for Dictionary type

	t, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading raw response from MW server: %v\n", err)
		fmt.Printf("Server responded: %d\n", resp.StatusCode)
		panic(err)
	}

	//err = mw.judge()
	err = json.Unmarshal(t, &ds.Suggestions)
	if err == nil {
		return &ds, errors.New("lookup failed for word: " + headword)
	}
	err = json.Unmarshal(t, &ds.Entries)
	if err != nil {
		panic("unexpected response from server: " + err.Error())
	}
	return &ds, nil
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

type SimpleHomonymJSON struct {
	Headword            string
	SimpleHomonymGroups []SimpleHomonymEntry
}

/*context for Definitions*/
type SimpleHomonymEntry struct {
	SenseIndex   int
	HomonymSense string
	PartOfSpeech string
	Definitions  []string
}

func (h *SimpleHomonymJSON) Print() {
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
{{range $HOMONYMGROUPS := .SimpleHomonymGroups}}
	┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	┃ {{.PartOfSpeech}} ({{.HomonymSense}}) 
	┠────────────────────────────────────────{{range $DEF := .Definitions}}
	┃ - {{.}}{{end}}
{{end}}` + "\n"
	ht := template.New("SimpleHomonymGrouping")
	t, err := ht.Parse(hjTemplate)
	if err != nil {
		fmt.Printf("template error: %v\n", err)
	}
	err = t.Execute(os.Stdout, h)
	if err != nil {
		fmt.Printf("template execute error: %v\n", err)
	}

}

func (r *DefinitionSet) PrintSuggestions() {
	// termwidth := 80
	var outputColumns int = 3
	sugs := *r.Suggestions //create new array and  copy

	/*to make nice columns*/
	maxSize := 0
	formattedSugs := make([]string, len(sugs))
	for _, s := range sugs {
		if len(s) > maxSize {
			maxSize = len(s)
		}
	}
	/*need too loop again bc of maxSize change*/
	for j, s := range sugs {
		if len(sugs) > 9 && j > 9 { //you can't possibly get >100 suggestions..
			formattedSugs[j] = s
		} else {
			formattedSugs[j] = s + " "
		}
		for i := 0; i < maxSize-len(s); i++ {
			formattedSugs[j] = " " + formattedSugs[j]
		}
	}

	/*do the printing*/
	c := 0
	for n := 1; n < 1+len(*r.Suggestions); n += outputColumns {
		for i := 0; i < outputColumns; i++ {
			if c == len(sugs)-1 {
				fmt.Println()
				return
			}
			// fmt.Printf("%d) %s\t", i+n, sugs[i+n])
			// fmt.Printf("%d) %s\t", i+n, formattedSugs[i+n])
			fmt.Printf("%s\t", formattedSugs[i+n])
			c++
		}
		fmt.Println()
	}
}

/*Parse the raw Merriam-Webster response for the text required to define a word's possible meanings.*/
func (r *DefinitionSet) GetSimpleHomonymJSON() SimpleHomonymJSON {
	var oj SimpleHomonymJSON
	oj.Headword = r.Headword
	for i, e := range r.Entries {
		hEntry := &SimpleHomonymEntry{SenseIndex: i + 1, HomonymSense: e.Meta.Id, PartOfSpeech: e.Fl, Definitions: e.Shortdef}
		oj.SimpleHomonymGroups = append(oj.SimpleHomonymGroups, *hEntry)
	}
	return oj
}

func (r *DefinitionSet) doForEntries() {
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
