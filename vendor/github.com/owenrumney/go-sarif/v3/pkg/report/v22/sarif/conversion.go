package sarif

// Conversion - Describes how a converter transformed the output of a static analysis tool from the analysis tool's native output format into the SARIF format.
type Conversion struct {
	// The locations of the analysis tool's per-run log files.
	AnalysisToolLogFiles []*ArtifactLocation `json:"analysisToolLogFiles"`

	// An invocation object that describes the invocation of the converter.
	Invocation *Invocation `json:"invocation,omitempty"`

	// Key/value pairs that provide additional information about the conversion.
	Properties *PropertyBag `json:"properties,omitempty"`

	// A tool object that describes the converter.
	Tool *Tool `json:"tool,omitempty"`
}

// NewConversion - creates a new
func NewConversion() *Conversion {
	return &Conversion{
		AnalysisToolLogFiles: make([]*ArtifactLocation, 0),
	}
}

// WithAnalysisToolLogFiles - add a AnalysisToolLogFiles to the Conversion
func (a *Conversion) WithAnalysisToolLogFiles(analysisToolLogFiles []*ArtifactLocation) *Conversion {
	a.AnalysisToolLogFiles = analysisToolLogFiles
	return a
}

// AddAnalysisToolLogFile - add a single AnalysisToolLogFile to the Conversion
func (a *Conversion) AddAnalysisToolLogFile(analysisToolLogFile *ArtifactLocation) *Conversion {
	a.AnalysisToolLogFiles = append(a.AnalysisToolLogFiles, analysisToolLogFile)
	return a
}

// WithInvocation - add a Invocation to the Conversion
func (i *Conversion) WithInvocation(invocation *Invocation) *Conversion {
	i.Invocation = invocation
	return i
}

// WithProperties - add a Properties to the Conversion
func (p *Conversion) WithProperties(properties *PropertyBag) *Conversion {
	p.Properties = properties
	return p
}

// WithTool - add a Tool to the Conversion
func (t *Conversion) WithTool(tool *Tool) *Conversion {
	t.Tool = tool
	return t
}
