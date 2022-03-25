package linter

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/PrinceMerluza/devcenter-content-linter/blueprintrepo"
)

type NotContainsCondition struct {
	Path        string
	NotContains *[]string
}

func (condition *NotContainsCondition) Validate() *ConditionResult {
	ret := &ConditionResult{
		FileHighlights: &[]FileHighlight{},
	}
	ret.IsSuccess = true

	file, err := os.Open(condition.Path)
	if err != nil {
		ret.Error = err
		ret.IsSuccess = false
		return ret
	}
	defer file.Close()

	for _, contains := range *condition.NotContains {
		if strings.TrimSpace(contains) == "" {
			ret.Error = errors.New("value of notcontains is empty")
			ret.IsSuccess = false
			break
		}

		scanner := bufio.NewScanner(file)
		lineNumber := 0
		for scanner.Scan() {
			lineNumber++
			lineString := scanner.Text()

			matched, err := regexp.MatchString(contains, lineString)
			if err != nil {
				ret.Error = err
				ret.IsSuccess = false
				return ret
			}

			if matched {
				ret.IsSuccess = false
				*ret.FileHighlights = append(*ret.FileHighlights, FileHighlight{
					Path:        blueprintrepo.GetRelPath(condition.Path),
					LineNumber:  lineNumber,
					LineContent: lineString,
					LineCount:   1,
				})
			}
		}

		if err := scanner.Err(); err != nil {
			ret.Error = err
			ret.IsSuccess = false
			return ret
		}
	}

	return ret
}
