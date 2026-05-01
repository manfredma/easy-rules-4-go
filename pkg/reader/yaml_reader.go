package reader

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/manfredma/easy-rules-4-go/pkg/api"
	exprpkg "github.com/manfredma/easy-rules-4-go/pkg/expr"
)

type YamlRuleFactory struct{}

func NewYamlRuleFactory() *YamlRuleFactory { return &YamlRuleFactory{} }

func (f *YamlRuleFactory) CreateRulesFrom(path string) (*api.Rules, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var defs []RuleDefinition
	if err := yaml.Unmarshal(data, &defs); err != nil {
		return nil, err
	}
	return buildRules(defs)
}

func buildRules(defs []RuleDefinition) (*api.Rules, error) {
	rules := api.NewRules()
	for _, def := range defs {
		b := exprpkg.NewExprRuleBuilder().
			Name(def.Name).
			Description(def.Description).
			Priority(def.Priority).
			When(def.Condition)
		for _, a := range def.Actions {
			b.Then(a.Expr, a.OutputKey)
		}
		rule, err := b.Build()
		if err != nil {
			return nil, err
		}
		rules.Register(rule)
	}
	return rules, nil
}
