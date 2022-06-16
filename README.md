# Define

## Description

Stop going to your browser just to get a GOOD definition for a word.

## Status

Working enough

## Usage

```
$ define -h
Usage of define:
  -debug
    	Print debug output
  -dict-api-key string
    	Overwrite any configuration of MW_DICTIONARY_API_KEY
  -f string
    	override location to config file (default is .MW-api-keys in home directory) (default "/home/zeebrow/.define.conf")
  -q	Use an alternative dictionary with more reliable api responses
  -stdin
    	Read a json from stdin instead of calling api
  -thes-api-key string
    	Overwrite any configuration of MW_THESAURUS_API_KEY
  -version
    	print the version and exit.
```

### Example Output

```
$ define monster
monster:
	noun (monster:1)
	----
	(1/3)	an animal of strange or terrifying shape
	(2/3)	one unusually large for its kind
	(3/3)	an animal or plant of abnormal form or structure

	adjective (monster:2)
	---------
	(1/1)	enormous or impressive especially in size, extent, or numbers

	noun (Gila monster)
	----
	(1/1)	a large, stout, venomous lizard (Heloderma suspectum) that has rough, bumpy, black and orange, pinkish, or yellowish skin, a thick tail, and venom glands in the lower lip and that is found especially in arid regions of the southwestern U.S. and northwestern Mexico

	noun (green-eyed monster)
	----
	(1/1)	jealousy imagined as a monster that attacks people —usually used with the

	noun (Loch Ness Monster)
	----
	(1/1)	a large, long-necked creature that has long been reported to exist in the waters of Loch Ness in Scotland

	noun (party monster)
	----
	(1/1)	a person known for frequent often wild partying : party animal

```

## Install

```
go build -ldflags "-X config.MWDictionaryApiKey=XZY -X config.MWThesaurusApiKey"
```

### Build args

```
-X main.MWDictionaryApiKey=XZY
-X main.MWThesaurusApiKey=ABC
```
### Environmentals

```
export MW_CONFIG_FILEPATH=/path/to/api/keys/dict
export MW_DICTIONARY_API_KEY=XYZ
export MW_THESAURUS_API_KEY=ABC
```

`MW_CONFIG_FILEPATH` should be a json file structured like this:
```json
{
    "dictionary": "dictionary api key",
    "thesaurus": "thesaurus api key"
}
~        
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


