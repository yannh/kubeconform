#!/usr/bin/env bats

@test "Pass when parsing a valid Kubernetes config YAML file" {
  run bin/kubeconform -file fixtures/valid.yaml -summary
  [ "$status" -eq 0 ]
  [ "$output" = "Run summary - Valid: 1, Invalid: 0, Errors: 0 Skipped: 0" ]
}