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

func textConversionModerator(action string, index int, delimiter int) error {
	toDeleteIndex := index
	for index != 0 && text[index-1] == "" {
		index -= 1
	}
	if index-delimiter < 0 {
		return errors.New("the amount of words that you want to " + action + " is not enough")
	}
	for j := index - 1; j+delimiter >= index; j-- {
		if action == "low" {
			text[j] = strings.ToLower(text[j])
		} else if action == "cap" {
			temp, err := capitalize(text[j])
			if err != nil {
				return err
			}
			text[j] = temp
		} else if action == "up" {
			text[j] = strings.ToUpper(text[j])
		}
		text = deleteActions(text, toDeleteIndex)
	}
	return nil
}

func readFile(inputFileName string) (string, error) {
	content, err := ioutil.ReadFile(inputFileName)
	return string(content), err
}

func splitString(input string) []string {
	split := strings.Split(input, " ")
	//Check if the next element of the slice is the second part of the action, Example : ["(cap," "333)"] and concatenate them
	for index, word := range split {
		if word == "(low," || word == "(up," || word == "(cap," {
			if index != len(split)-1 && isTheSecondPart(split[index+1]) { // input = onePunctFunc(input)
				// input = groupPunctFunc(input)
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

func toDecimal(index int, base int) error {

	var temp int64
	var err error
	for index != -1 && text[index-1] == "" {
		index -= 1
	}
	if index-1 < 0 {
		return errors.New("invalid syntax")
	}
	temp, err = strconv.ParseInt(text[index-1], base, 64)
	if err != nil {
		return err
	}
	text[index-1] = strconv.Itoa(int(temp))
	return nil
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
	//split string
	text = splitString(input)
	//matching user prompts
	for i := 0; i < len(text); i++ {
		if len(text) <= 1 {
			fmt.Fprintln(os.Stderr, "insert at least one word")
			break
		}
		if text[i] == "(low)" {
			err = textConversionModerator("low", i, 1)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
		} else if text[i] == "(up)" {
			err = textConversionModerator("up", i, 1)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
		} else if text[i] == "(cap)" {
			err = textConversionModerator("cap", i, 1)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
		} else if text[i] == "(bin)" {
			err = toDecimal(i, 2)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			text = deleteActions(text, i)
		} else if text[i] == "(hex)" {
			err = toDecimal(i, 16)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			text = deleteActions(text, i)
		} else if lowRegex.FindStringSubmatch(text[i]) != nil {
			match = lowRegex.FindStringSubmatch(text[i])
			delimiter, _ = strconv.Atoi(match[2])
			err = textConversionModerator("low", i, 1)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
		} else if upRegex.FindStringSubmatch(text[i]) != nil {
			match = upRegex.FindStringSubmatch(text[i])
			delimiter, _ = strconv.Atoi(match[2])
			err = textConversionModerator("up", i, 1)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
		} else if capRegex.FindStringSubmatch(text[i]) != nil {
			match = capRegex.FindStringSubmatch(text[i])
			delimiter, _ = strconv.Atoi(match[2])
			err = textConversionModerator("cap", i, 1)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
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
