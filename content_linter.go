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

// type ruleGroup struct {
// 	id string
// 	description string
// 	rules rules
// }

// type rules struct {

// }

func cloneBlueprint(tmpDirectory string, url string) (dirPath string, err error) {
	// Clone the blueprint into the temporary directory
	_, err = exec.Command("git", "-C", tmpDirectory, "clone", url).Output()
	if err != nil {
		return "", err
	}

	files, err := os.ReadDir(tmpDirectory)
	if err != nil {
		return "", err
	}

	if len(files) < 1 {
		return "", errors.New("can't find cloned repo directory")
	}

	logger.Println("Successfully cloned blueprint")

	return filepath.Join(tmpDirectory, files[0].Name()), nil
}

func loadRuleConfig(filePath string) (err error) {
	ctx := context.Background()
	rs := &jsonschema.Schema{}

	// Load schema file
	schemaData, err := os.ReadFile("./schemas/linter-rules.schema.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(schemaData, rs); err != nil {
		return err
	}

	// Load the rule config file
	ruleData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	errs, err := rs.ValidateBytes(ctx, ruleData)
	if err != nil {
		return err
	}

	// Log schema validation errors
	for _, err = range errs {
		log.Print(err.Error())
	}
	if len(errs) > 0 {
		return errors.New("rule file syntax is invalid")
	}

	return nil
}
