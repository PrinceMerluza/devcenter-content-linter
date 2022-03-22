package linter

import (
	"os"
)

type PathExistsCondition struct {
	Path string
}

func (condition *PathExistsCondition) Validate() *ConditionResult {
	ret := &ConditionResult{}
	ret.IsSuccess = true

	if _, err := os.Stat(condition.Path); os.IsNotExist(err) {
		ret.IsSuccess = false
	}

	return ret
}
