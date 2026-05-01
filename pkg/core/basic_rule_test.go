package core

import (
	"testing"

	"github.com/manfredma/easy-rules-4-go/pkg/api"
)

func TestBasicRule_Defaults(t *testing.T) {
	r := &BasicRule{}
	if r.GetName() != api.DefaultRuleName {
		t.Errorf("expected %s, got %s", api.DefaultRuleName, r.GetName())
	}
	if r.GetDescription() != api.DefaultRuleDescription {
		t.Errorf("expected %s, got %s", api.DefaultRuleDescription, r.GetDescription())
	}
	if r.GetPriority() != api.DefaultRulePriority {
		t.Errorf("expected %d, got %d", api.DefaultRulePriority, r.GetPriority())
	}
}

func TestBasicRule_CustomValues(t *testing.T) {
	r := NewBasicRule("my-rule", "desc", 5)
	if r.GetName() != "my-rule" {
		t.Errorf("expected my-rule, got %s", r.GetName())
	}
	if r.GetPriority() != 5 {
		t.Errorf("expected 5, got %d", r.GetPriority())
	}
}

func TestBasicRule_EvaluateReturnsFalse(t *testing.T) {
	r := &BasicRule{}
	if r.Evaluate(api.NewFacts()) {
		t.Error("default Evaluate should return false")
	}
}

func TestBasicRule_ExecuteReturnsNil(t *testing.T) {
	r := &BasicRule{}
	if err := r.Execute(api.NewFacts()); err != nil {
		t.Errorf("default Execute should return nil, got %v", err)
	}
}

type EmbedRule struct {
	BasicRule
}

func (e *EmbedRule) Evaluate(_ *api.Facts) bool { return true }
func (e *EmbedRule) Execute(_ *api.Facts) error  { return nil }

func TestBasicRule_Embed(t *testing.T) {
	r := &EmbedRule{}
	r.BasicRule = *NewBasicRule("embed", "embed rule", 1)
	if r.GetName() != "embed" {
		t.Errorf("expected embed, got %s", r.GetName())
	}
	if !r.Evaluate(api.NewFacts()) {
		t.Error("expected true from embedded evaluate")
	}
}
