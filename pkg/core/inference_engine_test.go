package core

import (
	"testing"

	"github.com/manfredma/easy-rules-4-go/pkg/api"
)

func TestInferenceEngine_LoopsUntilNoCandidate(t *testing.T) {
	facts := api.NewFacts()
	facts.Put("counter", 0)

	rule := NewRuleBuilder().
		Name("increment").
		When(func(f *api.Facts) bool { return f.Get("counter").(int) < 3 }).
		Then(func(f *api.Facts) error {
			f.Put("counter", f.Get("counter").(int)+1)
			return nil
		}).
		Build()

	rules := api.NewRules(rule)
	engine := NewInferenceRulesEngine(api.NewEngineParameters())
	engine.Fire(rules, facts)

	if facts.Get("counter").(int) != 3 {
		t.Errorf("expected counter=3, got %v", facts.Get("counter"))
	}
}

func TestInferenceEngine_StopsWhenNoCandidates(t *testing.T) {
	count := 0
	rules := api.NewRules(
		NewRuleBuilder().Name("r1").When(api.ConditionFalse).Then(func(_ *api.Facts) error { count++; return nil }).Build(),
	)
	engine := NewInferenceRulesEngine(api.NewEngineParameters())
	engine.Fire(rules, api.NewFacts())
	if count != 0 {
		t.Errorf("expected 0 executions, got %d", count)
	}
}

func TestInferenceEngine_GetParameters(t *testing.T) {
	params := api.NewEngineParameters()
	params.PriorityThreshold = 7
	engine := NewInferenceRulesEngine(params)
	if engine.GetParameters().PriorityThreshold != 7 {
		t.Error("expected PriorityThreshold=7")
	}
}

func TestInferenceEngine_GetRuleListeners(t *testing.T) {
	engine := NewInferenceRulesEngine(api.NewEngineParameters())
	if engine.GetRuleListeners() != nil {
		t.Error("expected nil rule listeners initially")
	}
	engine.RegisterRuleListener(&api.DefaultRuleListener{})
	if len(engine.GetRuleListeners()) != 1 {
		t.Error("expected 1 rule listener after register")
	}
}

func TestInferenceEngine_GetEngineListeners(t *testing.T) {
	engine := NewInferenceRulesEngine(api.NewEngineParameters())
	if engine.GetEngineListeners() != nil {
		t.Error("expected nil engine listeners initially")
	}
	engine.RegisterEngineListener(&api.DefaultRulesEngineListener{})
	if len(engine.GetEngineListeners()) != 1 {
		t.Error("expected 1 engine listener after register")
	}
}

func TestInferenceEngine_Check(t *testing.T) {
	rules := api.NewRules(
		NewRuleBuilder().Name("yes").When(api.ConditionTrue).Build(),
		NewRuleBuilder().Name("no").When(api.ConditionFalse).Build(),
	)
	engine := NewInferenceRulesEngine(api.NewEngineParameters())
	result := engine.Check(rules, api.NewFacts())
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}
