package main

import (
	"reflect"
	"testing"
)

func TestSkipKindMaps(t *testing.T) {
	for _, testCase := range []struct {
		name         string
		csvSkipKinds string
		expect       map[string]bool
	}{
		{
			"nothing to skip",
			"",
			map[string]bool{},
		},
		{
			"a single kind to skip",
			"somekind",
			map[string]bool{
				"somekind": true,
			},
		},
		{
			"multiple kinds to skip",
			"somekind,anotherkind,yetsomeotherkind",
			map[string]bool{
				"somekind":         true,
				"anotherkind":      true,
				"yetsomeotherkind": true,
			},
		},
	} {
		got := skipKindsMap(testCase.csvSkipKinds)
		if !reflect.DeepEqual(got, testCase.expect) {
			t.Errorf("%s - got %+v, expected %+v", testCase.name, got, testCase.expect)
		}
	}
}
