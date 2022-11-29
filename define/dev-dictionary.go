package define

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type APIResp struct {
	Word     string    `json:"word"`
	Meanings []Meaning `json:"meanings"`
}

type Defn struct {
	Definition string   `json:"definition"`
	Synonyms   []string `json:"synonyms"`
	Example    string   `json:"example"`
}

type Meaning struct {
	PartOfSpeech string `json:"partOfSpeech"`
	Definitions  []Defn `json:"definitions"`
}

func devGetDef(word string) APIResp {
	var dictionaryapi_response []APIResp
	res, err := http.Get("https://api.dictionaryapi.dev/api/v2/entries/en/" + word)
	if err != nil {
		fmt.Println(err)
	}
	resp_code := res.StatusCode
	if resp_code != 200 {
		fmt.Printf("ERROR request returned response code %d. Probably not a word: %s\n", resp_code, word)
		os.Exit(1)
	}
	b, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(b, &dictionaryapi_response)
	if err != nil {
		fmt.Println(err)
	}

	return dictionaryapi_response[0]
}

func devAllDefns(d Defn) {
	const defn_template = ` {{.}}
	`
	const synonyms_template = `
		Synonyms:{{if .}}{{range $SYNONS := .}}
		- {{.}}{{end}}
		{{else}} None{{end}}
	`

	dt := template.New("dt")
	_dt, err := dt.Parse(defn_template)
	if err != nil {
		fmt.Printf("Defn template error: %v", err)
	}
	err = _dt.Execute(os.Stdout, d.Definition)
	if err != nil {
		fmt.Printf("Defn template execute error: %v", err)
	}

	st := template.New("st")
	t, err := st.Parse(synonyms_template)
	if err != nil {
		fmt.Printf("Syn template error: %v", err)
	}
	err = t.Execute(os.Stdout, d.Synonyms)
	if err != nil {
		fmt.Printf("Syn template execute error: %v", err)
	}
}

func DevPrintMeanings(word string) {
	const tab = "  "
	resp := devGetDef(word)
	uline := strings.Repeat("-", len(resp.Word))
	fmt.Printf("\n%s\n", resp.Word)
	fmt.Printf("%s\n\n", uline)

	for i := 0; i < len(resp.Meanings); i++ {
		fmt.Printf("%s%s:\n", tab, resp.Meanings[i].PartOfSpeech)
		fmt.Printf("%s%v\n", tab, strings.Repeat("-", len(resp.Meanings[i].PartOfSpeech)+1))
		num_defns := len(resp.Meanings[i].Definitions)
		for d := 0; d < num_defns; d++ {
			fmt.Printf("%2s(%d/%d)", tab, d+1, num_defns)
			devAllDefns(resp.Meanings[i].Definitions[d])
		}

		fmt.Println()
	}
}
