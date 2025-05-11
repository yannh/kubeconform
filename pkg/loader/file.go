package loader

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/yannh/kubeconform/pkg/cache"
	"io"
	gourl "net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// FileLoader loads json file url.
type FileLoader struct {
	cache cache.Cache
}

func (l FileLoader) Load(url string) (any, error) {
	path, err := l.ToFile(url)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			msg := fmt.Sprintf("could not open file %s", path)
			return nil, NewNotFoundError(errors.New(msg))
		}
		return nil, err
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return jsonschema.UnmarshalJSON(bytes.NewReader(content))
}

// ToFile is helper method to convert file url to file path.
func (l FileLoader) ToFile(url string) (string, error) {
	u, err := gourl.Parse(url)
	if err != nil {
		return "", err
	}
	if u.Scheme != "file" {
		return url, nil
	}
	path := u.Path
	if runtime.GOOS == "windows" {
		path = strings.TrimPrefix(path, "/")
		path = filepath.FromSlash(path)
	}
	return path, nil
}

func NewFileLoader() *FileLoader {
	return &FileLoader{}
}
