package core

import "github.com/manfredma/easy-rules-4-go/pkg/api"

type InferenceRulesEngine struct {
	delegate *DefaultRulesEngine
}

func NewInferenceRulesEngine(params *api.EngineParameters) *InferenceRulesEngine {
	return &InferenceRulesEngine{delegate: NewDefaultRulesEngine(params)}
}

func (e *InferenceRulesEngine) GetParameters() *api.EngineParameters {
	return e.delegate.GetParameters()
}

func (e *InferenceRulesEngine) GetRuleListeners() []api.RuleListener {
	return e.delegate.GetRuleListeners()
}

func (e *InferenceRulesEngine) GetEngineListeners() []api.RulesEngineListener {
	return e.delegate.GetEngineListeners()
}

func (e *InferenceRulesEngine) RegisterRuleListener(l api.RuleListener) {
	e.delegate.RegisterRuleListener(l)
}

func (e *InferenceRulesEngine) RegisterEngineListener(l api.RulesEngineListener) {
	e.delegate.RegisterEngineListener(l)
}

func (e *InferenceRulesEngine) Fire(rules *api.Rules, facts *api.Facts) {
	for {
		candidates := e.selectCandidates(rules, facts)
		if candidates.IsEmpty() {
			break
		}
		e.delegate.Fire(candidates, facts)
	}
}

func (e *InferenceRulesEngine) Check(rules *api.Rules, facts *api.Facts) map[api.Rule]bool {
	return e.delegate.Check(rules, facts)
}

func (e *InferenceRulesEngine) selectCandidates(rules *api.Rules, facts *api.Facts) *api.Rules {
	candidates := api.NewRules()
	for _, rule := range rules.Slice() {
		if rule.Evaluate(facts) {
			candidates.Register(rule)
		}
	}
	return candidates
}
