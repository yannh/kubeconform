package resource

import (
	"os"
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
