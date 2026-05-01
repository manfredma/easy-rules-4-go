package expr

import (
	"github.com/manfredma/easy-rules-4-go/pkg/api"
	"github.com/manfredma/easy-rules-4-go/pkg/core"
)

type ExprRule struct {
	core.BasicRule
	condition api.Condition
	actions   []api.Action
}

func (r *ExprRule) Evaluate(facts *api.Facts) bool {
	return r.condition(facts)
}

func (r *ExprRule) Execute(facts *api.Facts) error {
	for _, action := range r.actions {
		if err := action(facts); err != nil {
			return err
		}
	}
	return nil
}

type ExprRuleBuilder struct {
	name        string
	description string
	priority    int
	condition   string
	actions     []struct{ expr, key string }
}

func NewExprRuleBuilder() *ExprRuleBuilder {
	return &ExprRuleBuilder{
		name:        api.DefaultRuleName,
		description: api.DefaultRuleDescription,
		priority:    api.DefaultRulePriority,
	}
}

func (b *ExprRuleBuilder) Name(name string) *ExprRuleBuilder {
	b.name = name
	return b
}

func (b *ExprRuleBuilder) Description(desc string) *ExprRuleBuilder {
	b.description = desc
	return b
}

func (b *ExprRuleBuilder) Priority(priority int) *ExprRuleBuilder {
	b.priority = priority
	return b
}

func (b *ExprRuleBuilder) When(expression string) *ExprRuleBuilder {
	b.condition = expression
	return b
}

// Then adds an action: expression result is stored in outputKey of facts.
func (b *ExprRuleBuilder) Then(expression, outputKey string) *ExprRuleBuilder {
	b.actions = append(b.actions, struct{ expr, key string }{expression, outputKey})
	return b
}

func (b *ExprRuleBuilder) Build() (*ExprRule, error) {
	cond, err := NewExprCondition(b.condition)
	if err != nil {
		return nil, err
	}
	var actions []api.Action
	for _, a := range b.actions {
		action, err := NewExprAction(a.expr, a.key)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}
	r := &ExprRule{condition: cond, actions: actions}
	r.BasicRule = *core.NewBasicRule(b.name, b.description, b.priority)
	return r, nil
}
