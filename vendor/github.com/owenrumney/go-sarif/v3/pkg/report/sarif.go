package report

import (
	v210 "github.com/owenrumney/go-sarif/v3/pkg/report/v210/sarif"
	v22 "github.com/owenrumney/go-sarif/v3/pkg/report/v22/sarif"
)

type Version string

func NewV210Report() *v210.Report {
	return v210.NewReport()
}

func NewV22Report() *v22.Report {
	return v22.NewReport()
}
