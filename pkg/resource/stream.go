package resource

import (
	"bufio"
	"context"
	"io"
	"strings"
)

func yamlSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Return nothing if at end of file and no data passed
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := strings.Index(string(data), "---\n"); i >= 0 {
		return i + 4, data[0:i], nil
	}

	// If at end of file with data return the data
	if atEOF {
		return len(data), data, nil
	}

	return
}

// FromStream reads resources from a byte stream, usually here stdin
func FromStream(ctx context.Context, path string, r io.Reader) (<-chan Resource, <-chan error) {
	resources := make(chan Resource)
	errors := make(chan error)
	stop := false

	go func() {
		<-ctx.Done()
		stop = true
	}()

	go func() {
		scanner := bufio.NewScanner(r)
		scanner.Split(yamlSplit)

		for res := scanner.Scan(); res != false; res = scanner.Scan() {
			if stop == true {
				break
			}
			resources <- Resource{Path: path, Bytes: scanner.Bytes()}
		}

		close(resources)
		close(errors)
	}()

	return resources, errors
}
