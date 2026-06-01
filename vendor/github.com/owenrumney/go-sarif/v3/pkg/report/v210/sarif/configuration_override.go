package sarif

// ConfigurationOverride - Information about how a specific rule or notification was reconfigured at runtime.
type ConfigurationOverride struct {
	// Specifies how the rule or notification was configured during the scan.
	Configuration *ReportingConfiguration `json:"configuration,omitempty"`

	// A reference used to locate the descriptor whose configuration was overridden.
	Descriptor *ReportingDescriptorReference `json:"descriptor,omitempty"`

	// Key/value pairs that provide additional information about the configuration override.
	Properties *PropertyBag `json:"properties,omitempty"`
}

// NewConfigurationOverride - creates a new
func NewConfigurationOverride() *ConfigurationOverride {
	return &ConfigurationOverride{}
}

// WithConfiguration - add a Configuration to the ConfigurationOverride
func (c *ConfigurationOverride) WithConfiguration(configuration *ReportingConfiguration) *ConfigurationOverride {
	c.Configuration = configuration
	return c
}

// WithDescriptor - add a Descriptor to the ConfigurationOverride
func (d *ConfigurationOverride) WithDescriptor(descriptor *ReportingDescriptorReference) *ConfigurationOverride {
	d.Descriptor = descriptor
	return d
}

// WithProperties - add a Properties to the ConfigurationOverride
func (p *ConfigurationOverride) WithProperties(properties *PropertyBag) *ConfigurationOverride {
	p.Properties = properties
	return p
}
