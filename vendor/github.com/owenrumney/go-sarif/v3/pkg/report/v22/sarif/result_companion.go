package sarif

// NewRuleResult - creates a new result with the ruleID set
func NewRuleResult(ruleID string) *Result {
	return NewResult().WithRuleID(ruleID)
}
