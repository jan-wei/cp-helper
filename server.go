package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const cpTemplate = `#include <bits/stdc++.h>

using namespace std;

int main() {
	ios::sync_with_stdio(false);
	cin.tie(NULL);
	
	
	
	return 0;
}
`

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

func enableCORS(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
    (*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST")
}

func CreateNewProblem(w http.ResponseWriter, r *http.Request) {
	// set CORS headers
    enableCORS(&w, r)
    
    if (*r).Method == "OPTIONS" {
		return
	}
	
	// parse request body
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Failed to read request body.", http.StatusBadRequest)
		return
	}
	
	// parse data into struct
	var problem Problem
	err = json.Unmarshal(reqBody, &problem)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Failed to parse request body.", http.StatusBadRequest)
		return
	}
	
	// create a new test case file to hold the test cases
	const cpTestFilesDir = "/home/janwei97/cp-files/test/"
	testFileName := problem.Number + "_test.json"
	testFilePath := cpTestFilesDir + testFileName
	newTestFile, err := os.Create(testFilePath)
	if err != nil {
		log.Printf("Error creating test file: %v", err)
		http.Error(w, "Failed to create test file.", http.StatusInternalServerError)
		return
	}
	
	defer newTestFile.Close()
	
	// write the request body to the new test file
	_, err = newTestFile.Write(reqBody)
	if err != nil {
		log.Printf("Error writing to test file: %v", err)
		http.Error(w, "Failed to write to test file.", http.StatusInternalServerError)
		return
	}
	
	// create the problem's template
	titleComment := "// Problem Title : " + problem.Title
	linkComment := "// Link : " + problem.Link
	authorComment := "// Author : janwei25"
	
	now := time.Now()
	formattedTime := fmt.Sprintf("%02d-%02d-%d %02d:%02d:%02d", now.Month(), now.Day(), now.Year(), now.Hour(), now.Minute(), now.Second())
	
	dateComment := "// Date : " + formattedTime
	
	comments := titleComment + "\n" + linkComment + "\n" + authorComment + "\n" + dateComment + "\n\n"
	
	// add comment to cpTemplate
	fullTemplate := comments + cpTemplate
	
	// create a new source file to write the template to
	// this file is where the problem would be coded in
	const cpSrcFilesDir = "/home/janwei97/cp-files/src/"
	srcFileName := problem.Number + ".cpp"
	srcFilePath := cpSrcFilesDir + srcFileName
	newSrcFile, err := os.Create(srcFilePath)
	if err != nil {
		log.Printf("Error creating source file: %v", err)
		http.Error(w, "Failed to create source file.", http.StatusInternalServerError)
		return
	}
	
	defer newSrcFile.Close()
	
	// write the template to the new source file
	_, err = newSrcFile.Write([]byte(fullTemplate))
	if err != nil {
		log.Printf("Error writing to source file: %v", err)
		http.Error(w, "Failed to write to source file.", http.StatusInternalServerError)
		return
	}
	
	// invoke new src file with geany
	command := exec.Command("geany", srcFilePath)
	err = command.Run()
	if err != nil {
		log.Printf("Error opening source file with Geany: %v", err)
		http.Error(w, "Failed to open source file with Geany.", http.StatusInternalServerError)
		return
	}
	
	// request successful
	w.WriteHeader(http.StatusOK)
	return
}

func main() {
	http.HandleFunc("/problem/new", CreateNewProblem)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
