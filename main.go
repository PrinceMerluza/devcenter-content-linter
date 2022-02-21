package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
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

func evaluateContent(paramsData *paramBlueprint) (err error) {
	// Create temporary directory for blueprints
	tempDir, err := ioutil.TempDir(".", ".tmp-content-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// Clone blueprint
	if _, err = cloneBlueprint(tempDir, paramsData.repoPath); err != nil {
		return err
	}

	// Load the configuration file
	if err = loadRuleConfig(paramsData.rulePath); err != nil {
		return err
	}

	return nil
}

func main() {
	// Get CLI parameter values
	paramsData, err := getParams()
	if err != nil {
		log.Fatal(err)
	}

	// Evaluate the content against the rule
	if err = evaluateContent(&paramsData); err != nil {
		log.Fatal(err)
	}

	logger.Println("Success")
}
