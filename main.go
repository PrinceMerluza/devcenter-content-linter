package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

var logger = log.New(os.Stdout, "", 0)

type paramBlueprint struct {
	repoPath string
	rulePath string
	fileDir  string
}

func getParams() (paramBlueprint, error) {
	if len(os.Args) < 3 {
		return paramBlueprint{}, errors.New("blueprint repository and config file is required")
	}

	repoPath := os.Args[1]
	rulePath := os.Args[2]
	fileDir := "."

	if len(os.Args) > 3 {
		fileDir = os.Args[3]
	}

	return paramBlueprint{repoPath, rulePath, fileDir}, nil
}

func prepareFiles(paramsData *paramBlueprint) (r *EvaluationData, errors []error) {
	var wg sync.WaitGroup
	var err error
	r = &EvaluationData{}

	wg.Add(2)
	go func() {
		defer wg.Done()

		// For local copy of content repo
		if !IsURL(paramsData.repoPath) {
			r.ContentPath = paramsData.repoPath
			return
		}

		// For remote repo, clone the blueprint
		tempDir, err := createTempDir(paramsData.fileDir)
		if err != nil {
			errors = append(errors, err)
		}
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

func evaluate(data *EvaluationData) (*EvaluationResult, error) {
	// Evaluate the content
	finalResult, err := data.Evaluate()
	if err != nil {
		return nil, err
	}

	// Sort the slices
	sort.SliceStable(finalResult.SuccessResults, func(i, j int) bool {
		return finalResult.SuccessResults[i].Id < finalResult.SuccessResults[j].Id
	})
	sort.SliceStable(finalResult.FailureResults, func(i, j int) bool {
		return finalResult.FailureResults[i].Id < finalResult.FailureResults[j].Id
	})
	sort.SliceStable(finalResult.ErrorResults, func(i, j int) bool {
		return finalResult.ErrorResults[i].Id < finalResult.ErrorResults[j].Id
	})

	return finalResult, err
}

func printResults(finalResult *EvaluationResult) {
	for _, result := range finalResult.SuccessResults {
		fmt.Printf("\n--- SUCCESS --\n")

		fmt.Printf("%s \nLevel: %s \nDescription: %s",
			result.Id, result.Level, result.Description)

		if result.FileHighlights != nil {
			for _, fileHighlight := range result.FileHighlights {
				fmt.Printf("\nFile: %v \nLine #%v \n%v", fileHighlight.Path, fileHighlight.LineNumber, fileHighlight.LineContent)
			}
		}
		fmt.Println("--")
	}

	for _, result := range finalResult.FailureResults {
		fmt.Printf("\n--- FAILED --\n")

		fmt.Printf("%s \nLevel: %s \nDescription: %s",
			result.Id, result.Level, result.Description)

		if result.FileHighlights != nil {
			for _, fileHighlight := range result.FileHighlights {
				fmt.Printf("\nFile: %v \nLine #%v \n%v", fileHighlight.Path, fileHighlight.LineNumber, fileHighlight.LineContent)
			}
		}
		fmt.Println("--")
	}

	for _, result := range finalResult.ErrorResults {
		fmt.Printf("\n--- ERROR --\n")

		fmt.Printf("%s \nLevel: %s \nDescription: %s",
			result.Id, result.Level, result.Description)

		if result.FileHighlights != nil {
			for _, fileHighlight := range result.FileHighlights {
				fmt.Printf("\nFile: %v \nLine #%v \n%v", fileHighlight.Path, fileHighlight.LineNumber, fileHighlight.LineContent)
			}
		}
		fmt.Println("--")
	}
}

func exportJsonResult(finalResult *EvaluationResult, filename string) error {
	data, err := json.Marshal(finalResult)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func createTempDir(root string) (string, error) {
	// Create temporary directory for blueprints
	tempDir, err := ioutil.TempDir(root, ".tmp-content-*")
	if err != nil {
		log.Print("Can't create temporary directory")
		return "", err
	}

	return tempDir, nil
}

func main() {
	// Get CLI parameter values
	paramsData, err := getParams()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare the repo and the rule file
	data, errs := prepareFiles(&paramsData)
	for _, err := range errs {
		log.Print(err.Error())
	}
	if len(errs) > 0 {
		log.Fatal("Error when preparing necessary files")
	}

	// Evalue the repository
	finalResult, err := evaluate(data)
	if err != nil {
		log.Fatal(err)
	}

	// Output
	printResults(finalResult)
	err = exportJsonResult(finalResult, filepath.Join(paramsData.fileDir, "result.json"))
	if err != nil {
		log.Fatal(err)
	}
}
