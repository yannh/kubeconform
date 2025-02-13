package registry

import (
    "fmt"
    "testing"
    "os"
)


// Test generated using Keploy
func TestDownloadSchema_SchemaPathError(t *testing.T) {
    // Arrange
    pathTemplate := "./schemas/%s/%s/%s.json"
    strict := false
    debug := false
    registry, err := newLocalRegistry(pathTemplate, strict, debug)
    if err != nil {
        t.Fatalf("Failed to create LocalRegistry: %v", err)
    }

    // Mock schemaPathFunc to return an error
    registry.schemaPathFunc = func(pathTemplate, resourceKind, resourceAPIVersion, k8sVersion string, strict bool) (string, error) {
        return "", fmt.Errorf("mock schemaPath error")
    }

    // Act
    filePath, content, err := registry.DownloadSchema("Pod", "v1", "1.21")

    // Assert
    if err != nil {
        t.Errorf("Expected nil error, got %v", err)
    }
    if filePath != "" {
        t.Errorf("Expected empty filePath, got %s", filePath)
    }
    if len(content) != 0 {
        t.Errorf("Expected empty content, got %v", content)
    }
}


// Test generated using Keploy
func TestDownloadSchema_Success(t *testing.T) {
    // Arrange
    pathTemplate := "./schemas/%s/%s/%s.json"
    strict := false
    debug := true
    registry, err := newLocalRegistry(pathTemplate, strict, debug)
    if err != nil {
        t.Fatalf("Failed to create LocalRegistry: %v", err)
    }

    // Mock schemaPathFunc to return a valid file path
    registry.schemaPathFunc = func(pathTemplate, resourceKind, resourceAPIVersion, k8sVersion string, strict bool) (string, error) {
        return "./valid_schema.json", nil
    }

    // Create a valid schema file
    validFile := "./valid_schema.json"
    expectedContent := []byte(`{"type": "object"}`)
    os.WriteFile(validFile, expectedContent, 0644)
    defer os.Remove(validFile)

    // Act
    filePath, content, err := registry.DownloadSchema("Pod", "v1", "1.21")

    // Assert
    if err != nil {
        t.Errorf("Expected nil error, got %v", err)
    }
    if filePath != validFile {
        t.Errorf("Expected filePath '%s', got %s", validFile, filePath)
    }
    if string(content) != string(expectedContent) {
        t.Errorf("Expected content '%s', got %s", string(expectedContent), string(content))
    }
}


// Test generated using Keploy
func TestDownloadSchema_FileNotFound(t *testing.T) {
    // Arrange
    pathTemplate := "./schemas/%s/%s/%s.json"
    strict := false
    debug := true
    registry, err := newLocalRegistry(pathTemplate, strict, debug)
    if err != nil {
        t.Fatalf("Failed to create LocalRegistry: %v", err)
    }

    // Mock schemaPathFunc to return a valid file path
    registry.schemaPathFunc = func(pathTemplate, resourceKind, resourceAPIVersion, k8sVersion string, strict bool) (string, error) {
        return "./nonexistent_file.json", nil
    }

    // Act
    filePath, content, err := registry.DownloadSchema("Pod", "v1", "1.21")

    // Assert
    if err == nil || err.Error() != "could not open file ./nonexistent_file.json" {
        t.Errorf("Expected NotFoundError, got %v", err)
    }
    if filePath != "./nonexistent_file.json" {
        t.Errorf("Expected filePath './nonexistent_file.json', got %s", filePath)
    }
    if content != nil {
        t.Errorf("Expected nil content, got %v", content)
    }
}
