package expr

import (
	"testing"

	"github.com/manfredma/easy-rules-4-go/pkg/api"
)

func TestExprCondition_TrueExpression(t *testing.T) {
	facts := api.NewFacts()
	facts.Put("rain", true)

	cond, err := NewExprCondition("rain == true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cond(facts) {
		t.Error("expected condition to be true")
	}
}

func TestExprCondition_FalseExpression(t *testing.T) {
	facts := api.NewFacts()
	facts.Put("rain", false)

	cond, err := NewExprCondition("rain == true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cond(facts) {
		t.Error("expected condition to be false")
	}
}

func TestExprCondition_InvalidExpression(t *testing.T) {
	_, err := NewExprCondition("((( unmatched")
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestExprAction_ModifyFacts(t *testing.T) {
	facts := api.NewFacts()
	facts.Put("num", 0)

	action, err := NewExprAction("num + 1", "num")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = action(facts)
	if facts.Get("num").(int) != 1 {
		t.Errorf("expected num=1, got %v", facts.Get("num"))
	}
}

func TestExprRule_EvaluateAndExecute(t *testing.T) {
	facts := api.NewFacts()
	facts.Put("temperature", 35)
	facts.Put("cooled", false)

	rule, err := NewExprRuleBuilder().
		Name("heat-rule").
		When("temperature > 30").
		Then("!cooled", "cooled").
		Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !rule.Evaluate(facts) {
		t.Error("expected evaluate to be true")
	}
	_ = rule.Execute(facts)
	if facts.Get("cooled") != true {
		t.Errorf("expected cooled=true, got %v", facts.Get("cooled"))
	}
}

func TestExprRule_InvalidCondition(t *testing.T) {
	_, err := NewExprRuleBuilder().
		Name("bad").
		When("((( unmatched").
		Build()
	if err == nil {
		t.Error("expected error for invalid condition")
	}
}
