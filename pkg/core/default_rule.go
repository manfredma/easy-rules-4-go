package core

import "github.com/manfredma/easy-rules-4-go/pkg/api"

type DefaultRule struct {
	BasicRule
	condition api.Condition
	actions   []api.Action
}

func (r *DefaultRule) Evaluate(facts *api.Facts) bool {
	return r.condition(facts)
}

func (r *DefaultRule) Execute(facts *api.Facts) error {
	for _, action := range r.actions {
		if err := action(facts); err != nil {
			return err
		}
	}
	return nil
}
