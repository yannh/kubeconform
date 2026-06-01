package sarif

// WebRequest - Describes an HTTP request.
type WebRequest struct {
	// The request headers.
	Headers map[string]string `json:"headers,omitempty"`

	// The request parameters.
	Parameters map[string]string `json:"parameters,omitempty"`

	// The body of the request.
	Body *ArtifactContent `json:"body,omitempty"`

	// The index within the run.webRequests array of the request object associated with this result.
	Index int `json:"index"`

	// The HTTP method. Well-known values are 'GET', 'PUT', 'POST', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS', 'TRACE', 'CONNECT'.
	Method *string `json:"method,omitempty"`

	// Key/value pairs that provide additional information about the request.
	Properties *PropertyBag `json:"properties,omitempty"`

	// The request protocol. Example: 'http'.
	Protocol *string `json:"protocol,omitempty"`

	// The target of the request.
	Target *string `json:"target,omitempty"`

	// The request version. Example: '1.1'.
	Version *string `json:"version,omitempty"`
}

// NewWebRequest - creates a new
func NewWebRequest() *WebRequest {
	return &WebRequest{
		Index: -1,
	}
}

// AddHeader - add a single Header to the WebRequest
func (h *WebRequest) AddHeader(key, header string) *WebRequest {
	h.Headers[key] = header
	return h
}

// WithHeaders - add a Headers to the WebRequest
func (h *WebRequest) WithHeaders(headers map[string]string) *WebRequest {
	h.Headers = headers
	return h
}

// AddParameter - add a single Parameter to the WebRequest
func (p *WebRequest) AddParameter(key, parameter string) *WebRequest {
	p.Parameters[key] = parameter
	return p
}

// WithParameters - add a Parameters to the WebRequest
func (p *WebRequest) WithParameters(parameters map[string]string) *WebRequest {
	p.Parameters = parameters
	return p
}

// WithBody - add a Body to the WebRequest
func (b *WebRequest) WithBody(body *ArtifactContent) *WebRequest {
	b.Body = body
	return b
}

// WithIndex - add a Index to the WebRequest
func (i *WebRequest) WithIndex(index int) *WebRequest {
	i.Index = index
	return i
}

// WithMethod - add a Method to the WebRequest
func (m *WebRequest) WithMethod(method string) *WebRequest {
	m.Method = &method
	return m
}

// WithProperties - add a Properties to the WebRequest
func (p *WebRequest) WithProperties(properties *PropertyBag) *WebRequest {
	p.Properties = properties
	return p
}

// WithProtocol - add a Protocol to the WebRequest
func (p *WebRequest) WithProtocol(protocol string) *WebRequest {
	p.Protocol = &protocol
	return p
}

// WithTarget - add a Target to the WebRequest
func (t *WebRequest) WithTarget(target string) *WebRequest {
	t.Target = &target
	return t
}

// WithVersion - add a Version to the WebRequest
func (v *WebRequest) WithVersion(version string) *WebRequest {
	v.Version = &version
	return v
}
