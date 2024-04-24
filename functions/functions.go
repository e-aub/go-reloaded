package functions

import (
	"regexp"
	"unicode"
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
	// var result string
	var err error
	// switch action {
	// case "low":
	// case "up" :
	// case "cap" :
	// case "hex" :
	// case "bin" :
	// }
	return word, err
}

func toLower(text string) string {
	runeText := []rune(text)
	var result string
	for _, letter := range runeText {
		if letter >= 'A' && letter <= 'Z' {
			result += string(unicode.ToLower(letter))
		} else {
			result += string(letter)
		}
	}
	return result
}

func toUpper(text string) string {
	runeText := []rune(text)
	var result string
	for _, letter := range runeText {
		if letter >= 'a' && letter <= 'z' {
			result += string(unicode.ToUpper(letter))
		} else {
			result += string(letter)
		}
	}
	return result
}
