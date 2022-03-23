package linter

import (
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/PrinceMerluza/devcenter-content-linter/config"
	"github.com/PrinceMerluza/devcenter-content-linter/utils"
)

type ContainsCondition struct {
	Path        string
	ContainsArr *[]config.ContainsCondition
}

func (condition *ContainsCondition) Validate() *ConditionResult {
	ret := &ConditionResult{
		FileHighlights: &[]FileHighlight{},
		IsSuccess:      true,
	}

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
			// Get the index of matchign string
			index := strings.Index(dataString, contains.Value)
			if index < 0 {
				ret.IsSuccess = false
				break
			}

			lineNumber := strings.Count(dataString[:index], "\n") + 1

			lineContent, err := utils.GetStringAtLine(dataString, lineNumber)
			if err != nil {
				ret.Error = err
				ret.IsSuccess = false
				break
			}

			*ret.FileHighlights = append(*ret.FileHighlights, FileHighlight{
				Path:        condition.Path,
				LineNumber:  lineNumber,
				LineContent: lineContent,
				LineCount:   1,
			})
		case "regex":
			re, err := regexp.Compile(contains.Value)
			if err != nil {
				ret.Error = err
				ret.IsSuccess = false
				break
			}

			loc := re.FindStringIndex(dataString)
			if loc == nil {
				ret.IsSuccess = false
				break
			}

			match := dataString[loc[0]:loc[1]]
			lineIndex := strings.Count(dataString[:loc[0]], "\n") + 1
			lineCount := strings.Count(dataString[loc[0]:loc[1]], "\n")

			*ret.FileHighlights = append(*ret.FileHighlights, FileHighlight{
				Path:        condition.Path,
				LineNumber:  lineIndex,
				LineContent: match,
				LineCount:   lineCount,
			})
		default:
			ret.Error = errors.New("unknown contains type")
			ret.IsSuccess = false
		}
	}

	return ret
}
