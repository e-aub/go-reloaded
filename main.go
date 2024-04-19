package main

import (
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
	inputFileName          string
	outputFileName         string
	input                  []byte
	err                    error
	lowRegex               = regexp.MustCompile(`(\(low, (\d+)\))`)
	upRegex                = regexp.MustCompile(`(\(up, (\d+)\))`)
	capRegex               = regexp.MustCompile(`(\(cap, (\d+)\))`)
	onePunctuationRegex    = regexp.MustCompile(`(\s*([.,!?:;]+)+\s*)`)
	groupPunctuationsRegex = regexp.MustCompile(`([.,!?:;\s]+[.,!?:;])`)
	vowelsRegex            = regexp.MustCompile(`((\s+[aA])\s([aeiouh]))`)
	quotesRegex            = regexp.MustCompile(`('\s*([-.,!?:;]*\w+(?:[-.,!?:;\s]+\w+)+[-.,!?:;]*)\s*')`)
	delimiter              = 1
	match                  []string
)

func main() {
	//Check for valid arguments
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Invalid arguments Usage : go run main.go <input.txt> <output.txt>")
		return
	}
	inputFileName = os.Args[1]
	outputFileName = os.Args[2]
	//Read file content
	input, err = ioutil.ReadFile(inputFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading content in the file\nError: %v\n", err)
		return
	}
	text := splitString(input)
	//extracting regex match
	for i := 0; i <= len(text)-1; i++ {
		if len(text) <= 1 {
			break
		}
		if text[i] == "(low)" {
			text[i-1] = strings.ToLower(text[i-1])
			text = deleteActions(text, i)
		} else if lowRegex.FindStringSubmatch(text[i]) != nil {
			match = lowRegex.FindStringSubmatch(text[i])
			delimiter, _ = strconv.Atoi(match[2])
			if i-delimiter < 0 {
				fmt.Fprintf(os.Stderr, "The amount of words that you want to lowercase is not enough")
				continue
			}
			for j := i - 1; j+delimiter >= i; j-- {
				text[j] = strings.ToLower(text[j])
			}
			text = deleteActions(text, i)
			delimiter = 1
		} else if text[i] == "(up)" {
			text[i-1] = strings.ToUpper(text[i-1])
			text = deleteActions(text, i)
		} else if upRegex.FindStringSubmatch(text[i]) != nil {
			match = upRegex.FindStringSubmatch(text[i])
			delimiter, _ = strconv.Atoi(match[2])
			if i-delimiter < 0 {
				fmt.Fprintf(os.Stderr, "The amount of words that you want to uppercase is not enough")
				continue
			}
			for j := i - 1; j+delimiter >= i; j-- {
				text[j] = strings.ToUpper(text[j])
			}
			text = deleteActions(text, i)
			delimiter = 1
		} else if text[i] == "(cap)" {
			text[i-1] = capitalize(text[i-1])
			text = deleteActions(text, i)
		} else if capRegex.FindStringSubmatch(text[i]) != nil {
			match = capRegex.FindStringSubmatch(text[i])
			delimiter, _ = strconv.Atoi(match[2])
			if i-delimiter < 0 {
				fmt.Fprintf(os.Stderr, "The amount of words that you want to capitalize is not enough")
				continue
			}
			for j := i - 1; j+delimiter >= i; j-- {
				text[j] = capitalize(text[j])
			}
			delimiter = 1
			text = deleteActions(text, i)
		} else if text[i] == "(hex)" {
			coverted, err := strconv.ParseInt(text[i-1], 16, 0)
			if err != nil {
				panic(err)
			}
			text[i-1] = strconv.Itoa(int(coverted))
			text = deleteActions(text, i)
		} else if text[i] == "(bin)" {
			coverted, _ := strconv.ParseInt(text[i-1], 2, 0)
			text[i-1] = strconv.Itoa(int(coverted))
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
	fmt.Println(strText)
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

func splitString(input []byte) []string {
	split := strings.Split(string(input), " ")
	for index, word := range split {
		//Check if the next element of the slice is the second part of the action, Example : ["(cap," "333)"] and concatenate them
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

func capitalize(s string) string {
	rs := []rune(s)

	if len(rs) == 0 {
		return ""
	}
	rs[0] = unicode.ToUpper(rs[0])
	for i := 1; i < len(rs); i++ {
		rs[i] = unicode.ToLower(rs[i])
	}
	return string(rs)
}

func deleteActions(text []string, toDelete int) []string {
	text = append(text[:toDelete], text[toDelete+1:]...)
	return text
}

func onePunctFunc(text string) string {
	text = onePunctuationRegex.ReplaceAllString(text, "$2 ")
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
