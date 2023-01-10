package config

import (
	"reflect"
	"testing"
)

func TestSkipKindMaps(t *testing.T) {
	for _, testCase := range []struct {
		name         string
		csvSkipKinds string
		expect       map[string]struct{}
	}{
		{
			"nothing to skip",
			"",
			map[string]struct{}{},
		},
		{
			"a single kind to skip",
			"somekind",
			map[string]struct{}{
				"somekind": {},
			},
		},
		{
			"multiple kinds to skip",
			"somekind,anotherkind,yetsomeotherkind",
			map[string]struct{}{
				"somekind":         {},
				"anotherkind":      {},
				"yetsomeotherkind": {},
			},
		},
	} {
		got := splitCSV(testCase.csvSkipKinds)
		if !reflect.DeepEqual(got, testCase.expect) {
			t.Errorf("%s - got %+v, expected %+v", testCase.name, got, testCase.expect)
		}
	}
}

func TestFromFlags(t *testing.T) {
	testCases := []struct {
		args []string
		conf Config
	}{
		{
			[]string{},
			Config{
				Files:             []string{},
				KubernetesVersion: "master",
				NumberOfWorkers:   4,
				OutputFormat:      "text",
				SchemaLocations:   nil,
				SkipKinds:         map[string]struct{}{},
				RejectKinds:       map[string]struct{}{},
			},
		},
		{
			[]string{"-h"},
			Config{
				Files:             []string{},
				Help:              true,
				KubernetesVersion: "master",
				NumberOfWorkers:   4,
				OutputFormat:      "text",
				SchemaLocations:   nil,
				SkipKinds:         map[string]struct{}{},
				RejectKinds:       map[string]struct{}{},
			},
		},
		{
			[]string{"-v"},
			Config{
				Files:             []string{},
				Version:           true,
				KubernetesVersion: "master",
				NumberOfWorkers:   4,
				OutputFormat:      "text",
				SchemaLocations:   nil,
				SkipKinds:         map[string]struct{}{},
				RejectKinds:       map[string]struct{}{},
			},
		},
		{
			[]string{"-skip", "a,b,c"},
			Config{
				Files:             []string{},
				KubernetesVersion: "master",
				NumberOfWorkers:   4,
				OutputFormat:      "text",
				SchemaLocations:   nil,
				SkipKinds:         map[string]struct{}{"a": {}, "b": {}, "c": {}},
				RejectKinds:       map[string]struct{}{},
			},
		},
		{
			[]string{"-summary", "-verbose", "file1", "file2"},
			Config{
				Files:             []string{"file1", "file2"},
				KubernetesVersion: "master",
				NumberOfWorkers:   4,
				OutputFormat:      "text",
				SchemaLocations:   nil,
				SkipKinds:         map[string]struct{}{},
				RejectKinds:       map[string]struct{}{},
				Summary:           true,
				Verbose:           true,
			},
		},
		{
			[]string{"-cache", "cache", "-ignore-missing-schemas", "-kubernetes-version", "1.16.0", "-n", "2", "-output", "json",
				"-schema-location", "folder", "-schema-location", "anotherfolder", "-skip", "kinda,kindb", "-strict",
				"-reject", "kindc,kindd", "-summary", "-debug", "-verbose", "file1", "file2"},
			Config{
				Cache:                "cache",
				Debug:                true,
				Files:                []string{"file1", "file2"},
				IgnoreMissingSchemas: true,
				KubernetesVersion:    "1.16.0",
				NumberOfWorkers:      2,
				OutputFormat:         "json",
				SchemaLocations:      []string{"folder", "anotherfolder"},
				SkipKinds:            map[string]struct{}{"kinda": {}, "kindb": {}},
				RejectKinds:          map[string]struct{}{"kindc": {}, "kindd": {}},
				Strict:               true,
				Summary:              true,
				Verbose:              true,
			},
		},
		{
			[]string{"file1,file2,file3"},
			Config{
				Files:             []string{"file1", "file2", "file3"},
				KubernetesVersion: "master",
				NumberOfWorkers:   4,
				OutputFormat:      "text",
				SchemaLocations:   nil,
				SkipKinds:         map[string]struct{}{},
				RejectKinds:       map[string]struct{}{},
			},
		},
	}

	for i, testCase := range testCases {
		cfg, _, _ := FromFlags("kubeconform", testCase.args)
		if reflect.DeepEqual(cfg, testCase.conf) != true {
			t.Errorf("test %d: failed parsing config - expected , got: \n%+v\n%+v", i, testCase.conf, cfg)
		}
	}
}
