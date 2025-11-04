package sarif

// ThreadFlowLocation - A location visited by an analysis tool while simulating or monitoring the execution of a program.
type ThreadFlowLocation struct {
	// A dictionary, each of whose keys specifies a variable or expression, the associated value of which represents the variable or expression value. For an annotation of kind 'continuation', for example, this dictionary might hold the current assumed values of a set of global variables.
	State map[string]MultiformatMessageString `json:"state,omitempty"`

	// An integer representing the temporal order in which execution reached this location.
	ExecutionOrder int `json:"executionOrder"`

	// The Coordinated Universal Time (UTC) date and time at which this location was executed.
	ExecutionTimeUtc *string `json:"executionTimeUtc,omitempty"`

	// Specifies the importance of this location in understanding the code flow in which it occurs. The order from most to least important is "essential", "important", "unimportant". Default: "important".
	Importance string `json:"importance"`

	// The index within the run threadFlowLocations array.
	Index int `json:"index"`

	// A set of distinct strings that categorize the thread flow location. Well-known kinds include 'acquire', 'release', 'enter', 'exit', 'call', 'return', 'branch', 'implicit', 'false', 'true', 'caution', 'danger', 'unknown', 'unreachable', 'taint', 'function', 'handler', 'lock', 'memory', 'resource', 'scope' and 'value'.
	Kinds []string `json:"kinds"`

	// The code location.
	Location *Location `json:"location,omitempty"`

	// The name of the module that contains the code that is executing.
	Module *string `json:"module,omitempty"`

	// An integer representing a containment hierarchy within the thread flow.
	NestingLevel *int `json:"nestingLevel,omitempty"`

	// Key/value pairs that provide additional information about the threadflow location.
	Properties *PropertyBag `json:"properties,omitempty"`

	// The call stack leading to this location.
	Stack *Stack `json:"stack,omitempty"`

	// An array of references to rule or taxonomy reporting descriptors that are applicable to the thread flow location.
	Taxa []*ReportingDescriptorReference `json:"taxa"`

	// A web request associated with this thread flow location.
	WebRequest *WebRequest `json:"webRequest,omitempty"`

	// A web response associated with this thread flow location.
	WebResponse *WebResponse `json:"webResponse,omitempty"`
}

// NewThreadFlowLocation - creates a new
func NewThreadFlowLocation() *ThreadFlowLocation {
	return &ThreadFlowLocation{
		ExecutionOrder: -1,
		Importance:     "important",
		Index:          -1,
		Kinds:          make([]string, 0),
		Taxa:           make([]*ReportingDescriptorReference, 0),
	}
}

// AddState - add a single State to the ThreadFlowLocation
func (s *ThreadFlowLocation) AddState(key string, state MultiformatMessageString) *ThreadFlowLocation {
	s.State[key] = state
	return s
}

// WithState - add a State to the ThreadFlowLocation
func (s *ThreadFlowLocation) WithState(state map[string]MultiformatMessageString) *ThreadFlowLocation {
	s.State = state
	return s
}

// WithExecutionOrder - add a ExecutionOrder to the ThreadFlowLocation
func (e *ThreadFlowLocation) WithExecutionOrder(executionOrder int) *ThreadFlowLocation {
	e.ExecutionOrder = executionOrder
	return e
}

// WithExecutionTimeUtc - add a ExecutionTimeUtc to the ThreadFlowLocation
func (e *ThreadFlowLocation) WithExecutionTimeUtc(executionTimeUtc string) *ThreadFlowLocation {
	e.ExecutionTimeUtc = &executionTimeUtc
	return e
}

// WithImportance - add a Importance to the ThreadFlowLocation
func (i *ThreadFlowLocation) WithImportance(importance string) *ThreadFlowLocation {
	i.Importance = importance
	return i
}

// WithIndex - add a Index to the ThreadFlowLocation
func (i *ThreadFlowLocation) WithIndex(index int) *ThreadFlowLocation {
	i.Index = index
	return i
}

// WithKinds - add a Kinds to the ThreadFlowLocation
func (k *ThreadFlowLocation) WithKinds(kinds []string) *ThreadFlowLocation {
	k.Kinds = kinds
	return k
}

// AddKind - add a single Kind to the ThreadFlowLocation
func (k *ThreadFlowLocation) AddKind(kind string) *ThreadFlowLocation {
	k.Kinds = append(k.Kinds, kind)
	return k
}

// WithLocation - add a Location to the ThreadFlowLocation
func (l *ThreadFlowLocation) WithLocation(location *Location) *ThreadFlowLocation {
	l.Location = location
	return l
}

// WithModule - add a Module to the ThreadFlowLocation
func (m *ThreadFlowLocation) WithModule(module string) *ThreadFlowLocation {
	m.Module = &module
	return m
}

// WithNestingLevel - add a NestingLevel to the ThreadFlowLocation
func (n *ThreadFlowLocation) WithNestingLevel(nestingLevel int) *ThreadFlowLocation {
	n.NestingLevel = &nestingLevel
	return n
}

// WithProperties - add a Properties to the ThreadFlowLocation
func (p *ThreadFlowLocation) WithProperties(properties *PropertyBag) *ThreadFlowLocation {
	p.Properties = properties
	return p
}

// WithStack - add a Stack to the ThreadFlowLocation
func (s *ThreadFlowLocation) WithStack(stack *Stack) *ThreadFlowLocation {
	s.Stack = stack
	return s
}

// WithTaxa - add a Taxa to the ThreadFlowLocation
func (t *ThreadFlowLocation) WithTaxa(taxa []*ReportingDescriptorReference) *ThreadFlowLocation {
	t.Taxa = taxa
	return t
}

// AddTaxa - add a single Taxa to the ThreadFlowLocation
func (t *ThreadFlowLocation) AddTaxa(taxa *ReportingDescriptorReference) *ThreadFlowLocation {
	t.Taxa = append(t.Taxa, taxa)
	return t
}

// WithWebRequest - add a WebRequest to the ThreadFlowLocation
func (w *ThreadFlowLocation) WithWebRequest(webRequest *WebRequest) *ThreadFlowLocation {
	w.WebRequest = webRequest
	return w
}

// WithWebResponse - add a WebResponse to the ThreadFlowLocation
func (w *ThreadFlowLocation) WithWebResponse(webResponse *WebResponse) *ThreadFlowLocation {
	w.WebResponse = webResponse
	return w
}
