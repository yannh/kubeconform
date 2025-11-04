package sarif

// NewRule creates a new Rule and returns a pointer to it
func NewRule(ruleID string) *ReportingDescriptor {
	return NewReportingDescriptor().WithID(ruleID)
}

// WithDescription specifies short description for a rule and returns the updated rule.
// Short description should be a single sentence that is understandable when visible space is limited to a single line
// of text.
func (rule *ReportingDescriptor) WithDescription(description string) *ReportingDescriptor {
	rule.ShortDescription = NewMultiformatMessageString().WithText(description)
	return rule
}

// WithMarkdownHelp specifies a help text  for a rule and returns the updated rule
func (rule *ReportingDescriptor) WithMarkdownHelp(markdown string) *ReportingDescriptor {
	if rule.Help == nil {
		rule.Help = NewMultiformatMessageString()
	}
	rule.Help.Text = &markdown
	rule.Help.Markdown = &markdown
	return rule
}
