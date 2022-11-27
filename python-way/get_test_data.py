#!/usr/bin/python3
import click
import json
import logging

from definition import Definition, Sense

logger = logging.getLogger()
logger.setLevel(logging.INFO)
if logger.getEffectiveLevel() == 10:
    # debug format
    sh = logging.StreamHandler()
    formatter = logging.Formatter('%(levelname)s:\t%(filename)s line (%(lineno)d): %(message)s')
    sh.setFormatter(formatter)
    logger.addHandler(sh)


# global
WORD = None
test_data = None # "cache"
please_help = False
please_list = False
registered_functions = {} 

pass_definition = click.make_pass_decorator(Definition)

@click.group()
@click.option("--word", help="word to define")
@click.pass_context
def cli(ctx, word):
    ctx.obj = Definition(word)

################################################
# Commands
################################################

@cli.command()
@pass_definition
def isword(definition):
    definition: Definition = definition
    """
    Determine if the word exists in Merriam-Webster's dictionary.
    If not, print out the suggestions returned.
    """
    is_a_word = False
    for t in definition.data:
        if type(t) == type("asdf"):
            break
        else:
            print("is a word")
            is_a_word = True
            return
    print("not a word. suggestions:")
    print(definition.data)
    return

@cli.command()
@pass_definition
def get_hwi(definition):
    definition: Definition = definition
    for defn in definition.data:
        if 'hwi' in defn.keys():
            print(json.dumps(defn['hwi'], indent=2))
            print("-------")

@cli.command()
@pass_definition
def testcase(definition):
    definition: Definition = definition
    """
    doc string is arg help
    """
    print("Executing test case logic")
    return 

@cli.command()
@pass_definition
def printout_data(definition):
    definition: Definition = definition
    """
    simply print the data
    """
    print(definition.data)

@cli.command()
@pass_definition
def get_sense_indeces(definition):
    definition: Definition = definition
    click.secho(f"found {len(definition.sseqs)} sseqs for word '{definition.word}'", fg='green')
    for sseq in definition.sseqs:
        for sense in sseq.senses:
            s = click.style(f"{sense.sense_number}-{sense.sense_letter}", fg="red")
            click.echo(s)

@cli.command()
@pass_definition
def get_defining_text(definition):#@@@
    definition: Definition = definition
    click.secho(f"found {len(definition.sseqs)} sseqs for word '{definition.word}'", fg='green')
    for sseq in definition.sseqs:
        for sn, sense in sseq.sense_map.items():
            senseno = click.style(sn, fg='yellow')
            defining_text = click.style(sense.dt, fg='blue')
            click.echo(f"{senseno}:\n{defining_text}")
    
@cli.command()
@pass_definition
def get_shortdefs(definition):#@@@
    definition: Definition = definition
    for i, d in enumerate(definition.shortdefs):
        click.echo(f"{i+1}) {d}")

cli()
