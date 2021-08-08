# Define

## Description

Stop going to your browser just to get a GOOD definition for a word.

## Status

Working enough

## Usage

```
define [word]

# long output opened similar to $MANPAGER
define --more [word]
```

## Features 

* short, quick CLI lookup
* long, "more info" pager reader
* spelling suggestions

## Want

### long-term

* `--more` syntax highlighting
* other dictionary integrations: slang, programming, pop culture, etc.
* link to more info when unable to find match
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
* you can't configure where the api keys are (must be in same dir as executable)
* I'm not worthy to parse the responses from Merriam-Webster. 

### Why it's hard

Arrays of arrays of arrays. JSON objects are represented as `["key": {}]` arrays. Everything is an array. To the point where it makes me think Python would be a much better choice for this, overall:

They have a duck-typed api? I should use a duck-typed program to interpret it.

Go is super fast, and I need to learn it - but it's verbose. Python isn't fast, but it's flexible enough to get me somewhere maintainable.

I think there's a way to get the best of both worlds.

The basic usage will be handled by go, anything more and go will pipe the json to a python script.

