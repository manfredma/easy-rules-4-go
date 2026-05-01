package composite

import (
	"testing"

	"github.com/manfredma/easy-rules-4-go/pkg/api"
	"github.com/manfredma/easy-rules-4-go/pkg/core"
)

func TestConditionalRuleGroup_GatePassesExecutesOthers(t *testing.T) {
	count := 0
	g := NewConditionalRuleGroup("g", "", 1)
	g.AddRule(core.NewRuleBuilder().Name("gate").Priority(0).When(api.ConditionTrue).Then(func(_ *api.Facts) error { count++; return nil }).Build())
	g.AddRule(core.NewRuleBuilder().Name("r1").Priority(1).When(api.ConditionTrue).Then(func(_ *api.Facts) error { count++; return nil }).Build())
	g.AddRule(core.NewRuleBuilder().Name("r2").Priority(2).When(api.ConditionTrue).Then(func(_ *api.Facts) error { count++; return nil }).Build())

	facts := api.NewFacts()
	if !g.Evaluate(facts) {
		t.Error("expected true when gate passes")
	}
	_ = g.Execute(facts)
	if count != 3 {
		t.Errorf("expected 3 (gate + r1 + r2), got %d", count)
	}
}

func TestConditionalRuleGroup_GateFailsSkipsAll(t *testing.T) {
	count := 0
	g := NewConditionalRuleGroup("g", "", 1)
	g.AddRule(core.NewRuleBuilder().Name("gate").Priority(0).When(api.ConditionFalse).Then(func(_ *api.Facts) error { count++; return nil }).Build())
	g.AddRule(core.NewRuleBuilder().Name("r1").Priority(1).When(api.ConditionTrue).Then(func(_ *api.Facts) error { count++; return nil }).Build())

	if g.Evaluate(api.NewFacts()) {
		t.Error("expected false when gate fails")
	}
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}
}

func TestConditionalRuleGroup_OnlyGateFired_WhenOthersFalse(t *testing.T) {
	count := 0
	g := NewConditionalRuleGroup("g", "", 1)
	g.AddRule(core.NewRuleBuilder().Name("gate").Priority(0).When(api.ConditionTrue).Then(func(_ *api.Facts) error { count++; return nil }).Build())
	g.AddRule(core.NewRuleBuilder().Name("r1").Priority(1).When(api.ConditionFalse).Then(func(_ *api.Facts) error { count++; return nil }).Build())

	facts := api.NewFacts()
	g.Evaluate(facts)
	_ = g.Execute(facts)
	if count != 1 {
		t.Errorf("expected only gate to fire (count=1), got %d", count)
	}
}

func TestConditionalRuleGroup_PanicsOnTwoRulesWithSameHighestPriority(t *testing.T) {
	g := NewConditionalRuleGroup("g", "", 1)
	g.AddRule(core.NewRuleBuilder().Name("r1").Priority(0).When(api.ConditionTrue).Build())
	g.AddRule(core.NewRuleBuilder().Name("r2").Priority(0).When(api.ConditionTrue).Build())

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for duplicate highest priority")
		}
	}()
	g.Evaluate(api.NewFacts())
}
