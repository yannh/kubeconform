// Derived from https://github.com/instrumenta/openapi2jsonschema
// Go port of openapi2jsonschema.py.
package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// OrderedMap preserves key insertion order, matching Python dict semantics
// (insertion-ordered) used by PyYAML SafeLoader and json.dumps.
type OrderedMap struct {
	keys   []string
	values map[string]any
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{values: map[string]any{}}
}

func (m *OrderedMap) Has(k string) bool {
	_, ok := m.values[k]
	return ok
}

func (m *OrderedMap) Get(k string) (any, bool) {
	v, ok := m.values[k]
	return v, ok
}

func (m *OrderedMap) Set(k string, v any) {
	if _, ok := m.values[k]; !ok {
		m.keys = append(m.keys, k)
	}
	m.values[k] = v
}

func (m *OrderedMap) Keys() []string { return m.keys }

func (m *OrderedMap) Len() int { return len(m.keys) }

// MarshalJSON emits keys in insertion order with no HTML escaping, matching
// Python's json.dumps default.
func (m *OrderedMap) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, k := range m.keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		kb, err := encodeJSON(k)
		if err != nil {
			return nil, err
		}
		buf.Write(kb)
		buf.WriteByte(':')
		vb, err := encodeJSON(m.values[k])
		if err != nil {
			return nil, err
		}
		buf.Write(vb)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func encodeJSON(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	// Encoder always appends a trailing newline; strip it.
	out := buf.Bytes()
	if n := len(out); n > 0 && out[n-1] == '\n' {
		out = out[:n-1]
	}
	return out, nil
}

// yamlToData converts a yaml.Node into Go values made of *OrderedMap, []any,
// and primitive types, preserving mapping key order.
func yamlToData(n *yaml.Node) (any, error) {
	switch n.Kind {
	case yaml.DocumentNode:
		if len(n.Content) == 0 {
			return nil, nil
		}
		return yamlToData(n.Content[0])
	case yaml.MappingNode:
		m := NewOrderedMap()
		for i := 0; i < len(n.Content); i += 2 {
			kNode := n.Content[i]
			vNode := n.Content[i+1]
			var key string
			if err := kNode.Decode(&key); err != nil {
				key = kNode.Value
			}
			v, err := yamlToData(vNode)
			if err != nil {
				return nil, err
			}
			m.Set(key, v)
		}
		return m, nil
	case yaml.SequenceNode:
		out := make([]any, 0, len(n.Content))
		for _, c := range n.Content {
			v, err := yamlToData(c)
			if err != nil {
				return nil, err
			}
			out = append(out, v)
		}
		return out, nil
	case yaml.ScalarNode:
		// Decode preserves YAML scalar typing (bool, int, float, string, null).
		var v any
		if err := n.Decode(&v); err != nil {
			return nil, err
		}
		return v, nil
	case yaml.AliasNode:
		return yamlToData(n.Alias)
	}
	return nil, nil
}

// additionalProperties recreates the kubectl behaviour: any object with a
// "properties" key gets `additionalProperties: false` set (unless the caller
// asks to skip it for the root).
func additionalProperties(data any, skip bool) any {
	m, ok := data.(*OrderedMap)
	if !ok {
		if list, ok := data.([]any); ok {
			for i := range list {
				list[i] = additionalProperties(list[i], false)
			}
		}
		return data
	}
	if m.Has("properties") && !skip {
		if !m.Has("additionalProperties") {
			m.Set("additionalProperties", false)
		}
	}
	for _, k := range m.Keys() {
		v, _ := m.Get(k)
		m.Set(k, additionalProperties(v, false))
	}
	return m
}

// replaceIntOrString replaces `format: int-or-string` with the canonical
// oneOf{string,integer} expansion.
func replaceIntOrString(data any) any {
	m, ok := data.(*OrderedMap)
	if !ok {
		if list, ok := data.([]any); ok {
			out := make([]any, len(list))
			for i, x := range list {
				out[i] = replaceIntOrString(x)
			}
			return out
		}
		return data
	}
	out := NewOrderedMap()
	for _, k := range m.Keys() {
		v, _ := m.Get(k)
		switch vv := v.(type) {
		case *OrderedMap:
			if f, ok := vv.Get("format"); ok {
				if s, ok := f.(string); ok && s == "int-or-string" {
					replacement := NewOrderedMap()
					strType := NewOrderedMap()
					strType.Set("type", "string")
					intType := NewOrderedMap()
					intType.Set("type", "integer")
					replacement.Set("oneOf", []any{strType, intType})
					out.Set(k, replacement)
					continue
				}
			}
			out.Set(k, replaceIntOrString(vv))
		case []any:
			rep := make([]any, len(vv))
			for i, x := range vv {
				rep[i] = replaceIntOrString(x)
			}
			out.Set(k, rep)
		default:
			out.Set(k, v)
		}
	}
	return out
}

func writeSchemaFile(schema any, filename string, denyRootAdditionalProperties bool) error {
	schema = additionalProperties(schema, !denyRootAdditionalProperties)
	schema = replaceIntOrString(schema)

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(schema); err != nil {
		return err
	}
	// Encoder appends "\n"; Python's `print(s, file=f)` also appends one.
	// json.dumps itself does not, so the encoder's newline matches print().

	// Treat the input as a path component only — guard against directory
	// traversal in user-supplied filenames.
	filename = filepath.Base(filename)
	if err := os.WriteFile(filename, buf.Bytes(), 0o644); err != nil {
		return err
	}
	fmt.Printf("JSON schema written to %s\n", filename)
	return nil
}

func formatFilename(format, kind, group, fullgroup, version string) string {
	r := strings.NewReplacer(
		"{kind}", kind,
		"{group}", group,
		"{fullgroup}", fullgroup,
		"{version}", version,
	)
	return strings.ToLower(r.Replace(format)) + ".json"
}

func openInput(path string) (io.ReadCloser, error) {
	if strings.HasPrefix(path, "http") {
		client := http.DefaultClient
		if os.Getenv("DISABLE_SSL_CERT_VALIDATION") != "" {
			client = &http.Client{Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}}
		}
		resp, err := client.Get(path)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode >= 400 {
			resp.Body.Close()
			return nil, fmt.Errorf("HTTP %d fetching %s", resp.StatusCode, path)
		}
		return resp.Body, nil
	}
	return os.Open(path)
}

// process consumes one CRD source (file or URL) and writes one JSON schema per
// version found, mirroring the original Python control flow.
func process(path string, filenameFormat string, denyRoot bool) error {
	r, err := openInput(path)
	if err != nil {
		return err
	}
	defer r.Close()

	dec := yaml.NewDecoder(r)
	var defs []*OrderedMap
	for {
		var node yaml.Node
		if err := dec.Decode(&node); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		v, err := yamlToData(&node)
		if err != nil {
			return err
		}
		m, ok := v.(*OrderedMap)
		if !ok || m == nil {
			continue
		}
		if items, ok := m.Get("items"); ok {
			if list, ok := items.([]any); ok {
				for _, it := range list {
					if im, ok := it.(*OrderedMap); ok {
						defs = append(defs, im)
					}
				}
			}
		}
		kind, ok := m.Get("kind")
		if !ok {
			continue
		}
		if ks, _ := kind.(string); ks != "CustomResourceDefinition" {
			continue
		}
		defs = append(defs, m)
	}

	for _, y := range defs {
		spec, ok := getMap(y, "spec")
		if !ok {
			continue
		}
		names, _ := getMap(spec, "names")
		kind, _ := getString(names, "kind")
		fullgroup, _ := getString(spec, "group")
		group := strings.SplitN(fullgroup, ".", 2)[0]

		versions, hasVersions := spec.Get("versions")
		versionList, _ := versions.([]any)

		if hasVersions && len(versionList) > 0 {
			for _, vAny := range versionList {
				vMap, ok := vAny.(*OrderedMap)
				if !ok {
					continue
				}
				vName, _ := getString(vMap, "name")
				schemaMap, hasSchema := getMap(vMap, "schema")
				if hasSchema {
					if root, ok := schemaMap.Get("openAPIV3Schema"); ok {
						filename := formatFilename(filenameFormat, kind, group, fullgroup, vName)
						if err := writeSchemaFile(root, filename, denyRoot); err != nil {
							return err
						}
						continue
					}
				}
				validation, hasValidation := getMap(spec, "validation")
				if hasValidation {
					if root, ok := validation.Get("openAPIV3Schema"); ok {
						filename := formatFilename(filenameFormat, kind, group, fullgroup, vName)
						if err := writeSchemaFile(root, filename, denyRoot); err != nil {
							return err
						}
					}
				}
			}
		} else if validation, hasValidation := getMap(spec, "validation"); hasValidation {
			if root, ok := validation.Get("openAPIV3Schema"); ok {
				vName, _ := getString(spec, "version")
				filename := formatFilename(filenameFormat, kind, group, fullgroup, vName)
				if err := writeSchemaFile(root, filename, denyRoot); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func getMap(m *OrderedMap, key string) (*OrderedMap, bool) {
	if m == nil {
		return nil, false
	}
	v, ok := m.Get(key)
	if !ok {
		return nil, false
	}
	mm, ok := v.(*OrderedMap)
	return mm, ok
}

func getString(m *OrderedMap, key string) (string, bool) {
	if m == nil {
		return "", false
	}
	v, ok := m.Get(key)
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Missing FILE parameter.\nUsage: %s [FILE]\n", os.Args[0])
		os.Exit(1)
	}
	filenameFormat := os.Getenv("FILENAME_FORMAT")
	if filenameFormat == "" {
		filenameFormat = "{kind}_{version}"
	}
	denyRoot := os.Getenv("DENY_ROOT_ADDITIONAL_PROPERTIES") != ""

	for _, arg := range os.Args[1:] {
		if err := process(arg, filenameFormat, denyRoot); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}
