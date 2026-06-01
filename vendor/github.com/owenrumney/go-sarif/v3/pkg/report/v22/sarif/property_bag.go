package sarif

// PropertyBag - Key/value pairs that provide additional information about the object.
type PropertyBag struct {
	// AdditionalProperties - additional properties
	Properties Properties `json:"properties,omitempty"`

	// A set of distinct strings that provide additional information.
	Tags []string `json:"tags"`
}

// NewPropertyBag - creates a new
func NewPropertyBag() *PropertyBag {
	return &PropertyBag{
		Properties: make(Properties),
		Tags:       make([]string, 0),
	}
}

// AddProperty - add a property to the properties
func (p *PropertyBag) Add(key string, value interface{}) *PropertyBag {
	p.Properties[key] = value
	return p
}

// WithTags - add a Tags to the PropertyBag
func (t *PropertyBag) WithTags(tags []string) *PropertyBag {
	t.Tags = tags
	return t
}

// AddTag - add a single Tag to the PropertyBag
func (t *PropertyBag) AddTag(tag string) *PropertyBag {
	t.Tags = append(t.Tags, tag)
	return t
}
