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
	Number       string     `json:"number"`
	Title        string     `json:"title"`
	Link         string     `json:"link"`
	NumTestCases int        `json:"numTestCases"`
	TestCases    []TestCase `json:"testCases"`
}

func compile(problemNumber string) error {
	fullSrcPath := "\\Users\\janwe\\Documents\\src\\codeforces\\problems\\" + problemNumber + ".cpp"
	fullBinPath := "\\Users\\janwe\\Documents\\src\\codeforces\\bin\\" + problemNumber

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

func retrieveTestCases(problemNumber string) (*Problem, error) {
	fullTestPath := "\\Users\\janwe\\Documents\\src\\codeforces\\tests\\" + problemNumber + "_test.json"

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

	return &problem, nil
}

func runTestCases(problem *Problem) error {
	fullBinPath := "\\Users\\janwe\\Documents\\src\\codeforces\\bin\\" + problem.Number

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
	if len(os.Args) < 2 {
		fmt.Println("judge: Please specify a problem to judge")
		os.Exit(1)
	}

	input := os.Args[1]

	compileFile := true
	if len(os.Args) > 2 && os.Args[2] == "-nc" {
		compileFile = false
	}

	var err error

	if compileFile {
		err = compile(input)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	problem, err := retrieveTestCases(input)
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
