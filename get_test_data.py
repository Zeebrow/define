#!/usr/bin/python3
import requests
import json
import argparse

"""
when you need to filter definition data, add a function to do so and wrap it with the `@command` decorator.

what it does:
Turn a python function into a cli argument by simply adding the `command` decorator

how it works:
- only supports 'boolean' flags (either run or don't run the function)
- function docstring = help

what it needs:
- automatic completion generation:
    - a runcom or executable
- convert function name underscores to hyphens

"""

# global
WORD = None
test_data = None # "cache"
please_help = False
please_list = False
registered_functions = {} 
def register_function(fname, f):
    registered_functions[fname] = f

class MyBoolean(argparse.Action):
    # pretty much copy-pasted from the argparse module
    def __init__(
        self,
                 option_strings,
                 dest,
                 default=None,
                 type=None,
                 choices=None,
                 required=False,
                 help=None,
                 metavar=None
                 ):

        _option_strings = []
        for option_string in option_strings:
            _option_strings.append(option_string)

        if help is not None and default is not None and default is not argparse.SUPPRESS:
            help += " (default: %(default)s)"

        super().__init__(
            option_strings=_option_strings,
            dest=dest,
            nargs=0,
            default=default,
            type=type,
            choices=choices,
            required=required,
            help=help,
            metavar=metavar)

    def __call__(self, parser, namespace, values, option_string=None):
        if self.dest in registered_functions.keys():
            exec(f"{registered_functions[self.dest](test_data)}")
        if option_string in self.option_strings:
            setattr(namespace, self.dest, not option_string.startswith('--no-'))

    def format_usage(self):
        return ' | '.join(self.option_strings)


DICT_URL = lambda w, k: f"https://www.dictionaryapi.com/api/v3/references/collegiate/json/{w}?key={k}"
THES_URL = lambda w, k: f"https://www.dictionaryapi.com/api/v3/references/thesaurus/json/{w}?key={k}"
API_KEYS_FILE="/home/zeebrow/.local/etc/define.conf"

def get_api_key(_type='dictionary'):
    conf_file = None
    with open(API_KEYS_FILE, 'r') as f:
        conf_file = json.load(f)
    return conf_file[_type]

def define(word: str):
    print("defining...")
    r = requests.get(DICT_URL(word, get_api_key()), headers={"Content-Type": "application/json"})
    return r.json()

parser = argparse.ArgumentParser(add_help=False)
dostuff_group = parser.add_mutually_exclusive_group()
dostuff_group.add_argument("-w", "--word", dest='__word', action='store', required=False, help="word to lookup")
parser.add_argument("-l", "--list", dest='__list', action='store_true', required=False, help="list available tests")
# note, silently ignored kwargs...
parser.add_argument("-h", "--help", dest='__help', action='store_true', required=False, help="print what options are available and their descriptions")
WORD = parser.parse_known_args()[0].__word
__list = parser.parse_known_args()[0].__list
__help = parser.parse_known_args()[0].__help
if __list:
    please_list = True
if __help:
    please_help = True


def command(f):
    parser.add_argument(
        f"--{f.__name__}",
        required=False,
        help="".join([f.__doc__ if f.__doc__ is not None else f.__name__]),
        action=MyBoolean,
    )
    register_function(f.__name__, f)

    global test_data
    if not (please_list or please_help) and test_data is None:
        test_data = define(WORD)

    def wrapped_arg_function(*args, **kwargs):
        f(test_data)

    return wrapped_arg_function

@command
def isword(data):
    """
    Determine if the word exists in Merriam-Webster's dictionary.
    If not, print out the suggestions returned.
    """
    is_a_word = False
    for t in data:
        if type(t) == type("asdf"):
            continue
        else:
            print("is a word")
            is_a_word = True
            return
    print("not a word. suggestions:")
    print(data)
    return

@command
def get_hwi(data):
    for defn in data:
        if 'hwi' in defn.keys():
            print(json.dumps(defn['hwi'], indent=2))
            print("-------")

@command
def testcase(data):
    """
    doc string is arg help
    """
    print("Executing test case logic")
    return 

@command
def printout_data(data):
    """
    simply print the data
    """
    print(data)

@command
def tcdata_01(data):#
    """
    print definitions with more than one `meta['id']` 
    """
    meta_ids = []
    for defn in data:
        if len(defn['meta']['id'].split(':')) > 1:
            meta_ids.append(defn['meta']['id'])
    print(meta_ids)


if __name__ == '__main__':
    args = parser.parse_args()
    if please_list:
        for r in registered_functions.keys():
            print(f"--{r}")
        exit()
    if please_help:
        parser.print_help()
        exit()
