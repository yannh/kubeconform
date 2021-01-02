package resource

import (
	"os"
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
