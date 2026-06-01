package sarif

// ExternalProperties - The top-level element of an external property file.
type ExternalProperties struct {
	// Addresses that will be merged with a separate run.
	Addresses []*Address `json:"addresses"`

	// An array of artifact objects that will be merged with a separate run.
	Artifacts []*Artifact `json:"artifacts,omitempty"`

	// A conversion object that will be merged with a separate run.
	Conversion *Conversion `json:"conversion,omitempty"`

	// The analysis tool object that will be merged with a separate run.
	Driver *ToolComponent `json:"driver,omitempty"`

	// Tool extensions that will be merged with a separate run.
	Extensions []*ToolComponent `json:"extensions"`

	// Key/value pairs that provide additional information that will be merged with a separate run.
	ExternalizedProperties *PropertyBag `json:"externalizedProperties,omitempty"`

	// An array of graph objects that will be merged with a separate run.
	Graphs []*Graph `json:"graphs"`

	// A stable, unique identifier for this external properties object, in the form of a GUID.
	GuID *string `json:"guid,omitempty"`

	// Describes the invocation of the analysis tool that will be merged with a separate run.
	Invocations []*Invocation `json:"invocations"`

	// An array of logical locations such as namespaces, types or functions that will be merged with a separate run.
	LogicalLocations []*LogicalLocation `json:"logicalLocations"`

	// Tool policies that will be merged with a separate run.
	Policies []*ToolComponent `json:"policies"`

	// Key/value pairs that provide additional information about the external properties.
	Properties *PropertyBag `json:"properties,omitempty"`

	// An array of result objects that will be merged with a separate run.
	Results []*Result `json:"results"`

	// A stable, unique identifier for the run associated with this external properties object, in the form of a GUID.
	RunGuID *string `json:"runGuid,omitempty"`

	// The URI of the JSON schema corresponding to the version of the external property file format.
	Schema *string `json:"schema,omitempty"`

	// Tool taxonomies that will be merged with a separate run.
	Taxonomies []*ToolComponent `json:"taxonomies"`

	// An array of threadFlowLocation objects that will be merged with a separate run.
	ThreadFlowLocations []*ThreadFlowLocation `json:"threadFlowLocations"`

	// Tool translations that will be merged with a separate run.
	Translations []*ToolComponent `json:"translations"`

	// The SARIF format version of this external properties object.
	Version *string `json:"version,omitempty"`

	// Requests that will be merged with a separate run.
	WebRequests []*WebRequest `json:"webRequests"`

	// Responses that will be merged with a separate run.
	WebResponses []*WebResponse `json:"webResponses"`
}

// NewExternalProperties - creates a new
func NewExternalProperties() *ExternalProperties {
	return &ExternalProperties{
		Addresses:           make([]*Address, 0),
		Artifacts:           make([]*Artifact, 0),
		Extensions:          make([]*ToolComponent, 0),
		Graphs:              make([]*Graph, 0),
		Invocations:         make([]*Invocation, 0),
		LogicalLocations:    make([]*LogicalLocation, 0),
		Policies:            make([]*ToolComponent, 0),
		Results:             make([]*Result, 0),
		Taxonomies:          make([]*ToolComponent, 0),
		ThreadFlowLocations: make([]*ThreadFlowLocation, 0),
		Translations:        make([]*ToolComponent, 0),
		WebRequests:         make([]*WebRequest, 0),
		WebResponses:        make([]*WebResponse, 0),
	}
}

// WithAddresses - add a Addresses to the ExternalProperties
func (a *ExternalProperties) WithAddresses(addresses []*Address) *ExternalProperties {
	a.Addresses = addresses
	return a
}

// AddAddresse - add a single Addresse to the ExternalProperties
func (a *ExternalProperties) AddAddresse(addresse *Address) *ExternalProperties {
	a.Addresses = append(a.Addresses, addresse)
	return a
}

// WithArtifacts - add a Artifacts to the ExternalProperties
func (a *ExternalProperties) WithArtifacts(artifacts []*Artifact) *ExternalProperties {
	a.Artifacts = artifacts
	return a
}

// AddArtifact - add a single Artifact to the ExternalProperties
func (a *ExternalProperties) AddArtifact(artifact *Artifact) *ExternalProperties {
	a.Artifacts = append(a.Artifacts, artifact)
	return a
}

// WithConversion - add a Conversion to the ExternalProperties
func (c *ExternalProperties) WithConversion(conversion *Conversion) *ExternalProperties {
	c.Conversion = conversion
	return c
}

// WithDriver - add a Driver to the ExternalProperties
func (d *ExternalProperties) WithDriver(driver *ToolComponent) *ExternalProperties {
	d.Driver = driver
	return d
}

// WithExtensions - add a Extensions to the ExternalProperties
func (e *ExternalProperties) WithExtensions(extensions []*ToolComponent) *ExternalProperties {
	e.Extensions = extensions
	return e
}

// AddExtension - add a single Extension to the ExternalProperties
func (e *ExternalProperties) AddExtension(extension *ToolComponent) *ExternalProperties {
	e.Extensions = append(e.Extensions, extension)
	return e
}

// WithExternalizedProperties - add a ExternalizedProperties to the ExternalProperties
func (e *ExternalProperties) WithExternalizedProperties(externalizedProperties *PropertyBag) *ExternalProperties {
	e.ExternalizedProperties = externalizedProperties
	return e
}

// WithGraphs - add a Graphs to the ExternalProperties
func (g *ExternalProperties) WithGraphs(graphs []*Graph) *ExternalProperties {
	g.Graphs = graphs
	return g
}

// AddGraph - add a single Graph to the ExternalProperties
func (g *ExternalProperties) AddGraph(graph *Graph) *ExternalProperties {
	g.Graphs = append(g.Graphs, graph)
	return g
}

// WithGuID - add a GuID to the ExternalProperties
func (g *ExternalProperties) WithGuID(guid string) *ExternalProperties {
	g.GuID = &guid
	return g
}

// WithInvocations - add a Invocations to the ExternalProperties
func (i *ExternalProperties) WithInvocations(invocations []*Invocation) *ExternalProperties {
	i.Invocations = invocations
	return i
}

// AddInvocation - add a single Invocation to the ExternalProperties
func (i *ExternalProperties) AddInvocation(invocation *Invocation) *ExternalProperties {
	i.Invocations = append(i.Invocations, invocation)
	return i
}

// WithLogicalLocations - add a LogicalLocations to the ExternalProperties
func (l *ExternalProperties) WithLogicalLocations(logicalLocations []*LogicalLocation) *ExternalProperties {
	l.LogicalLocations = logicalLocations
	return l
}

// AddLogicalLocation - add a single LogicalLocation to the ExternalProperties
func (l *ExternalProperties) AddLogicalLocation(logicalLocation *LogicalLocation) *ExternalProperties {
	l.LogicalLocations = append(l.LogicalLocations, logicalLocation)
	return l
}

// WithPolicies - add a Policies to the ExternalProperties
func (p *ExternalProperties) WithPolicies(policies []*ToolComponent) *ExternalProperties {
	p.Policies = policies
	return p
}

// AddPolicie - add a single Policie to the ExternalProperties
func (p *ExternalProperties) AddPolicie(policie *ToolComponent) *ExternalProperties {
	p.Policies = append(p.Policies, policie)
	return p
}

// WithProperties - add a Properties to the ExternalProperties
func (p *ExternalProperties) WithProperties(properties *PropertyBag) *ExternalProperties {
	p.Properties = properties
	return p
}

// WithResults - add a Results to the ExternalProperties
func (r *ExternalProperties) WithResults(results []*Result) *ExternalProperties {
	r.Results = results
	return r
}

// AddResult - add a single Result to the ExternalProperties
func (r *ExternalProperties) AddResult(result *Result) *ExternalProperties {
	r.Results = append(r.Results, result)
	return r
}

// WithRunGuID - add a RunGuID to the ExternalProperties
func (r *ExternalProperties) WithRunGuID(runGuid string) *ExternalProperties {
	r.RunGuID = &runGuid
	return r
}

// WithSchema - add a Schema to the ExternalProperties
func (s *ExternalProperties) WithSchema(schema string) *ExternalProperties {
	s.Schema = &schema
	return s
}

// WithTaxonomies - add a Taxonomies to the ExternalProperties
func (t *ExternalProperties) WithTaxonomies(taxonomies []*ToolComponent) *ExternalProperties {
	t.Taxonomies = taxonomies
	return t
}

// AddTaxonomie - add a single Taxonomie to the ExternalProperties
func (t *ExternalProperties) AddTaxonomie(taxonomie *ToolComponent) *ExternalProperties {
	t.Taxonomies = append(t.Taxonomies, taxonomie)
	return t
}

// WithThreadFlowLocations - add a ThreadFlowLocations to the ExternalProperties
func (t *ExternalProperties) WithThreadFlowLocations(threadFlowLocations []*ThreadFlowLocation) *ExternalProperties {
	t.ThreadFlowLocations = threadFlowLocations
	return t
}

// AddThreadFlowLocation - add a single ThreadFlowLocation to the ExternalProperties
func (t *ExternalProperties) AddThreadFlowLocation(threadFlowLocation *ThreadFlowLocation) *ExternalProperties {
	t.ThreadFlowLocations = append(t.ThreadFlowLocations, threadFlowLocation)
	return t
}

// WithTranslations - add a Translations to the ExternalProperties
func (t *ExternalProperties) WithTranslations(translations []*ToolComponent) *ExternalProperties {
	t.Translations = translations
	return t
}

// AddTranslation - add a single Translation to the ExternalProperties
func (t *ExternalProperties) AddTranslation(translation *ToolComponent) *ExternalProperties {
	t.Translations = append(t.Translations, translation)
	return t
}

// WithVersion - add a Version to the ExternalProperties
func (v *ExternalProperties) WithVersion(version string) *ExternalProperties {
	v.Version = &version
	return v
}

// WithWebRequests - add a WebRequests to the ExternalProperties
func (w *ExternalProperties) WithWebRequests(webRequests []*WebRequest) *ExternalProperties {
	w.WebRequests = webRequests
	return w
}

// AddWebRequest - add a single WebRequest to the ExternalProperties
func (w *ExternalProperties) AddWebRequest(webRequest *WebRequest) *ExternalProperties {
	w.WebRequests = append(w.WebRequests, webRequest)
	return w
}

// WithWebResponses - add a WebResponses to the ExternalProperties
func (w *ExternalProperties) WithWebResponses(webResponses []*WebResponse) *ExternalProperties {
	w.WebResponses = webResponses
	return w
}

// AddWebResponse - add a single WebResponse to the ExternalProperties
func (w *ExternalProperties) AddWebResponse(webResponse *WebResponse) *ExternalProperties {
	w.WebResponses = append(w.WebResponses, webResponse)
	return w
}
