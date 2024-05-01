package functions

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	// Matches one punctuation or one punctuation and a space before it
	onePunctuationRegex = regexp.MustCompile(` *([.,!?:;])`)
	// Matches groups of consecutive punctuation marks separated by spaces
	groupPunctuationsRegex = regexp.MustCompile(`( *[.,!?:; ]+[.,!?:;])`)
	// Matches aA followed by a vowel
	vowelsRegex = regexp.MustCompile(`(([\W]+[aA]) +([aAeEiIoOuUhH]))`)
	// Matches aA followed by a vowel in the beginning
	vowelsRegex2 = regexp.MustCompile(`((^[aA]) +([aAeEiIoOuUhH]))`)
)

// function the check if the output file has a valid extension
func IsValidExtension(outputFileName string) error {
	if !regexp.MustCompile(`\.txt$`).Match([]byte(outputFileName)) {
		return errors.New("enter a valid file extension (.txt)")
	}
	return nil
}

func SplitKeepSeparator(text, pattern string) []string {
	regex := regexp.MustCompile(pattern)
	//Extract indices of patterns
	indices := regex.FindAllStringIndex(text, -1)
	var parts []string
	start := 0
	for _, indexPair := range indices {
		// and append text from the start to the index
		parts = append(parts, text[start:indexPair[1]])
		// Change the start to where we end last time
		start = indexPair[1]
	}
	// from the last matched index to the last index of the text
	parts = append(parts, text[start:])
	// Delete the last string if it is empty
	if parts[len(parts)-1] == "" {
		parts = parts[0 : len(parts)-1]
	}
	return parts
}

// Perform actions on text
func ActionsModerator(word string, action string) (string, error) {
	var result string
	var err error
	switch action {
	case "low":
		result = toLower(word)
	case "up":
		result = toUpper(word)
	case "cap":
		result = capitalize(word)
	case "hex":
		word = strings.Trim(word, " ")
		temp, err := strconv.ParseInt(word, 16, 64)
		if err != nil {
			return word, err
		}
		result = strconv.Itoa(int(temp))
	case "bin":
		word = strings.Trim(word, " ")
		temp, err := strconv.ParseInt(word, 2, 64)
		if err != nil {
			return word, err
		}
		result = strconv.Itoa(int(temp))
	}
	return result, err
}

func toLower(text string) string {
	var result string
	for _, letter := range text {
		if letter >= 'A' && letter <= 'Z' {
			result += string(unicode.ToLower(letter))
		} else {
			result += string(letter)
		}
	}
	return result
}

func toUpper(text string) string {
	var result string
	for _, letter := range text {
		if letter >= 'a' && letter <= 'z' {
			result += string(unicode.ToUpper(letter))
		} else {
			result += string(letter)
		}
	}
	return result
}

func capitalize(text string) string {
	text = strings.ToLower(text)
	runeTxt := []rune(text)
	prevIsletter := false
	for index, letter := range runeTxt {
		if letter >= 'a' && letter <= 'z' {
			if !prevIsletter {
				runeTxt[index] = unicode.ToUpper(letter)
				prevIsletter = true
			}
		} else {
			prevIsletter = false
			continue
		}
	}
	return string(runeTxt)
}

func PunctFunc(text string) string {
	//add space after punc
	text = onePunctuationRegex.ReplaceAllString(text, "$1 ")
	//delete spaces between punctuations
	text = groupPunctuationsRegex.ReplaceAllStringFunc(text, func(match string) string {
		withoutSpaces := strings.ReplaceAll(match, " ", "")
		return withoutSpaces
	})
	return text
}

func VowelFix(text string) string {
	// add n to a if a is followed by a vowel
	text = vowelsRegex2.ReplaceAllString(text, "${2}n ${3}")
	text = vowelsRegex.ReplaceAllString(text, "${2}n ${3}")
	return text
}

// Clean whitespaces
func CleanSpaces(text string) string {
	var result []string
	var strRes string
	text = strings.TrimSpace(text)
	temp := strings.Split(text, " ")
	for _, word := range temp {
		if word != "" && word != " " {
			result = append(result, word)
		} else {
			continue
		}
	}
	strRes = strings.Join(result, " ")
	strRes = regexp.MustCompile(` \n+|\n+`).ReplaceAllString(strRes, "\n")
	strRes = regexp.MustCompile(`\n+ +`).ReplaceAllString(strRes, "\n")
	return strRes

}

func Quotes(text string) string {
	// add spaces after and before quote to become easy to recognize
	text = regexp.MustCompile(` +`).ReplaceAllString(text, " ")
	text = regexp.MustCompile(`''`).ReplaceAllString(text, " ' ' ")
	text = regexp.MustCompile(` '`).ReplaceAllString(text, "  '  ")
	text = regexp.MustCompile(`' `).ReplaceAllString(text, "  '  ")
	text = regexp.MustCompile(`^'`).ReplaceAllString(text, "  '  ")
	text = regexp.MustCompile(`'$`).ReplaceAllString(text, "  '  ")
	//match text between quotes surround it by quotes
	text = regexp.MustCompile(` +' +([-.,!?:; ]*\w*(?:[-.,!?:; ]*'*\w+)+[-.,!?:;]*) +' `).ReplaceAllString(text, " '$1' ")
	// fix punctuation if ruined by quotes process
	text = regexp.MustCompile(`(') +([-.,!?:;]) `).ReplaceAllString(text, "${1}${2} ")
	// add space if there are two attached closing quotes
	text = regexp.MustCompile(`''`).ReplaceAllString(text, "' '")
	return text
}
