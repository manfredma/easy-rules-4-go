package api

import "testing"

func TestFacts_PutAndGet(t *testing.T) {
	facts := NewFacts()
	facts.Put("rain", true)
	val := facts.Get("rain")
	if val != true {
		t.Errorf("expected true, got %v", val)
	}
}

func TestFacts_PutOverwrite(t *testing.T) {
	facts := NewFacts()
	facts.Put("rain", true)
	facts.Put("rain", false)
	if facts.Get("rain") != false {
		t.Error("expected overwrite to false")
	}
}

func TestFacts_Remove(t *testing.T) {
	facts := NewFacts()
	facts.Put("rain", true)
	facts.Remove("rain")
	if facts.Get("rain") != nil {
		t.Error("expected nil after remove")
	}
}

func TestFacts_AsMap(t *testing.T) {
	facts := NewFacts()
	facts.Put("a", 1)
	facts.Put("b", 2)
	m := facts.AsMap()
	if m["a"] != 1 || m["b"] != 2 {
		t.Error("AsMap mismatch")
	}
}

func TestFacts_Clear(t *testing.T) {
	facts := NewFacts()
	facts.Put("a", 1)
	facts.Clear()
	if facts.Get("a") != nil {
		t.Error("expected nil after clear")
	}
}

func TestFacts_String(t *testing.T) {
	facts := NewFacts()
	facts.Put("x", 42)
	s := facts.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}

func TestFact_String(t *testing.T) {
	f := &Fact{Name: "rain", Value: true}
	s := f.String()
	if s == "" {
		t.Error("expected non-empty Fact.String()")
	}
}

func TestEngineParameters_Defaults(t *testing.T) {
	p := NewEngineParameters()
	if p.SkipOnFirstAppliedRule {
		t.Error("SkipOnFirstAppliedRule should default to false")
	}
	if p.PriorityThreshold == 0 {
		t.Error("PriorityThreshold should be MaxInt")
	}
}

func TestDefaultRuleListener_NoOps(t *testing.T) {
	l := &DefaultRuleListener{}
	facts := NewFacts()
	rule := &mockRule{name: "mock", priority: 0}
	// None of these should panic
	if !l.BeforeEvaluate(rule, facts) {
		t.Error("BeforeEvaluate should return true")
	}
	l.AfterEvaluate(rule, facts, true)
	l.OnEvaluationError(rule, facts, nil)
	l.BeforeExecute(rule, facts)
	l.OnSuccess(rule, facts)
	l.OnFailure(rule, facts, nil)
}

func TestDefaultRulesEngineListener_NoOps(t *testing.T) {
	l := &DefaultRulesEngineListener{}
	facts := NewFacts()
	rules := NewRules()
	// None of these should panic
	l.BeforeEvaluate(rules, facts)
	l.AfterExecute(rules, facts)
}
