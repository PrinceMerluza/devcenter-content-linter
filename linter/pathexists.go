package linter

import (
	"os"

	"github.com/PrinceMerluza/devcenter-content-linter/logger"
)

type PathExistsCondition struct {
	Path string
}

func (condition *PathExistsCondition) Validate() *ConditionResult {
	ret := &ConditionResult{}
	ret.IsSuccess = true

	if _, err := os.Stat(condition.Path); err != nil {
		logger.Tracef("Path %s does not exist \n", condition.Path)
		ret.IsSuccess = false
	}

	return ret
}
