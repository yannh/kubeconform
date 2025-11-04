package sarif

// ReportingConfiguration - Information about a rule or notification that can be configured at runtime.
type ReportingConfiguration struct {
	// Specifies whether the report may be produced during the scan.
	Enabled bool `json:"enabled"`

	// Specifies the failure level for the report.
	Level string `json:"level,omitempty"`

	// Contains configuration information specific to a report.
	Parameters *PropertyBag `json:"parameters,omitempty"`

	// Key/value pairs that provide additional information about the reporting configuration.
	Properties *PropertyBag `json:"properties,omitempty"`

	// Specifies the relative priority of the report. Used for analysis output only.
	Rank float64 `json:"rank"`
}

// NewReportingConfiguration - creates a new
func NewReportingConfiguration() *ReportingConfiguration {
	return &ReportingConfiguration{
		Enabled: true,
		Level:   "warning",
		Rank:    -1.000000,
	}
}

// WithEnabled - add a Enabled to the ReportingConfiguration
func (e *ReportingConfiguration) WithEnabled(enabled bool) *ReportingConfiguration {
	e.Enabled = enabled
	return e
}

// WithLevel - add a Level to the ReportingConfiguration
func (l *ReportingConfiguration) WithLevel(level string) *ReportingConfiguration {
	l.Level = level
	return l
}

// WithParameters - add a Parameters to the ReportingConfiguration
func (p *ReportingConfiguration) WithParameters(parameters *PropertyBag) *ReportingConfiguration {
	p.Parameters = parameters
	return p
}

// WithProperties - add a Properties to the ReportingConfiguration
func (p *ReportingConfiguration) WithProperties(properties *PropertyBag) *ReportingConfiguration {
	p.Properties = properties
	return p
}

// WithRank - add a Rank to the ReportingConfiguration
func (r *ReportingConfiguration) WithRank(rank float64) *ReportingConfiguration {
	r.Rank = rank
	return r
}
