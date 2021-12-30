# fuzzylinks
Link extractor and wordlist generator for fuzzing, written in go

This is based on Tomnomnom's Anew https://github.com/tomnomnom/anew and inspired by GerbenJavado's LinkFinder https://github.com/GerbenJavado/LinkFinder

I'm sure there are other tools that do a very similar and likely better job as fuzzylinks, but I couldn't find one that suited my needs and focused on wordlist generation for fuzzing. Alternatives I've tried have missed links for the sake of being cleaner with fewer false positives - fuzzylinks takes the opposite approach and attempts to extract as many potential links as possible.
In addition to link extraction, fuzzylinks will also split the links down into its component parts and include this in the output file too. This is so that when passed to a fuzzer, directories/endpoints that are only referred to within a link are included. 

fuzzylinks will append any new links or compoentns of links it finds to the file passed as the first argument. The file handling is based on Anew and acts in the same manner. 

This is my first go project so I would very much appreciate any constructive criticism and issues/pull requests for improvements.

Cheers!

N25


## Usage
```
cat inputFile | fuzzylinks [-q|-v] outputFile

  -q quiet mode prevents writing to stdout

  -v verbose mode - adds some print lines to stdout about the links found and possible link splits
```

## ToDo
- Create fileWrite function to reduce code duplication between the initial link matching and link splitting
- Add additinal combinations of split links rather than only a complete split into component parts
