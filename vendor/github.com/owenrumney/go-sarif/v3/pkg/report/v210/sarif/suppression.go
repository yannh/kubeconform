package sarif

// Suppression - A suppression that is relevant to a result.
type Suppression struct {
	// A stable, unique identifier for the suprression in the form of a GUID.
	GuID *string `json:"guid,omitempty"`

	// A string representing the justification for the suppression.
	Justification *string `json:"justification,omitempty"`

	// A string that indicates where the suppression is persisted.
	Kind *string `json:"kind,omitempty"`

	// Identifies the location associated with the suppression.
	Location *Location `json:"location,omitempty"`

	// Key/value pairs that provide additional information about the suppression.
	Properties *PropertyBag `json:"properties,omitempty"`

	// A string that indicates the review status of the suppression.
	Status *string `json:"status,omitempty"`
}

// NewSuppression - creates a new
func NewSuppression() *Suppression {
	return &Suppression{}
}

// WithGuID - add a GuID to the Suppression
func (g *Suppression) WithGuID(guid string) *Suppression {
	g.GuID = &guid
	return g
}

// WithJustification - add a Justification to the Suppression
func (j *Suppression) WithJustification(justification string) *Suppression {
	j.Justification = &justification
	return j
}

// WithKind - add a Kind to the Suppression
func (k *Suppression) WithKind(kind string) *Suppression {
	k.Kind = &kind
	return k
}

// WithLocation - add a Location to the Suppression
func (l *Suppression) WithLocation(location *Location) *Suppression {
	l.Location = location
	return l
}

// WithProperties - add a Properties to the Suppression
func (p *Suppression) WithProperties(properties *PropertyBag) *Suppression {
	p.Properties = properties
	return p
}

// WithStatus - add a Status to the Suppression
func (s *Suppression) WithStatus(status string) *Suppression {
	s.Status = &status
	return s
}
