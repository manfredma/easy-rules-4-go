package composite

import (
	"sort"

	"github.com/manfredma/easy-rules-4-go/pkg/api"
)

type ConditionalRuleGroup struct {
	CompositeRule
	conditionalRule       api.Rule
	successfulEvaluations []api.Rule
}

func NewConditionalRuleGroup(name, description string, priority int) *ConditionalRuleGroup {
	g := &ConditionalRuleGroup{}
	g.CompositeRule = *NewCompositeRule(name, description, priority)
	return g
}

func (g *ConditionalRuleGroup) Evaluate(facts *api.Facts) bool {
	g.successfulEvaluations = nil
	g.conditionalRule = g.getRuleWithHighestPriority()
	if g.conditionalRule.Evaluate(facts) {
		for _, r := range g.rules {
			if r != g.conditionalRule && r.Evaluate(facts) {
				g.successfulEvaluations = append(g.successfulEvaluations, r)
			}
		}
		return true
	}
	return false
}

func (g *ConditionalRuleGroup) Execute(facts *api.Facts) error {
	if err := g.conditionalRule.Execute(facts); err != nil {
		return err
	}
	sorted := make([]api.Rule, len(g.successfulEvaluations))
	copy(sorted, g.successfulEvaluations)
	sort.SliceStable(sorted, func(i, j int) bool {
		pi, pj := sorted[i].GetPriority(), sorted[j].GetPriority()
		if pi != pj {
			return pi < pj
		}
		return sorted[i].GetName() < sorted[j].GetName()
	})
	for _, r := range sorted {
		if err := r.Execute(facts); err != nil {
			return err
		}
	}
	return nil
}

func (g *ConditionalRuleGroup) getRuleWithHighestPriority() api.Rule {
	if len(g.rules) == 0 {
		panic("ConditionalRuleGroup has no rules")
	}
	sorted := make([]api.Rule, len(g.rules))
	copy(sorted, g.rules)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].GetPriority() < sorted[j].GetPriority()
	})
	highest := sorted[0]
	if len(sorted) > 1 && sorted[1].GetPriority() == highest.GetPriority() {
		panic("ConditionalRuleGroup: only one rule can have the highest priority")
	}
	return highest
}
