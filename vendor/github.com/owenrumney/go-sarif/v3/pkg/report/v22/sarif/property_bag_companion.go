package sarif

import "encoding/json"

// MarshalJSON - custom JSON marshaller for PropertyBag
func (p PropertyBag) MarshalJSON() ([]byte, error) {
	// type Alias PropertyBag
	aux := make(map[string]interface{})
	for k, v := range p.Properties {
		aux[k] = v
	}
	if len(p.Tags) > 0 {
		aux["tags"] = p.Tags
	}
	return json.Marshal(aux)
}

// UnmarshalJSON - custom JSON unmarshaller for PropertyBag
func (p *PropertyBag) UnmarshalJSON(data []byte) error {
	// type Alias PropertyBag
	aux := struct {
		Tags []string `json:"tags,omitempty"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	p.Tags = aux.Tags

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	delete(raw, "tags")
	p.Properties = raw
	return nil
}
