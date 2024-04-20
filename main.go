package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

//Declaring variables

var (
	text                   []string
	base                   int
	inputFileName          string
	outputFileName         string
	input                  string
	err                    error
	lowRegex               = regexp.MustCompile(`(\(low, (\d+)\))`)
	upRegex                = regexp.MustCompile(`(\(up, (\d+)\))`)
	capRegex               = regexp.MustCompile(`(\(cap, (\d+)\))`)
	onePunctuationRegex    = regexp.MustCompile(`(\s*([.,!?:;]+)+\s*)`)
	lastPunctuationRegex   = regexp.MustCompile(`(\s*([.,!?:;]+)+\s*$)`)
	groupPunctuationsRegex = regexp.MustCompile(`([.,!?:;\s]+[.,!?:;])`)
	vowelsRegex            = regexp.MustCompile(`((\s*[aA])\s([aeiouh]))`)
	quotesRegex            = regexp.MustCompile(`('\s*([-.,!?:;]*\w+(?:[-.,!?:;\s]+\w+)+[-.,!?:;]*)\s*')`)
	delimiter              = 1
	match                  []string
)

func readFile(inputFileName string) (string, error) {
	content, err := ioutil.ReadFile(inputFileName)
	return string(content), err
}

func splitString(input string) []string {
	split := strings.Split(input, " ")
	//Check if the next element of the slice is the second part of the action, Example : ["(cap," "333)"] and concatenate them
	for index, word := range split {
		if word == "(low," || word == "(up," || word == "(cap," {
			if index != len(split)-1 && isTheSecondPart(split[index+1]) {
				split[index] = split[index] + " " + split[index+1]
				split = append(split[:index+1], split[index+2:]...)
			}
		} else {
			continue
		}
	}
	return split
}

func isTheSecondPart(s string) bool {
	if s == "" {
		return false
	}
	isNumber := true
	for i := 0; i <= len(s)-2; i++ {
		if s[i] >= '0' && s[i] <= '9' {
			isNumber = true
		} else {
			isNumber = false
			break
		}
	}
	if isNumber && s[len(s)-1] == ')' {
		return true
	}
	return false
}

func capitalize(s string) (string, error) {
	var err error
	rs := []rune(s)

	if len(rs) == 0 {
		return "", err
	}
	if !(rs[0] >= 'a' && rs[0] <= 'z' || rs[0] >= 'A' && rs[0] <= 'Z') {
		err = errors.New("capitalize : not valid syntax")
	}
	rs[0] = unicode.ToUpper(rs[0])
	for i := 1; i < len(rs); i++ {
		rs[i] = unicode.ToLower(rs[i])
	}
	return string(rs), err
}

func deleteActions(text []string, toDelete int) []string {
	text = append(text[:toDelete], text[toDelete+1:]...)
	return text
}

func onePunctFunc(text string) string {
	text = onePunctuationRegex.ReplaceAllString(text, "$2 ")
	text = lastPunctuationRegex.ReplaceAllString(text, "$2")
	return text
}

func groupPunctFunc(text string) string {
	text = groupPunctuationsRegex.ReplaceAllStringFunc(text, func(match string) string {
		withoutSpaces := strings.ReplaceAll(match, " ", "")
		return withoutSpaces
	})
	return text
}

func vowelFix(text string) string {
	text = vowelsRegex.ReplaceAllString(text, "${2}n ${3}")
	return text
}

func writeToOutput(toWrite string, outputFileName string) {
	// Write output
	err = ioutil.WriteFile(outputFileName, []byte(toWrite), 0777)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while writing content in the file\nError: %v\n", err)
		return
	}
}

func main() {
	//Check for valid arguments
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Invalid arguments Usage : go run main.go <input.txt> <output.txt>")
		return
	}
	inputFileName = os.Args[1]
	outputFileName = os.Args[2]
	//Read file content
	input, err = readFile(inputFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading content in the file\nError: %v\n", err)
		return
	}
	input = onePunctFunc(input)
	input = groupPunctFunc(input)
	//split string
	text = splitString(input)
	//matching user prompts
	for i := 0; i < len(text); i++ {
		if len(text) <= 1 {
			break
		}
		if text[i] == "(low)" {
			k := i
			if k-1 < 0 {
				fmt.Fprintf(os.Stderr, "The amount of words that you want to lowercase is not enough")
				continue
			}
			for text[k-1] == "" {
				k -= 1
			}
			text[k-1] = strings.ToLower(text[k-1])
			text = deleteActions(text, i)
		} else if lowRegex.FindStringSubmatch(text[i]) != nil {
			match = lowRegex.FindStringSubmatch(text[i])
			delimiter, _ = strconv.Atoi(match[2])
			k := i
			for text[k-1] == "" {
				k -= 1
			}
			if k-delimiter < 0 {
				fmt.Fprintf(os.Stderr, "The amount of words that you want to lowercase is not enough")
				continue
			}
			for j := k - 1; j+delimiter >= k; j-- {
				text[j] = strings.ToLower(text[j])
			}
			text = deleteActions(text, i)
			delimiter = 1
		} else if text[i] == "(up)" {
			k := i
			if k-1 < 0 {
				fmt.Fprintf(os.Stderr, "The amount of words that you want to uppercase is not enough")
				continue
			}
			for text[k-1] == "" {
				k -= 1
			}
			text[k-1] = strings.ToUpper(text[k-1])
			text = deleteActions(text, i)
		} else if upRegex.FindStringSubmatch(text[i]) != nil {
			match = upRegex.FindStringSubmatch(text[i])
			delimiter, _ = strconv.Atoi(match[2])
			k := i
			for text[k-1] == "" {
				k -= 1
			}
			if k-delimiter < 0 {
				fmt.Fprintf(os.Stderr, "The amount of words that you want to uppercase is not enough")
				continue
			}
			for j := k - 1; j+delimiter >= k; j-- {
				text[j] = strings.ToUpper(text[j])
			}
			text = deleteActions(text, i)
			delimiter = 1
		} else if text[i] == "(cap)" {
			k := i
			if k-1 < 0 {
				fmt.Fprintf(os.Stderr, "The amount of words that you want to capitalize is not enough")
				continue
			}
			for text[k-1] == "" {
				k -= 1
			}
			text[k-1], err = capitalize(text[k-1])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			text = deleteActions(text, i)
		} else if capRegex.FindStringSubmatch(text[i]) != nil {
			match = capRegex.FindStringSubmatch(text[i])
			delimiter, _ = strconv.Atoi(match[2])
			k := i
			for text[k-1] == "" {
				k -= 1
			}
			if k-delimiter < 0 {
				fmt.Fprintf(os.Stderr, "The amount of words that you want to capitalize is not enough")
				continue
			}
			for j := k - 1; j+delimiter >= k; j-- {
				text[j], err = capitalize(text[j])
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					break
				}
			}
			delimiter = 1
			text = deleteActions(text, i)
		} else if text[i] == "(hex)" {
			k := i
			if k-1 < 0 {
				fmt.Fprintf(os.Stderr, "The amount of words that you want to convert to decimal is not enough")
				continue
			}
			for text[k-1] == "" {
				k -= 1
			}

			coverted, err := strconv.ParseInt(text[k-1], 16, 0)
			if err != nil {
				panic(err)
			}
			text[k-1] = strconv.Itoa(int(coverted))
			text = deleteActions(text, i)
		} else if text[i] == "(bin)" {
			k := i
			if k-1 < 0 {
				fmt.Fprintf(os.Stderr, "The amount of words that you want to convert to decimal is not enough")
				continue
			}
			for text[k-1] == "" {
				k -= 1
			}

			coverted, _ := strconv.ParseInt(text[k-1], 2, 0)
			text[k-1] = strconv.Itoa(int(coverted))
			text = deleteActions(text, i)
		}
	}
	// join slice
	strText := strings.Join(text, " ")
	// Handle space before and after punctuations
	strText = onePunctFunc(strText)
	//Handle successive punctuations
	strText = groupPunctFunc(strText)
	// add n to a if a *vowel*
	strText = vowelFix(strText)
	// single quotes fix
	strText = quotesRegex.ReplaceAllString(strText, "'$2'")
	writeToOutput(strText, outputFileName)
}
