package expr

import (
	"github.com/expr-lang/expr"
	"github.com/manfredma/easy-rules-4-go/pkg/api"
)

func NewExprCondition(expression string) (api.Condition, error) {
	program, err := expr.Compile(expression)
	if err != nil {
		return nil, err
	}
	return func(facts *api.Facts) bool {
		result, err := expr.Run(program, facts.AsMap())
		if err != nil {
			return false
		}
		v, _ := result.(bool)
		return v
	}, nil
}
