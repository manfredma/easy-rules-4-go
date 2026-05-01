package core

import "github.com/manfredma/easy-rules-4-go/pkg/api"

type RuleBuilder struct {
	name        string
	description string
	priority    int
	condition   api.Condition
	actions     []api.Action
}

func NewRuleBuilder() *RuleBuilder {
	return &RuleBuilder{
		name:        api.DefaultRuleName,
		description: api.DefaultRuleDescription,
		priority:    api.DefaultRulePriority,
		condition:   api.ConditionFalse,
	}
}

func (b *RuleBuilder) Name(name string) *RuleBuilder {
	b.name = name
	return b
}

func (b *RuleBuilder) Description(desc string) *RuleBuilder {
	b.description = desc
	return b
}

func (b *RuleBuilder) Priority(priority int) *RuleBuilder {
	b.priority = priority
	return b
}

func (b *RuleBuilder) When(condition api.Condition) *RuleBuilder {
	b.condition = condition
	return b
}

func (b *RuleBuilder) Then(action api.Action) *RuleBuilder {
	b.actions = append(b.actions, action)
	return b
}

func (b *RuleBuilder) Build() api.Rule {
	r := &DefaultRule{
		condition: b.condition,
		actions:   b.actions,
	}
	r.BasicRule = *NewBasicRule(b.name, b.description, b.priority)
	return r
}
