package composite

import (
	"testing"

	"github.com/manfredma/easy-rules-4-go/pkg/api"
	"github.com/manfredma/easy-rules-4-go/pkg/core"
)

func TestActivationRuleGroup_FiresFirstMatchingRule(t *testing.T) {
	fired := ""
	g := NewActivationRuleGroup("g", "", 1)
	g.AddRule(core.NewRuleBuilder().Name("r1").Priority(1).When(api.ConditionTrue).Then(func(_ *api.Facts) error { fired = "r1"; return nil }).Build())
	g.AddRule(core.NewRuleBuilder().Name("r2").Priority(2).When(api.ConditionTrue).Then(func(_ *api.Facts) error { fired = "r2"; return nil }).Build())

	facts := api.NewFacts()
	if !g.Evaluate(facts) {
		t.Error("expected true")
	}
	_ = g.Execute(facts)
	if fired != "r1" {
		t.Errorf("expected r1 (highest priority), got %s", fired)
	}
}

func TestActivationRuleGroup_ReturnsFalseWhenNoneMatch(t *testing.T) {
	g := NewActivationRuleGroup("g", "", 1)
	g.AddRule(falseRule("r1"))
	if g.Evaluate(api.NewFacts()) {
		t.Error("expected false")
	}
}

func TestActivationRuleGroup_SkipsLowerPriorityRules(t *testing.T) {
	count := 0
	g := NewActivationRuleGroup("g", "", 1)
	g.AddRule(core.NewRuleBuilder().Name("r1").Priority(1).When(api.ConditionTrue).Then(func(_ *api.Facts) error { count++; return nil }).Build())
	g.AddRule(core.NewRuleBuilder().Name("r2").Priority(2).When(api.ConditionTrue).Then(func(_ *api.Facts) error { count++; return nil }).Build())

	facts := api.NewFacts()
	g.Evaluate(facts)
	_ = g.Execute(facts)
	if count != 1 {
		t.Errorf("expected only 1 rule fired (XOR), got %d", count)
	}
}
