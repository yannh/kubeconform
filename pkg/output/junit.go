package output

// References:
// https://github.com/windyroad/JUnit-Schema/blob/master/JUnit.xsd
// https://llg.cubic.org/docs/junit/
// https://github.com/jstemmer/go-junit-report/blob/master/formatter/formatter.go
// https://github.com/junit-team/junit5/blob/main/platform-tests/src/test/resources/jenkins-junit.xsd

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"time"

	"github.com/yannh/kubeconform/pkg/validator"
)

type TestSuiteCollection struct {
	XMLName  xml.Name    `xml:"testsuites"`
	Name     string      `xml:"name,attr"`
	Time     float64     `xml:"time,attr"`
	Tests    int         `xml:"tests,attr"`
	Failures int         `xml:"failures,attr"`
	Disabled int         `xml:"disabled,attr"`
	Errors   int         `xml:"errors,attr"`
	Suites   []TestSuite `xml:"testsuite"`
}

type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type TestSuite struct {
	XMLName  xml.Name   `xml:"testsuite"`
	Cases    []TestCase `xml:"testcase"`
	Name     string     `xml:"name,attr"`
	Id       int        `xml:"id,attr"`
	Tests    int        `xml:"tests,attr"`
	Failures int        `xml:"failures,attr"`
	Errors   int        `xml:"errors,attr"`
	Disabled int        `xml:"disabled,attr"`
	Skipped  int        `xml:"skipped,attr"`
}

type TestCase struct {
	XMLName   xml.Name         `xml:"testcase"`
	Name      string           `xml:"name,attr"`
	ClassName string           `xml:"classname,attr"`
	Time      int              `xml:"time,attr"` // Optional, but for Buildkite support  https://github.com/yannh/kubeconform/issues/127
	Skipped   *TestCaseSkipped `xml:"skipped,omitempty"`
	Error     *TestCaseError   `xml:"error,omitempty"`
	Failure   []TestCaseError  `xml:"failure,omitempty"`
}

type TestCaseSkipped struct {
	Message string `xml:"message,attr"`
}

type TestCaseError struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Content string `xml:",chardata"`
}

type junito struct {
	id          int
	w           io.Writer
	withSummary bool
	verbose     bool
	suitesIndex map[string]int // map filename to index in suites
	suites      []TestSuite
	startTime   time.Time
}

func junitOutput(w io.Writer, withSummary bool, isStdin, verbose bool) Output {
	return &junito{
		id:          0,
		w:           w,
		withSummary: withSummary,
		verbose:     verbose,
		suites:      []TestSuite{},
		suitesIndex: make(map[string]int),
		startTime:   time.Now(),
	}
}

// Write adds a result to the report.
func (o *junito) Write(result validator.Result) error {
	var suite TestSuite
	i, found := o.suitesIndex[result.Resource.Path]

	if !found {
		o.id++
		suite = TestSuite{
			Name:  result.Resource.Path,
			Id:    o.id,
			Tests: 0, Failures: 0, Errors: 0, Disabled: 0, Skipped: 0,
			Cases: make([]TestCase, 0),
		}
		o.suites = append(o.suites, suite)
		i = len(o.suites) - 1
		o.suitesIndex[result.Resource.Path] = i
	}

	sig, _ := result.Resource.Signature()
	var objectName string
	if len(sig.Namespace) > 0 {
		objectName = fmt.Sprintf("%s/%s", sig.Namespace, sig.Name)
	} else {
		objectName = sig.Name
	}
	typeName := fmt.Sprintf("%s@%s", sig.Kind, sig.Version)
	testCase := TestCase{ClassName: typeName, Name: objectName}

	switch result.Status {
	case validator.Valid:
	case validator.Invalid:
		o.suites[i].Failures++
		failure := TestCaseError{Message: result.Err.Error()}
		testCase.Failure = append(testCase.Failure, failure)
	case validator.Error:
		o.suites[i].Errors++
		testCase.Error = &TestCaseError{Message: result.Err.Error()}
	case validator.Skipped:
		testCase.Skipped = &TestCaseSkipped{}
		o.suites[i].Skipped++
	case validator.Empty:
		return nil
	}

	o.suites[i].Tests++
	o.suites[i].Cases = append(o.suites[i].Cases, testCase)

	return nil
}

// Flush outputs the results as XML
func (o *junito) Flush() error {
	runtime := time.Now().Sub(o.startTime)

	totalValid := 0
	totalInvalid := 0
	totalErrors := 0
	totalSkipped := 0

	for _, suite := range o.suites {
		for _, tCase := range suite.Cases {
			if tCase.Error != nil {
				totalErrors++
			} else if tCase.Skipped != nil {
				totalSkipped++
			} else if len(tCase.Failure) > 0 {
				totalInvalid++
			} else {
				totalValid++
			}
		}
	}

	root := TestSuiteCollection{
		Name:     "kubeconform",
		Time:     runtime.Seconds(),
		Tests:    totalValid + totalInvalid + totalErrors + totalSkipped,
		Failures: totalInvalid,
		Errors:   totalErrors,
		Disabled: totalSkipped,
		Suites:   o.suites,
	}

	// 2-space indentation
	content, err := xml.MarshalIndent(root, "", "  ")

	if err != nil {
		return err
	}

	writer := bufio.NewWriter(o.w)
	writer.Write(content)
	writer.WriteByte('\n')
	writer.Flush()

	return nil
}
