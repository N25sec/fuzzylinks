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
	var quietMode bool
	var verboseMode bool
	var stripSlash bool
	var addSplits bool
	flag.BoolVar(&addSplits, "s", false, "add lines to output for each component of found links")
	flag.BoolVar(&stripSlash, "r", false, "remove leading slashes, for more consistent use with fuzzers")
	flag.BoolVar(&quietMode, "q", false, "quiet mode prevents writing to stdout")
	flag.BoolVar(&verboseMode, "v", false, "verbose mode")
	flag.Parse()
	fileOut := flag.Arg(0)
	lines := make(map[string]bool)
	var f io.WriteCloser
	if fileOut != "" {
		readFile, err := os.Open(fileOut)
		if err == nil {
			scanLines := bufio.NewScanner(readFile)
			for scanLines.Scan() {
				lines[scanLines.Text()] = true
			}
			readFile.Close()
		}
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
		re := regexp.MustCompile(`(((?:[a-zA-Z]{1,10}:\/\/)[^"'<>)(\s\n\/]{1,}\.[a-zA-Z0-9]{2,}[^"'<>)(\s\n]{0,})|((?:\/|\\){1,}[a-zA-Z0-9]{1,}[^"'<>)(\s\n]{0,})|((?:\.{1,2}(?:\/|\\)){1,}[a-zA-Z0-9]{1,}[^"'<>)(\s\n]{0,})|((?:"|'){1,1}[a-zA-Z0-9]{1,}[^"'<>)(\s\n](?:\\|\/){1,}[^"'<>)(\s\n]{0,}))`)
		match := re.FindStringSubmatch(line)

		if match == nil {
			continue
		} else {
			line = strings.Trim(match[0], "\"'")
			if stripSlash {
				line = strings.TrimLeft(line, "/\\")
			}
			if verboseMode && !quietMode {
				fmt.Printf(" MATCH >>> %s\n", line)
			}
			if lines[line] {
				continue
			}
			lines[line] = true
			if !quietMode {
				fmt.Println(line)
			}
			if fileOut != "" {
				fmt.Fprintf(f, "%s\n", line)
			}
			if addSplits {
				lSplit := regexp.MustCompile(`(?:\/|\\|\?|=|&)`)
				lineSplit := lSplit.Split(line, -1)
				if lineSplit != nil {
					for i := 0; i < len(lineSplit); i++ {
						splitLine := strings.Trim(lineSplit[i], "\"'")
						if stripSlash {
							splitLine = strings.TrimLeft(splitLine, "/\\")
						}
						if verboseMode && !quietMode {
							fmt.Printf(" SPLIT >>> %s\n", splitLine)
						}
						if strings.HasSuffix(splitLine, ":") {
							continue
						}
						if lines[splitLine] {
							continue
						}
						lines[splitLine] = true
						if !quietMode {
							fmt.Println(splitLine)
						}
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
}
