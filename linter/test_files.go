package linter

import (
	"path/filepath"

	"github.com/PrinceMerluza/devcenter-content-linter/logger"
)

var (
	emptyFile       string = "./test/empty.md"
	containsFile    string = "./test/contains.md"
	notContainsFile string = "./test/notcontains.md"
	refExists       string = "./test/refExists.md"
	refExists2      string = "./test/refExists2.md"
	incorrectPath   string = "./aasifGJASDIOOJ123LKRJAWSLIEUWE/qadGHQAWIUEHAWE"
)

func relPath(path string) string {
	relPath, err := filepath.Rel(".", path)
	if err != nil {
		logger.Fatal(err)
	}

	return relPath
}
