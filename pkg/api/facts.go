package api

import "fmt"

type Fact struct {
	Name  string
	Value any
}

func (f *Fact) String() string {
	return fmt.Sprintf("Fact{name=%s, value=%v}", f.Name, f.Value)
}

type Facts struct {
	facts map[string]any
}

func NewFacts() *Facts {
	return &Facts{facts: make(map[string]any)}
}

func (f *Facts) Put(name string, value any) {
	f.facts[name] = value
}

func (f *Facts) Get(name string) any {
	return f.facts[name]
}

func (f *Facts) Remove(name string) {
	delete(f.facts, name)
}

func (f *Facts) AsMap() map[string]any {
	result := make(map[string]any, len(f.facts))
	for k, v := range f.facts {
		result[k] = v
	}
	return result
}

func (f *Facts) Clear() {
	f.facts = make(map[string]any)
}

func (f *Facts) String() string {
	return fmt.Sprintf("%v", f.facts)
}
