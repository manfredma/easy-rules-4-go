package api

import "math"

type EngineParameters struct {
	SkipOnFirstAppliedRule      bool
	SkipOnFirstNonTriggeredRule bool
	SkipOnFirstFailedRule       bool
	PriorityThreshold           int
}

func NewEngineParameters() *EngineParameters {
	return &EngineParameters{
		PriorityThreshold: math.MaxInt,
	}
}

type RulesEngine interface {
	GetParameters() *EngineParameters
	GetRuleListeners() []RuleListener
	GetEngineListeners() []RulesEngineListener
	Fire(rules *Rules, facts *Facts)
	Check(rules *Rules, facts *Facts) map[Rule]bool
}
