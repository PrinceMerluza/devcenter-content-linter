package main

import "strings"

func NewBoolPtr(val bool) *bool {
	return &val
}

func IsURL(path string) bool {
	return strings.HasPrefix(strings.ToLower(path), "http")
}
