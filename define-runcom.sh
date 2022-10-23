#!/bin/bash
mkdir -p "$XDG_DATA_HOME/define"
DEFINED_WORDS="$XDG_DATA_HOME/define/defined_words.txt"
[ ! -f "$DEFINED_WORDS" ] && touch "$DEFINED_WORDS"

# skip if first arg is a flag
[ "${1:0:1}" == '-' ] && define "$@" && exit

define "$@"
rc=$?
if [ "$rc" -eq 0 ]; then 
  # add to list if grep returns 1
  grep -q "$@" "$DEFINED_WORDS" || echo "$@" >> "$DEFINED_WORDS"
fi
