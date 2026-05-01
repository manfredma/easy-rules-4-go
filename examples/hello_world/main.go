package main

import (
	"fmt"

	"github.com/manfredma/easy-rules-4-go/pkg/api"
	"github.com/manfredma/easy-rules-4-go/pkg/core"
)

func main() {
	rule := core.NewRuleBuilder().
		Name("hello world rule").
		Description("say hello world").
		When(api.ConditionTrue).
		Then(func(_ *api.Facts) error {
			fmt.Println("Hello World!")
			return nil
		}).
		Build()

	rules := api.NewRules(rule)
	engine := core.NewDefaultRulesEngine(api.NewEngineParameters())
	engine.Fire(rules, api.NewFacts())
}
