#!/usr/bin/env bats

@test "Pass when parsing a valid Kubernetes config YAML file" {
  run bin/kubeconform -file fixtures/valid.yaml -summary
  [ "$status" -eq 0 ]
  [ "$output" = "Run summary - Valid: 1, Invalid: 0, Errors: 0 Skipped: 0" ]
}

@test "Pass when parsing a Kubernetes file with string and integer quantities" {
  run bin/kubeconform -verbose -file fixtures/quantity.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "fixtures/quantity.yaml - LimitRange is valid" ]
}

@test "Pass when parsing a valid Kubernetes config file with null arrays" {
  run bin/kubeconform -verbose -file fixtures/null_string.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "fixtures/null_string.yaml - Service is valid" ]
}

@test "Pass when parsing a multi-document config file" {
  run bin/kubeconform -summary -file fixtures/multi_valid.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "Run summary - Valid: 6, Invalid: 0, Errors: 0 Skipped: 0" ]
}

@test "Fail when parsing a multi-document config file with one invalid resource" {
  run bin/kubeconform -file fixtures/multi_invalid.yaml
  [ "$status" -eq 1 ]
}

@test "Fail when parsing an invalid Kubernetes config file" {
  run bin/kubeconform -file fixtures/invalid.yaml
  [ "$status" -eq 1 ]
}

@test "Return relevant error for non-existent file" {
  run bin/kubeconform -file fixtures/not-here
  [ "$status" -eq 1 ]
  [ $(expr "$output" : "^failed opening fixtures/not-here") -ne 0 ]
}
