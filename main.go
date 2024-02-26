package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	// Regular expressions to match transformation actions with delimiters
	var (
		capRegexN = regexp.MustCompile(`(\(cap, )(\d+)\)`)
		lowRegexN = regexp.MustCompile(`(\(low, )(\d+)\)`)
		upRegexN  = regexp.MustCompile(`(\(up, )(\d+)\)`)
		cap       = "(cap)"
		low       = "(low)"
		up        = "(up)"
		hex       = "(hex)"
		bin       = "(bin)"
		delimeter int
		textSlice []string
		err       error
	)
	const oneWord = 1 // Constant for single word delimiter

	// Check if correct number of arguments is provided
	if len(os.Args) != 3 {
		fmt.Println("Enter valid arguments (go run main.go <input.txt> <output.txt>)")
		os.Exit(1)
	}
	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Read input file
	text := ReadFile(inputFile)
	// Split text into words
	textSlice = SplitString(text)

	// Iterate over words in text and Check if word is a transformation action
	for index, word := range textSlice {
		switch {
		case capRegexN.MatchString(word) || word == cap:
			// Process capitalization transformation
			if capRegexN.MatchString(word) {
				matches := capRegexN.FindStringSubmatch(word)
				delimeter, err = strconv.Atoi(matches[2])
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			} else if word == cap {
				delimeter = oneWord
			}
			textSlice = ProcessTransformation(cap, delimeter, index, textSlice)

		case lowRegexN.MatchString(word) || word == low:
			// Process lowercase transformation
			if lowRegexN.MatchString(word) {
				matches := lowRegexN.FindStringSubmatch(word)
				delimeter, err = strconv.Atoi(matches[2])
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			} else if word == low {
				delimeter = oneWord
			}
			textSlice = ProcessTransformation(low, delimeter, index, textSlice)

		case upRegexN.MatchString(word) || word == up:
			// Process uppercase transformation
			if upRegexN.MatchString(word) {
				matches := upRegexN.FindStringSubmatch(word)
				delimeter, err = strconv.Atoi(matches[2])
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			} else if word == up {
				delimeter = oneWord
			}
			textSlice = ProcessTransformation(up, delimeter, index, textSlice)

		case word == hex:
			// Process hexadecimal transformation
			textSlice = ProcessTransformation(hex, 1, index, textSlice)

		case word == bin:
			// Process binary transformation
			textSlice = ProcessTransformation(bin, 1, index, textSlice)

		default:
			// Skip non-transformation words
			continue
		}
	}

	// Merge modified text into single string
	result := MergeString(textSlice)
	// Perform punctuation and vowel transformations
	result = PunctuationsAndVowels(result)
	// Write result to output file
	WriteFile(outputFile, result)
}

// ReadFile reads content from the input file
func ReadFile(inputFile string) string {
	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(input) == 0 {
		fmt.Fprintln(os.Stderr, "File is empty")
		os.Exit(1)
	}
	return string(input)
}

// SplitString splits the text into words
func SplitString(text string) []string {
	var result []string
	textSlice := strings.Fields(text)
	for i, word := range textSlice {
		if ParatheseExist(word) {
			textSlice[i-1] = textSlice[i-1] + " " + word
			textSlice[i] = ""
		} else {
			continue
		}
	}
	for _, word := range textSlice {
		if word != "" {
			result = append(result, word)
		} else {
			continue
		}
	}
	return result
}

// ParatheseExist checks if parentheses exist in the word
func ParatheseExist(s string) bool {
	for _, letter := range s {
		if letter == '(' {
			return false
		}
		if letter == ')' {
			return true
		}
	}
	return false
}

// ProcessTransformation processes various transformations
func ProcessTransformation(action string, delimeter int, index int, textSlice []string) []string {
	if index-delimeter < 0 {
		fmt.Printf("There is not %d words before\n", delimeter)
		os.Exit(1)
	}
	if action == "(cap)" {
		textSlice[index] = ""
		for delimeter != 0 {
			textSlice[index-delimeter] = Capitalize(textSlice[index-delimeter])
			delimeter--
		}
	} else if action == "(low)" {
		textSlice[index] = ""
		for delimeter != 0 {
			textSlice[index-delimeter] = strings.ToLower(textSlice[index-delimeter])
			delimeter--
		}
	} else if action == "(up)" {
		textSlice[index] = ""
		for delimeter != 0 {
			textSlice[index-delimeter] = strings.ToUpper(textSlice[index-delimeter])
			delimeter--
		}
	} else if action == "(hex)" || action == "(bin)" {
		textSlice[index] = ""
		var decimalValue int64
		var base int
		if action == "(hex)" {
			base = 16
		} else {
			base = 2
		}
		decimalValue, err := strconv.ParseInt(textSlice[index-delimeter], base, 64)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		textSlice[index-delimeter] = strconv.Itoa(int(decimalValue))
	}
	return textSlice
}

// Capitalize capitalizes the first letter of the word
func Capitalize(s string) string {
	var result string
	if IsLower(rune(s[0])) {
		result = result + string(s[0]-32) + s[1:]
	} else {
		result = s
	}
	return result
}

// IsLower checks if the letter is lowercase
func IsLower(let rune) bool {
	if let >= 'a' && let <= 'z' {
		return true
	}
	return false
}

// PunctuationsAndVowels processes punctuation and vowels
func PunctuationsAndVowels(text string) string {
	vowels := regexp.MustCompile(`(\s[Aa])(\s+[aAeEiIoOuUhH])`)
	punctMark := regexp.MustCompile(`('\s+)([^']+)(\s+')`)
	beforePunct := regexp.MustCompile(`\s+([,.:;?!])`)
	afterPunct := regexp.MustCompile(`([,.:;!?]+)\s*`)
	endOfLine := regexp.MustCompile(`([,.:;!?])\s+$`)
	result := vowels.ReplaceAllString(text, "${1}n${2}")
	result = beforePunct.ReplaceAllString(result, "$1")
	result = afterPunct.ReplaceAllString(result, "$1 ")
	result = endOfLine.ReplaceAllString(result, "$1")
	result = punctMark.ReplaceAllString(result, "'${2}'")
	return result
}

// MergeString merges words into a single string
func MergeString(textSlice []string) string {
	var result string
	for _, word := range textSlice {
		if word != "" {
			result = result + word + " "
		} else {
			continue
		}
	}
	return strings.TrimSpace(result)
}

// WriteFile writes output to the output file
func WriteFile(outputFile string, output string) {
	err := ioutil.WriteFile(outputFile, []byte(output), 0700)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
