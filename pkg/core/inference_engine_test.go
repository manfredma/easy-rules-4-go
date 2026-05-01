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
