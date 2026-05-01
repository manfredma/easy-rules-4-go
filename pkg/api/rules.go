package api

import "sort"

type Rules struct {
	rules []Rule
}

func NewRules(rules ...Rule) *Rules {
	r := &Rules{}
	r.Register(rules...)
	return r
}

func (rs *Rules) Register(rules ...Rule) {
	for _, r := range rules {
		rs.rules = append(rs.rules, r)
	}
	rs.sort()
}

func (rs *Rules) Unregister(name string) {
	for i, r := range rs.rules {
		if r.GetName() == name {
			rs.rules = append(rs.rules[:i], rs.rules[i+1:]...)
			return
		}
	}
}

func (rs *Rules) IsEmpty() bool {
	return len(rs.rules) == 0
}

func (rs *Rules) Size() int {
	return len(rs.rules)
}

func (rs *Rules) Slice() []Rule {
	result := make([]Rule, len(rs.rules))
	copy(result, rs.rules)
	return result
}

func (rs *Rules) sort() {
	sort.SliceStable(rs.rules, func(i, j int) bool {
		pi, pj := rs.rules[i].GetPriority(), rs.rules[j].GetPriority()
		if pi != pj {
			return pi < pj
		}
		return rs.rules[i].GetName() < rs.rules[j].GetName()
	})
}
