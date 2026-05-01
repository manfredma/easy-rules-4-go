package composite

import (
	"sort"

	"github.com/manfredma/easy-rules-4-go/pkg/api"
	"github.com/manfredma/easy-rules-4-go/pkg/core"
)

type CompositeRule struct {
	core.BasicRule
	rules []api.Rule
}

func NewCompositeRule(name, description string, priority int) *CompositeRule {
	c := &CompositeRule{}
	c.BasicRule = *core.NewBasicRule(name, description, priority)
	return c
}

func (c *CompositeRule) AddRule(rule api.Rule) {
	c.rules = append(c.rules, rule)
	sort.SliceStable(c.rules, func(i, j int) bool {
		pi, pj := c.rules[i].GetPriority(), c.rules[j].GetPriority()
		if pi != pj {
			return pi < pj
		}
		return c.rules[i].GetName() < c.rules[j].GetName()
	})
}

func (c *CompositeRule) RemoveRule(name string) {
	for i, r := range c.rules {
		if r.GetName() == name {
			c.rules = append(c.rules[:i], c.rules[i+1:]...)
			return
		}
	}
}

func (c *CompositeRule) GetRules() []api.Rule {
	result := make([]api.Rule, len(c.rules))
	copy(result, c.rules)
	return result
}
