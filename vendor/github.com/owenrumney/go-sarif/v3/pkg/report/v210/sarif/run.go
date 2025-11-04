package sarif

// Run - Describes a single run of an analysis tool, and contains the reported output of that run.
type Run struct {
	// The artifact location specified by each uriBaseId symbol on the machine where the tool originally ran.
	OriginalUriBaseIds map[string]ArtifactLocation `json:"originalUriBaseIds,omitempty"`

	// Addresses associated with this run instance, if any.
	Addresses []*Address `json:"addresses"`

	// An array of artifact objects relevant to the run.
	Artifacts []*Artifact `json:"artifacts,omitempty"`

	// Automation details that describe this run.
	AutomationDetails *RunAutomationDetails `json:"automationDetails,omitempty"`

	// The 'guid' property of a previous SARIF 'run' that comprises the baseline that was used to compute result 'baselineState' properties for the run.
	BaselineGuID *string `json:"baselineGuid,omitempty"`

	// Specifies the unit in which the tool measures columns.
	ColumnKind *string `json:"columnKind,omitempty"`

	// A conversion object that describes how a converter transformed an analysis tool's native reporting format into the SARIF format.
	Conversion *Conversion `json:"conversion,omitempty"`

	// Specifies the default encoding for any artifact object that refers to a text file.
	DefaultEncoding *string `json:"defaultEncoding,omitempty"`

	// Specifies the default source language for any artifact object that refers to a text file that contains source code.
	DefaultSourceLanguage *string `json:"defaultSourceLanguage,omitempty"`

	// References to external property files that should be inlined with the content of a root log file.
	ExternalPropertyFileReferences *ExternalPropertyFileReferences `json:"externalPropertyFileReferences,omitempty"`

	// An array of zero or more unique graph objects associated with the run.
	Graphs []*Graph `json:"graphs"`

	// Describes the invocation of the analysis tool.
	Invocations []*Invocation `json:"invocations"`

	// The language of the messages emitted into the log file during this run (expressed as an ISO 639-1 two-letter lowercase culture code) and an optional region (expressed as an ISO 3166-1 two-letter uppercase subculture code associated with a country or region). The casing is recommended but not required (in order for this data to conform to RFC5646).
	Language string `json:"language"`

	// An array of logical locations such as namespaces, types or functions.
	LogicalLocations []*LogicalLocation `json:"logicalLocations"`

	// An ordered list of character sequences that were treated as line breaks when computing region information for the run.
	NewlineSequences []string `json:"newlineSequences"`

	// Contains configurations that may potentially override both reportingDescriptor.defaultConfiguration (the tool's default severities) and invocation.configurationOverrides (severities established at run-time from the command line).
	Policies []*ToolComponent `json:"policies"`

	// Key/value pairs that provide additional information about the run.
	Properties *PropertyBag `json:"properties,omitempty"`

	// An array of strings used to replace sensitive information in a redaction-aware property.
	RedactionTokens []string `json:"redactionTokens"`

	// The set of results contained in an SARIF log. The results array can be omitted when a run is solely exporting rules metadata. It must be present (but may be empty) if a log file represents an actual scan.
	Results []*Result `json:"results"`

	// Automation details that describe the aggregate of runs to which this run belongs.
	RunAggregates []*RunAutomationDetails `json:"runAggregates"`

	// A specialLocations object that defines locations of special significance to SARIF consumers.
	SpecialLocations *SpecialLocations `json:"specialLocations,omitempty"`

	// An array of toolComponent objects relevant to a taxonomy in which results are categorized.
	Taxonomies []*ToolComponent `json:"taxonomies"`

	// An array of threadFlowLocation objects cached at run level.
	ThreadFlowLocations []*ThreadFlowLocation `json:"threadFlowLocations"`

	// Information about the tool or tool pipeline that generated the results in this run. A run can only contain results produced by a single tool or tool pipeline. A run can aggregate results from multiple log files, as long as context around the tool run (tool command-line arguments and the like) is identical for all aggregated files.
	Tool *Tool `json:"tool,omitempty"`

	// The set of available translations of the localized data provided by the tool.
	Translations []*ToolComponent `json:"translations"`

	// Specifies the revision in version control of the artifacts that were scanned.
	VersionControlProvenance []*VersionControlDetails `json:"versionControlProvenance"`

	// An array of request objects cached at run level.
	WebRequests []*WebRequest `json:"webRequests"`

	// An array of response objects cached at run level.
	WebResponses []*WebResponse `json:"webResponses"`
}

// NewRun - creates a new
func NewRun() *Run {
	return &Run{
		Addresses:                make([]*Address, 0),
		Artifacts:                make([]*Artifact, 0),
		Graphs:                   make([]*Graph, 0),
		Invocations:              make([]*Invocation, 0),
		Language:                 "en-US",
		LogicalLocations:         make([]*LogicalLocation, 0),
		NewlineSequences:         []string{"\r\n", "\n"},
		Policies:                 make([]*ToolComponent, 0),
		RedactionTokens:          make([]string, 0),
		Results:                  make([]*Result, 0),
		RunAggregates:            make([]*RunAutomationDetails, 0),
		Taxonomies:               make([]*ToolComponent, 0),
		ThreadFlowLocations:      make([]*ThreadFlowLocation, 0),
		Translations:             make([]*ToolComponent, 0),
		VersionControlProvenance: make([]*VersionControlDetails, 0),
		WebRequests:              make([]*WebRequest, 0),
		WebResponses:             make([]*WebResponse, 0),
	}
}

// AddOriginalUriBaseId - add a single OriginalUriBaseId to the Run
func (o *Run) AddOriginalUriBaseId(key string, originalUriBaseId ArtifactLocation) *Run {
	o.OriginalUriBaseIds[key] = originalUriBaseId
	return o
}

// WithOriginalUriBaseIds - add a OriginalUriBaseIds to the Run
func (o *Run) WithOriginalUriBaseIds(originalUriBaseIds map[string]ArtifactLocation) *Run {
	o.OriginalUriBaseIds = originalUriBaseIds
	return o
}

// WithAddresses - add a Addresses to the Run
func (a *Run) WithAddresses(addresses []*Address) *Run {
	a.Addresses = addresses
	return a
}

// AddAddresse - add a single Addresse to the Run
func (a *Run) AddAddresse(addresse *Address) *Run {
	a.Addresses = append(a.Addresses, addresse)
	return a
}

// WithArtifacts - add a Artifacts to the Run
func (a *Run) WithArtifacts(artifacts []*Artifact) *Run {
	a.Artifacts = artifacts
	return a
}

// AddArtifact - add a single Artifact to the Run
func (a *Run) AddArtifact(artifact *Artifact) *Run {
	a.Artifacts = append(a.Artifacts, artifact)
	return a
}

// WithAutomationDetails - add a AutomationDetails to the Run
func (a *Run) WithAutomationDetails(automationDetails *RunAutomationDetails) *Run {
	a.AutomationDetails = automationDetails
	return a
}

// WithBaselineGuID - add a BaselineGuID to the Run
func (b *Run) WithBaselineGuID(baselineGuid string) *Run {
	b.BaselineGuID = &baselineGuid
	return b
}

// WithColumnKind - add a ColumnKind to the Run
func (c *Run) WithColumnKind(columnKind string) *Run {
	c.ColumnKind = &columnKind
	return c
}

// WithConversion - add a Conversion to the Run
func (c *Run) WithConversion(conversion *Conversion) *Run {
	c.Conversion = conversion
	return c
}

// WithDefaultEncoding - add a DefaultEncoding to the Run
func (d *Run) WithDefaultEncoding(defaultEncoding string) *Run {
	d.DefaultEncoding = &defaultEncoding
	return d
}

// WithDefaultSourceLanguage - add a DefaultSourceLanguage to the Run
func (d *Run) WithDefaultSourceLanguage(defaultSourceLanguage string) *Run {
	d.DefaultSourceLanguage = &defaultSourceLanguage
	return d
}

// WithExternalPropertyFileReferences - add a ExternalPropertyFileReferences to the Run
func (e *Run) WithExternalPropertyFileReferences(externalPropertyFileReferences *ExternalPropertyFileReferences) *Run {
	e.ExternalPropertyFileReferences = externalPropertyFileReferences
	return e
}

// WithGraphs - add a Graphs to the Run
func (g *Run) WithGraphs(graphs []*Graph) *Run {
	g.Graphs = graphs
	return g
}

// AddGraph - add a single Graph to the Run
func (g *Run) AddGraph(graph *Graph) *Run {
	g.Graphs = append(g.Graphs, graph)
	return g
}

// WithInvocations - add a Invocations to the Run
func (i *Run) WithInvocations(invocations []*Invocation) *Run {
	i.Invocations = invocations
	return i
}

// AddInvocation - add a single Invocation to the Run
func (i *Run) AddInvocation(invocation *Invocation) *Run {
	i.Invocations = append(i.Invocations, invocation)
	return i
}

// WithLanguage - add a Language to the Run
func (l *Run) WithLanguage(language string) *Run {
	l.Language = language
	return l
}

// WithLogicalLocations - add a LogicalLocations to the Run
func (l *Run) WithLogicalLocations(logicalLocations []*LogicalLocation) *Run {
	l.LogicalLocations = logicalLocations
	return l
}

// AddLogicalLocation - add a single LogicalLocation to the Run
func (l *Run) AddLogicalLocation(logicalLocation *LogicalLocation) *Run {
	l.LogicalLocations = append(l.LogicalLocations, logicalLocation)
	return l
}

// WithNewlineSequences - add a NewlineSequences to the Run
func (n *Run) WithNewlineSequences(newlineSequences []string) *Run {
	n.NewlineSequences = newlineSequences
	return n
}

// AddNewlineSequence - add a single NewlineSequence to the Run
func (n *Run) AddNewlineSequence(newlineSequence string) *Run {
	n.NewlineSequences = append(n.NewlineSequences, newlineSequence)
	return n
}

// WithPolicies - add a Policies to the Run
func (p *Run) WithPolicies(policies []*ToolComponent) *Run {
	p.Policies = policies
	return p
}

// AddPolicie - add a single Policie to the Run
func (p *Run) AddPolicie(policie *ToolComponent) *Run {
	p.Policies = append(p.Policies, policie)
	return p
}

// WithProperties - add a Properties to the Run
func (p *Run) WithProperties(properties *PropertyBag) *Run {
	p.Properties = properties
	return p
}

// WithRedactionTokens - add a RedactionTokens to the Run
func (r *Run) WithRedactionTokens(redactionTokens []string) *Run {
	r.RedactionTokens = redactionTokens
	return r
}

// AddRedactionToken - add a single RedactionToken to the Run
func (r *Run) AddRedactionToken(redactionToken string) *Run {
	r.RedactionTokens = append(r.RedactionTokens, redactionToken)
	return r
}

// WithResults - add a Results to the Run
func (r *Run) WithResults(results []*Result) *Run {
	r.Results = results
	return r
}

// AddResult - add a single Result to the Run
func (r *Run) AddResult(result *Result) *Run {
	r.Results = append(r.Results, result)
	return r
}

// WithRunAggregates - add a RunAggregates to the Run
func (r *Run) WithRunAggregates(runAggregates []*RunAutomationDetails) *Run {
	r.RunAggregates = runAggregates
	return r
}

// AddRunAggregate - add a single RunAggregate to the Run
func (r *Run) AddRunAggregate(runAggregate *RunAutomationDetails) *Run {
	r.RunAggregates = append(r.RunAggregates, runAggregate)
	return r
}

// WithSpecialLocations - add a SpecialLocations to the Run
func (s *Run) WithSpecialLocations(specialLocations *SpecialLocations) *Run {
	s.SpecialLocations = specialLocations
	return s
}

// WithTaxonomies - add a Taxonomies to the Run
func (t *Run) WithTaxonomies(taxonomies []*ToolComponent) *Run {
	t.Taxonomies = taxonomies
	return t
}

// AddTaxonomie - add a single Taxonomie to the Run
func (t *Run) AddTaxonomie(taxonomie *ToolComponent) *Run {
	t.Taxonomies = append(t.Taxonomies, taxonomie)
	return t
}

// WithThreadFlowLocations - add a ThreadFlowLocations to the Run
func (t *Run) WithThreadFlowLocations(threadFlowLocations []*ThreadFlowLocation) *Run {
	t.ThreadFlowLocations = threadFlowLocations
	return t
}

// AddThreadFlowLocation - add a single ThreadFlowLocation to the Run
func (t *Run) AddThreadFlowLocation(threadFlowLocation *ThreadFlowLocation) *Run {
	t.ThreadFlowLocations = append(t.ThreadFlowLocations, threadFlowLocation)
	return t
}

// WithTool - add a Tool to the Run
func (t *Run) WithTool(tool *Tool) *Run {
	t.Tool = tool
	return t
}

// WithTranslations - add a Translations to the Run
func (t *Run) WithTranslations(translations []*ToolComponent) *Run {
	t.Translations = translations
	return t
}

// AddTranslation - add a single Translation to the Run
func (t *Run) AddTranslation(translation *ToolComponent) *Run {
	t.Translations = append(t.Translations, translation)
	return t
}

// WithVersionControlProvenance - add a VersionControlProvenance to the Run
func (v *Run) WithVersionControlProvenance(versionControlProvenance []*VersionControlDetails) *Run {
	v.VersionControlProvenance = versionControlProvenance
	return v
}

// AddVersionControlProvenance - add a single VersionControlProvenance to the Run
func (v *Run) AddVersionControlProvenance(versionControlProvenance *VersionControlDetails) *Run {
	v.VersionControlProvenance = append(v.VersionControlProvenance, versionControlProvenance)
	return v
}

// WithWebRequests - add a WebRequests to the Run
func (w *Run) WithWebRequests(webRequests []*WebRequest) *Run {
	w.WebRequests = webRequests
	return w
}

// AddWebRequest - add a single WebRequest to the Run
func (w *Run) AddWebRequest(webRequest *WebRequest) *Run {
	w.WebRequests = append(w.WebRequests, webRequest)
	return w
}

// WithWebResponses - add a WebResponses to the Run
func (w *Run) WithWebResponses(webResponses []*WebResponse) *Run {
	w.WebResponses = webResponses
	return w
}

// AddWebResponse - add a single WebResponse to the Run
func (w *Run) AddWebResponse(webResponse *WebResponse) *Run {
	w.WebResponses = append(w.WebResponses, webResponse)
	return w
}
