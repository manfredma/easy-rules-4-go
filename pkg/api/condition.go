package api

type Condition func(facts *Facts) bool

var ConditionTrue Condition = func(_ *Facts) bool { return true }
var ConditionFalse Condition = func(_ *Facts) bool { return false }
