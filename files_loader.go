package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/qri-io/jsonschema"
)

// Clone the blueprint to a given directory
// Returns relative directory path (to the path param) to the repo contents
func cloneBlueprint(tmpDirectory string, url string) (dirPath string, err error) {
	logger.Print("Cloning blueprint...")

	// Clone the blueprint into the temporary directory
	_, err = exec.Command("git", "-C", tmpDirectory, "clone", url).Output()
	if err != nil {
		return
	}

	files, err := os.ReadDir(tmpDirectory)
	if err != nil {
		return
	}

	if len(files) < 1 {
		err = errors.New("can't find cloned repo directory")
		return
	}

	logger.Println("Successfully cloned blueprint")
	dirPath = filepath.Join(tmpDirectory, files[0].Name())

	return
}

// Loads the rule json file in the given path
// Returns the unmarshalled RuleData
func loadRuleConfig(filePath string) (retData *RuleData, err error) {
	ctx := context.Background()
	rs := &jsonschema.Schema{}
	retData = &RuleData{}

	logger.Print("Processing rule configuration...")

	if err = json.Unmarshal([]byte(RuleSetSchema), rs); err != nil {
		return
	}

	// Load the rule config file
	rawRules, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	// Verify if rule config file follows schema
	errs, err := rs.ValidateBytes(ctx, rawRules)
	if err != nil {
		return
	}
	for _, err = range errs {
		log.Print(err.Error())
	}
	if len(errs) > 0 {
		err = errors.New("rule file syntax is invalid")
		return
	}

	// Marshall the rule json file
	if err = json.Unmarshal(rawRules, retData); err != nil {
		return
	}

	logger.Print("Successfully processed rule configuration")
	return
}
