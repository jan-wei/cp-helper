package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type TestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type Problem struct {
	Platform     string     `json:"platform"`
	Number       string     `json:"number"`
	Title        string     `json:"title"`
	Link         string     `json:"link"`
	NumTestCases int        `json:"numTestCases"`
	TestCases    []TestCase `json:"testCases"`
}

func compile(platform, problemNumber string) error {
	fullSrcPath := fmt.Sprintf("\\Users\\janwe\\Documents\\src\\%s\\problems\\%s.cpp", platform, problemNumber)
	fullBinPath := fmt.Sprintf("\\Users\\janwe\\Documents\\src\\%s\\bin\\%s", platform, problemNumber)

	args := []string{"-std=c++17", fullSrcPath, "-o", fullBinPath}

	cmd := exec.Command("g++", args...)

	fmt.Printf("Compiling %s...\n", fullSrcPath)
	// compile the source file
	compilerOutput, err := cmd.CombinedOutput()

	// dump compiler output (if any)
	outputString := string(compilerOutput)
	if len(outputString) > 0 {
		fmt.Println(string(compilerOutput))
	}

	if err != nil {
		return errors.New("judge: Compile error")
	}

	fmt.Printf("Compile success!\n\n")

	return nil
}

func retrieveTestCases(platform, problemNumber string) (*Problem, error) {
	fullTestPath := fmt.Sprintf("\\Users\\janwe\\Documents\\src\\%s\\tests\\%s_test.json", platform, problemNumber)

	fmt.Printf("Retrieving test cases from %s...\n", fullTestPath)

	testFile, err := os.Open(fullTestPath)
	if err != nil {
		return nil, err
	}

	defer testFile.Close()

	// get the file size (needed in order to allocate byte array)
	fileStat, err := testFile.Stat()
	if err != nil {
		return nil, err
	}

	// read the file into a byte array
	data := make([]byte, fileStat.Size())
	_, err = testFile.Read(data)
	if err != nil {
		return nil, err
	}

	// unmarshal it into a go struct
	var problem Problem
	err = json.Unmarshal(data, &problem)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Test cases retrieved!\n\n")

	// populate platform
	problem.Platform = platform

	return &problem, nil
}

func runTestCases(problem *Problem) error {
	fullBinPath := fmt.Sprintf("\\Users\\janwe\\Documents\\src\\%s\\bin\\%s", problem.Platform, problem.Number)

	fmt.Println("Running test cases...\n")

	for i := 0; i < problem.NumTestCases; i++ {
		fmt.Printf("TEST #%d\n", i+1)
		fmt.Println("Input:")
		fmt.Printf("%s\n", problem.TestCases[i].Input)

		cmd := exec.Command(fullBinPath)

		// pipe test case input into the command
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}

		// write test case input to pipe
		_, err = io.WriteString(stdin, problem.TestCases[i].Input)
		if err != nil {
			return err
		}

		// execute the program
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(output))
			return errors.New("judge: Runtime error")
		}

		stdin.Close()

		// print output and expected answer
		fmt.Println("Output:")
		fmt.Printf("%s\n", string(output))
		fmt.Println("Answer:")
		fmt.Printf("%s\n", problem.TestCases[i].Output)
	}

	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("judge: Please specify the platform and problem to judge")
		os.Exit(1)
	}

	platform := os.Args[1]
	input := os.Args[2]

	if platform != "codeforces" && platform != "atcoder" {
		fmt.Println("judge: Platform not supported")
		os.Exit(1)
	}

	compileFile := true
	if len(os.Args) > 3 && os.Args[3] == "-nc" {
		compileFile = false
	}

	var err error

	if compileFile {
		err = compile(platform, input)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	problem, err := retrieveTestCases(platform, input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = runTestCases(problem)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
