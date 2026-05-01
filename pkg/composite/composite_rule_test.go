package composite

import (
	"errors"
	"testing"

	"github.com/manfredma/easy-rules-4-go/pkg/api"
	"github.com/manfredma/easy-rules-4-go/pkg/core"
)

func TestCompositeRule_AddAndGetRules(t *testing.T) {
	c := NewCompositeRule("c", "desc", 1)
	c.AddRule(trueRule("r1"))
	c.AddRule(trueRule("r2"))
	rules := c.GetRules()
	if len(rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(rules))
	}
}

func TestCompositeRule_RemoveRule(t *testing.T) {
	c := NewCompositeRule("c", "desc", 1)
	c.AddRule(trueRule("r1"))
	c.AddRule(trueRule("r2"))
	c.RemoveRule("r1")
	rules := c.GetRules()
	if len(rules) != 1 {
		t.Errorf("expected 1 rule after remove, got %d", len(rules))
	}
	if rules[0].GetName() != "r2" {
		t.Errorf("expected r2 to remain, got %s", rules[0].GetName())
	}
}

func TestCompositeRule_RemoveNonExistentRule(t *testing.T) {
	c := NewCompositeRule("c", "desc", 1)
	c.AddRule(trueRule("r1"))
	c.RemoveRule("doesNotExist") // should not panic or remove anything
	if len(c.GetRules()) != 1 {
		t.Error("expected 1 rule, removal of nonexistent should be a no-op")
	}
}

func TestCompositeRule_AddRuleSortsbyPriority(t *testing.T) {
	c := NewCompositeRule("c", "desc", 1)
	c.AddRule(core.NewRuleBuilder().Name("high").Priority(10).When(api.ConditionTrue).Build())
	c.AddRule(core.NewRuleBuilder().Name("low").Priority(1).When(api.ConditionTrue).Build())
	rules := c.GetRules()
	if rules[0].GetName() != "low" {
		t.Errorf("expected low-priority rule first, got %s", rules[0].GetName())
	}
}

func TestUnitRuleGroup_ExecuteReturnsErrorOnFailure(t *testing.T) {
	g := NewUnitRuleGroup("g", "", 1)
	g.AddRule(core.NewRuleBuilder().Name("fail").When(api.ConditionTrue).Then(func(_ *api.Facts) error {
		return errors.New("unit fail")
	}).Build())
	err := g.Execute(api.NewFacts())
	if err == nil {
		t.Error("expected error from failing rule in UnitRuleGroup")
	}
}

func TestActivationRuleGroup_ExecuteReturnsErrorOnFailure(t *testing.T) {
	g := NewActivationRuleGroup("g", "", 1)
	g.AddRule(core.NewRuleBuilder().Name("fail").Priority(1).When(api.ConditionTrue).Then(func(_ *api.Facts) error {
		return errors.New("activation fail")
	}).Build())
	facts := api.NewFacts()
	g.Evaluate(facts)
	err := g.Execute(facts)
	if err == nil {
		t.Error("expected error from failing rule in ActivationRuleGroup")
	}
}

func TestConditionalRuleGroup_ExecuteReturnsErrorOnFailure(t *testing.T) {
	g := NewConditionalRuleGroup("g", "", 1)
	g.AddRule(core.NewRuleBuilder().Name("gate").Priority(0).When(api.ConditionTrue).Then(func(_ *api.Facts) error {
		return errors.New("gate fail")
	}).Build())
	facts := api.NewFacts()
	g.Evaluate(facts)
	err := g.Execute(facts)
	if err == nil {
		t.Error("expected error from failing gate in ConditionalRuleGroup")
	}
}
