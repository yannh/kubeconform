package resource_test

import (
	"bytes"
	"context"
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
