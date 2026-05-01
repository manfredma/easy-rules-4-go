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
