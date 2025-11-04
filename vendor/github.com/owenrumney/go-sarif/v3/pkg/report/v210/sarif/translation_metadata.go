package sarif

// TranslationMetadata - Provides additional metadata related to translation.
type TranslationMetadata struct {
	// The absolute URI from which the translation metadata can be downloaded.
	DownloadURI *string `json:"downloadUri,omitempty"`

	// A comprehensive description of the translation metadata.
	FullDescription *MultiformatMessageString `json:"fullDescription,omitempty"`

	// The full name associated with the translation metadata.
	FullName *string `json:"fullName,omitempty"`

	// The absolute URI from which information related to the translation metadata can be downloaded.
	InformationURI *string `json:"informationUri,omitempty"`

	// The name associated with the translation metadata.
	Name *string `json:"name,omitempty"`

	// Key/value pairs that provide additional information about the translation metadata.
	Properties *PropertyBag `json:"properties,omitempty"`

	// A brief description of the translation metadata.
	ShortDescription *MultiformatMessageString `json:"shortDescription,omitempty"`
}

// NewTranslationMetadata - creates a new
func NewTranslationMetadata() *TranslationMetadata {
	return &TranslationMetadata{}
}

// WithDownloadURI - add a DownloadURI to the TranslationMetadata
func (d *TranslationMetadata) WithDownloadURI(downloadUri string) *TranslationMetadata {
	d.DownloadURI = &downloadUri
	return d
}

// WithFullDescription - add a FullDescription to the TranslationMetadata
func (f *TranslationMetadata) WithFullDescription(fullDescription *MultiformatMessageString) *TranslationMetadata {
	f.FullDescription = fullDescription
	return f
}

// WithFullName - add a FullName to the TranslationMetadata
func (f *TranslationMetadata) WithFullName(fullName string) *TranslationMetadata {
	f.FullName = &fullName
	return f
}

// WithInformationURI - add a InformationURI to the TranslationMetadata
func (i *TranslationMetadata) WithInformationURI(informationUri string) *TranslationMetadata {
	i.InformationURI = &informationUri
	return i
}

// WithName - add a Name to the TranslationMetadata
func (n *TranslationMetadata) WithName(name string) *TranslationMetadata {
	n.Name = &name
	return n
}

// WithProperties - add a Properties to the TranslationMetadata
func (p *TranslationMetadata) WithProperties(properties *PropertyBag) *TranslationMetadata {
	p.Properties = properties
	return p
}

// WithShortDescription - add a ShortDescription to the TranslationMetadata
func (s *TranslationMetadata) WithShortDescription(shortDescription *MultiformatMessageString) *TranslationMetadata {
	s.ShortDescription = shortDescription
	return s
}
