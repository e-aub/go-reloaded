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
	splitRegex       = `\((low|up|cap), (\d+)\)|\(cap\)|\(up\)|\(low\)|\(hex\)|\(bin\)`          // Regular expression for splitting text with the pattern
	withoutNumRgx    = regexp.MustCompile(`\((low|up|cap|hex|bin)\)`)                            //Regular expression to match actions without numbers
	bigWithoutNumRgx = regexp.MustCompile(`((?:\b[\w_-]+\b[\W\s]*){1})\((low|up|cap|hex|bin)\)`) // Regular expression to match actions without numbers along with 1 word before
	withNumRgx       = regexp.MustCompile(`\(((low|up|cap)), (\d+)\)`)                           // Regular expression to match actions with numbers
	bigWithNumRgx    = `((?:\b[\w]+\b\W*\s*){%s})\(((?:%s),\s*\d+)\)`                            // Regular expression template for matching actions with numbers along with preceding words
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
	// add n to a if the next is a vowel
	text = []byte(functions.VowelFix(string(text)))
	sliced := functions.SplitKeepSeparator(string(text), splitRegex) // Split text using custom function
	for index, line := range sliced {                                // Iterate over each line in the sliced text
		if withoutNumRgx.Match([]byte(line)) { // Handle actions without numbers
			if match := bigWithoutNumRgx.FindStringSubmatch(line); match != nil { // If there is one word before, handle it
				word := match[1]
				action := match[2]
				word, err := functions.ActionsModerator(word, action) // Perform action on the word
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					continue
				}
				sliced[index] = bigWithoutNumRgx.ReplaceAllString(sliced[index], word) // Replace original word and action with processed word
			} else {
				action := withoutNumRgx.FindStringSubmatch(line)[1] // if there isn't a word before Prompt user to enter one word before the action
				fmt.Fprintf(os.Stderr, "enter one word before (%s)\n", action)

			}
		} else if match := withNumRgx.FindStringSubmatch(line); match != nil { // Handle actions with numbers
			asciiDelimiter := match[3]
			action := match[2]
			bigWithNumberCompiledRegex := regexp.MustCompile(fmt.Sprintf(bigWithNumRgx, asciiDelimiter, action)) //Compile regular expression pattern with delimiter and action
			if bigMatch := bigWithNumberCompiledRegex.FindStringSubmatch(line); bigMatch != nil {                // If there are preceding words, handle them
				words, err := functions.ActionsModerator(bigMatch[1], action) // Perform action on words
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					continue
				}
				sliced[index] = bigWithNumberCompiledRegex.ReplaceAllString(line, words) //Replace original words and action with processed words
			} else {
				fmt.Fprintln(os.Stderr, "enter a valid words number to "+action) // Prompt user to enter a valid number of words before the action
			}

		}
	}
	//Join lines together
	result := strings.Join(sliced, "")
	//Remove extra spaces
	result = functions.CleanSpaces(result)
	// Handle single punctuations
	result = functions.OnePunctFunc(result)
	//Handle successive punctuations
	result = functions.GroupPunctFunc(result)
	// single quotes fix
	result = functions.QuotesFix(result)
	//Write output the result and add a new line at the end
	err = ioutil.WriteFile(outputFileName, []byte(result+"\n"), 0777)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while writing content in the file\nError: %v\n", err)
		return
	}
}
