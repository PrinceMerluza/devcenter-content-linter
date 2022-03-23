package utils

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetStringAtLine(data string, line int) (string, error) {
	lines := strings.Count(data, "\n")

	if lines < line {
		return "", errors.New("line number out of range")
	}

	lines = 0
	for _, c := range data {
		if c == '\n' {
			lines++
		}

		if lines == line {
			break
		}
	}
	lastIndex := strings.Index(data[lines:], "\n")

	return data[lines:lastIndex], nil
}

func NewBoolPtr(val bool) *bool {
	return &val
}

func IsURL(path string) bool {
	_, err := url.ParseRequestURI(path)
	return err == nil
}

func CloneRepoTemp(repoUrl string) (string, error) {
	tmpPath, err := os.MkdirTemp("", "gc-content")
	if err != nil {
		fmt.Println("Error creating temp dir:", err)
		return "", err
	}

	fmt.Println("Cloning blueprint...")

	// Clone the blueprint into the temporary directory
	_, err = exec.Command("git", "-C", tmpPath, "clone", repoUrl).Output()
	if err != nil {
		fmt.Println("Error cloning repo:", err)
		return "", err
	}

	files, err := os.ReadDir(tmpPath)
	if err != nil {
		return "", err
	}

	if len(files) < 1 {
		err = errors.New("can't find cloned repo directory")
		return "", err
	}

	fmt.Println("Successfully cloned blueprint")
	dirPath := filepath.Join(tmpPath, files[0].Name())

	return dirPath, nil
}
