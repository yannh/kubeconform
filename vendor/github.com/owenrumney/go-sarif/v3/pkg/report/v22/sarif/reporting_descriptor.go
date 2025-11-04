package sarif

// ReportingDescriptor - Metadata that describes a specific report produced by the tool, as part of the analysis it provides or its runtime reporting.
type ReportingDescriptor struct {
	// A set of name/value pairs with arbitrary names. Each value is a multiformatMessageString object, which holds message strings in plain text and (optionally) Markdown format. The strings can include placeholders, which can be used to construct a message in combination with an arbitrary number of additional string arguments.
	MessageStrings map[string]MultiformatMessageString `json:"messageStrings,omitempty"`

	// Default reporting configuration information.
	DefaultConfiguration *ReportingConfiguration `json:"defaultConfiguration,omitempty"`

	// An array of unique identifies in the form of a GUID by which this report was known in some previous version of the analysis tool.
	DeprecatedGuids []string `json:"deprecatedGuids,omitempty"`

	// An array of stable, opaque identifiers by which this report was known in some previous version of the analysis tool.
	DeprecatedIds []string `json:"deprecatedIds,omitempty"`

	// An array of readable identifiers by which this report was known in some previous version of the analysis tool.
	DeprecatedNames []string `json:"deprecatedNames,omitempty"`

	// A description of the report. Should, as far as possible, provide details sufficient to enable resolution of any problem indicated by the result.
	FullDescription *MultiformatMessageString `json:"fullDescription,omitempty"`

	// A unique identifier for the reporting descriptor in the form of a GUID.
	Guid *Guid `json:"guid,omitempty"`

	// Provides the primary documentation for the report, useful when there is no online documentation.
	Help *MultiformatMessageString `json:"help,omitempty"`

	// A URI where the primary documentation for the report can be found.
	HelpURI *string `json:"helpUri,omitempty"`

	// A stable, opaque identifier for the report.
	ID *string `json:"id,omitempty"`

	// A report identifier that is understandable to an end user.
	Name *string `json:"name,omitempty"`

	// Key/value pairs that provide additional information about the report.
	Properties *PropertyBag `json:"properties,omitempty"`

	// An array of objects that describe relationships between this reporting descriptor and others.
	Relationships []*ReportingDescriptorRelationship `json:"relationships"`

	// A concise description of the report. Should be a single sentence that is understandable when visible space is limited to a single line of text.
	ShortDescription *MultiformatMessageString `json:"shortDescription,omitempty"`
}

// NewReportingDescriptor - creates a new
func NewReportingDescriptor() *ReportingDescriptor {
	return &ReportingDescriptor{
		DeprecatedGuids: make([]string, 0),
		DeprecatedIds:   make([]string, 0),
		DeprecatedNames: make([]string, 0),
		Relationships:   make([]*ReportingDescriptorRelationship, 0),
	}
}

// AddMessageString - add a single MessageString to the ReportingDescriptor
func (m *ReportingDescriptor) AddMessageString(key string, messageString MultiformatMessageString) *ReportingDescriptor {
	m.MessageStrings[key] = messageString
	return m
}

// WithMessageStrings - add a MessageStrings to the ReportingDescriptor
func (m *ReportingDescriptor) WithMessageStrings(messageStrings map[string]MultiformatMessageString) *ReportingDescriptor {
	m.MessageStrings = messageStrings
	return m
}

// WithDefaultConfiguration - add a DefaultConfiguration to the ReportingDescriptor
func (d *ReportingDescriptor) WithDefaultConfiguration(defaultConfiguration *ReportingConfiguration) *ReportingDescriptor {
	d.DefaultConfiguration = defaultConfiguration
	return d
}

// WithDeprecatedGuids - add a DeprecatedGuids to the ReportingDescriptor
func (d *ReportingDescriptor) WithDeprecatedGuids(deprecatedGuids []string) *ReportingDescriptor {
	d.DeprecatedGuids = deprecatedGuids
	return d
}

// AddDeprecatedGuid - add a single DeprecatedGuid to the ReportingDescriptor
func (d *ReportingDescriptor) AddDeprecatedGuid(deprecatedGuid string) *ReportingDescriptor {
	d.DeprecatedGuids = append(d.DeprecatedGuids, deprecatedGuid)
	return d
}

// WithDeprecatedIds - add a DeprecatedIds to the ReportingDescriptor
func (d *ReportingDescriptor) WithDeprecatedIds(deprecatedIds []string) *ReportingDescriptor {
	d.DeprecatedIds = deprecatedIds
	return d
}

// AddDeprecatedId - add a single DeprecatedId to the ReportingDescriptor
func (d *ReportingDescriptor) AddDeprecatedId(deprecatedId string) *ReportingDescriptor {
	d.DeprecatedIds = append(d.DeprecatedIds, deprecatedId)
	return d
}

// WithDeprecatedNames - add a DeprecatedNames to the ReportingDescriptor
func (d *ReportingDescriptor) WithDeprecatedNames(deprecatedNames []string) *ReportingDescriptor {
	d.DeprecatedNames = deprecatedNames
	return d
}

// AddDeprecatedName - add a single DeprecatedName to the ReportingDescriptor
func (d *ReportingDescriptor) AddDeprecatedName(deprecatedName string) *ReportingDescriptor {
	d.DeprecatedNames = append(d.DeprecatedNames, deprecatedName)
	return d
}

// WithFullDescription - add a FullDescription to the ReportingDescriptor
func (f *ReportingDescriptor) WithFullDescription(fullDescription *MultiformatMessageString) *ReportingDescriptor {
	f.FullDescription = fullDescription
	return f
}

// WithGuid - add a Guid to the ReportingDescriptor
func (g *ReportingDescriptor) WithGuid(guid *Guid) *ReportingDescriptor {
	g.Guid = guid
	return g
}

// WithHelp - add a Help to the ReportingDescriptor
func (h *ReportingDescriptor) WithHelp(help *MultiformatMessageString) *ReportingDescriptor {
	h.Help = help
	return h
}

// WithHelpURI - add a HelpURI to the ReportingDescriptor
func (h *ReportingDescriptor) WithHelpURI(helpUri string) *ReportingDescriptor {
	h.HelpURI = &helpUri
	return h
}

// WithID - add a ID to the ReportingDescriptor
func (i *ReportingDescriptor) WithID(id string) *ReportingDescriptor {
	i.ID = &id
	return i
}

// WithName - add a Name to the ReportingDescriptor
func (n *ReportingDescriptor) WithName(name string) *ReportingDescriptor {
	n.Name = &name
	return n
}

// WithProperties - add a Properties to the ReportingDescriptor
func (p *ReportingDescriptor) WithProperties(properties *PropertyBag) *ReportingDescriptor {
	p.Properties = properties
	return p
}

// WithRelationships - add a Relationships to the ReportingDescriptor
func (r *ReportingDescriptor) WithRelationships(relationships []*ReportingDescriptorRelationship) *ReportingDescriptor {
	r.Relationships = relationships
	return r
}

// AddRelationship - add a single Relationship to the ReportingDescriptor
func (r *ReportingDescriptor) AddRelationship(relationship *ReportingDescriptorRelationship) *ReportingDescriptor {
	r.Relationships = append(r.Relationships, relationship)
	return r
}

// WithShortDescription - add a ShortDescription to the ReportingDescriptor
func (s *ReportingDescriptor) WithShortDescription(shortDescription *MultiformatMessageString) *ReportingDescriptor {
	s.ShortDescription = shortDescription
	return s
}
