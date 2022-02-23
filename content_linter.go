package main

import (
	"fmt"
	"sync"
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
	Id        string
	Rule      *Rule
	isSuccess bool
}

func (input *EvaluationData) evaluate() *EvaluationResult {
	var wg sync.WaitGroup
	contentPath := input.ContentPath
	ruleData := input.RuleData
	finalResult := &EvaluationResult{}

	for id, ruleGroup := range *ruleData.RuleGroups {
		idCpy := id
		ruleGroupCpy := ruleGroup

		wg.Add(1)
		go func() {
			defer wg.Done()

			finalResult.results = append(finalResult.results, ruleGroupCpy.evaluate(idCpy, contentPath)...)
		}()
	}

	wg.Wait()

	return finalResult
}

func (ruleGroup *RuleGroup) evaluate(groupId string, path string) []*RuleResult {
	var wg sync.WaitGroup
	results := []*RuleResult{}

	for id, rule := range *ruleGroup.Rules {
		ruleIdFull := fmt.Sprintf("%s_%s", groupId, id)
		ruleCpy := rule

		wg.Add(1)
		go func() {
			defer wg.Done()

			results = append(results, ruleCpy.evaluate(ruleIdFull, path))
		}()
	}

	wg.Wait()
	return results
}

func (rule *Rule) evaluate(ruleId string, path string) *RuleResult {
	ret := &RuleResult{
		Id: ruleId,
	}

	return ret
}
