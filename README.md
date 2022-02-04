# fuzzylinks
Link extractor and wordlist generator for fuzzing, written in go

This is based on Tomnomnom's Anew https://github.com/tomnomnom/anew and inspired by GerbenJavado's LinkFinder https://github.com/GerbenJavado/LinkFinder

fuzzylinks will append any new links or compoentns of links it finds to the file passed as the first argument. This file handling process is based on Anew and acts in the same manner. 

I'm sure there are other tools that do a very similar and likely better job as fuzzylinks, but I created this to learn the Go language and to suit my needs by focusing on wordlist generation for fuzzing. Alternatives I've tried have missed links for the sake of being cleaner with fewer false positives - fuzzylinks takes the opposite approach and attempts to extract as many potential links as possible.
In addition to link extraction, fuzzylinks can also split the links down into its component parts and include this in the output file too, using `-s`. This is so that when passed to a fuzzer, directories/endpoints that are only referred to within a link are included. Currently, this only includes each individual component of the link on it's own line with no combinations of link components for now.


## Installation
`go install github.com/N25sec/fuzzylinks@latest`

## Usage
`cat file.txt | fuzzylinks output.txt`
```
cat inputFile | fuzzylinks [-s|-r|-q|-v] outputFile

  -s (split) add lines to output for each component of found links
  
  -r remove leading slashes for more consistency when using with a fuzzing tool

  -q quiet mode prevents writing to stdout

  -v verbose mode - adds some print lines to stdout about the links found and possible link splits
```

## ToDo
- Create fileWrite function to reduce code duplication between the initial link matching and link splitting
- Add additinal combinations of split links rather than only a complete split into component parts


This is my first go project and I would very much appreciate any constructive criticism and issues/pull requests for improvements.

