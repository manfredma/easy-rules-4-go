package main

import (
	"fmt"

	"github.com/manfredma/easy-rules-4-go/pkg/api"
	"github.com/manfredma/easy-rules-4-go/pkg/core"
)

type WeatherRule struct {
	core.BasicRule
}

func (r *WeatherRule) Evaluate(facts *api.Facts) bool {
	rain, _ := facts.Get("rain").(bool)
	return rain
}

func (r *WeatherRule) Execute(_ *api.Facts) error {
	fmt.Println("It rains, take an umbrella!")
	return nil
}

func main() {
	rule := &WeatherRule{}
	rule.BasicRule = *core.NewBasicRule("weather rule", "if it rains then take an umbrella", 1)

	facts := api.NewFacts()
	facts.Put("rain", true)

	rules := api.NewRules(rule)
	engine := core.NewDefaultRulesEngine(api.NewEngineParameters())
	engine.Fire(rules, facts)
}
