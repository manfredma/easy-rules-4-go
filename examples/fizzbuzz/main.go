package main

import (
	"fmt"

	"github.com/manfredma/easy-rules-4-go/pkg/api"
	"github.com/manfredma/easy-rules-4-go/pkg/core"
	easyrules_expr "github.com/manfredma/easy-rules-4-go/pkg/expr"
)

func main() {
	// fizzBuzzRule has highest priority (lowest number = first evaluated)
	fizzBuzzRule, _ := easyrules_expr.NewExprRuleBuilder().
		Name("fizzbuzz rule").
		Priority(0).
		When(`number % 3 == 0 && number % 5 == 0`).
		Then(`"fizzbuzz"`, "print").
		Build()

	// fizzRule only fires if print not already set
	fizzRule, _ := easyrules_expr.NewExprRuleBuilder().
		Name("fizz rule").
		Priority(1).
		When(`number % 3 == 0 && print == nil`).
		Then(`"fizz"`, "print").
		Build()

	// buzzRule only fires if print not already set
	buzzRule, _ := easyrules_expr.NewExprRuleBuilder().
		Name("buzz rule").
		Priority(2).
		When(`number % 5 == 0 && print == nil`).
		Then(`"buzz"`, "print").
		Build()

	// numberRule always fires last and prints the accumulated result
	numberRule := core.NewRuleBuilder().
		Name("number rule").
		Priority(3).
		When(api.ConditionTrue).
		Then(func(facts *api.Facts) error {
			if facts.Get("print") == nil {
				fmt.Println(facts.Get("number"))
			} else {
				fmt.Println(facts.Get("print"))
				facts.Remove("print")
			}
			return nil
		}).
		Build()

	rules := api.NewRules(fizzBuzzRule, fizzRule, buzzRule, numberRule)
	engine := core.NewDefaultRulesEngine(api.NewEngineParameters())

	for i := 1; i <= 20; i++ {
		facts := api.NewFacts()
		facts.Put("number", i)
		engine.Fire(rules, facts)
	}
}
