package sarif

// Location - A location within a programming artifact.
type Location struct {
	// A set of regions relevant to the location.
	Annotations []*Region `json:"annotations"`

	// Value that distinguishes this location from all other locations within a single result object.
	ID int `json:"id"`

	// The logical locations associated with the result.
	LogicalLocations []*LogicalLocation `json:"logicalLocations"`

	// A message relevant to the location.
	Message *Message `json:"message,omitempty"`

	// Identifies the artifact and region.
	PhysicalLocation *PhysicalLocation `json:"physicalLocation,omitempty"`

	// Key/value pairs that provide additional information about the location.
	Properties *PropertyBag `json:"properties,omitempty"`

	// An array of objects that describe relationships between this location and others.
	Relationships []*LocationRelationship `json:"relationships"`
}

// NewLocation - creates a new
func NewLocation() *Location {
	return &Location{
		Annotations:      make([]*Region, 0),
		ID:               -1,
		LogicalLocations: make([]*LogicalLocation, 0),
		Relationships:    make([]*LocationRelationship, 0),
	}
}

// WithAnnotations - add a Annotations to the Location
func (a *Location) WithAnnotations(annotations []*Region) *Location {
	a.Annotations = annotations
	return a
}

// AddAnnotation - add a single Annotation to the Location
func (a *Location) AddAnnotation(annotation *Region) *Location {
	a.Annotations = append(a.Annotations, annotation)
	return a
}

// WithID - add a ID to the Location
func (i *Location) WithID(id int) *Location {
	i.ID = id
	return i
}

// WithLogicalLocations - add a LogicalLocations to the Location
func (l *Location) WithLogicalLocations(logicalLocations []*LogicalLocation) *Location {
	l.LogicalLocations = logicalLocations
	return l
}

// AddLogicalLocation - add a single LogicalLocation to the Location
func (l *Location) AddLogicalLocation(logicalLocation *LogicalLocation) *Location {
	l.LogicalLocations = append(l.LogicalLocations, logicalLocation)
	return l
}

// WithMessage - add a Message to the Location
func (m *Location) WithMessage(message *Message) *Location {
	m.Message = message
	return m
}

// WithPhysicalLocation - add a PhysicalLocation to the Location
func (p *Location) WithPhysicalLocation(physicalLocation *PhysicalLocation) *Location {
	p.PhysicalLocation = physicalLocation
	return p
}

// WithProperties - add a Properties to the Location
func (p *Location) WithProperties(properties *PropertyBag) *Location {
	p.Properties = properties
	return p
}

// WithRelationships - add a Relationships to the Location
func (r *Location) WithRelationships(relationships []*LocationRelationship) *Location {
	r.Relationships = relationships
	return r
}

// AddRelationship - add a single Relationship to the Location
func (r *Location) AddRelationship(relationship *LocationRelationship) *Location {
	r.Relationships = append(r.Relationships, relationship)
	return r
}
