package linter

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type RefExistsCondition struct {
	Path              string
	ReferencePatterns *[]string
}

func (condition *RefExistsCondition) Validate() *ConditionResult {
	ret := &ConditionResult{
		FileHighlights: &[]FileHighlight{},
	}
	ret.IsSuccess = true

	file, err := os.Open(condition.Path)
	if err != nil {
		ret.IsSuccess = false
		ret.Error = err
		return ret
	}
	defer file.Close()

	for _, pattern := range *condition.ReferencePatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			ret.Error = err
			ret.IsSuccess = false
			return ret
		}

		scanner := bufio.NewScanner(file)
		lineNumber := 0
		for scanner.Scan() {
			lineNumber++
			lineString := scanner.Text()

			subMatch := re.FindStringSubmatch(lineString)
			if subMatch == nil {
				continue
			}

			if len(subMatch) <= 1 {
				ret.Error = errors.New("no matching group found. Regex may be incorrect")
				ret.IsSuccess = false
				return ret
			}

			// NOTE: The second submatch(1st matching group) is always used to get the path
			// condition.Path is always a file so need to get the directory, before adding relative path
			pathToCheck := filepath.Join(condition.Path, "..", subMatch[1])

			fmt.Println(pathToCheck)

			if _, err := os.Stat(pathToCheck); err != nil {
				ret.IsSuccess = false
			}
			*ret.FileHighlights = append(*ret.FileHighlights, FileHighlight{
				Path:        condition.Path,
				LineNumber:  lineNumber,
				LineContent: lineString,
				LineCount:   1,
			})
		}

		if err := scanner.Err(); err != nil {
			ret.Error = err
			ret.IsSuccess = false
			return ret
		}
	}

	return ret
}
