package linter

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/PrinceMerluza/devcenter-content-linter/config"
)

type ValidationData struct {
	RuleSetName string
	Description string
	ContentPath string
	RuleData    *config.RuleSet
}

type ValidationResult struct {
	SuccessResults *[]RuleResult `json:"success"`
	FailureResults *[]RuleResult `json:"failed"`
}

type RuleResult struct {
	Id             string           `json:"id"`
	Level          config.Level     `json:"level"`
	Description    string           `json:"description"`
	IsSuccess      bool             `json:"-"`
	FileHighlights *[]FileHighlight `json:"fileHighlights,omitempty"`
	Error          *ValidationError `json:"error,omitempty"`
}

type ConditionResult struct {
	IsSuccess      bool
	FileHighlights *[]FileHighlight
	Error          error
}

type FileHighlight struct {
	Path        string `json:"path"`
	LineNumber  int    `json:"lineNumber"`
	LineContent string `json:"lineContent"`
}

type ValidationError struct {
	RuleId string
	Err    error
}

type Validator interface {
	Validate() *ConditionResult
}

// Validate the content
func (input *ValidationData) Validate() (*ValidationResult, error) {
	if input == nil {
		return nil, errors.New("nil validation data")
	}

	rulesCount := 0
	contentPath := input.ContentPath
	ruleData := input.RuleData
	finalResult := &ValidationResult{
		SuccessResults: &[]RuleResult{},
		FailureResults: &[]RuleResult{},
	}
	ch := make(chan *RuleResult)

	if _, err := os.Stat(contentPath); os.IsNotExist(err) {
		return nil, err
	}

	for id, ruleGroup := range *ruleData.RuleGroups {
		rulesCount += len(*ruleGroup.Rules)
		if err := validateRuleGroup(&ruleGroup, ch, id, contentPath); err != nil {
			return finalResult, err
		}
	}

	for i := 0; i < rulesCount; i++ {
		ruleResult := <-ch

		if ruleResult.IsSuccess {
			*finalResult.SuccessResults = append(*finalResult.SuccessResults, *ruleResult)
			continue
		}
		if !ruleResult.IsSuccess {
			*finalResult.FailureResults = append(*finalResult.FailureResults, *ruleResult)
			continue
		}
	}

	return finalResult, nil
}

// Evaluate the rulegroup. Channel should be passed where the RuleResults will
// be sent to.
func validateRuleGroup(ruleGroup *config.RuleGroup, ch chan *RuleResult, groupId string, path string) error {
	if ch == nil {
		return fmt.Errorf("%s: channel is missing", groupId)
	}

	if len(groupId) <= 0 {
		return fmt.Errorf("%s: group id is blank", groupId)
	}

	if len(path) <= 0 {
		return fmt.Errorf("%s: path is blank", groupId)
	}

	for id, rule := range *ruleGroup.Rules {
		ruleIdFull := fmt.Sprintf("%s_%v", groupId, id)
		ruleCpy := rule

		go func() {
			ch <- validateRule(&ruleCpy, ruleIdFull, path)
		}()
	}

	return nil
}

// Evaluate the specific rule and get the RuleResult. Path is the root of
// content files
func validateRule(rule *config.Rule, ruleId string, contentPath string) *RuleResult {
	ret := &RuleResult{
		Id:          ruleId,
		Level:       rule.Level,
		Description: rule.Description,
	}

	var filePath *string
	if rule.File != nil {
		tmpPath := path.Join(contentPath, *rule.File)
		filePath = &tmpPath
	}

	for _, condition := range *rule.Conditions {
		condResult := validateCondition(&condition, filePath)
		if condResult == nil {
			ret.Error = &ValidationError{
				RuleId: ruleId,
				Err:    errors.New("unexpected error. No result from condition"),
			}
			break
		}

		ret.IsSuccess = condResult.IsSuccess
		ret.FileHighlights = condResult.FileHighlights

		if condResult.Error != nil {
			ret.Error = &ValidationError{
				RuleId: ruleId,
				Err:    condResult.Error,
			}
		}

		// Short circuit failing conditions
		if !ret.IsSuccess {
			break
		}
	}

	return ret
}

// Evaluate the condition. Any failure in any type of condition will short circuit the evaluation.
func validateCondition(condition *config.Condition, filePath *string) *ConditionResult {
	var ret *ConditionResult
	var validator Validator

	// PathExists Condition
	if condition.PathExists != nil {
		validator = &PathExistsCondition{
			Path: *condition.PathExists,
		}
	}

	// Contains Conditions
	if condition.Contains != nil {
		validator = &ContainsCondition{
			Path:        *filePath,
			ContainsArr: condition.Contains,
		}
	}

	// Not Contains Condition
	if condition.NotContains != nil {
		validator = &NotContainsCondition{
			Path:        *filePath,
			NotContains: condition.NotContains,
		}
	}

	// Check reference Exist Condition
	if condition.CheckReferenceExist != nil {
		validator = &RefExistsCondition{
			Path:              *filePath,
			ReferencePatterns: condition.CheckReferenceExist,
		}
	}

	ret = validator.Validate()
	return ret
}
