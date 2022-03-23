package utils

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetStringAtLine(s string, line int) (string, error) {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	// adjust line number to 0-based index
	aLine := line - 1
	if aLine > len(lines) {
		return "", errors.New("line number out of range")
	}

	return lines[aLine], nil
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
