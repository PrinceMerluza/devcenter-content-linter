package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
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

func evaluate() error {
	// Get CLI parameter values
	paramsData, err := getParams()
	if err != nil {
		return err
	}

	// Create temporary directory for blueprints
	tempDir, err := ioutil.TempDir(".", ".tmp-content-*")
	if err != nil {
		log.Print("Can't create temporary directory")
		return err
	}
	defer os.RemoveAll(tempDir)

	// Clone blueprint and load rule config
	data, errs := prepareFiles(&paramsData, tempDir)
	for _, err := range errs {
		log.Print(err.Error())
	}
	if len(errs) > 0 {
		return errors.New("Error when preparing necessary files")
	}

	// Evaluate the content
	finalResult, err := data.Evaluate()
	if err != nil {
		return err
	}

	sort.SliceStable(finalResult.results, func(i, j int) bool {
		return finalResult.results[i].Id < finalResult.results[j].Id
	})

	for _, result := range finalResult.results {
		fmt.Printf("\n-----\n")

		if result.Error != nil || result.IsSuccess == nil {
			fmt.Printf("%s \n Error: %v \n", result.Id, result.Error)
			continue
		}

		fmt.Printf("%s \nLevel: %s \nDescription: %s \nSuccess: %v",
			result.Id, result.Rule.Level, result.Rule.Description, *result.IsSuccess)

		if result.FileHighlights != nil {
			for _, fileHighlight := range result.FileHighlights {
				fmt.Printf("\nFile: %v \nLine #%v \n%v", fileHighlight.Path, fileHighlight.LineNumber, fileHighlight.LineContent)
			}
		}
	}

	return err
}

func main() {
	if err := evaluate(); err != nil {
		log.Fatal(err)
	}
}
