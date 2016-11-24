package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func usage() string {
	return `Usage:
	teamcity-go-test -test <binary> [-parallelism n]

	Test names must be listed one per line on stdin.
`
}

func main() {
	testBinary := flag.String("test", "", "executable containing the tests to run")
	parallelism := flag.Int("parallelism", 1, "number of tests to execute in parallel")
	flag.Parse()

	if testBinary == nil || *testBinary == "" {
		fmt.Fprint(os.Stderr, usage())
		os.Exit(1)
	}

	testNames := make([]string, 0, 0)
	stdInReader := bufio.NewReader(os.Stdin)

	for {
		line, err := stdInReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if strings.TrimSpace(line) != "" {
					testNames = append(testNames, line)
				}
				break
			}
			fmt.Fprintf(os.Stderr, "error reading stdin: %s", err)
			os.Exit(1)
		}

		if strings.TrimSpace(line) != "" {
			testNames = append(testNames, line)
		}
	}

	testQueue := make(chan string)
	messages := make(chan string)
	completed := make(chan struct{})

	for i := 0; i < *parallelism; i++ {
		go runWorker(testQueue, messages, completed, *testBinary)
	}

	go func() {
		for _, testName := range testNames {
			testQueue <- strings.TrimSpace(testName)
		}
	}()

	resultsCount := 0
	for {
		select {
		case message := <-messages:
			fmt.Printf("%s", message)
		case <-completed:
			resultsCount++
		}

		if resultsCount == len(testNames) {
			break
		}
	}
}

func runWorker(inputQueue <-chan string, messages chan<- string, done chan<- struct{}, binaryName string) {
	for {
		select {
		case testName := <-inputQueue:
			test := NewTeamCityTest(testName)
			messages <- fmt.Sprintf("%s", test.FormatStartNotice())
			runTest(test, binaryName)
			messages <- test.FormatTestOutput()
			done <- struct{}{}
		}
	}
}

func runTest(test *TeamCityTest, binaryName string) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	cmd := exec.Command(binaryName, "-test.v", "-test.run", fmt.Sprintf("^%s$", test.Name))
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	// Not sure what to do with errors here other than report them out to the runner.
	cmd.Run()

	test.ParseTestRunnerOutput(out.String(), errOut.String())
}
