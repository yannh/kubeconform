package sarif

// Notification - Describes a condition relevant to the tool itself, as opposed to being relevant to a target being analyzed by the tool.
type Notification struct {
	// A reference used to locate the rule descriptor associated with this notification.
	AssociatedRule *ReportingDescriptorReference `json:"associatedRule,omitempty"`

	// A reference used to locate the descriptor relevant to this notification.
	Descriptor *ReportingDescriptorReference `json:"descriptor,omitempty"`

	// The runtime exception, if any, relevant to this notification.
	Exception *Exception `json:"exception,omitempty"`

	// A value specifying the severity level of the notification.
	Level string `json:"level"`

	// The locations relevant to this notification.
	Locations []*Location `json:"locations"`

	// A message that describes the condition that was encountered.
	Message *Message `json:"message,omitempty"`

	// Key/value pairs that provide additional information about the notification.
	Properties *PropertyBag `json:"properties,omitempty"`

	// A set of locations relevant to this notification.
	RelatedLocations []*Location `json:"relatedLocations"`

	// The thread identifier of the code that generated the notification.
	ThreadID *int `json:"threadId,omitempty"`

	// The Coordinated Universal Time (UTC) date and time at which the analysis tool generated the notification.
	TimeUtc *string `json:"timeUtc,omitempty"`
}

// NewNotification - creates a new
func NewNotification() *Notification {
	return &Notification{
		Level:            "warning",
		Locations:        make([]*Location, 0),
		RelatedLocations: make([]*Location, 0),
	}
}

// WithAssociatedRule - add a AssociatedRule to the Notification
func (a *Notification) WithAssociatedRule(associatedRule *ReportingDescriptorReference) *Notification {
	a.AssociatedRule = associatedRule
	return a
}

// WithDescriptor - add a Descriptor to the Notification
func (d *Notification) WithDescriptor(descriptor *ReportingDescriptorReference) *Notification {
	d.Descriptor = descriptor
	return d
}

// WithException - add a Exception to the Notification
func (e *Notification) WithException(exception *Exception) *Notification {
	e.Exception = exception
	return e
}

// WithLevel - add a Level to the Notification
func (l *Notification) WithLevel(level string) *Notification {
	l.Level = level
	return l
}

// WithLocations - add a Locations to the Notification
func (l *Notification) WithLocations(locations []*Location) *Notification {
	l.Locations = locations
	return l
}

// AddLocation - add a single Location to the Notification
func (l *Notification) AddLocation(location *Location) *Notification {
	l.Locations = append(l.Locations, location)
	return l
}

// WithMessage - add a Message to the Notification
func (m *Notification) WithMessage(message *Message) *Notification {
	m.Message = message
	return m
}

// WithProperties - add a Properties to the Notification
func (p *Notification) WithProperties(properties *PropertyBag) *Notification {
	p.Properties = properties
	return p
}

// WithRelatedLocations - add a RelatedLocations to the Notification
func (r *Notification) WithRelatedLocations(relatedLocations []*Location) *Notification {
	r.RelatedLocations = relatedLocations
	return r
}

// AddRelatedLocation - add a single RelatedLocation to the Notification
func (r *Notification) AddRelatedLocation(relatedLocation *Location) *Notification {
	r.RelatedLocations = append(r.RelatedLocations, relatedLocation)
	return r
}

// WithThreadID - add a ThreadID to the Notification
func (t *Notification) WithThreadID(threadId int) *Notification {
	t.ThreadID = &threadId
	return t
}

// WithTimeUtc - add a TimeUtc to the Notification
func (t *Notification) WithTimeUtc(timeUtc string) *Notification {
	t.TimeUtc = &timeUtc
	return t
}
