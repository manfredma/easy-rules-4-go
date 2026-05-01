package composite

import "github.com/manfredma/easy-rules-4-go/pkg/api"

type ActivationRuleGroup struct {
	CompositeRule
	selectedRule api.Rule
}

func NewActivationRuleGroup(name, description string, priority int) *ActivationRuleGroup {
	g := &ActivationRuleGroup{}
	g.CompositeRule = *NewCompositeRule(name, description, priority)
	return g
}

func (g *ActivationRuleGroup) Evaluate(facts *api.Facts) bool {
	for _, r := range g.rules {
		if r.Evaluate(facts) {
			g.selectedRule = r
			return true
		}
	}
	return false
}

func (g *ActivationRuleGroup) Execute(facts *api.Facts) error {
	if g.selectedRule != nil {
		return g.selectedRule.Execute(facts)
	}
	return nil
}
