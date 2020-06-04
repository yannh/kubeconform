#!/usr/bin/env bats

@test "Pass when parsing a valid Kubernetes config YAML file" {
  run bin/kubeconform -file fixtures/valid.yaml -summary
  [ "$status" -eq 0 ]
  [ "$output" = "Summary: 1 resource found in 1 file - Valid: 1, Invalid: 0, Errors: 0 Skipped: 0" ]
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
  [ "$output" = "Summary: 6 resources found in 1 file - Valid: 6, Invalid: 0, Errors: 0 Skipped: 0" ]
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

@test "Fail when parsing a config with additional properties and strict set" {
  run bin/kubeconform -strict -k8sversion 1.16.0 -file fixtures/extra_property.yaml
  [ "$status" -eq 1 ]
}

@test "Fail when parsing a config with CRD" {
  run bin/kubeconform -file fixtures/test_crd.yaml
  [ "$status" -eq 1 ]
}

@test "Pass when parsing a config with CRD and ignoring missing schemas" {
  run bin/kubeconform -file fixtures/test_crd.yaml -ignore-missing-schemas
  [ "$status" -eq 0 ]
}

@test "Succeed parsing a CRD when additional schema passed" {
  run bin/kubeconform -file fixtures/test_crd.yaml -schema fixtures/crd_schema.yaml
  [ "$status" -eq 0 ]
}
