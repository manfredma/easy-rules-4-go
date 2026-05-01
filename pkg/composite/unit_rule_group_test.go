package composite

import (
	"fmt"
	"testing"

	"github.com/manfredma/easy-rules-4-go/pkg/api"
	"github.com/manfredma/easy-rules-4-go/pkg/core"
)

func trueRule(name string) api.Rule {
	return core.NewRuleBuilder().Name(name).When(api.ConditionTrue).Then(func(_ *api.Facts) error { return nil }).Build()
}

func falseRule(name string) api.Rule {
	return core.NewRuleBuilder().Name(name).When(api.ConditionFalse).Then(func(_ *api.Facts) error { return nil }).Build()
}

func TestUnitRuleGroup_AllTrueEvaluatesTrue(t *testing.T) {
	g := NewUnitRuleGroup("g", "", 1)
	g.AddRule(trueRule("r1"))
	g.AddRule(trueRule("r2"))
	if !g.Evaluate(api.NewFacts()) {
		t.Error("expected true when all rules are true")
	}
}

func TestUnitRuleGroup_OneFalseEvaluatesFalse(t *testing.T) {
	g := NewUnitRuleGroup("g", "", 1)
	g.AddRule(trueRule("r1"))
	g.AddRule(falseRule("r2"))
	if g.Evaluate(api.NewFacts()) {
		t.Error("expected false when one rule is false")
	}
}

func TestUnitRuleGroup_ExecutesAllRules(t *testing.T) {
	count := 0
	g := NewUnitRuleGroup("g", "", 1)
	for i := 0; i < 3; i++ {
		idx := i
		g.AddRule(core.NewRuleBuilder().Name(fmt.Sprintf("r%d", idx)).When(api.ConditionTrue).Then(func(_ *api.Facts) error { count++; return nil }).Build())
	}
	_ = g.Execute(api.NewFacts())
	if count != 3 {
		t.Errorf("expected 3 executions, got %d", count)
	}
}

func TestUnitRuleGroup_EmptyEvaluatesFalse(t *testing.T) {
	g := NewUnitRuleGroup("g", "", 1)
	if g.Evaluate(api.NewFacts()) {
		t.Error("expected false for empty group")
	}
}
