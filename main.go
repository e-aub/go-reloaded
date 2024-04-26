package main

import (
	"fmt"
	"go-reloaded/functions"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var (
	splitRegex            = `\((low|up|cap), (\d+)\)|\(cap\)|\(up\)|\(low\)|\(hex\)|\(bin\)`
	withoutNumberRegex    = regexp.MustCompile(`\((low|up|cap|hex|bin)\)`)
	bigWithoutNumberRegex = regexp.MustCompile(`((?:\b[\w_-]+\b[\W\s]*){1})\((low|up|cap|hex|bin)\)`)
	withNumberRegex       = regexp.MustCompile(`\(((low|up|cap)), (\d+)\)`)
	bigWithNumberRegex    = `((?:\b[\w]+\b\W*\s*){%s})\(((?:%s),\s*\d+)\)`
)

func main() {
	//Check for valid arguments
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Invalid arguments Usage : go run main.go <input.txt> <output.txt>")
		return
	}
	inputFileName := os.Args[1]
	outputFileName := os.Args[2]
	//Read file content
	text, err := ioutil.ReadFile(inputFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading content in the file\nError: %v\n", err)
		return
	}
	sliced := functions.SplitKeepSeparator(string(text), splitRegex)
	for index, line := range sliced {
		if withoutNumberRegex.Match([]byte(line)) {
			if match := bigWithoutNumberRegex.FindStringSubmatch(line); match != nil {
				word := match[1]
				action := match[2]
				word, err := functions.ActionsModerator(word, action)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					continue
				}
				sliced[index] = bigWithoutNumberRegex.ReplaceAllString(sliced[index], word)
			} else {
				fmt.Fprintln(os.Stderr, "enter at least one word before")

			}
		} else if match := withNumberRegex.FindStringSubmatch(line); match != nil {
			asciiDelimiter := match[3]
			action := match[2]
			bigWithNumberCompiledRegex := regexp.MustCompile(fmt.Sprintf(bigWithNumberRegex, asciiDelimiter, action))
			if bigMatch := bigWithNumberCompiledRegex.FindStringSubmatch(line); bigMatch != nil {
				words, err := functions.ActionsModerator(bigMatch[1], action)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					continue
				}
				sliced[index] = bigWithNumberCompiledRegex.ReplaceAllString(line, words)
			} else {
				fmt.Fprintln(os.Stderr, "enter a valid words number to "+action)
			}

		}
	}
	result := strings.Join(sliced, "")
	result = functions.OnePunctFunc(result)
	// //Handle successive punctuations
	result = functions.GroupPunctFunc(result)
	// // add n to a if a *vowel*
	result = functions.VowelFix(result)
	// // single quotes fix
	result = functions.QuotesFix(result)
	fmt.Println(result)
	//Write output the result
	err = ioutil.WriteFile(outputFileName, []byte(result+"\n"), 0777)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while writing content in the file\nError: %v\n", err)
		return
	}
}
