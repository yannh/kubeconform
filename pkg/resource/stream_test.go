package resource_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/yannh/kubeconform/pkg/resource"
)

func TestFromStream(t *testing.T) {
	type have struct {
		Path   string
		Reader io.Reader
	}

	type want struct {
		Resources []resource.Resource
		Errors    []error
	}

	testCases := []struct {
		Have have
		Want want
	}{
		{
			Have: have{
				Path: "myfile",
				Reader: strings.NewReader(`---
apiVersion: v1
kind: ReplicationController
`),
			},
			Want: want{
				Resources: []resource.Resource{
					{
						Path: "myfile",
						Bytes: []byte(`---
apiVersion: v1
kind: ReplicationController
`),
					},
				},
				Errors: []error{},
			},
		},
		{
			Have: have{
				Path: "myfile",
				Reader: strings.NewReader(`apiVersion: v1
---
apiVersion: v2
`),
			},
			Want: want{
				Resources: []resource.Resource{
					{
						Path:  "myfile",
						Bytes: []byte(`apiVersion: v1`),
					},
					{
						Path: "myfile",
						Bytes: []byte(`apiVersion: v2
`),
					},
				},
				Errors: []error{},
			},
		},
		{
			Have: have{
				Path: "myfile",
				Reader: strings.NewReader(`apiVersion: v1
kind: ReplicationController
---
apiVersion: v1
kind: Deployment
---
apiVersion: v2
kind: CronJob
`),
			},
			Want: want{
				Resources: []resource.Resource{
					{
						Path: "myfile",
						Bytes: []byte(`apiVersion: v1
kind: ReplicationController`),
					},
					{
						Path: "myfile",
						Bytes: []byte(`apiVersion: v1
kind: Deployment`),
					},
					{
						Path: "myfile",
						Bytes: []byte(`apiVersion: v2
kind: CronJob
`),
					},
				},
				Errors: []error{},
			},
		},
		{
			Have: have{
				Path: "myfile",
				Reader: strings.NewReader(`apiVersion: v1
kind: ReplicationController
---
apiVersion: v1
kind: Deployment
`),
			},
			Want: want{
				Resources: []resource.Resource{
					{
						Path: "myfile",
						Bytes: []byte(`apiVersion: v1
kind: ReplicationController`),
					},
					{
						Path: "myfile",
						Bytes: []byte(`apiVersion: v1
kind: Deployment
`),
					},
				},
				Errors: []error{},
			},
		},
	}

	for testi, testCase := range testCases {
		ctx := context.Background()
		resChan, errChan := resource.FromStream(ctx, testCase.Have.Path, testCase.Have.Reader)
		var wg sync.WaitGroup

		wg.Add(2)
		go func() {
			res := []resource.Resource{}
			for r := range resChan {
				res = append(res, r)
			}

			if len(testCase.Want.Resources) != len(res) {
				t.Errorf("test %d - expected %d resources, got %d", testi, len(testCase.Want.Resources), len(res))
			}
			for i, v := range res {
				if !bytes.Equal(v.Bytes, testCase.Want.Resources[i].Bytes) {
					t.Errorf("test %d - for resource %d, got '%s', expected '%s'", testi, i, string(res[i].Bytes), string(testCase.Want.Resources[i].Bytes))
				}
			}

			wg.Done()
		}()

		go func() {
			errs := []error{}
			for e := range errChan {
				errs = append(errs, e)
			}
			if reflect.DeepEqual(testCase.Want.Errors, errs) == false {
				t.Errorf("expected error %+s, got %+s", testCase.Want.Errors, errs)
			}
			wg.Done()
		}()

		wg.Wait()
	}
}

// TestFromStreamResourceBytesNotAliased is a regression test for a data race in
// FromStream: each Resource was created with Bytes set to scanner.Bytes(), which
// aliases the bufio.Scanner's internal buffer. Resources are emitted on a channel
// and consumed concurrently with the producer goroutine; the scanner refills that
// buffer for later documents, overwriting the slices earlier Resources still point
// at. Callers that retain or parse resources as they arrive (e.g.
// ValidateWithContext) then read corrupted bytes — a different document spliced
// into the one being parsed — failing a different resource on every run.
//
// The stream is deliberately larger than the scanner's 4MB initial buffer so it
// must shift/refill mid-stream (the condition that clobbers earlier tokens). Each
// document carries a unique marker; after draining every resource, all markers
// must still be intact. Run with -race to also catch the underlying data race.
func TestFromStreamResourceBytesNotAliased(t *testing.T) {
	const (
		docs       = 400
		fillerSize = 16 * 1024 // docs*fillerSize (~6.4MB) exceeds the 4MB scanner buffer
	)

	var stream strings.Builder
	for i := 0; i < docs; i++ {
		fmt.Fprintf(
			&stream,
			"---\napiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: cm-%d\ndata:\n  filler: %q\n",
			i, strings.Repeat("x", fillerSize),
		)
	}

	resChan, errChan := resource.FromStream(context.Background(), "stream", strings.NewReader(stream.String()))

	go func() {
		for range errChan { //nolint:revive // drain the error channel
		}
	}()

	// Drain every resource; the scanner keeps refilling its buffer as we do,
	// clobbering any earlier Resource.Bytes that still alias it.
	var all []resource.Resource
	for res := range resChan {
		all = append(all, res)
	}

	if len(all) != docs {
		t.Fatalf("expected %d resources, got %d", docs, len(all))
	}

	for i, res := range all {
		marker := fmt.Sprintf("name: cm-%d\n", i)
		if !bytes.Contains(res.Bytes, []byte(marker)) {
			t.Fatalf(
				"resource %d: bytes were overwritten (Resource.Bytes aliased the scanner buffer, "+
					"clobbered when the scanner refilled); expected them to contain %q, got first 120 bytes:\n%s",
				i, marker, string(res.Bytes[:min(120, len(res.Bytes))]),
			)
		}
	}
}
