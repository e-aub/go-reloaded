package main

import (
	"fmt"
	"goreloaded/functions"
	"os"
	"regexp"
)

func main() {
	text := "hello world (low, 2)  how are! you bro (cap, 5)"
	// actionAndDelimiterRegex := `\((low|up|cap), (\d+)\)`
	splitRegex := `\((low|up|cap), (\d+)\)|\(cap\)|\(up\)|\(low\)|\(hex\)|\(bin\)`
	withoutNumberRegex := regexp.MustCompile(`((?:\b[\w_-]+\b[\W\s]*){1})\((low|up|cap|hex|bin)\)`)
	withNumberRegex := regexp.MustCompile(`\(((low|up|cap)), (\d+)\)`)
	checkWordsBeforeRegex := `((?:\b[\w]+\b\W*\s*){%s})\(((?:%s),\s*\d+)\)`
	sliced := functions.SplitKeepSeparator(text, splitRegex)

	fmt.Println(sliced)

	for index, line := range sliced {
		if match := withoutNumberRegex.FindStringSubmatch(line); match != nil {
			word := match[1]
			action := match[2]
			word, err := functions.ActionsModerator(word, action)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			sliced[index] = withoutNumberRegex.ReplaceAllString(sliced[index], "|"+word+"|")
		} else if match := withNumberRegex.FindStringSubmatch(line); match != nil {
			asciiDelimiter := match[3]
			action := match[2]
			checkWordsBeforeCompiledRegex := regexp.MustCompile(fmt.Sprintf(checkWordsBeforeRegex, asciiDelimiter, action))
			if bigMatch := checkWordsBeforeCompiledRegex.FindStringSubmatch(line); bigMatch != nil {
				words, err := functions.ActionsModerator(bigMatch[1], action)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					continue
				}
				sliced[index] = checkWordsBeforeCompiledRegex.ReplaceAllString(line, words)
			} else {
				fmt.Fprintln(os.Stderr, "enter a valid words number to "+action)
			}

		}
	}
	fmt.Println(sliced)
}
