package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type inputBlueprint struct {
	repoPath string
	rulePath string
}

func getParams() (inputBlueprint, error) {
	if len(os.Args) < 3 {
		return inputBlueprint{}, errors.New("blueprint repository and config file is required")
	}

	repoPath := os.Args[1]
	rulePath := os.Args[2]

	return inputBlueprint{repoPath, rulePath}, nil
}

func cloneBlueprint(url string) {
	_, err := exec.Command("git", "clone", url).Output()
	if err != nil {
		exitGracefully(err)
	}

	fmt.Println("Successfully cloned blueprint")
}

func exitGracefully(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
