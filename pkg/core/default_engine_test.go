package core

import (
	"errors"
	"testing"

	"github.com/manfredma/easy-rules-4-go/pkg/api"
)

func makeCountRule(name string, priority int, result *int, increment int) api.Rule {
	return NewRuleBuilder().
		Name(name).
		Priority(priority).
		When(api.ConditionTrue).
		Then(func(f *api.Facts) error { *result += increment; return nil }).
		Build()
}

func TestDefaultEngine_FiresMatchingRules(t *testing.T) {
	count := 0
	rules := api.NewRules(
		makeCountRule("r1", 1, &count, 1),
		makeCountRule("r2", 2, &count, 10),
	)
	engine := NewDefaultRulesEngine(api.NewEngineParameters())
	engine.Fire(rules, api.NewFacts())
	if count != 11 {
		t.Errorf("expected 11, got %d", count)
	}
}

func TestDefaultEngine_SkipsNonMatchingRules(t *testing.T) {
	count := 0
	rules := api.NewRules(
		NewRuleBuilder().Name("r1").When(api.ConditionFalse).Then(func(_ *api.Facts) error { count++; return nil }).Build(),
	)
	engine := NewDefaultRulesEngine(api.NewEngineParameters())
	engine.Fire(rules, api.NewFacts())
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}
}

func TestDefaultEngine_SkipOnFirstAppliedRule(t *testing.T) {
	count := 0
	rules := api.NewRules(
		makeCountRule("r1", 1, &count, 1),
		makeCountRule("r2", 2, &count, 10),
	)
	params := api.NewEngineParameters()
	params.SkipOnFirstAppliedRule = true
	engine := NewDefaultRulesEngine(params)
	engine.Fire(rules, api.NewFacts())
	if count != 1 {
		t.Errorf("expected 1 (stop after first applied), got %d", count)
	}
}

func TestDefaultEngine_SkipOnFirstNonTriggeredRule(t *testing.T) {
	count := 0
	rules := api.NewRules(
		NewRuleBuilder().Name("r1").Priority(1).When(api.ConditionFalse).Then(func(_ *api.Facts) error { count++; return nil }).Build(),
		makeCountRule("r2", 2, &count, 10),
	)
	params := api.NewEngineParameters()
	params.SkipOnFirstNonTriggeredRule = true
	engine := NewDefaultRulesEngine(params)
	engine.Fire(rules, api.NewFacts())
	if count != 0 {
		t.Errorf("expected 0 (stop after first non-triggered), got %d", count)
	}
}

func TestDefaultEngine_SkipOnFirstFailedRule(t *testing.T) {
	count := 0
	rules := api.NewRules(
		NewRuleBuilder().Name("r1").Priority(1).When(api.ConditionTrue).Then(func(_ *api.Facts) error { return errors.New("fail") }).Build(),
		makeCountRule("r2", 2, &count, 10),
	)
	params := api.NewEngineParameters()
	params.SkipOnFirstFailedRule = true
	engine := NewDefaultRulesEngine(params)
	engine.Fire(rules, api.NewFacts())
	if count != 0 {
		t.Errorf("expected 0 (stop after first failed), got %d", count)
	}
}

func TestDefaultEngine_PriorityThreshold(t *testing.T) {
	count := 0
	rules := api.NewRules(
		makeCountRule("r1", 1, &count, 1),
		makeCountRule("r2", 10, &count, 100),
	)
	params := api.NewEngineParameters()
	params.PriorityThreshold = 5
	engine := NewDefaultRulesEngine(params)
	engine.Fire(rules, api.NewFacts())
	if count != 1 {
		t.Errorf("expected 1 (r2 skipped by threshold), got %d", count)
	}
}

func TestDefaultEngine_Check(t *testing.T) {
	rules := api.NewRules(
		NewRuleBuilder().Name("yes").When(api.ConditionTrue).Build(),
		NewRuleBuilder().Name("no").When(api.ConditionFalse).Build(),
	)
	engine := NewDefaultRulesEngine(api.NewEngineParameters())
	result := engine.Check(rules, api.NewFacts())
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
	for r, v := range result {
		if r.GetName() == "yes" && !v {
			t.Error("expected yes=true")
		}
		if r.GetName() == "no" && v {
			t.Error("expected no=false")
		}
	}
}

func TestDefaultEngine_RuleListener(t *testing.T) {
	var beforeEvalCalled, onSuccessCalled bool

	listener := &testRuleListener{
		beforeEvaluate: func(_ api.Rule, _ *api.Facts) bool { beforeEvalCalled = true; return true },
		onSuccess:      func(_ api.Rule, _ *api.Facts) { onSuccessCalled = true },
	}

	rules := api.NewRules(NewRuleBuilder().Name("r1").When(api.ConditionTrue).Then(func(_ *api.Facts) error { return nil }).Build())
	engine := NewDefaultRulesEngine(api.NewEngineParameters())
	engine.RegisterRuleListener(listener)
	engine.Fire(rules, api.NewFacts())

	if !beforeEvalCalled {
		t.Error("expected BeforeEvaluate to be called")
	}
	if !onSuccessCalled {
		t.Error("expected OnSuccess to be called")
	}
}

type testRuleListener struct {
	api.DefaultRuleListener
	beforeEvaluate func(api.Rule, *api.Facts) bool
	onSuccess      func(api.Rule, *api.Facts)
	onFailure      func(api.Rule, *api.Facts, error)
}

func (l *testRuleListener) BeforeEvaluate(r api.Rule, f *api.Facts) bool {
	if l.beforeEvaluate != nil {
		return l.beforeEvaluate(r, f)
	}
	return true
}

func (l *testRuleListener) OnSuccess(r api.Rule, f *api.Facts) {
	if l.onSuccess != nil {
		l.onSuccess(r, f)
	}
}

func (l *testRuleListener) OnFailure(r api.Rule, f *api.Facts, err error) {
	if l.onFailure != nil {
		l.onFailure(r, f, err)
	}
}
