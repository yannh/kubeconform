package resource

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

type MockFileInfo struct {
	fileName string
}

func NewMockFileInfo(filename string) *MockFileInfo {
	return &MockFileInfo{fileName: filename}
}
func (m *MockFileInfo) Name() string       { return m.fileName }
func (m *MockFileInfo) Size() int64        { return 0 }           // length in bytes for regular files; system-dependent for others
func (m *MockFileInfo) Mode() os.FileMode  { return 0 }           // file mode bits
func (m *MockFileInfo) ModTime() time.Time { return time.Time{} } // modification time
func (m *MockFileInfo) IsDir() bool        { return false }       // abbreviation for Mode().IsDir()
func (m *MockFileInfo) Sys() interface{}   { return nil }         // underlying data source (can return nil)

func TestIsYamlFile(t *testing.T) {
	for i, testCase := range []struct {
		filename string
		expect   bool
	}{
		{
			"file.yaml",
			true,
		},
		{
			"/path/to/my/file.yaml",
			true,
		},
		{
			"file.yml",
			true,
		},
		{
			"/path/to/my/file.yml",
			true,
		},
		{
			"file.notyaml",
			false,
		},
		{
			"/path/to/my/file.notyaml",
			false,
		},
		{
			"/path/to/my/file",
			false,
		},
	} {
		if got := isYAMLFile(NewMockFileInfo(testCase.filename)); got != testCase.expect {
			t.Errorf("test %d: for filename %s, expected %t, got %t", i+1, testCase.filename, testCase.expect, got)
		}
	}
}

func TestIsJSONFile(t *testing.T) {
	for i, testCase := range []struct {
		filename string
		expect   bool
	}{
		{
			"file.json",
			true,
		},
		{
			"/path/to/my/file.json",
			true,
		},
		{
			"file.notjson",
			false,
		},
		{
			"/path/to/my/file",
			false,
		},
	} {
		if got := isJSONFile(NewMockFileInfo(testCase.filename)); got != testCase.expect {
			t.Errorf("test %d: for filename %s, expected %t, got %t", i+1, testCase.filename, testCase.expect, got)
		}
	}
}

func TestFindResourcesInReader(t *testing.T) {
	maxResourceSize := 4 * 1024 * 1024   // 4MB ought to be enough for everybody
	buf := make([]byte, maxResourceSize) // We reuse this to avoid multiple large memory allocations

	for i, testCase := range []struct {
		filePath string
		yamlData string
		res      []Resource
		errs     []error
	}{
		{
			"manifest.yaml",
			``,
			[]Resource{
				{
					Path:  "manifest.yaml",
					Bytes: nil,
					sig:   nil,
				},
			},
			nil,
		},
		{
			"manifest.yaml",
			`---
foo: bar
`,
			[]Resource{
				{
					Path:  "manifest.yaml",
					Bytes: []byte("---\nfoo: bar\n"),
					sig:   nil,
				},
			},
			nil,
		},
		{
			"manifest.yaml",
			`---
foo: bar
---
lorem: ipsum
`,
			[]Resource{
				{
					Path:  "manifest.yaml",
					Bytes: []byte("---\nfoo: bar"),
					sig:   nil,
				},
				{
					Path:  "manifest.yaml",
					Bytes: []byte("lorem: ipsum\n"),
					sig:   nil,
				},
			},
			nil,
		},
	} {
		res := make(chan Resource)
		errs := make(chan error)
		receivedResources := []Resource{}
		receivedErrs := []error{}
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			for {
				select {
				case receivedResource, ok := <-res:
					if ok {
						receivedResources = append(receivedResources, receivedResource)
						continue
					}
					res = nil

				case receivedErr, ok := <-errs:
					if ok {
						receivedErrs = append(receivedErrs, receivedErr)
						continue
					}
					errs = nil
				}

				if res == nil && errs == nil {
					break
				}
			}
			wg.Done()
		}()

		r := strings.NewReader(testCase.yamlData)
		findResourcesInReader(testCase.filePath, r, res, errs, buf)
		close(res)
		close(errs)
		wg.Wait()

		if len(receivedResources) != len(testCase.res) {
			t.Errorf("test %d: expected %d resources, received %d: %+v", i, len(testCase.res), len(receivedResources), receivedResources)
			continue
		}

		for j, r := range receivedResources {
			if r.Path != testCase.res[j].Path {
				t.Errorf("test %d, resource %d, expected path %s, received %s", i, j, testCase.res[j].Path, r.Path)
			}

			if string(r.Bytes) != string(testCase.res[j].Bytes) {
				t.Errorf("test %d, resource %d, expected Bytes %s, received %s", i, j, string(testCase.res[j].Bytes), string(r.Bytes))
			}
		}
	}
}

// Test generated using Keploy
func TestFindFilesInFolders_BasicFunctionality(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tempDir := t.TempDir()
	os.WriteFile(filepath.Join(tempDir, "file1.yaml"), []byte{}, 0644)
	os.WriteFile(filepath.Join(tempDir, "file2.json"), []byte{}, 0644)
	os.WriteFile(filepath.Join(tempDir, "ignored.txt"), []byte{}, 0644)

	ignorePatterns := []string{".*ignored.*"}
	files, errors := findFilesInFolders(ctx, []string{tempDir}, ignorePatterns)

	receivedFiles := []string{}
	for f := range files {
		receivedFiles = append(receivedFiles, f)
	}

	if len(receivedFiles) != 2 {
		t.Errorf("expected 2 files, got %d: %v", len(receivedFiles), receivedFiles)
	}

	select {
	case err := <-errors:
		t.Errorf("unexpected error: %v", err)
	default:
	}
}

// Test generated using Keploy
func TestFindResourcesInFile_ErrorHandling(t *testing.T) {
	resources := make(chan Resource)
	errors := make(chan error)
	buf := make([]byte, 4*1024*1024)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for err := range errors {
			if err == nil {
				t.Errorf("expected an error, got nil")
			}
		}
	}()

	findResourcesInFile("nonexistent.yaml", resources, errors, buf)
	close(resources)
	close(errors)
	wg.Wait()
}

// Test generated using Keploy
func TestDiscoveryError_ErrorMethod(t *testing.T) {
	underlyingErr := "file not found"
	de := DiscoveryError{
		Path: "/path/to/file.yaml",
		Err:  fmt.Errorf(underlyingErr),
	}

	if de.Error() != underlyingErr {
		t.Errorf("expected error message '%s', got '%s'", underlyingErr, de.Error())
	}
}

// Test generated using Keploy
func TestIsIgnored_MultiplePatterns(t *testing.T) {
	for i, testCase := range []struct {
		path            string
		ignorePatterns  []string
		expectedIgnored bool
		expectedError   bool
	}{
		{
			path:            "/path/to/ignored/file.yaml",
			ignorePatterns:  []string{".*ignored.*", ".*file.*"},
			expectedIgnored: true,
			expectedError:   false,
		},
		{
			path:            "/path/to/not_ignored/file.yaml",
			ignorePatterns:  []string{".*ignored.*", ".*not.*"},
			expectedIgnored: true,
			expectedError:   false,
		},
		{
			path:            "/path/to/file.yaml",
			ignorePatterns:  []string{"[invalid_regex"},
			expectedIgnored: false,
			expectedError:   true,
		},
	} {
		ignored, err := isIgnored(testCase.path, testCase.ignorePatterns)
		if ignored != testCase.expectedIgnored {
			t.Errorf("test %d: expected ignored=%t, got %t", i+1, testCase.expectedIgnored, ignored)
		}
		if (err != nil) != testCase.expectedError {
			t.Errorf("test %d: expected error=%t, got %v", i+1, testCase.expectedError, err)
		}
	}
}
