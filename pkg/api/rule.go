package api

const (
	DefaultRuleName        = "rule"
	DefaultRuleDescription = "description"
	DefaultRulePriority    = int(^uint(0)>>1) - 1 // math.MaxInt - 1
)

type Rule interface {
	GetName() string
	GetDescription() string
	GetPriority() int
	Evaluate(facts *Facts) bool
	Execute(facts *Facts) error
}
