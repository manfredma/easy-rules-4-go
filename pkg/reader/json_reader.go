package reader

import (
	"encoding/json"
	"os"

	"github.com/manfredma/easy-rules-4-go/pkg/api"
)

type JsonRuleFactory struct{}

func NewJsonRuleFactory() *JsonRuleFactory { return &JsonRuleFactory{} }

func (f *JsonRuleFactory) CreateRulesFrom(path string) (*api.Rules, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var defs []RuleDefinition
	if err := json.Unmarshal(data, &defs); err != nil {
		return nil, err
	}
	return buildRules(defs)
}
