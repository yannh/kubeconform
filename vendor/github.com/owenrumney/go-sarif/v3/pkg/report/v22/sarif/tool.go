package sarif

// Tool - The analysis tool that was run.
type Tool struct {
	// The analysis tool that was run.
	Driver *ToolComponent `json:"driver,omitempty"`

	// Tool extensions that contributed to or reconfigured the analysis tool that was run.
	Extensions []*ToolComponent `json:"extensions"`

	// Key/value pairs that provide additional information about the tool.
	Properties *PropertyBag `json:"properties,omitempty"`
}

// NewTool - creates a new
func NewTool() *Tool {
	return &Tool{
		Extensions: make([]*ToolComponent, 0),
	}
}

// WithDriver - add a Driver to the Tool
func (d *Tool) WithDriver(driver *ToolComponent) *Tool {
	d.Driver = driver
	return d
}

// WithExtensions - add a Extensions to the Tool
func (e *Tool) WithExtensions(extensions []*ToolComponent) *Tool {
	e.Extensions = extensions
	return e
}

// AddExtension - add a single Extension to the Tool
func (e *Tool) AddExtension(extension *ToolComponent) *Tool {
	e.Extensions = append(e.Extensions, extension)
	return e
}

// WithProperties - add a Properties to the Tool
func (p *Tool) WithProperties(properties *PropertyBag) *Tool {
	p.Properties = properties
	return p
}
