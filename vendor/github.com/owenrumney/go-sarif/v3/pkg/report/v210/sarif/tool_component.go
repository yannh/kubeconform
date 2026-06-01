package sarif

// ToolComponent - A component, such as a plug-in or the driver, of the analysis tool that was run.
type ToolComponent struct {
	// A dictionary, each of whose keys is a resource identifier and each of whose values is a multiformatMessageString object, which holds message strings in plain text and (optionally) Markdown format. The strings can include placeholders, which can be used to construct a message in combination with an arbitrary number of additional string arguments.
	GlobalMessageStrings map[string]MultiformatMessageString `json:"globalMessageStrings,omitempty"`

	// The component which is strongly associated with this component. For a translation, this refers to the component which has been translated. For an extension, this is the driver that provides the extension's plugin model.
	AssociatedComponent *ToolComponentReference `json:"associatedComponent,omitempty"`

	// The kinds of data contained in this object.
	Contents []string `json:"contents"`

	// The binary version of the tool component's primary executable file expressed as four non-negative integers separated by a period (for operating systems that express file versions in this way).
	DottedQuadFileVersion *string `json:"dottedQuadFileVersion,omitempty"`

	// The absolute URI from which the tool component can be downloaded.
	DownloadURI *string `json:"downloadUri,omitempty"`

	// A comprehensive description of the tool component.
	FullDescription *MultiformatMessageString `json:"fullDescription,omitempty"`

	// The name of the tool component along with its version and any other useful identifying information, such as its locale.
	FullName *string `json:"fullName,omitempty"`

	// A unique identifier for the tool component in the form of a GUID.
	GuID *string `json:"guid,omitempty"`

	// The absolute URI at which information about this version of the tool component can be found.
	InformationURI *string `json:"informationUri,omitempty"`

	// Specifies whether this object contains a complete definition of the localizable and/or non-localizable data for this component, as opposed to including only data that is relevant to the results persisted to this log file.
	IsComprehensive bool `json:"isComprehensive"`

	// The language of the messages emitted into the log file during this run (expressed as an ISO 639-1 two-letter lowercase language code) and an optional region (expressed as an ISO 3166-1 two-letter uppercase subculture code associated with a country or region). The casing is recommended but not required (in order for this data to conform to RFC5646).
	Language string `json:"language"`

	// The semantic version of the localized strings defined in this component; maintained by components that provide translations.
	LocalizedDataSemanticVersion *string `json:"localizedDataSemanticVersion,omitempty"`

	// An array of the artifactLocation objects associated with the tool component.
	Locations []*ArtifactLocation `json:"locations"`

	// The minimum value of localizedDataSemanticVersion required in translations consumed by this component; used by components that consume translations.
	MinimumRequiredLocalizedDataSemanticVersion *string `json:"minimumRequiredLocalizedDataSemanticVersion,omitempty"`

	// The name of the tool component.
	Name *string `json:"name,omitempty"`

	// An array of reportingDescriptor objects relevant to the notifications related to the configuration and runtime execution of the tool component.
	Notifications []*ReportingDescriptor `json:"notifications"`

	// The organization or company that produced the tool component.
	Organization *string `json:"organization,omitempty"`

	// A product suite to which the tool component belongs.
	Product *string `json:"product,omitempty"`

	// A localizable string containing the name of the suite of products to which the tool component belongs.
	ProductSuite *string `json:"productSuite,omitempty"`

	// Key/value pairs that provide additional information about the tool component.
	Properties *PropertyBag `json:"properties,omitempty"`

	// A string specifying the UTC date (and optionally, the time) of the component's release.
	ReleaseDateUtc *string `json:"releaseDateUtc,omitempty"`

	// An array of reportingDescriptor objects relevant to the analysis performed by the tool component.
	Rules []*ReportingDescriptor `json:"rules"`

	// The tool component version in the format specified by Semantic Versioning 2.0.
	SemanticVersion *string `json:"semanticVersion,omitempty"`

	// A brief description of the tool component.
	ShortDescription *MultiformatMessageString `json:"shortDescription,omitempty"`

	// An array of toolComponentReference objects to declare the taxonomies supported by the tool component.
	SupportedTaxonomies []*ToolComponentReference `json:"supportedTaxonomies"`

	// An array of reportingDescriptor objects relevant to the definitions of both standalone and tool-defined taxonomies.
	Taxa []*ReportingDescriptor `json:"taxa"`

	// Translation metadata, required for a translation, not populated by other component types.
	TranslationMetadata *TranslationMetadata `json:"translationMetadata,omitempty"`

	// The tool component version, in whatever format the component natively provides.
	Version *string `json:"version,omitempty"`
}

// NewToolComponent - creates a new
func NewToolComponent() *ToolComponent {
	return &ToolComponent{
		Contents:            []string{"localizedData", "nonLocalizedData"},
		IsComprehensive:     false,
		Language:            "en-US",
		Locations:           make([]*ArtifactLocation, 0),
		Notifications:       make([]*ReportingDescriptor, 0),
		Rules:               make([]*ReportingDescriptor, 0),
		SupportedTaxonomies: make([]*ToolComponentReference, 0),
		Taxa:                make([]*ReportingDescriptor, 0),
	}
}

// AddGlobalMessageString - add a single GlobalMessageString to the ToolComponent
func (g *ToolComponent) AddGlobalMessageString(key string, globalMessageString MultiformatMessageString) *ToolComponent {
	g.GlobalMessageStrings[key] = globalMessageString
	return g
}

// WithGlobalMessageStrings - add a GlobalMessageStrings to the ToolComponent
func (g *ToolComponent) WithGlobalMessageStrings(globalMessageStrings map[string]MultiformatMessageString) *ToolComponent {
	g.GlobalMessageStrings = globalMessageStrings
	return g
}

// WithAssociatedComponent - add a AssociatedComponent to the ToolComponent
func (a *ToolComponent) WithAssociatedComponent(associatedComponent *ToolComponentReference) *ToolComponent {
	a.AssociatedComponent = associatedComponent
	return a
}

// WithContents - add a Contents to the ToolComponent
func (c *ToolComponent) WithContents(contents []string) *ToolComponent {
	c.Contents = contents
	return c
}

// AddContent - add a single Content to the ToolComponent
func (c *ToolComponent) AddContent(content string) *ToolComponent {
	c.Contents = append(c.Contents, content)
	return c
}

// WithDottedQuadFileVersion - add a DottedQuadFileVersion to the ToolComponent
func (d *ToolComponent) WithDottedQuadFileVersion(dottedQuadFileVersion string) *ToolComponent {
	d.DottedQuadFileVersion = &dottedQuadFileVersion
	return d
}

// WithDownloadURI - add a DownloadURI to the ToolComponent
func (d *ToolComponent) WithDownloadURI(downloadUri string) *ToolComponent {
	d.DownloadURI = &downloadUri
	return d
}

// WithFullDescription - add a FullDescription to the ToolComponent
func (f *ToolComponent) WithFullDescription(fullDescription *MultiformatMessageString) *ToolComponent {
	f.FullDescription = fullDescription
	return f
}

// WithFullName - add a FullName to the ToolComponent
func (f *ToolComponent) WithFullName(fullName string) *ToolComponent {
	f.FullName = &fullName
	return f
}

// WithGuID - add a GuID to the ToolComponent
func (g *ToolComponent) WithGuID(guid string) *ToolComponent {
	g.GuID = &guid
	return g
}

// WithInformationURI - add a InformationURI to the ToolComponent
func (i *ToolComponent) WithInformationURI(informationUri string) *ToolComponent {
	i.InformationURI = &informationUri
	return i
}

// WithIsComprehensive - add a IsComprehensive to the ToolComponent
func (i *ToolComponent) WithIsComprehensive(isComprehensive bool) *ToolComponent {
	i.IsComprehensive = isComprehensive
	return i
}

// WithLanguage - add a Language to the ToolComponent
func (l *ToolComponent) WithLanguage(language string) *ToolComponent {
	l.Language = language
	return l
}

// WithLocalizedDataSemanticVersion - add a LocalizedDataSemanticVersion to the ToolComponent
func (l *ToolComponent) WithLocalizedDataSemanticVersion(localizedDataSemanticVersion string) *ToolComponent {
	l.LocalizedDataSemanticVersion = &localizedDataSemanticVersion
	return l
}

// WithLocations - add a Locations to the ToolComponent
func (l *ToolComponent) WithLocations(locations []*ArtifactLocation) *ToolComponent {
	l.Locations = locations
	return l
}

// AddLocation - add a single Location to the ToolComponent
func (l *ToolComponent) AddLocation(location *ArtifactLocation) *ToolComponent {
	l.Locations = append(l.Locations, location)
	return l
}

// WithMinimumRequiredLocalizedDataSemanticVersion - add a MinimumRequiredLocalizedDataSemanticVersion to the ToolComponent
func (m *ToolComponent) WithMinimumRequiredLocalizedDataSemanticVersion(minimumRequiredLocalizedDataSemanticVersion string) *ToolComponent {
	m.MinimumRequiredLocalizedDataSemanticVersion = &minimumRequiredLocalizedDataSemanticVersion
	return m
}

// WithName - add a Name to the ToolComponent
func (n *ToolComponent) WithName(name string) *ToolComponent {
	n.Name = &name
	return n
}

// WithNotifications - add a Notifications to the ToolComponent
func (n *ToolComponent) WithNotifications(notifications []*ReportingDescriptor) *ToolComponent {
	n.Notifications = notifications
	return n
}

// AddNotification - add a single Notification to the ToolComponent
func (n *ToolComponent) AddNotification(notification *ReportingDescriptor) *ToolComponent {
	n.Notifications = append(n.Notifications, notification)
	return n
}

// WithOrganization - add a Organization to the ToolComponent
func (o *ToolComponent) WithOrganization(organization string) *ToolComponent {
	o.Organization = &organization
	return o
}

// WithProduct - add a Product to the ToolComponent
func (p *ToolComponent) WithProduct(product string) *ToolComponent {
	p.Product = &product
	return p
}

// WithProductSuite - add a ProductSuite to the ToolComponent
func (p *ToolComponent) WithProductSuite(productSuite string) *ToolComponent {
	p.ProductSuite = &productSuite
	return p
}

// WithProperties - add a Properties to the ToolComponent
func (p *ToolComponent) WithProperties(properties *PropertyBag) *ToolComponent {
	p.Properties = properties
	return p
}

// WithReleaseDateUtc - add a ReleaseDateUtc to the ToolComponent
func (r *ToolComponent) WithReleaseDateUtc(releaseDateUtc string) *ToolComponent {
	r.ReleaseDateUtc = &releaseDateUtc
	return r
}

// WithRules - add a Rules to the ToolComponent
func (r *ToolComponent) WithRules(rules []*ReportingDescriptor) *ToolComponent {
	r.Rules = rules
	return r
}

// AddRule - add a single Rule to the ToolComponent
func (r *ToolComponent) AddRule(rule *ReportingDescriptor) *ToolComponent {
	r.Rules = append(r.Rules, rule)
	return r
}

// WithSemanticVersion - add a SemanticVersion to the ToolComponent
func (s *ToolComponent) WithSemanticVersion(semanticVersion string) *ToolComponent {
	s.SemanticVersion = &semanticVersion
	return s
}

// WithShortDescription - add a ShortDescription to the ToolComponent
func (s *ToolComponent) WithShortDescription(shortDescription *MultiformatMessageString) *ToolComponent {
	s.ShortDescription = shortDescription
	return s
}

// WithSupportedTaxonomies - add a SupportedTaxonomies to the ToolComponent
func (s *ToolComponent) WithSupportedTaxonomies(supportedTaxonomies []*ToolComponentReference) *ToolComponent {
	s.SupportedTaxonomies = supportedTaxonomies
	return s
}

// AddSupportedTaxonomie - add a single SupportedTaxonomie to the ToolComponent
func (s *ToolComponent) AddSupportedTaxonomie(supportedTaxonomie *ToolComponentReference) *ToolComponent {
	s.SupportedTaxonomies = append(s.SupportedTaxonomies, supportedTaxonomie)
	return s
}

// WithTaxa - add a Taxa to the ToolComponent
func (t *ToolComponent) WithTaxa(taxa []*ReportingDescriptor) *ToolComponent {
	t.Taxa = taxa
	return t
}

// AddTaxa - add a single Taxa to the ToolComponent
func (t *ToolComponent) AddTaxa(taxa *ReportingDescriptor) *ToolComponent {
	t.Taxa = append(t.Taxa, taxa)
	return t
}

// WithTranslationMetadata - add a TranslationMetadata to the ToolComponent
func (t *ToolComponent) WithTranslationMetadata(translationMetadata *TranslationMetadata) *ToolComponent {
	t.TranslationMetadata = translationMetadata
	return t
}

// WithVersion - add a Version to the ToolComponent
func (v *ToolComponent) WithVersion(version string) *ToolComponent {
	v.Version = &version
	return v
}
