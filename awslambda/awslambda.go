package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Zeebrow/define/define"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, word string) {
	MW_DICT_API_KEY := os.Getenv("MW_DICT_API_KEY")
	if MW_DICT_API_KEY == "" {
		fmt.Println("MW_DICT_API_KEY is not set")
		os.Exit(1)
	}
	mw := define.NewApi(MW_DICT_API_KEY)
	mw.Define(word)
	if mw.Response != nil {
		homonyms := mw.Response.GroupByHomonym()
		outpoutJson, err := json.MarshalIndent(homonyms, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(outpoutJson))
	} else {
		fmt.Printf("'%s' isn't a word. Did you mean one of these?\ni%v", word, mw.Suggestions)
	}
}

func main() {
	lambda.Start(HandleRequest)
	// HandleRequest(nil, "obese")
}
