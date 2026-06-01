package sarif

// VersionControlDetails - Specifies the information necessary to retrieve a desired revision from a version control system.
type VersionControlDetails struct {
	// A Coordinated Universal Time (UTC) date and time that can be used to synchronize an enlistment to the state of the repository at that time.
	AsOfTimeUtc *string `json:"asOfTimeUtc,omitempty"`

	// The name of a branch containing the revision.
	Branch *string `json:"branch,omitempty"`

	// The location in the local file system to which the root of the repository was mapped at the time of the analysis.
	MappedTo *ArtifactLocation `json:"mappedTo,omitempty"`

	// Key/value pairs that provide additional information about the version control details.
	Properties *PropertyBag `json:"properties,omitempty"`

	// The absolute URI of the repository.
	RepositoryURI *string `json:"repositoryUri,omitempty"`

	// A string that uniquely and permanently identifies the revision within the repository.
	RevisionID *string `json:"revisionId,omitempty"`

	// A tag that has been applied to the revision.
	RevisionTag *string `json:"revisionTag,omitempty"`
}

// NewVersionControlDetails - creates a new
func NewVersionControlDetails() *VersionControlDetails {
	return &VersionControlDetails{}
}

// WithAsOfTimeUtc - add a AsOfTimeUtc to the VersionControlDetails
func (a *VersionControlDetails) WithAsOfTimeUtc(asOfTimeUtc string) *VersionControlDetails {
	a.AsOfTimeUtc = &asOfTimeUtc
	return a
}

// WithBranch - add a Branch to the VersionControlDetails
func (b *VersionControlDetails) WithBranch(branch string) *VersionControlDetails {
	b.Branch = &branch
	return b
}

// WithMappedTo - add a MappedTo to the VersionControlDetails
func (m *VersionControlDetails) WithMappedTo(mappedTo *ArtifactLocation) *VersionControlDetails {
	m.MappedTo = mappedTo
	return m
}

// WithProperties - add a Properties to the VersionControlDetails
func (p *VersionControlDetails) WithProperties(properties *PropertyBag) *VersionControlDetails {
	p.Properties = properties
	return p
}

// WithRepositoryURI - add a RepositoryURI to the VersionControlDetails
func (r *VersionControlDetails) WithRepositoryURI(repositoryUri string) *VersionControlDetails {
	r.RepositoryURI = &repositoryUri
	return r
}

// WithRevisionID - add a RevisionID to the VersionControlDetails
func (r *VersionControlDetails) WithRevisionID(revisionId string) *VersionControlDetails {
	r.RevisionID = &revisionId
	return r
}

// WithRevisionTag - add a RevisionTag to the VersionControlDetails
func (r *VersionControlDetails) WithRevisionTag(revisionTag string) *VersionControlDetails {
	r.RevisionTag = &revisionTag
	return r
}
