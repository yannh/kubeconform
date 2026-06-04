package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

// roundTrip serialises and parses through encoding/json so the resulting
// generic structure can be compared with reflect.DeepEqual regardless of how
// the producer chose to nest *OrderedMap vs map[string]any.
func roundTrip(t *testing.T, v any) any {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return out
}

func mapOf(kv ...any) *OrderedMap {
	m := NewOrderedMap()
	for i := 0; i < len(kv); i += 2 {
		m.Set(kv[i].(string), kv[i+1])
	}
	return m
}

func TestAdditionalProperties(t *testing.T) {
	cases := []struct {
		name   string
		input  *OrderedMap
		expect string
	}{
		{
			name:   "object with properties gets additionalProperties:false",
			input:  mapOf("something", mapOf("properties", NewOrderedMap())),
			expect: `{"something":{"properties":{},"additionalProperties":false}}`,
		},
		{
			name:   "object without properties is left alone",
			input:  mapOf("something", mapOf("somethingelse", NewOrderedMap())),
			expect: `{"something":{"somethingelse":{}}}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := additionalProperties(tc.input, false)
			b, err := json.Marshal(got)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(b) != tc.expect {
				t.Fatalf("got %s want %s", b, tc.expect)
			}
		})
	}
}

func TestReplaceIntOrString(t *testing.T) {
	cases := []struct {
		name   string
		input  *OrderedMap
		expect string
	}{
		{
			name:   "int-or-string is expanded to oneOf",
			input:  mapOf("something", mapOf("format", "int-or-string")),
			expect: `{"something":{"oneOf":[{"type":"string"},{"type":"integer"}]}}`,
		},
		{
			name:   "other formats are left alone",
			input:  mapOf("something", mapOf("format", "string")),
			expect: `{"something":{"format":"string"}}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := replaceIntOrString(tc.input)
			b, err := json.Marshal(got)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(b) != tc.expect {
				t.Fatalf("got %s want %s", b, tc.expect)
			}
		})
	}
}

func TestFormatFilename(t *testing.T) {
	got := formatFilename("{kind}_{version}", "Prometheus", "monitoring", "monitoring.coreos.com", "v1")
	if got != "prometheus_v1.json" {
		t.Fatalf("got %s", got)
	}
	got = formatFilename("{kind}-{group}-{version}", "Prometheus", "monitoring", "monitoring.coreos.com", "v1")
	if got != "prometheus-monitoring-v1.json" {
		t.Fatalf("got %s", got)
	}
}

// TestAdditionalProperties_PreservesUnrelatedKeys asserts the recursive walk
// does not lose sibling keys when injecting additionalProperties.
func TestAdditionalProperties_PreservesUnrelatedKeys(t *testing.T) {
	input := mapOf(
		"properties", mapOf("name", mapOf("type", "string")),
		"required", []any{"name"},
		"type", "object",
	)
	got := additionalProperties(input, false)
	parsed := roundTrip(t, got)
	// `additionalProperties: false` is only added at the level that has a
	// "properties" key — the nested {type: string} does not.
	want := map[string]any{
		"properties":           map[string]any{"name": map[string]any{"type": "string"}},
		"required":             []any{"name"},
		"type":                 "object",
		"additionalProperties": false,
	}
	if !reflect.DeepEqual(parsed, want) {
		t.Fatalf("got %v want %v", parsed, want)
	}
}
