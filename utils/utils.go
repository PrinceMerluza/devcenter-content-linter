package utils

import (
	"bufio"
	"errors"
	"net/url"
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
