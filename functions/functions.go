package functions

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	onePunctuationRegex    = regexp.MustCompile(`(\s*([.,!?:;]+)+\s*)`)  // Matches one or more punctuation marks surrounded by whitespace
	lastPunctuationRegex   = regexp.MustCompile(`(\s*([.,!?:;]+)+\s*$)`) // Matches punctuation marks at the end of a string
	groupPunctuationsRegex = regexp.MustCompile(`([.,!?:;\s]+[.,!?:;])`) // Matches groups of consecutive punctuation marks
	vowelsRegex            = regexp.MustCompile(`((\s+[aA])\s+([aAeEiIoOuUhH]))`)
	vowelsRegex2           = regexp.MustCompile(`((^[aA])\s+([aAeEiIoOuUhH]))`)
	quotesRegex            = regexp.MustCompile(`('\s*([-.,!?:;]*\w+(?:[-.,!?:;\s]*\w+)+[-.,!?:;]*)\s*')`)
)

func SplitKeepSeparator(text, pattern string) []string {
	regex := regexp.MustCompile(pattern)
	indices := regex.FindAllStringIndex(text, -1)
	var parts []string
	start := 0
	for _, indexPair := range indices {
		parts = append(parts, text[start:indexPair[1]])
		start = indexPair[1]
	}
	parts = append(parts, text[start:])
	if parts[len(parts)-1] == "" {
		parts = parts[0 : len(parts)-1]
	}
	return parts
}

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

func OnePunctFunc(text string) string {
	text = onePunctuationRegex.ReplaceAllString(text, "$2 ")
	text = lastPunctuationRegex.ReplaceAllString(text, "$2")
	return text
}

func GroupPunctFunc(text string) string {
	text = groupPunctuationsRegex.ReplaceAllStringFunc(text, func(match string) string {
		withoutSpaces := strings.ReplaceAll(match, " ", "")
		return withoutSpaces
	})
	return text
}

func VowelFix(text string) string {
	text = vowelsRegex2.ReplaceAllString(text, "${2}n ${3}")
	text = vowelsRegex.ReplaceAllString(text, "${2}n ${3}")
	return text
}

func QuotesFix(text string) string {
	result := quotesRegex.ReplaceAllString(text, "'$2'")
	return result
}

func CleanSpaces(text string) string {
	var result []string
	text = strings.TrimSpace(text)
	temp := strings.Split(text, " ")
	for _, word := range temp {
		if word != "" && word != " " {
			result = append(result, word)
		} else {
			continue
		}
	}
	return strings.Join(result, " ")

}
