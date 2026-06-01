package sarif

// ResultProvenance - Contains information about how and when a result was detected.
type ResultProvenance struct {
	// An array of physicalLocation objects which specify the portions of an analysis tool's output that a converter transformed into the result.
	ConversionSources []*PhysicalLocation `json:"conversionSources"`

	// A GUID-valued string equal to the automationDetails.guid property of the run in which the result was first detected.
	FirstDetectionRunGuID *string `json:"firstDetectionRunGuid,omitempty"`

	// The Coordinated Universal Time (UTC) date and time at which the result was first detected. See "Date/time properties" in the SARIF spec for the required format.
	FirstDetectionTimeUtc *string `json:"firstDetectionTimeUtc,omitempty"`

	// The index within the run.invocations array of the invocation object which describes the tool invocation that detected the result.
	InvocationIndex int `json:"invocationIndex"`

	// A GUID-valued string equal to the automationDetails.guid property of the run in which the result was most recently detected.
	LastDetectionRunGuID *string `json:"lastDetectionRunGuid,omitempty"`

	// The Coordinated Universal Time (UTC) date and time at which the result was most recently detected. See "Date/time properties" in the SARIF spec for the required format.
	LastDetectionTimeUtc *string `json:"lastDetectionTimeUtc,omitempty"`

	// Key/value pairs that provide additional information about the result.
	Properties *PropertyBag `json:"properties,omitempty"`
}

// NewResultProvenance - creates a new
func NewResultProvenance() *ResultProvenance {
	return &ResultProvenance{
		ConversionSources: make([]*PhysicalLocation, 0),
		InvocationIndex:   -1,
	}
}

// WithConversionSources - add a ConversionSources to the ResultProvenance
func (c *ResultProvenance) WithConversionSources(conversionSources []*PhysicalLocation) *ResultProvenance {
	c.ConversionSources = conversionSources
	return c
}

// AddConversionSource - add a single ConversionSource to the ResultProvenance
func (c *ResultProvenance) AddConversionSource(conversionSource *PhysicalLocation) *ResultProvenance {
	c.ConversionSources = append(c.ConversionSources, conversionSource)
	return c
}

// WithFirstDetectionRunGuID - add a FirstDetectionRunGuID to the ResultProvenance
func (f *ResultProvenance) WithFirstDetectionRunGuID(firstDetectionRunGuid string) *ResultProvenance {
	f.FirstDetectionRunGuID = &firstDetectionRunGuid
	return f
}

// WithFirstDetectionTimeUtc - add a FirstDetectionTimeUtc to the ResultProvenance
func (f *ResultProvenance) WithFirstDetectionTimeUtc(firstDetectionTimeUtc string) *ResultProvenance {
	f.FirstDetectionTimeUtc = &firstDetectionTimeUtc
	return f
}

// WithInvocationIndex - add a InvocationIndex to the ResultProvenance
func (i *ResultProvenance) WithInvocationIndex(invocationIndex int) *ResultProvenance {
	i.InvocationIndex = invocationIndex
	return i
}

// WithLastDetectionRunGuID - add a LastDetectionRunGuID to the ResultProvenance
func (l *ResultProvenance) WithLastDetectionRunGuID(lastDetectionRunGuid string) *ResultProvenance {
	l.LastDetectionRunGuID = &lastDetectionRunGuid
	return l
}

// WithLastDetectionTimeUtc - add a LastDetectionTimeUtc to the ResultProvenance
func (l *ResultProvenance) WithLastDetectionTimeUtc(lastDetectionTimeUtc string) *ResultProvenance {
	l.LastDetectionTimeUtc = &lastDetectionTimeUtc
	return l
}

// WithProperties - add a Properties to the ResultProvenance
func (p *ResultProvenance) WithProperties(properties *PropertyBag) *ResultProvenance {
	p.Properties = properties
	return p
}
