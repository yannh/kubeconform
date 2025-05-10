package loader

import (
	"fmt"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/yannh/kubeconform/pkg/cache"
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
	if l.cache != nil {
		if cached, err := l.cache.Get(path); err == nil {
			return cached, nil
		}
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s, err := jsonschema.UnmarshalJSON(f)
	if err != nil {
		return nil, err
	}

	if l.cache != nil {
		if err = l.cache.Set(path, s); err != nil {
			return nil, fmt.Errorf("failed to write cache to disk: %s", err)
		}
	}

	return s, nil
}

// ToFile is helper method to convert file url to file path.
func (l FileLoader) ToFile(url string) (string, error) {
	u, err := gourl.Parse(url)
	if err != nil {
		return "", err
	}
	if u.Scheme != "file" {
		return "", fmt.Errorf("invalid file url: %s", u)
	}
	path := u.Path
	if runtime.GOOS == "windows" {
		path = strings.TrimPrefix(path, "/")
		path = filepath.FromSlash(path)
	}
	return path, nil
}

func NewFileLoader(cache cache.Cache) *FileLoader {
	return &FileLoader{
		cache: cache,
	}
}
