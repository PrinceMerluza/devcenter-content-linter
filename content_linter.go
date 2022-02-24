package main

import (
	"errors"
	"fmt"
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
	Type  string
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
	IsSuccess     bool
	FileHighlight *FileHighlight
	Error         *EvaluationError
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
func (rule *Rule) Evaluate(ruleId string, path string) *RuleResult {
	ret := &RuleResult{
		Id: ruleId,
	}

	return ret
}
