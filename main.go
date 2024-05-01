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
	// Regular expression for splitting text with the pattern
	splitRegex = `\((low|up|cap), (\d+)\)|\(cap\)|\(up\)|\(low\)|\(hex\)|\(bin\)`
	//Regular expression to match actions without numbers
	withoutNumRgx = regexp.MustCompile(`\((low|up|cap|hex|bin)\)`)
	// Regular expression to match actions without numbers along with 1 word before
	bigWithoutNumRgx = regexp.MustCompile(`((?:\b[\w_-]+\b[\W ]*){1})\((low|up|cap|hex|bin)\)`)
	// Regular expression to match actions with numbers
	withNumRgx = regexp.MustCompile(`\(((low|up|cap)), (\d+)\)`)
	// Regular expression template for matching actions with numbers along with preceding words
	bigWithNumRgx = `((?:\b[\w]+\b\W* *){%s})\(((?:%s), *\d+)\)`
)

func main() {
	//Check for valid arguments
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "\x1b[31mInvalid arguments Usage : go run main.go <input.txt> <output.txt>\x1b[0m")
		return
	}
	inputFileName := os.Args[1]
	outputFileName := os.Args[2]
	//Check output file extension

	if err := functions.IsValidExtension(outputFileName); err != nil {
		fmt.Fprintln(os.Stderr, "\x1b[31m"+err.Error()+"\x1b[0m")
		os.Exit(1)
	}
	//Read file content
	text, err := ioutil.ReadFile(inputFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\x1b[31mError while reading content in the file\nError: %v\n\x1b[0m", err)
		return
	}
	// add n to a if the next is a vowel
	text = []byte(functions.VowelFix(string(text)))
	// Split text with actions
	sliced := functions.SplitKeepSeparator(string(text), splitRegex)
	for index, line := range sliced {
		// Handle actions without numbers
		if withoutNumRgx.Match([]byte(line)) {
			// If there is one word before, handle it
			if match := bigWithoutNumRgx.FindStringSubmatch(line); match != nil {
				word := match[1]
				action := match[2]
				// Perform action on the word
				word, err := functions.ActionsModerator(word, action)
				if err != nil {
					fmt.Fprintln(os.Stderr, "\x1b[31m"+err.Error()+"\x1b[0m")
					continue
				}
				// Replace original word and action with processed word
				sliced[index] = bigWithoutNumRgx.ReplaceAllString(sliced[index], " "+word+" ")
			} else {
				// if there isn't a word before Prompt user to enter one word before the action
				action := withoutNumRgx.FindStringSubmatch(line)[1]
				fmt.Fprintf(os.Stderr, "\x1b[31menter one word before (%s)\n\x1b[0m", action)

			}
			// Handle actions with numbers
		} else if match := withNumRgx.FindStringSubmatch(line); match != nil {
			asciiDelimiter := match[3]
			action := match[2]
			//Compile regular expression pattern with delimiter and action
			bigWithNumberCompiledRegex := regexp.MustCompile(fmt.Sprintf(bigWithNumRgx, asciiDelimiter, action))
			// If there are preceding words, handle them
			if bigMatch := bigWithNumberCompiledRegex.FindStringSubmatch(line); bigMatch != nil {
				// Perform action on words
				words, err := functions.ActionsModerator(bigMatch[1], action)
				if err != nil {
					fmt.Fprintln(os.Stderr, "\x1b[31m"+err.Error()+"\x1b[0m")
					continue
				}
				//Replace original words and action with processed words
				sliced[index] = bigWithNumberCompiledRegex.ReplaceAllString(line, " "+words+" ")
			} else {
				// Prompt user to enter a valid number of words before the action
				fmt.Fprintln(os.Stderr, "\x1b[31menter a valid words number to "+action+"\x1b[0m")
			}

		}
	}
	//Join lines together
	result := strings.Join(sliced, "")
	//Remove extra spaces
	result = functions.CleanSpaces(result)
	// Handle punctuations
	result = functions.PunctFunc(result)
	// single quotes fix
	result = functions.Quotes(result)
	//Remove extra spaces
	result = functions.CleanSpaces(result)
	// add n to a if the next is a vowel
	result = functions.VowelFix(result)
	//Write output the result and add a new line at the end
	err = ioutil.WriteFile(outputFileName, []byte(result+"\n"), 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\x1b[31mError while writing content in the file\nError: %v\n", err.Error()+"\x1b[0m")
		return
	}
}
