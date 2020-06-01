package output

import (
	"bytes"
	"testing"
)

func TestTextWrite(t *testing.T) {
	type result struct {
		fileName, kind, version  string
		err error
		skipped bool
	}

	for _, testCase := range []struct {
		name string
		withSummary bool
		verbose bool

		res  []result
		expect string
	} {
		{
			"a single deployment, no summary, no verbose",
			false,
			false,
			[]result{
				{
					"deployment.yml",
					"Deployment",
					"apps/v1",
					nil,
					false,
				},
			},
			"",
		},
		{
			"a single deployment, summary, no verbose",
			true,
			false,
			[]result{
				{
					"deployment.yml",
					"Deployment",
					"apps/v1",
					nil,
					false,
				},
			},
			"Run summary - Valid: 1, Invalid: 0, Errors: 0 Skipped: 0\n",
		},
		{
			"a single deployment, verbose, with summary",
			true,
			true,
			[]result{
				{
					"deployment.yml",
					"Deployment",
					"apps/v1",
					nil,
					false,
				},
			},
			`deployment.yml - Deployment is valid
Run summary - Valid: 1, Invalid: 0, Errors: 0 Skipped: 0
`,
		},
	} {
		w := new(bytes.Buffer)
		o := Text(w, testCase.withSummary, testCase.verbose)

		for _, res := range testCase.res {
			o.Write(res.fileName, res.kind, res.version, res.err, res.skipped)
		}
		o.Flush()

		if w.String() != testCase.expect {
			t.Errorf("%s - expected: %s, got: %s", testCase.name, testCase.expect, w)
		}
	}
}