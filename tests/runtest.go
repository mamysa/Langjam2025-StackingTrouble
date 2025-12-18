package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var TESTS_PASSING []string = []string{
	"assert-fail.bla",
	"arith-0.bla",
	"cmp-0.bla",
	"boolean-short-circuit.bla",
	"object-0.bla",
}

const (
	RED   string = "\033[1;31m"
	GREEN        = "\033[0;32m"
	RESET        = "\033[0;0m"
)

func runInterpreter(filename string) []string {
	stdoutBuf := new(bytes.Buffer)
	cmd := exec.Command("go", "run", "../main.go", filename)
	cmd.Stdout = stdoutBuf

	if err := cmd.Run(); err != nil {
		//fmt.Printf("Error running %s", err)
	}

	outputTrimmed := strings.TrimSpace(stdoutBuf.String())

	outSplit := strings.Split(outputTrimmed, "\n")
	return outSplit
}

func extractTestMessages(filename string) []string {
	f, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	fsplit := strings.Split(string(f), "\n")

	stdoutPrefix := "#? STDOUT"

	expectedOutputs := []string{}
	for _, l := range fsplit {
		lineTrimmed := strings.TrimSpace(l)
		if strings.HasPrefix(lineTrimmed, stdoutPrefix) {

			lineWithoutPrefix := strings.TrimPrefix(lineTrimmed, stdoutPrefix)
			lineWithoutPrefix = strings.TrimSpace(lineWithoutPrefix)
			expectedOutputs = append(expectedOutputs, lineWithoutPrefix)
		}
	}
	return expectedOutputs
}

func compare(filename string, expected, actual []string) {
	le := len(expected)
	la := len(actual)

	if le != la {
		fmt.Printf("%s[FAIL] %s %s: output size mismatch, expected %d, actual:%d\n", RED, filename, RESET, le, la)
		fmt.Printf("e: %+v\na: %+v\n", expected, actual)
		return
	}

	for i := 0; i < le; i++ {
		a := expected[i]
		b := actual[i]

		if a != b {
			fmt.Printf("%s[FAIL] %s %s: output mismatch at index %d\n", RED, filename, RESET, i)
			fmt.Printf("e: %+v\na: %+v\n", expected, actual)
			return
		}
	}

	fmt.Printf("%s[PASS] %s %s\n", GREEN, filename, RESET)
}

func main() {
	for _, passingTest := range TESTS_PASSING {
		expectedOutputs := extractTestMessages(passingTest)
		actualOutputs := runInterpreter(passingTest)
		compare(passingTest, expectedOutputs, actualOutputs)
	}
}
