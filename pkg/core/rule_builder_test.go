package core

import (
	"errors"
	"testing"

	"github.com/manfredma/easy-rules-4-go/pkg/api"
)

func TestRuleBuilder_BuildWithConditionAndAction(t *testing.T) {
	facts := api.NewFacts()
	facts.Put("x", 10)

	rule := NewRuleBuilder().
		Name("test-rule").
		Description("desc").
		Priority(1).
		When(func(f *api.Facts) bool {
			return f.Get("x").(int) > 5
		}).
		Then(func(f *api.Facts) error {
			f.Put("result", "big")
			return nil
		}).
		Build()

	if rule.GetName() != "test-rule" {
		t.Errorf("expected test-rule, got %s", rule.GetName())
	}
	if !rule.Evaluate(facts) {
		t.Error("expected condition to be true")
	}
	_ = rule.Execute(facts)
	if facts.Get("result") != "big" {
		t.Error("expected result=big after execute")
	}
}

func TestRuleBuilder_DefaultsWhenNotSet(t *testing.T) {
	rule := NewRuleBuilder().Build()
	if rule.GetName() != api.DefaultRuleName {
		t.Errorf("expected default name, got %s", rule.GetName())
	}
	if rule.GetPriority() != api.DefaultRulePriority {
		t.Errorf("expected default priority, got %d", rule.GetPriority())
	}
	if rule.Evaluate(api.NewFacts()) {
		t.Error("default condition should return false")
	}
}

func TestRuleBuilder_MultipleActions(t *testing.T) {
	facts := api.NewFacts()
	facts.Put("count", 0)

	rule := NewRuleBuilder().
		When(api.ConditionTrue).
		Then(func(f *api.Facts) error { f.Put("count", 1); return nil }).
		Then(func(f *api.Facts) error { f.Put("count", f.Get("count").(int)+1); return nil }).
		Build()

	_ = rule.Execute(facts)
	if facts.Get("count") != 2 {
		t.Errorf("expected count=2, got %v", facts.Get("count"))
	}
}

func TestRuleBuilder_ActionError(t *testing.T) {
	rule := NewRuleBuilder().
		When(api.ConditionTrue).
		Then(func(_ *api.Facts) error { return errors.New("boom") }).
		Build()

	err := rule.Execute(api.NewFacts())
	if err == nil || err.Error() != "boom" {
		t.Errorf("expected error boom, got %v", err)
	}
}
