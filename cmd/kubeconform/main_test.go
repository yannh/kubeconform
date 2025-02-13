package main

import (
	"context"
	"testing"

	"github.com/yannh/kubeconform/pkg/validator"
)

type MockOutput struct{}

func (m *MockOutput) Write(res validator.Result) error {
	return nil
}

func (m *MockOutput) Flush() error {
	return nil
}

// Test generated using Keploy
func TestProcessResults_SuccessWithValidResults(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	validationResults := make(chan validator.Result, 2)
	validationResults <- validator.Result{Status: validator.Valid}
	validationResults <- validator.Result{Status: validator.Valid}
	close(validationResults)

	outputMock := &MockOutput{}
	exitOnError := false

	resultChan := processResults(cancel, outputMock, validationResults, exitOnError)
	success := <-resultChan

	if !success {
		t.Errorf("Expected success to be true, got false")
	}

	if ctx.Err() != nil {
		t.Errorf("Context should not be canceled")
	}
}

// Test generated using Keploy
func TestProcessResults_FailureWithInvalidResults(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	validationResults := make(chan validator.Result, 2)
	validationResults <- validator.Result{Status: validator.Valid}
	validationResults <- validator.Result{Status: validator.Invalid}
	close(validationResults)

	outputMock := &MockOutput{}
	exitOnError := false

	resultChan := processResults(cancel, outputMock, validationResults, exitOnError)
	success := <-resultChan

	if success {
		t.Errorf("Expected success to be false, got true")
	}

	if ctx.Err() != nil {
		t.Errorf("Context should not be canceled when exitOnError is false")
	}
}

// Test generated using Keploy
func TestProcessResults_CancelOnInvalidWithExitOnError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	validationResults := make(chan validator.Result, 2)
	validationResults <- validator.Result{Status: validator.Valid}
	validationResults <- validator.Result{Status: validator.Invalid}
	close(validationResults)

	outputMock := &MockOutput{}
	exitOnError := true

	resultChan := processResults(cancel, outputMock, validationResults, exitOnError)
	success := <-resultChan

	if success {
		t.Errorf("Expected success to be false, got true")
	}

	select {
	case <-ctx.Done():
	default:
		t.Errorf("Expected context to be canceled when exitOnError is true")
	}
}
