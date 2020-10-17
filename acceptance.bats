#!/usr/bin/env bats

@test "Pass when parsing a valid Kubernetes config YAML file" {
  run bin/kubeconform -summary fixtures/valid.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "Summary: 1 resource found in 1 file - Valid: 1, Invalid: 0, Errors: 0 Skipped: 0" ]
}

@test "Pass when parsing a folder containing valid YAML files" {
  run bin/kubeconform -summary fixtures/folder
  [ "$status" -eq 0 ]
  [ "$output" = "Summary: 7 resources found in 2 files - Valid: 7, Invalid: 0, Errors: 0 Skipped: 0" ]
}

@test "Pass when parsing a Kubernetes file with string and integer quantities" {
  run bin/kubeconform -verbose fixtures/quantity.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "fixtures/quantity.yaml - LimitRange is valid" ]
}

@test "Pass when parsing a valid Kubernetes config file with null arrays" {
  run bin/kubeconform -verbose fixtures/null_string.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "fixtures/null_string.yaml - Service is valid" ]
}

@test "Pass when parsing a multi-document config file" {
  run bin/kubeconform -summary fixtures/multi_valid.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "Summary: 6 resources found in 1 file - Valid: 6, Invalid: 0, Errors: 0 Skipped: 0" ]
}

@test "Fail when parsing a multi-document config file with one invalid resource" {
  run bin/kubeconform fixtures/multi_invalid.yaml
  [ "$status" -eq 1 ]
}

@test "Fail when parsing an invalid Kubernetes config file" {
  run bin/kubeconform fixtures/invalid.yaml
  [ "$status" -eq 1 ]
}

@test "Return relevant error for non-existent file" {
  run bin/kubeconform fixtures/not-here
  [ "$status" -eq 1 ]
  [ "$output" = "fixtures/not-here -  failed validation: open fixtures/not-here: no such file or directory" ]
}

@test "Fail when parsing a config with additional properties and strict set" {
  run bin/kubeconform -strict -k8sversion 1.16.0 fixtures/extra_property.yaml
  [ "$status" -eq 1 ]
}

@test "Fail when parsing a config with CRD" {
  run bin/kubeconform fixtures/test_crd.yaml
  [ "$status" -eq 1 ]
}

@test "Pass when parsing a config with CRD and ignoring missing schemas" {
  run bin/kubeconform -ignore-missing-schemas fixtures/test_crd.yaml
  [ "$status" -eq 0 ]
}

