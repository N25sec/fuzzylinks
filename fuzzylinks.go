package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func main() {

	// quietMode is a flag that prevents writing extracted links to stdout.
	var quietMode bool

	// verboseMode adds additional stdout prints.
	var verboseMode bool

	// stripSlashes removed the leading slashes from matches and splits
	var stripSlash bool

	// flag.BoolVar sets up the flag parsing.
	flag.BoolVar(&stripSlash, "s", false, "strip leading slashes, for more consistent use with fuzzers")

	flag.BoolVar(&quietMode, "q", false, "quiet mode prevents writing to stdout")

	flag.BoolVar(&verboseMode, "v", false, "verbose mode")
	// flag.Parse() initiates the parsing of the passed cli arguments.
	flag.Parse()

	// fileOut is the output file that this tool's stdout is written to, the first argument.
	fileOut := flag.Arg(0)

	// `lines` is the map of lines taken from stdin.
	lines := make(map[string]bool)

	// Not sure what the below does yet.
	var f io.WriteCloser

	if fileOut != "" {
		// Read entire file into map, if the file exists.
		readFile, err := os.Open(fileOut)
		if err == nil {
			scanLines := bufio.NewScanner(readFile)

			for scanLines.Scan() {
				lines[scanLines.Text()] = true
			}
			readFile.Close()
		}
		// We don't care about the error here as if the file is a new one, it will be created later but would error here.
		f, err = os.OpenFile(fileOut, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open file >> %s\n", err)
			return
		}
		defer f.Close()
	}

	scanStdin := bufio.NewScanner(os.Stdin)

	for scanStdin.Scan() {
		line := scanStdin.Text()
		// Here we can do some logic on the lines read in to extract the links.
		re := regexp.MustCompile(`(((?:[a-zA-Z]{1,10}:\/\/)[^"'<>)(\s\n\/]{1,}\.[a-zA-Z0-9]{2,}[^"'<>)(\s\n]{0,})|((?:\/|\\){1,}[a-zA-Z0-9]{1,}[^"'<>)(\s\n]{0,})|((?:\.{1,2}(?:\/|\\)){1,}[a-zA-Z0-9]{1,}[^"'<>)(\s\n]{0,})|((?:"|'){1,1}[a-zA-Z0-9]{1,}[^"'<>)(\s\n](?:\\|\/){1,}[^"'<>)(\s\n]{0,}))`)
		match := re.FindStringSubmatch(line)

		if match == nil {
			continue
		} else {
			line = strings.Trim(match[0], "\"'")

			if stripSlash {
				line = strings.TrimLeft(line, "/\\")
			}

			// Print and write to file.
			if verboseMode && !quietMode {
				fmt.Printf(" MATCH >>> %s\n", line)
			}
			// Skip if line is already in the map.
			if lines[line] {
				continue
			}
			// Avoid Duplicates.
			lines[line] = true
			if !quietMode {
				fmt.Println(line)
			}
			if fileOut != "" {
				fmt.Fprintf(f, "%s\n", line)
			}
			// Split on slashes and query strings - makes a lot of noise but may be useful.
			lSplit := regexp.MustCompile(`(?:\/|\\|\?|=|&)`)
			lineSplit := lSplit.Split(line, -1)

			if lineSplit != nil {
				for i := 0; i < len(lineSplit); i++ {
					splitLine := strings.Trim(lineSplit[i], "\"'")

					if stripSlash {
						splitLine = strings.TrimLeft(splitLine, "/\\")
					}
					// Print and write to file.
					if verboseMode && !quietMode {
						fmt.Printf(" SPLIT >>> %s\n", splitLine)
					}
					if strings.HasSuffix(splitLine, ":") {
						continue
					}
					// Skip if splitLine is already in the map.
					if lines[splitLine] {
						continue
					}
					// Avoid Duplicates.
					lines[splitLine] = true

					if !quietMode {
						fmt.Println(splitLine)
					}
					// Wrtie splitLine to fileOut.
					if fileOut != "" {
						fmt.Fprintf(f, "%s\n", splitLine)
					}
				}
			} else {
				if !quietMode {
					fmt.Println("Line Split Failed")
				}
			}
		}
	}
}
