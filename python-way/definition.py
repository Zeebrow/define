from pathlib import Path
import json
import pickle
import logging
import requests
from typing import List, Dict
import click
import re

DICT_URL = lambda w, k: f"https://www.dictionaryapi.com/api/v3/references/collegiate/json/{w}?key={k}"
THES_URL = lambda w, k: f"https://www.dictionaryapi.com/api/v3/references/thesaurus/json/{w}?key={k}"
API_KEYS_FILE="/home/zeebrow/.local/etc/define.conf"

logger = logging.getLogger(__name__)

def get_api_key(_type='dictionary'):
    conf_file = None
    with open(API_KEYS_FILE, 'r') as f:
        conf_file = json.load(f)
    return conf_file[_type]


class Definition:
    def __init__(self, word):
        self.word = word
        cachefile = Path('cachefile')
        # create cachefile if it doesn't exist
        if not cachefile.exists():
            with cachefile.open("wb") as cf:
                pickled_obj_base = [{"word": {"example-raw-data": True}}]
                pickle.dump(pickled_obj_base, cf)
                logger.info(f"Created new cachefile '{cachefile.absolute()}'")
        else:
            logger.debug("reading from cache...")
        
        self.data = None
        #try to get definition from cache before making request
        with cachefile.open("rb") as cf:
            unpickled_obj = pickle.load(cf)
            json.dumps(unpickled_obj, indent=2)
            for defn in unpickled_obj:
                if word in defn.keys():
                    logger.debug(f"found cached definition for '{word}'...")
                    self.data = defn[word]
                    break

            if self.data == None:
                logger.debug(f"getting definition for '{word}'...")
                r = requests.get(DICT_URL(word, get_api_key()), headers={"Content-Type": "application/json"})
                with cachefile.open("wb") as cf:
                    unpickled_obj.append({word: r.json()})
                    pickle.dump(unpickled_obj, cf)
                self.data = r.json()

        self.shortdefs = self.data[0]['shortdef']
        self.meta = self.data[0]["meta"]
        self.defs = self.data[0]["def"]
        self.sseqs: List[Sseq] = []
        for d in self.defs:
            for k,v in d.items():
                if k == "sseq":
                    self.sseqs.append(Sseq(v))
                else:
                    logger.error(f"Not a sseq: {k}")
                    continue

class Sseq:
    """an array of senses"""
    def __init__(self, data: list):
        self.raw = data
        self.senses: List[Sense] = []
        self.sense_map: Dict[str, Sense] = {}
        for sseq in data:
            sense_num = 0
            sense_letter = ' '
            for sense in sseq:
                # we want to know the sense number (sn) before creating a sense object
                if sense[0] != "sense":
                    continue
                sn: str = sense[1]['sn']
                # sn will always be a single number, single letter, or string of 'number + " " + letter'
                if len(sn) == 3:
                    sense_num = int(sn[0])
                    sense_letter = sn[-1]
                elif len(sn) == 1:
                    if sn.isnumeric():
                        sense_num = int(sn)
                        # if there is only one sense, we assign a letter even though one will not be given 
                        sense_letter = 'a' 
                    else:
                        sense_letter = sn
                else:
                    raise Exception("invalid sense number "+ sn)
                new_sense = Sense(data=sense, sense_num=sense_num, sense_let=sense_letter)
                self.senses.append(new_sense)
                self.sense_map[str(sense_num) + " " + sense_letter] = new_sense

class Sense:
    def __init__(self, data: list, sense_num=None, sense_let=None):
        if data[0] != "sense":
            logger.error(f"wanted 'sense', got {data[0]}")
            raise TypeError("not a sense array")
        self.sense_number = sense_num
        self.sense_letter = sense_let
        self.dt =  DefiningText(data=data[1]['dt'], sense_index=str(self.sense_number)+self.sense_letter)

class Text:
    def __init__(self, data: List[str]):
        if data[0] != "text":
            raise Exception("oopsie!...")
        self.raw_text = data[1]
        self.definitions: List[str] = []
        self.parse_definitions()

    def parse_definitions(self):
        # https://dictionaryapi.com/products/json#sec-2.xreftokens
        defs = self.raw_text.split("{bc}")
        defs.pop(0)
        _defs = []
        for d in defs:
            #sx_re = re.compile(r'\{sx\|([a-z]*)\|\|\}')
            a = self.sx_strip(d)
            b = a.strip()
            _defs.append(b)
        self.definitions = _defs
    def sx_strip(self, text):
        b = text.replace('{sx|', '')
        c = b.replace('||}', '')
        return c


    def __repr__(self) -> str:
        return '\n'.join(self.definitions)

class DefiningText:
    def __init__(self, data: List, sense_index: str):
        self.sense_index = sense_index
        self.defining_texts: List[Text] = []
        for i in data:
            if i[0] == "text":
                self.defining_texts.append(Text(i))

    def __repr__(self) -> str:
        rtn = ''
        rtn += self.sense_index + '\n'
        for t in self.defining_texts:
            rtn += str(t)
        return rtn

