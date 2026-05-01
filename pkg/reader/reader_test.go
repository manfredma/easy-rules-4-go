package reader

import (
	"testing"

	"github.com/manfredma/easy-rules-4-go/pkg/api"
)

func TestYamlRuleFactory_LoadAndFire(t *testing.T) {
	factory := NewYamlRuleFactory()
	rules, err := factory.CreateRulesFrom("testdata/weather-rules.yml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rules.Size() != 1 {
		t.Errorf("expected 1 rule, got %d", rules.Size())
	}

	facts := api.NewFacts()
	facts.Put("rain", true)
	facts.Put("umbrella", false)

	rule := rules.Slice()[0]
	if rule.GetName() != "weather rule" {
		t.Errorf("expected 'weather rule', got %s", rule.GetName())
	}
	if !rule.Evaluate(facts) {
		t.Error("expected condition to be true")
	}
	_ = rule.Execute(facts)
	if facts.Get("umbrella") != true {
		t.Errorf("expected umbrella=true after execute, got %v", facts.Get("umbrella"))
	}
}

func TestJsonRuleFactory_LoadAndFire(t *testing.T) {
	factory := NewJsonRuleFactory()
	rules, err := factory.CreateRulesFrom("testdata/weather-rules.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rules.Size() != 1 {
		t.Errorf("expected 1 rule, got %d", rules.Size())
	}

	facts := api.NewFacts()
	facts.Put("rain", true)
	facts.Put("umbrella", false)

	rule := rules.Slice()[0]
	if !rule.Evaluate(facts) {
		t.Error("expected condition to be true")
	}
	_ = rule.Execute(facts)
	if facts.Get("umbrella") != true {
		t.Errorf("expected umbrella=true, got %v", facts.Get("umbrella"))
	}
}

func TestYamlRuleFactory_FileNotFound(t *testing.T) {
	factory := NewYamlRuleFactory()
	_, err := factory.CreateRulesFrom("testdata/nonexistent.yml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
