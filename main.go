package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var logger = log.New(os.Stdout, "", 0)

type paramBlueprint struct {
	repoPath string
	rulePath string
}

func getParams() (paramBlueprint, error) {
	if len(os.Args) < 3 {
		return paramBlueprint{}, errors.New("blueprint repository and config file is required")
	}

	repoPath := os.Args[1]
	rulePath := os.Args[2]

	return paramBlueprint{repoPath, rulePath}, nil
}

func prepareFiles(paramsData *paramBlueprint, tempDir string) (r *EvaluationData, errors []error) {
	var wg sync.WaitGroup
	var err error
	r = &EvaluationData{}

	wg.Add(2)
	go func() {
		defer wg.Done()

		if r.ContentPath, err = cloneBlueprint(tempDir, paramsData.repoPath); err != nil {
			errors = append(errors, err)
		}
	}()

	go func() {
		defer wg.Done()

		if r.RuleData, err = loadRuleConfig(paramsData.rulePath); err != nil {
			errors = append(errors, err)
		}
	}()

	wg.Wait()

	return r, errors
}

func main() {
	// Get CLI parameter values
	paramsData, err := getParams()
	if err != nil {
		log.Fatal(err)
	}

	// Create temporary directory for blueprints
	tempDir, err := ioutil.TempDir(".", ".tmp-content-*")
	if err != nil {
		log.Fatal("Can't create temporary directory")
	}
	defer os.RemoveAll(tempDir)

	// Clone blueprint and load rule config
	data, errs := prepareFiles(&paramsData, tempDir)
	for _, err := range errs {
		log.Print(err.Error())
	}
	if len(errs) > 0 {
		log.Fatal("Error when preparing necessary files")
	}

	// Evaluate the content
	finalResult, err := data.Evaluate()
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, result := range finalResult.results {
		if result.Error != nil || result.IsSuccess == nil {
			fmt.Printf(`Error on running test %s
				Error: %v
			`, result.Id, result.Error)
			continue
		}

		fmt.Printf(`%s
			Level: %s
			Description: %s
			Success: %v
		`, result.Id, result.Rule.Level, result.Rule.Description, *result.IsSuccess)
	}

	logger.Println("END")
}
