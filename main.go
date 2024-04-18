package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

//Declaring variables

var (
	inputFileName string
	// outputFileName string
	input []byte
	// output []byte
	err       error
	lowRegex  = regexp.MustCompile(`(\(low,\s(\d+)\))`)
	upRegex   = regexp.MustCompile(`(\(up,\s(\d+)\))`)
	capRegex  = regexp.MustCompile(`(\(cap,\s(\d+)\))`)
	delimiter = 1
	match     []string
)

func main() {
	//Check for valid arguments
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Invalid arguments Usage : go run main.go <input.txt> <output.txt>")
		return
	}
	inputFileName = os.Args[1]
	// outputFileName = os.Args[2]
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
		if text[i] == "low" {
			text[i-1] = strings.ToLower(text[i-1])
		} else if lowRegex.FindStringSubmatch(text[i]) != nil {
			match = lowRegex.FindStringSubmatch(text[i])
			delimiter, _ = strconv.Atoi(match[2])
			if i-delimiter < 0 {
				continue
			}
			for j := i - 1; j+delimiter >= i; j-- {
				text[j] = strings.ToLower(text[j])
			}
			delimiter = 1
		} else if text[i] == "(up)" {
			text[i-1] = strings.ToUpper(text[i-1])
		} else if upRegex.FindStringSubmatch(text[i]) != nil {
			match = upRegex.FindStringSubmatch(text[i])
			delimiter, _ = strconv.Atoi(match[2])
			if i-delimiter < 0 {
				continue
			}
			for j := i - 1; j+delimiter >= i; j-- {
				text[j] = strings.ToUpper(text[j])
			}
			delimiter = 1
		} else if text[i] == "(cap)" {
			text[i-1] = Capitalize(text[i-1])
		} else if capRegex.FindStringSubmatch(text[i]) != nil {
			match = capRegex.FindStringSubmatch(text[i])
			delimiter, _ = strconv.Atoi(match[2])
			if i-delimiter < 0 {
				continue
			}
			for j := i - 1; j+delimiter >= i; j-- {
				text[j] = Capitalize(text[j])
			}
			delimiter = 1
		} else if text[i] == "(hex)" {
			coverted, err := strconv.ParseInt(text[i-1], 16, 0)
			if err != nil {
				panic(err)
			}
			text[i-1] = strconv.Itoa(int(coverted))
		} else if text[i] == "(bin)" {
			fmt.Println("passed")
			coverted, _ := strconv.ParseInt(text[i-1], 2, 0)
			text[i-1] = strconv.Itoa(int(coverted))
		}
	}
	fmt.Println(text)

}

func valid(s string) bool {
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

func splitString([]byte) []string {
	split := strings.Split(string(input), " ")
	for index, word := range split {
		if word == "(low," || word == "(up," || word == "(cap," && index != len(split)-1 && valid(split[index+1]) {
			split[index] = split[index] + " " + split[index+1]
			split = append(split[:index+1], split[index+2:]...)
		} else {
			continue
		}
	}
	return split
}

//Write output
// err = ioutil.WriteFile(outputFileName, output, 0777)
// if err != nil {
// 	fmt.Fprintf(os.Stderr, "Error while writing content in the file\nError: %v\n", err)
// 	return
// }

func Capitalize(s string) string {
	sRune := []rune(s)
	for index, letter := range sRune {
		if letter >= 'A' && letter <= 'Z' {
			sRune[index] += 32
		}
	}
	if sRune[0] >= 'a' && sRune[0] <= 'z' {
		sRune[0] -= 32
	}
	for i := 0; i < len(sRune)-1; i++ {
		if !(sRune[i] >= 'a' && sRune[i] <= 'z' || sRune[i] >= '0' && sRune[i] <= '9' || sRune[i] >= 'A' && sRune[i] <= 'Z') {
			if sRune[i+1] >= 'a' && sRune[i+1] <= 'z' {
				sRune[i+1] -= 32
			}
		}
	}
	return string(sRune)
}
