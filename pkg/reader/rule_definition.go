package reader

type ActionDefinition struct {
	Expr      string `yaml:"expr" json:"expr"`
	OutputKey string `yaml:"outputKey" json:"outputKey"`
}

type RuleDefinition struct {
	Name        string             `yaml:"name" json:"name"`
	Description string             `yaml:"description" json:"description"`
	Priority    int                `yaml:"priority" json:"priority"`
	Condition   string             `yaml:"condition" json:"condition"`
	Actions     []ActionDefinition `yaml:"actions" json:"actions"`
}
