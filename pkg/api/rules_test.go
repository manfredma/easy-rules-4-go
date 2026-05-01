package api

import "testing"

type mockRule struct {
	name     string
	priority int
}

func (r *mockRule) GetName() string        { return r.name }
func (r *mockRule) GetDescription() string { return "" }
func (r *mockRule) GetPriority() int       { return r.priority }
func (r *mockRule) Evaluate(_ *Facts) bool { return false }
func (r *mockRule) Execute(_ *Facts) error { return nil }

func TestRules_RegisterAndSize(t *testing.T) {
	rules := NewRules()
	rules.Register(&mockRule{name: "r1", priority: 1})
	rules.Register(&mockRule{name: "r2", priority: 2})
	if rules.Size() != 2 {
		t.Errorf("expected 2, got %d", rules.Size())
	}
}

func TestRules_SortedByPriority(t *testing.T) {
	rules := NewRules()
	rules.Register(&mockRule{name: "r2", priority: 2})
	rules.Register(&mockRule{name: "r1", priority: 1})
	names := []string{}
	for _, r := range rules.Slice() {
		names = append(names, r.GetName())
	}
	if names[0] != "r1" || names[1] != "r2" {
		t.Errorf("expected [r1, r2], got %v", names)
	}
}

func TestRules_SamePrioritySortedByName(t *testing.T) {
	rules := NewRules()
	rules.Register(&mockRule{name: "z", priority: 1})
	rules.Register(&mockRule{name: "a", priority: 1})
	names := []string{}
	for _, r := range rules.Slice() {
		names = append(names, r.GetName())
	}
	if names[0] != "a" || names[1] != "z" {
		t.Errorf("expected [a, z], got %v", names)
	}
}

func TestRules_UnregisterByName(t *testing.T) {
	rules := NewRules()
	rules.Register(&mockRule{name: "r1", priority: 1})
	rules.Unregister("r1")
	if rules.Size() != 0 {
		t.Error("expected empty after unregister")
	}
}

func TestRules_IsEmpty(t *testing.T) {
	rules := NewRules()
	if !rules.IsEmpty() {
		t.Error("expected empty")
	}
	rules.Register(&mockRule{name: "r1", priority: 1})
	if rules.IsEmpty() {
		t.Error("expected not empty")
	}
}
