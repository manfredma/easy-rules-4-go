package api

type RuleListener interface {
	BeforeEvaluate(rule Rule, facts *Facts) bool
	AfterEvaluate(rule Rule, facts *Facts, evaluationResult bool)
	OnEvaluationError(rule Rule, facts *Facts, err error)
	BeforeExecute(rule Rule, facts *Facts)
	OnSuccess(rule Rule, facts *Facts)
	OnFailure(rule Rule, facts *Facts, err error)
}

type RulesEngineListener interface {
	BeforeEvaluate(rules *Rules, facts *Facts)
	AfterExecute(rules *Rules, facts *Facts)
}

// DefaultRuleListener provides no-op implementations of RuleListener.
type DefaultRuleListener struct{}

func (d *DefaultRuleListener) BeforeEvaluate(_ Rule, _ *Facts) bool        { return true }
func (d *DefaultRuleListener) AfterEvaluate(_ Rule, _ *Facts, _ bool)      {}
func (d *DefaultRuleListener) OnEvaluationError(_ Rule, _ *Facts, _ error) {}
func (d *DefaultRuleListener) BeforeExecute(_ Rule, _ *Facts)              {}
func (d *DefaultRuleListener) OnSuccess(_ Rule, _ *Facts)                  {}
func (d *DefaultRuleListener) OnFailure(_ Rule, _ *Facts, _ error)         {}

// DefaultRulesEngineListener provides no-op implementations of RulesEngineListener.
type DefaultRulesEngineListener struct{}

func (d *DefaultRulesEngineListener) BeforeEvaluate(_ *Rules, _ *Facts) {}
func (d *DefaultRulesEngineListener) AfterExecute(_ *Rules, _ *Facts)   {}
