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

func TestDefaultEngine_GetParameters(t *testing.T) {
	params := api.NewEngineParameters()
	params.PriorityThreshold = 42
	engine := NewDefaultRulesEngine(params)
	if engine.GetParameters().PriorityThreshold != 42 {
		t.Error("expected PriorityThreshold=42")
	}
}

func TestDefaultEngine_GetRuleListeners(t *testing.T) {
	engine := NewDefaultRulesEngine(api.NewEngineParameters())
	if engine.GetRuleListeners() != nil {
		t.Error("expected nil listeners initially")
	}
	engine.RegisterRuleListener(&testRuleListener{})
	if len(engine.GetRuleListeners()) != 1 {
		t.Error("expected 1 rule listener")
	}
}

func TestDefaultEngine_GetEngineListeners(t *testing.T) {
	engine := NewDefaultRulesEngine(api.NewEngineParameters())
	if engine.GetEngineListeners() != nil {
		t.Error("expected nil engine listeners initially")
	}
	engine.RegisterEngineListener(&testEngineListener{})
	if len(engine.GetEngineListeners()) != 1 {
		t.Error("expected 1 engine listener")
	}
}

func TestDefaultEngine_EngineListenerCallbacks(t *testing.T) {
	var beforeCalled, afterCalled bool
	listener := &testEngineListener{
		beforeEvaluate: func(_ *api.Rules, _ *api.Facts) { beforeCalled = true },
		afterExecute:   func(_ *api.Rules, _ *api.Facts) { afterCalled = true },
	}
	rules := api.NewRules(NewRuleBuilder().Name("r1").When(api.ConditionTrue).Then(func(_ *api.Facts) error { return nil }).Build())
	engine := NewDefaultRulesEngine(api.NewEngineParameters())
	engine.RegisterEngineListener(listener)
	engine.Fire(rules, api.NewFacts())
	if !beforeCalled {
		t.Error("expected engine BeforeEvaluate to be called")
	}
	if !afterCalled {
		t.Error("expected engine AfterExecute to be called")
	}
}

func TestDefaultEngine_OnFailureListener(t *testing.T) {
	var failureCalled bool
	listener := &testRuleListener{
		onFailure: func(_ api.Rule, _ *api.Facts, _ error) { failureCalled = true },
	}
	rules := api.NewRules(
		NewRuleBuilder().Name("r1").When(api.ConditionTrue).Then(func(_ *api.Facts) error { return errors.New("fail") }).Build(),
	)
	engine := NewDefaultRulesEngine(api.NewEngineParameters())
	engine.RegisterRuleListener(listener)
	engine.Fire(rules, api.NewFacts())
	if !failureCalled {
		t.Error("expected OnFailure to be called")
	}
}

func TestDefaultEngine_NilParams(t *testing.T) {
	// NewDefaultRulesEngine with nil params should not panic
	engine := NewDefaultRulesEngine(nil)
	if engine.GetParameters() == nil {
		t.Error("expected non-nil parameters")
	}
}

func TestDefaultEngine_PanicRecovery(t *testing.T) {
	// A rule whose Evaluate panics should trigger OnEvaluationError
	var evalErrCalled bool
	listener := &testRuleListener{
		evalErr: func(_ api.Rule, _ *api.Facts, _ error) { evalErrCalled = true },
	}
	panicRule := NewRuleBuilder().
		Name("panic-rule").
		When(func(_ *api.Facts) bool { panic("boom") }).
		Then(func(_ *api.Facts) error { return nil }).
		Build()
	rules := api.NewRules(panicRule)
	engine := NewDefaultRulesEngine(api.NewEngineParameters())
	engine.RegisterRuleListener(listener)
	engine.Fire(rules, api.NewFacts()) // must not panic
	if !evalErrCalled {
		t.Error("expected OnEvaluationError to be called on panic")
	}
}

func TestDefaultEngine_SkipOnFirstNonTriggeredAfterPanic(t *testing.T) {
	// panic in Evaluate with SkipOnFirstNonTriggeredRule=true should stop
	count := 0
	panicRule := NewRuleBuilder().
		Name("panic-rule").Priority(1).
		When(func(_ *api.Facts) bool { panic("boom") }).
		Then(func(_ *api.Facts) error { return nil }).
		Build()
	countRule := makeCountRule("r2", 2, &count, 1)
	params := api.NewEngineParameters()
	params.SkipOnFirstNonTriggeredRule = true
	engine := NewDefaultRulesEngine(params)
	engine.Fire(api.NewRules(panicRule, countRule), api.NewFacts())
	if count != 0 {
		t.Errorf("expected 0 (stop after panic with SkipOnFirstNonTriggeredRule), got %d", count)
	}
}

type testRuleListener struct {
	api.DefaultRuleListener
	beforeEvaluate func(api.Rule, *api.Facts) bool
	onSuccess      func(api.Rule, *api.Facts)
	onFailure      func(api.Rule, *api.Facts, error)
	evalErr        func(api.Rule, *api.Facts, error)
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

func (l *testRuleListener) OnEvaluationError(r api.Rule, f *api.Facts, err error) {
	if l.evalErr != nil {
		l.evalErr(r, f, err)
	}
}

type testEngineListener struct {
	api.DefaultRulesEngineListener
	beforeEvaluate func(*api.Rules, *api.Facts)
	afterExecute   func(*api.Rules, *api.Facts)
}

func (l *testEngineListener) BeforeEvaluate(rules *api.Rules, facts *api.Facts) {
	if l.beforeEvaluate != nil {
		l.beforeEvaluate(rules, facts)
	}
}

func (l *testEngineListener) AfterExecute(rules *api.Rules, facts *api.Facts) {
	if l.afterExecute != nil {
		l.afterExecute(rules, facts)
	}
}
