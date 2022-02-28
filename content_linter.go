package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

type RuleLevel string

const (
	Undefined RuleLevel = ""
	Warning   RuleLevel = "warning"
	Error     RuleLevel = "error"
)

type EvaluationData struct {
	RuleSetName string
	Description string
	ContentPath string
	RuleData    *RuleData
}

type RuleData struct {
	Name        string
	Description string
	RuleGroups  *map[string]RuleGroup
}

type RuleGroup struct {
	Description string
	Rules       *map[string]Rule
}

type Rule struct {
	Description string
	Path        *string
	Files       *[]string
	Conditions  *[]Condition
	Level       RuleLevel
}

type Condition struct {
	PathExists          *bool
	Contains            *[]ContainsCondition
	NotContains         *[]string
	MarkdownMeta        *map[string]string
	CheckReferenceExist *CheckReferenceExistCondition
}

type ContainsCondition struct {
	Type  string // static or regex
	Value string
}

type CheckReferenceExistCondition struct {
	Pattern    string
	MatchGroup int
}

type EvaluationResult struct {
	results []*RuleResult
}

type RuleResult struct {
	Id            string
	Rule          *Rule
	IsSuccess     *bool
	FileHighlight *FileHighlight
	Error         *EvaluationError
}

type ConditionResult struct {
	IsSuccess     *bool
	FileHighlight *FileHighlight
	Error         error
}

type FileHighlight struct {
	Path      string
	LineStart int
	LineEnd   int
}

type EvaluationError struct {
	RuleId string
	Err    error
}

func (e *EvaluationError) Error() string {
	return fmt.Sprintf("error on rule %s: %v", e.RuleId, e.Err)
}

// Evaluate the content
func (input *EvaluationData) Evaluate() (*EvaluationResult, error) {
	if input == nil {
		return nil, errors.New("nil evaluation Data")
	}

	rulesCount := 0
	contentPath := input.ContentPath
	ruleData := input.RuleData
	finalResult := &EvaluationResult{}
	ch := make(chan *RuleResult)

	for id, ruleGroup := range *ruleData.RuleGroups {
		rulesCount += len(*ruleGroup.Rules)
		if err := ruleGroup.Evaluate(ch, id, contentPath); err != nil {
			return finalResult, err
		}
	}

	for i := 0; i < rulesCount; i++ {
		ruleResult := <-ch
		finalResult.results = append(finalResult.results, ruleResult)
	}

	return finalResult, nil
}

// Evaluate the rulegroup. Channel should be passed where the RuleResults will
// be sent to.
func (ruleGroup *RuleGroup) Evaluate(ch chan *RuleResult, groupId string, path string) error {
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
		ruleIdFull := fmt.Sprintf("%s_%s", groupId, id)
		ruleCpy := rule

		go func() {
			ch <- ruleCpy.Evaluate(ruleIdFull, path)
		}()
	}

	return nil
}

// Evaluate the specific rule and get the RuleResult. Path is the root of
// content files
func (rule *Rule) Evaluate(ruleId string, contentPath string) *RuleResult {
	ret := &RuleResult{
		Id:   ruleId,
		Rule: rule,
	}

	// Short circuited evaluation for conditions
	for _, condition := range *rule.Conditions {
		condResult := condition.Evaluate(rule, contentPath)
		if condResult == nil {
			ret.Error = &EvaluationError{
				RuleId: ruleId,
				Err:    errors.New("unexpected error. No result from condition"),
			}
			break
		}
		if condResult.Error != nil {
			ret.IsSuccess = NewBoolPtr(false)
			ret.Error = &EvaluationError{
				RuleId: ruleId,
				Err:    condResult.Error,
			}
			break
		}
		if condResult.IsSuccess == nil {
			ret.Error = &EvaluationError{
				RuleId: ruleId,
				Err:    errors.New("unexpected error. Success status not able to be determined"),
			}
			break
		}

		ret.IsSuccess = condResult.IsSuccess
		ret.FileHighlight = condResult.FileHighlight
	}

	return ret
}

// Evaluate the condition. Any failure in any type of condition will short circuit the evaluation.
func (condition *Condition) Evaluate(rule *Rule, contentPath string) *ConditionResult {
	var ret *ConditionResult

	filePaths := []string{}

	// Determine the relative filepaths
	if rule.Path == nil && rule.Files == nil {
		ret.Error = errors.New("rules has no path or files in it")
	}
	if rule.Files != nil {
		for _, file := range *rule.Files {
			filePaths = append(filePaths, path.Join(contentPath, file))
		}
	}
	if rule.Path != nil {
		filePaths = append(filePaths, path.Join(contentPath, *rule.Path))
	}

	// PathExists Condition
	if condition.PathExists != nil && *condition.PathExists {
		ret = EvaluatePathExistCondition(&filePaths)
		if ret.IsSuccess != nil && !*ret.IsSuccess {
			return ret
		}
	}

	// Contains Conditions
	if condition.Contains != nil {
		ret = EvaluateContainsCondition(&filePaths, condition.Contains)
		if ret.IsSuccess != nil && !*ret.IsSuccess {
			return ret
		}
	}

	// Not Contains Condition
	if condition.NotContains != nil {
		ret = EvaluateNotContainsCondition(&filePaths, condition.NotContains)
		if ret.IsSuccess != nil && !*ret.IsSuccess {
			return ret
		}
	}

	return ret
}

func EvaluatePathExistCondition(filePaths *[]string) *ConditionResult {
	ret := &ConditionResult{}

	for _, path := range *filePaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			ret.IsSuccess = NewBoolPtr(false)
			break
		}

		ret.IsSuccess = NewBoolPtr(true)
	}

	return ret
}

func EvaluateContainsCondition(filePaths *[]string, arrContains *[]ContainsCondition) *ConditionResult {
	ret := &ConditionResult{}

	for _, path := range *filePaths {
		for _, contains := range *arrContains {
			fileData, err := os.ReadFile(path)
			if err != nil {
				ret.Error = err
				return ret
			}

			dataString := string(fileData[:])

			switch contains.Type {
			case "static":
				ret.IsSuccess = NewBoolPtr(strings.Contains(dataString, contains.Value))
			case "regex":
				matched, err := regexp.MatchString(contains.Value, dataString)
				if err != nil {
					ret.Error = err
					return ret
				}
				ret.IsSuccess = NewBoolPtr(matched)
			default:
				ret.Error = errors.New("unknown contains type")
				return ret
			}
		}
	}

	return ret
}

func EvaluateNotContainsCondition(filePaths *[]string, notContains *[]string) *ConditionResult {
	ret := &ConditionResult{}

	for _, path := range *filePaths {
		for _, contains := range *notContains {
			fileData, err := os.ReadFile(path)
			if err != nil {
				ret.Error = err
				return ret
			}

			dataString := string(fileData[:])

			matched, err := regexp.MatchString(contains, dataString)
			if err != nil {
				ret.Error = err
				return ret
			}
			ret.IsSuccess = NewBoolPtr(!matched)
		}
	}

	return ret
}
