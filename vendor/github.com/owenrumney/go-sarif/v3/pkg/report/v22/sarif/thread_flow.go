package sarif

// ThreadFlow - Describes a sequence of code locations that specify a path through a single thread of execution such as an operating system or fiber.
type ThreadFlow struct {
	// Values of relevant expressions at the start of the thread flow that may change during thread flow execution.
	InitialState map[string]MultiformatMessageString `json:"initialState,omitempty"`

	// Values of relevant expressions at the start of the thread flow that remain constant.
	ImmutableState map[string]MultiformatMessageString `json:"immutableState,omitempty"`

	// An string that uniquely identifies the threadFlow within the codeFlow in which it occurs.
	ID *string `json:"id,omitempty"`

	// A temporally ordered array of 'threadFlowLocation' objects, each of which describes a location visited by the tool while producing the result.
	Locations []*ThreadFlowLocation `json:"locations,omitempty"`

	// A message relevant to the thread flow.
	Message *Message `json:"message,omitempty"`

	// Key/value pairs that provide additional information about the thread flow.
	Properties *PropertyBag `json:"properties,omitempty"`
}

// NewThreadFlow - creates a new
func NewThreadFlow() *ThreadFlow {
	return &ThreadFlow{
		Locations: make([]*ThreadFlowLocation, 0),
	}
}

// AddInitialState - add a single InitialState to the ThreadFlow
func (i *ThreadFlow) AddInitialState(key string, initialState MultiformatMessageString) *ThreadFlow {
	i.InitialState[key] = initialState
	return i
}

// WithInitialState - add a InitialState to the ThreadFlow
func (i *ThreadFlow) WithInitialState(initialState map[string]MultiformatMessageString) *ThreadFlow {
	i.InitialState = initialState
	return i
}

// AddImmutableState - add a single ImmutableState to the ThreadFlow
func (i *ThreadFlow) AddImmutableState(key string, immutableState MultiformatMessageString) *ThreadFlow {
	i.ImmutableState[key] = immutableState
	return i
}

// WithImmutableState - add a ImmutableState to the ThreadFlow
func (i *ThreadFlow) WithImmutableState(immutableState map[string]MultiformatMessageString) *ThreadFlow {
	i.ImmutableState = immutableState
	return i
}

// WithID - add a ID to the ThreadFlow
func (i *ThreadFlow) WithID(id string) *ThreadFlow {
	i.ID = &id
	return i
}

// WithLocations - add a Locations to the ThreadFlow
func (l *ThreadFlow) WithLocations(locations []*ThreadFlowLocation) *ThreadFlow {
	l.Locations = locations
	return l
}

// AddLocation - add a single Location to the ThreadFlow
func (l *ThreadFlow) AddLocation(location *ThreadFlowLocation) *ThreadFlow {
	l.Locations = append(l.Locations, location)
	return l
}

// WithMessage - add a Message to the ThreadFlow
func (m *ThreadFlow) WithMessage(message *Message) *ThreadFlow {
	m.Message = message
	return m
}

// WithProperties - add a Properties to the ThreadFlow
func (p *ThreadFlow) WithProperties(properties *PropertyBag) *ThreadFlow {
	p.Properties = properties
	return p
}
