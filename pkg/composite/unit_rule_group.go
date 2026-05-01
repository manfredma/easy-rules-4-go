package composite

import "github.com/manfredma/easy-rules-4-go/pkg/api"

type UnitRuleGroup struct {
	CompositeRule
}

func NewUnitRuleGroup(name, description string, priority int) *UnitRuleGroup {
	g := &UnitRuleGroup{}
	g.CompositeRule = *NewCompositeRule(name, description, priority)
	return g
}

func (g *UnitRuleGroup) Evaluate(facts *api.Facts) bool {
	if len(g.rules) == 0 {
		return false
	}
	for _, r := range g.rules {
		if !r.Evaluate(facts) {
			return false
		}
	}
	return true
}

func (g *UnitRuleGroup) Execute(facts *api.Facts) error {
	for _, r := range g.rules {
		if err := r.Execute(facts); err != nil {
			return err
		}
	}
	return nil
}
