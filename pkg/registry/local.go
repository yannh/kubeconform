package registry

import (
    "errors"
    "fmt"
    "io"
    "log"
    "os"
)

type LocalRegistry struct {
    pathTemplate   string
    strict         bool
    debug          bool
    schemaPathFunc func(pathTemplate, resourceKind, resourceAPIVersion, k8sVersion string, strict bool) (string, error)
}

func newLocalRegistry(pathTemplate string, strict bool, debug bool) (*LocalRegistry, error) {
    return &LocalRegistry{
        pathTemplate:   pathTemplate,
        strict:         strict,
        debug:          debug,
        schemaPathFunc: schemaPath,
    }, nil
}

func (r LocalRegistry) DownloadSchema(resourceKind, resourceAPIVersion, k8sVersion string) (string, []byte, error) {
    schemaFile, err := r.schemaPathFunc(r.pathTemplate, resourceKind, resourceAPIVersion, k8sVersion, r.strict)
    if err != nil {
        return schemaFile, []byte{}, nil
    }
    f, err := os.Open(schemaFile)
    if err != nil {
        if os.IsNotExist(err) {
            msg := fmt.Sprintf("could not open file %s", schemaFile)
            if r.debug {
                log.Print(msg)
            }
            return schemaFile, nil, newNotFoundError(errors.New(msg))
        }

        msg := fmt.Sprintf("failed to open schema at %s: %s", schemaFile, err)
        if r.debug {
            log.Print(msg)
        }
        return schemaFile, nil, errors.New(msg)
    }
    defer f.Close()
    content, err := io.ReadAll(f)
    if err != nil {
        msg := fmt.Sprintf("failed to read schema at %s: %s", schemaFile, err)
        if r.debug {
            log.Print(msg)
        }
        return schemaFile, nil, err
    }
    if r.debug {
        log.Printf("using schema found at %s", schemaFile)
    }
    return schemaFile, content, nil
}
