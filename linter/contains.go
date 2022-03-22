package linter

import (
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/PrinceMerluza/devcenter-content-linter/config"
)

type ContainsCondition struct {
	Path        string
	ContainsArr *[]config.ContainsCondition
}

func (condition *ContainsCondition) Validate() *ConditionResult {
	ret := &ConditionResult{
		FileHighlights: &[]FileHighlight{},
	}
	ret.IsSuccess = true

	fileData, err := os.ReadFile(condition.Path)
	if err != nil {
		ret.Error = err
		ret.IsSuccess = false
		return ret
	}

	dataString := string(fileData[:])

	for _, contains := range *condition.ContainsArr {
		switch contains.Type {
		case "static":
			if !strings.Contains(dataString, contains.Value) {
				ret.IsSuccess = false
			}
		case "regex":
			re, err := regexp.Compile(contains.Value)
			if err != nil {
				ret.Error = err
				ret.IsSuccess = false
				return ret
			}

			matched := re.MatchString(dataString)
			if !matched {
				ret.IsSuccess = false
			}
		default:
			ret.Error = errors.New("unknown contains type")
			ret.IsSuccess = false
		}
	}

	return ret
}
