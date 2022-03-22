package linter

import (
	"bufio"
	"os"
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

			subMatch := re.FindSubmatch([]byte(lineString))
			if subMatch == nil {
				continue
			}

			// NOTE: The second submatch(1st matching group) is always used to get the path
			pathToCheck := string(subMatch[1])

			if _, err := os.Stat(pathToCheck); os.IsNotExist(err) {
				ret.IsSuccess = false
			}
			*ret.FileHighlights = append(*ret.FileHighlights, FileHighlight{
				Path:        condition.Path,
				LineNumber:  lineNumber,
				LineContent: lineString,
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
