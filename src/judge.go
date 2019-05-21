package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type TestCase struct {
	Input string 	`json:"input"`	
	Output string 	`json:"output"`
}

type Problem struct {
	Number string 			`json:"number"`
	Title string 			`json:"title"`
	Link string				`json:"link"`
	NumTestCases int 		`json:"numTestCases"`
	TestCases []TestCase 	`json:"testCases"`
}

func compile(problemNumber string) error {
	fullSrcPath := "/home/janwei97/cp-files/src/" + problemNumber + ".cpp"
	fullBinPath := "/home/janwei97/cp-files/bin/" + problemNumber
	
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
	fullTestPath := "/home/janwei97/cp-files/test/" + problemNumber + "_test.json"
	
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
	fullBinPath := "/home/janwei97/cp-files/bin/" + problem.Number
	numTestsPassed := 0
	
	fmt.Println("Running test cases...\n")
	
	for i := 0; i < problem.NumTestCases; i++ {
		fmt.Printf("\033[96mTEST #%d\033[0m\n", i + 1)
		fmt.Println("\033[93mInput:\033[0m")
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
		
		// trim trailing '\n' characters for both output and expected answer
		outputString := string(output)
		outputString = strings.TrimSuffix(outputString, "\n")
		problem.TestCases[i].Output = strings.TrimSuffix(problem.TestCases[i].Output, "\n")
		
		// print output and expected answer
		fmt.Println("\033[93mOutput:\033[0m")
		fmt.Printf("%s\n\n", outputString)
		fmt.Println("\033[93mAnswer:\033[0m")
		fmt.Printf("%s\n\n", problem.TestCases[i].Output)
		
		fmt.Printf("\033[95mVerdict: \033[0m")
		// perform checks between output and expected answer
		if(outputString == problem.TestCases[i].Output) {
			// correct
			numTestsPassed++
			fmt.Println("\033[92mCorrect\033[0m")
		} else {
			fmt.Println("\033[91mWrong Answer\033[0m")
		}
		
		fmt.Printf("\n")
	}
	
	// print the number of test cases passed
	fmt.Printf("%d/%d test cases passed\n\n", numTestsPassed, problem.NumTestCases)
	
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
