# Define

## Description

Stop going to your browser just to get a GOOD definition for a word.

## Status

Working enough

## Usage

```
zeebrow@zeebrow-pc:define(master)$ ./define 
Specify a word to look up!
Usage of define:
  -a	Save output to file
  -debug
    	Print debug output
  -dict-api-key string
    	Overwrite any configuration of MW_DICTIONARY_API_KEY
  -m	Print more detailed definitions
  -s	Save output to file
  -thes-api-key string
    	Overwrite any configuration of MW_THESAURUS_API_KEY
  -w string
    	Save output to file
  -x	Print potentially offensive definitions
```

## Install

```
go build -ldflags "-X config.MWDictionaryApiKey=XZY -X config.MWThesaurusApiKey"
```

### Build args

```
-X config.MWDictionaryApiKey=XZY
-X config.MWThesaurusApiKey=ABC
```
### Environmentals

```
export MW_DICTIONARY_API_KEY=XYZ
export MW_THESAURUS_API_KEY=ABC
```

## Features 

* short, quick CLI lookup
* long, "more info" pager reader
* spelling suggestions

## Want

### long-term

* client-server model
* `--more` syntax highlighting
* other dictionary integrations: slang, programming, pop culture, etc.
* local caching of content looked up in the past

### short-term

* `-m` formatting
* Install for cross-platform
* Sort-by-PoS 
* deal with `{things}` in def/senses, e.g.

```
zeebrow@zeebrow-pc:define(master)$ ./ezget tool | jq '.[0].def[0].sseq[2]'
[
  [
    "sense",
    {
      "sn": "3 a",
      "dt": [
        [
          "text",
          "{bc}one who is used or manipulated by another "
```

## Why it sucks

* you need your own Merriam-Webster api key to get access to `-m` better definitions
* output format :pukeemoji:

### Why it's hard

Arrays of arrays of arrays. JSON objects are represented as `["key": {}]` arrays. Everything is an array. To the point where it makes me think Python would be a much better choice for some things.

Go is super fast, and I need to learn it. Python isn't fast, but it's flexiblity could go a long way for specific parts here.


