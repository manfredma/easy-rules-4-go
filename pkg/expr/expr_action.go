package expr

import (
	"github.com/expr-lang/expr"
	"github.com/manfredma/easy-rules-4-go/pkg/api"
)

// NewExprAction compiles an expression and stores its result back into facts under outputKey.
func NewExprAction(expression string, outputKey string) (api.Action, error) {
	program, err := expr.Compile(expression)
	if err != nil {
		return nil, err
	}
	return func(facts *api.Facts) error {
		result, err := expr.Run(program, facts.AsMap())
		if err != nil {
			return err
		}
		facts.Put(outputKey, result)
		return nil
	}, nil
}
