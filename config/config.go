package config

var (
	LoadedRuleSet *RuleSet
)

type Level string

const (
	Undefined Level = ""
	Warning   Level = "warning"
	Error     Level = "error"
)

type RuleSet struct {
	Name        string
	Description string
	RuleGroups  *map[string]RuleGroup
}

type RuleGroup struct {
	Description string
	Rules       *[]Rule
}

type Rule struct {
	Description string
	File        *string
	Conditions  *[]Condition
	Level       Level
}

type Condition struct {
	PathExists          *string
	Contains            *[]ContainsCondition
	NotContains         *[]string
	CheckReferenceExist *[]string
}

type ContainsCondition struct {
	Type  string // static or regex
	Value string
}
