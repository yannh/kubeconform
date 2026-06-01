package sarif

import "github.com/owenrumney/go-sarif/v3/pkg/report/utils"

// AddRule returns an existing ReportingDescriptor for the ruleID or creates a new ReportingDescriptor and returns a pointer to it
func (run *Run) AddRule(ruleID string) *ReportingDescriptor {
	for _, rule := range run.Tool.Driver.Rules {
		if *rule.ID == ruleID {
			return rule
		}
	}
	rule := NewRule(ruleID)
	run.Tool.Driver.Rules = append(run.Tool.Driver.Rules, rule)
	return rule
}

// GetRuleIndex returns the index of the rule with the given ID, or -1 if it does not exist
func (run *Run) GetRuleIndex(ruleID string) int {
	if run.Tool == nil || run.Tool.Driver == nil || run.Tool.Driver.Rules == nil || ruleID == "" {
		return -1
	}

	for i, rule := range run.Tool.Driver.Rules {
		if *rule.ID == ruleID {
			return i
		}
	}
	return -1
}

// AddDistinctArtifact will handle deduplication of simple artifact additions
func (run *Run) AddDistinctArtifact(uri string) *Artifact {
	for _, artifact := range run.Artifacts {
		if *artifact.Location.URI == uri {
			return artifact
		}
	}

	a := NewArtifact().WithLength(utils.DefaultLengthInt)
	a.WithLocation(NewSimpleArtifactLocation(uri))

	run.Artifacts = append(run.Artifacts, a)
	return a
}

// CreateResultForRule returns an existing Result or creates a new one and returns a pointer to it
func (run *Run) CreateResultForRule(ruleID string) *Result {
	ruleIndex := run.GetRuleIndex(ruleID)
	result := NewRuleResult(ruleID).WithRuleIndex(ruleIndex)
	run.AddResult(result)
	return result
}

// DedupeArtifacts will remove any duplicate artifacts from the run
func (run *Run) DedupeArtifacts() error {
	dupes := map[*Artifact]bool{}
	deduped := []*Artifact{}

	for _, a := range run.Artifacts {
		if _, ok := dupes[a]; !ok {
			dupes[a] = true
			deduped = append(deduped, a)
		}
	}
	run.Artifacts = deduped
	return nil
}
