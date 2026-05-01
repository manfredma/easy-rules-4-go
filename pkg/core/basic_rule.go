package core

import "github.com/manfredma/easy-rules-4-go/pkg/api"

type BasicRule struct {
	name        string
	description string
	priority    int
	prioritySet bool
}

func NewBasicRule(name, description string, priority int) *BasicRule {
	return &BasicRule{name: name, description: description, priority: priority, prioritySet: true}
}

func (r *BasicRule) GetName() string {
	if r.name == "" {
		return api.DefaultRuleName
	}
	return r.name
}

func (r *BasicRule) GetDescription() string {
	if r.description == "" {
		return api.DefaultRuleDescription
	}
	return r.description
}

func (r *BasicRule) GetPriority() int {
	if !r.prioritySet {
		return api.DefaultRulePriority
	}
	return r.priority
}

func (r *BasicRule) Evaluate(_ *api.Facts) bool { return false }
func (r *BasicRule) Execute(_ *api.Facts) error  { return nil }

func (r *BasicRule) SetName(name string)        { r.name = name }
func (r *BasicRule) SetDescription(desc string) { r.description = desc }
func (r *BasicRule) SetPriority(priority int)   { r.priority = priority; r.prioritySet = true }
