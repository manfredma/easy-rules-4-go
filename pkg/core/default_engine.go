package core

import (
	"fmt"

	"github.com/manfredma/easy-rules-4-go/pkg/api"
)

type DefaultRulesEngine struct {
	parameters      *api.EngineParameters
	ruleListeners   []api.RuleListener
	engineListeners []api.RulesEngineListener
}

func NewDefaultRulesEngine(params *api.EngineParameters) *DefaultRulesEngine {
	if params == nil {
		params = api.NewEngineParameters()
	}
	return &DefaultRulesEngine{parameters: params}
}

func (e *DefaultRulesEngine) GetParameters() *api.EngineParameters         { return e.parameters }
func (e *DefaultRulesEngine) GetRuleListeners() []api.RuleListener          { return e.ruleListeners }
func (e *DefaultRulesEngine) GetEngineListeners() []api.RulesEngineListener { return e.engineListeners }

func (e *DefaultRulesEngine) RegisterRuleListener(l api.RuleListener) {
	e.ruleListeners = append(e.ruleListeners, l)
}

func (e *DefaultRulesEngine) RegisterEngineListener(l api.RulesEngineListener) {
	e.engineListeners = append(e.engineListeners, l)
}

func (e *DefaultRulesEngine) Fire(rules *api.Rules, facts *api.Facts) {
	e.triggerEngineListenersBefore(rules, facts)
	e.doFire(rules, facts)
	e.triggerEngineListenersAfter(rules, facts)
}

func (e *DefaultRulesEngine) doFire(rules *api.Rules, facts *api.Facts) {
	if rules.IsEmpty() {
		return
	}
	for _, rule := range rules.Slice() {
		if rule.GetPriority() > e.parameters.PriorityThreshold {
			break
		}
		if !e.shouldBeEvaluated(rule, facts) {
			continue
		}
		evaluationResult := false
		var evalErr error
		func() {
			defer func() {
				if rec := recover(); rec != nil {
					evalErr = fmt.Errorf("panic: %v", rec)
				}
			}()
			evaluationResult = rule.Evaluate(facts)
		}()
		if evalErr != nil {
			e.triggerOnEvaluationError(rule, facts, evalErr)
			if e.parameters.SkipOnFirstNonTriggeredRule {
				break
			}
			continue
		}
		if evaluationResult {
			e.triggerAfterEvaluate(rule, facts, true)
			e.triggerBeforeExecute(rule, facts)
			err := rule.Execute(facts)
			if err != nil {
				e.triggerOnFailure(rule, facts, err)
				if e.parameters.SkipOnFirstFailedRule {
					break
				}
			} else {
				e.triggerOnSuccess(rule, facts)
				if e.parameters.SkipOnFirstAppliedRule {
					break
				}
			}
		} else {
			e.triggerAfterEvaluate(rule, facts, false)
			if e.parameters.SkipOnFirstNonTriggeredRule {
				break
			}
		}
	}
}

func (e *DefaultRulesEngine) Check(rules *api.Rules, facts *api.Facts) map[api.Rule]bool {
	e.triggerEngineListenersBefore(rules, facts)
	result := make(map[api.Rule]bool)
	for _, rule := range rules.Slice() {
		if e.shouldBeEvaluated(rule, facts) {
			result[rule] = rule.Evaluate(facts)
		}
	}
	e.triggerEngineListenersAfter(rules, facts)
	return result
}

func (e *DefaultRulesEngine) shouldBeEvaluated(rule api.Rule, facts *api.Facts) bool {
	for _, l := range e.ruleListeners {
		if !l.BeforeEvaluate(rule, facts) {
			return false
		}
	}
	return true
}

func (e *DefaultRulesEngine) triggerAfterEvaluate(rule api.Rule, facts *api.Facts, result bool) {
	for _, l := range e.ruleListeners {
		l.AfterEvaluate(rule, facts, result)
	}
}

func (e *DefaultRulesEngine) triggerOnEvaluationError(rule api.Rule, facts *api.Facts, err error) {
	for _, l := range e.ruleListeners {
		l.OnEvaluationError(rule, facts, err)
	}
}

func (e *DefaultRulesEngine) triggerBeforeExecute(rule api.Rule, facts *api.Facts) {
	for _, l := range e.ruleListeners {
		l.BeforeExecute(rule, facts)
	}
}

func (e *DefaultRulesEngine) triggerOnSuccess(rule api.Rule, facts *api.Facts) {
	for _, l := range e.ruleListeners {
		l.OnSuccess(rule, facts)
	}
}

func (e *DefaultRulesEngine) triggerOnFailure(rule api.Rule, facts *api.Facts, err error) {
	for _, l := range e.ruleListeners {
		l.OnFailure(rule, facts, err)
	}
}

func (e *DefaultRulesEngine) triggerEngineListenersBefore(rules *api.Rules, facts *api.Facts) {
	for _, l := range e.engineListeners {
		l.BeforeEvaluate(rules, facts)
	}
}

func (e *DefaultRulesEngine) triggerEngineListenersAfter(rules *api.Rules, facts *api.Facts) {
	for _, l := range e.engineListeners {
		l.AfterExecute(rules, facts)
	}
}
