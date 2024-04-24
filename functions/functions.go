package functions

import "regexp"

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
	var err error
	return word, err
}
