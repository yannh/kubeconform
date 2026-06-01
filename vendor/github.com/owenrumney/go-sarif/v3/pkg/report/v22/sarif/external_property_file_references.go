package sarif

// ExternalPropertyFileReferences - References to external property files that should be inlined with the content of a root log file.
type ExternalPropertyFileReferences struct {
	// An array of external property files containing run.addresses arrays to be merged with the root log file.
	Addresses []*ExternalPropertyFileReference `json:"addresses"`

	// An array of external property files containing run.artifacts arrays to be merged with the root log file.
	Artifacts []*ExternalPropertyFileReference `json:"artifacts"`

	// An external property file containing a run.conversion object to be merged with the root log file.
	Conversion *ExternalPropertyFileReference `json:"conversion,omitempty"`

	// An external property file containing a run.driver object to be merged with the root log file.
	Driver *ExternalPropertyFileReference `json:"driver,omitempty"`

	// An array of external property files containing run.extensions arrays to be merged with the root log file.
	Extensions []*ExternalPropertyFileReference `json:"extensions"`

	// An external property file containing a run.properties object to be merged with the root log file.
	ExternalizedProperties *ExternalPropertyFileReference `json:"externalizedProperties,omitempty"`

	// An array of external property files containing a run.graphs object to be merged with the root log file.
	Graphs []*ExternalPropertyFileReference `json:"graphs"`

	// An array of external property files containing run.invocations arrays to be merged with the root log file.
	Invocations []*ExternalPropertyFileReference `json:"invocations"`

	// An array of external property files containing run.logicalLocations arrays to be merged with the root log file.
	LogicalLocations []*ExternalPropertyFileReference `json:"logicalLocations"`

	// An array of external property files containing run.policies arrays to be merged with the root log file.
	Policies []*ExternalPropertyFileReference `json:"policies"`

	// Key/value pairs that provide additional information about the external property files.
	Properties *PropertyBag `json:"properties,omitempty"`

	// An array of external property files containing run.results arrays to be merged with the root log file.
	Results []*ExternalPropertyFileReference `json:"results"`

	// An array of external property files containing run.taxonomies arrays to be merged with the root log file.
	Taxonomies []*ExternalPropertyFileReference `json:"taxonomies"`

	// An array of external property files containing run.threadFlowLocations arrays to be merged with the root log file.
	ThreadFlowLocations []*ExternalPropertyFileReference `json:"threadFlowLocations"`

	// An array of external property files containing run.translations arrays to be merged with the root log file.
	Translations []*ExternalPropertyFileReference `json:"translations"`

	// An array of external property files containing run.requests arrays to be merged with the root log file.
	WebRequests []*ExternalPropertyFileReference `json:"webRequests"`

	// An array of external property files containing run.responses arrays to be merged with the root log file.
	WebResponses []*ExternalPropertyFileReference `json:"webResponses"`
}

// NewExternalPropertyFileReferences - creates a new
func NewExternalPropertyFileReferences() *ExternalPropertyFileReferences {
	return &ExternalPropertyFileReferences{
		Addresses:           make([]*ExternalPropertyFileReference, 0),
		Artifacts:           make([]*ExternalPropertyFileReference, 0),
		Extensions:          make([]*ExternalPropertyFileReference, 0),
		Graphs:              make([]*ExternalPropertyFileReference, 0),
		Invocations:         make([]*ExternalPropertyFileReference, 0),
		LogicalLocations:    make([]*ExternalPropertyFileReference, 0),
		Policies:            make([]*ExternalPropertyFileReference, 0),
		Results:             make([]*ExternalPropertyFileReference, 0),
		Taxonomies:          make([]*ExternalPropertyFileReference, 0),
		ThreadFlowLocations: make([]*ExternalPropertyFileReference, 0),
		Translations:        make([]*ExternalPropertyFileReference, 0),
		WebRequests:         make([]*ExternalPropertyFileReference, 0),
		WebResponses:        make([]*ExternalPropertyFileReference, 0),
	}
}

// WithAddresses - add a Addresses to the ExternalPropertyFileReferences
func (a *ExternalPropertyFileReferences) WithAddresses(addresses []*ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	a.Addresses = addresses
	return a
}

// AddAddresse - add a single Addresse to the ExternalPropertyFileReferences
func (a *ExternalPropertyFileReferences) AddAddresse(addresse *ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	a.Addresses = append(a.Addresses, addresse)
	return a
}

// WithArtifacts - add a Artifacts to the ExternalPropertyFileReferences
func (a *ExternalPropertyFileReferences) WithArtifacts(artifacts []*ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	a.Artifacts = artifacts
	return a
}

// AddArtifact - add a single Artifact to the ExternalPropertyFileReferences
func (a *ExternalPropertyFileReferences) AddArtifact(artifact *ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	a.Artifacts = append(a.Artifacts, artifact)
	return a
}

// WithConversion - add a Conversion to the ExternalPropertyFileReferences
func (c *ExternalPropertyFileReferences) WithConversion(conversion *ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	c.Conversion = conversion
	return c
}

// WithDriver - add a Driver to the ExternalPropertyFileReferences
func (d *ExternalPropertyFileReferences) WithDriver(driver *ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	d.Driver = driver
	return d
}

// WithExtensions - add a Extensions to the ExternalPropertyFileReferences
func (e *ExternalPropertyFileReferences) WithExtensions(extensions []*ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	e.Extensions = extensions
	return e
}

// AddExtension - add a single Extension to the ExternalPropertyFileReferences
func (e *ExternalPropertyFileReferences) AddExtension(extension *ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	e.Extensions = append(e.Extensions, extension)
	return e
}

// WithExternalizedProperties - add a ExternalizedProperties to the ExternalPropertyFileReferences
func (e *ExternalPropertyFileReferences) WithExternalizedProperties(externalizedProperties *ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	e.ExternalizedProperties = externalizedProperties
	return e
}

// WithGraphs - add a Graphs to the ExternalPropertyFileReferences
func (g *ExternalPropertyFileReferences) WithGraphs(graphs []*ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	g.Graphs = graphs
	return g
}

// AddGraph - add a single Graph to the ExternalPropertyFileReferences
func (g *ExternalPropertyFileReferences) AddGraph(graph *ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	g.Graphs = append(g.Graphs, graph)
	return g
}

// WithInvocations - add a Invocations to the ExternalPropertyFileReferences
func (i *ExternalPropertyFileReferences) WithInvocations(invocations []*ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	i.Invocations = invocations
	return i
}

// AddInvocation - add a single Invocation to the ExternalPropertyFileReferences
func (i *ExternalPropertyFileReferences) AddInvocation(invocation *ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	i.Invocations = append(i.Invocations, invocation)
	return i
}

// WithLogicalLocations - add a LogicalLocations to the ExternalPropertyFileReferences
func (l *ExternalPropertyFileReferences) WithLogicalLocations(logicalLocations []*ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	l.LogicalLocations = logicalLocations
	return l
}

// AddLogicalLocation - add a single LogicalLocation to the ExternalPropertyFileReferences
func (l *ExternalPropertyFileReferences) AddLogicalLocation(logicalLocation *ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	l.LogicalLocations = append(l.LogicalLocations, logicalLocation)
	return l
}

// WithPolicies - add a Policies to the ExternalPropertyFileReferences
func (p *ExternalPropertyFileReferences) WithPolicies(policies []*ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	p.Policies = policies
	return p
}

// AddPolicie - add a single Policie to the ExternalPropertyFileReferences
func (p *ExternalPropertyFileReferences) AddPolicie(policie *ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	p.Policies = append(p.Policies, policie)
	return p
}

// WithProperties - add a Properties to the ExternalPropertyFileReferences
func (p *ExternalPropertyFileReferences) WithProperties(properties *PropertyBag) *ExternalPropertyFileReferences {
	p.Properties = properties
	return p
}

// WithResults - add a Results to the ExternalPropertyFileReferences
func (r *ExternalPropertyFileReferences) WithResults(results []*ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	r.Results = results
	return r
}

// AddResult - add a single Result to the ExternalPropertyFileReferences
func (r *ExternalPropertyFileReferences) AddResult(result *ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	r.Results = append(r.Results, result)
	return r
}

// WithTaxonomies - add a Taxonomies to the ExternalPropertyFileReferences
func (t *ExternalPropertyFileReferences) WithTaxonomies(taxonomies []*ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	t.Taxonomies = taxonomies
	return t
}

// AddTaxonomie - add a single Taxonomie to the ExternalPropertyFileReferences
func (t *ExternalPropertyFileReferences) AddTaxonomie(taxonomie *ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	t.Taxonomies = append(t.Taxonomies, taxonomie)
	return t
}

// WithThreadFlowLocations - add a ThreadFlowLocations to the ExternalPropertyFileReferences
func (t *ExternalPropertyFileReferences) WithThreadFlowLocations(threadFlowLocations []*ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	t.ThreadFlowLocations = threadFlowLocations
	return t
}

// AddThreadFlowLocation - add a single ThreadFlowLocation to the ExternalPropertyFileReferences
func (t *ExternalPropertyFileReferences) AddThreadFlowLocation(threadFlowLocation *ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	t.ThreadFlowLocations = append(t.ThreadFlowLocations, threadFlowLocation)
	return t
}

// WithTranslations - add a Translations to the ExternalPropertyFileReferences
func (t *ExternalPropertyFileReferences) WithTranslations(translations []*ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	t.Translations = translations
	return t
}

// AddTranslation - add a single Translation to the ExternalPropertyFileReferences
func (t *ExternalPropertyFileReferences) AddTranslation(translation *ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	t.Translations = append(t.Translations, translation)
	return t
}

// WithWebRequests - add a WebRequests to the ExternalPropertyFileReferences
func (w *ExternalPropertyFileReferences) WithWebRequests(webRequests []*ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	w.WebRequests = webRequests
	return w
}

// AddWebRequest - add a single WebRequest to the ExternalPropertyFileReferences
func (w *ExternalPropertyFileReferences) AddWebRequest(webRequest *ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	w.WebRequests = append(w.WebRequests, webRequest)
	return w
}

// WithWebResponses - add a WebResponses to the ExternalPropertyFileReferences
func (w *ExternalPropertyFileReferences) WithWebResponses(webResponses []*ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	w.WebResponses = webResponses
	return w
}

// AddWebResponse - add a single WebResponse to the ExternalPropertyFileReferences
func (w *ExternalPropertyFileReferences) AddWebResponse(webResponse *ExternalPropertyFileReference) *ExternalPropertyFileReferences {
	w.WebResponses = append(w.WebResponses, webResponse)
	return w
}
