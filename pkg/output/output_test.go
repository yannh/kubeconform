package output

import (
    "testing"
)


// Test generated using Keploy
func TestNew_JSONOutput(t *testing.T) {
    writer := &mockWriter{}
    output, err := New(writer, "json", false, false, false)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if output == nil {
        t.Fatalf("Expected a valid Output implementation, got nil")
    }
}


// Test generated using Keploy
func TestNew_JUnitOutput(t *testing.T) {
    writer := &mockWriter{}
    output, err := New(writer, "junit", false, false, false)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if output == nil {
        t.Fatalf("Expected a valid Output implementation, got nil")
    }
}


// Test generated using Keploy
func TestNew_PrettyOutput(t *testing.T) {
    writer := &mockWriter{}
    output, err := New(writer, "pretty", false, false, false)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if output == nil {
        t.Fatalf("Expected a valid Output implementation, got nil")
    }
}


// Test generated using Keploy
func TestNew_UnsupportedFormat(t *testing.T) {
    writer := &mockWriter{}
    _, err := New(writer, "unsupported", false, false, false)
    if err == nil {
        t.Fatalf("Expected an error, got nil")
    }
    expectedError := "'outputFormat' must be 'json', 'junit', 'pretty', 'tap' or 'text'"
    if err.Error() != expectedError {
        t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
    }
}


// Test generated using Keploy
func TestNew_TapOutput(t *testing.T) {
    writer := &mockWriter{}
    output, err := New(writer, "tap", false, false, false)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if output == nil {
        t.Fatalf("Expected a valid Output implementation, got nil")
    }
}


// Test generated using Keploy
func TestNew_TextOutput(t *testing.T) {
    writer := &mockWriter{}
    output, err := New(writer, "text", false, false, false)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if output == nil {
        t.Fatalf("Expected a valid Output implementation, got nil")
    }
}
