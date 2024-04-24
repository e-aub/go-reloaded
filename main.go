package main

import (
	"fmt"
	"go-reloaded/functions"
	"os"
	"regexp"
)

func main() {
	text := "  (cap)"
	// actionAndDelimiterRegex := `\((low|up|cap), (\d+)\)`
	splitRegex := `\((low|up|cap), (\d+)\)|\(cap\)|\(up\)|\(low\)|\(hex\)|\(bin\)`
	withoutNumberRegex := regexp.MustCompile(`\((low|up|cap|hex|bin)\)`)
	bigWithoutNumberRegex := regexp.MustCompile(`((?:\b[\w_-]+\b[\W\s]*){1})\((low|up|cap|hex|bin)\)`)
	withNumberRegex := regexp.MustCompile(`\(((low|up|cap)), (\d+)\)`)
	bigWithNumberRegex := `((?:\b[\w]+\b\W*\s*){%s})\(((?:%s),\s*\d+)\)`
	sliced := functions.SplitKeepSeparator(text, splitRegex)

	fmt.Println(sliced)

	for index, line := range sliced {
		if withoutNumberRegex.Match([]byte(line)) {
			if match := bigWithoutNumberRegex.FindStringSubmatch(line); match != nil {
				word := match[1]
				action := match[2]
				word, err := functions.ActionsModerator(word, action)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					continue
				}
				sliced[index] = bigWithoutNumberRegex.ReplaceAllString(sliced[index], "|"+word+"|")
			} else {
				fmt.Fprintln(os.Stderr, "enter at least one word before")

			}
		} else if match := withNumberRegex.FindStringSubmatch(line); match != nil {
			asciiDelimiter := match[3]
			action := match[2]
			bigWithNumberCompiledRegex := regexp.MustCompile(fmt.Sprintf(bigWithNumberRegex, asciiDelimiter, action))
			if bigMatch := bigWithNumberCompiledRegex.FindStringSubmatch(line); bigMatch != nil {
				words, err := functions.ActionsModerator(bigMatch[1], action)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					continue
				}
				sliced[index] = bigWithNumberCompiledRegex.ReplaceAllString(line, words)
			} else {
				fmt.Fprintln(os.Stderr, "enter a valid words number to "+action)
			}

		}
	}
	fmt.Println(sliced)
}
