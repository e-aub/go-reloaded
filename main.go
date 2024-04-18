package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

//Declaring variables

var (
	inputFileName  string
	outputFileName string
	input          []byte
	output         []byte
	err            error
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
		fmt.Fprintf(os.Stderr, "Error while writing content in the file\nError: %v\n", err)
		return
	}
	output = input
	//Write output
	err = ioutil.WriteFile(outputFileName, output, 07)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while writing content in the file\nError: %v\n", err)
		return
	}
}
