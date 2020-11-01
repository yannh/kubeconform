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

@test "Pass when parsing a valid Kubernetes config file with int_to_string vars" {
  run bin/kubeconform -verbose fixtures/int_or_string.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "fixtures/int_or_string.yaml - Service heapster is valid" ]
}

@test "Pass when parsing a valid Kubernetes config JSON file" {
  run bin/kubeconform -kubernetes-version 1.17.1 -summary fixtures/valid.json
  [ "$status" -eq 0 ]
  [ "$output" = "Summary: 1 resource found in 1 file - Valid: 1, Invalid: 0, Errors: 0 Skipped: 0" ]
}

@test "Pass when parsing a valid Kubernetes config YAML file with generate name" {
  run bin/kubeconform -verbose fixtures/generate_name.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "fixtures/generate_name.yaml - Job pi-{{ generateName }} is valid" ]
}

@test "Pass when parsing a Kubernetes file with string and integer quantities" {
  run bin/kubeconform -verbose fixtures/quantity.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "fixtures/quantity.yaml - LimitRange mem-limit-range is valid" ]
}

@test "Pass when parsing a valid Kubernetes config file with null arrays" {
  run bin/kubeconform -verbose fixtures/null_string.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "fixtures/null_string.yaml - Service frontend is valid" ]
}

@test "Pass when parsing a valid Kubernetes config file with null strings" {
  run bin/kubeconform -summary fixtures/null_string.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "Summary: 1 resource found in 1 file - Valid: 1, Invalid: 0, Errors: 0 Skipped: 0" ]
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
  [ "$output" = "fixtures/not-here - failed validation: lstat fixtures/not-here: no such file or directory" ]
}

@test "Pass when parsing a blank config file" {
   run bin/kubeconform -summary fixtures/blank.yaml
   [ "$status" -eq 0 ]
   [ "$output" = "Summary: 0 resource found in 1 file - Valid: 0, Invalid: 0, Errors: 0 Skipped: 0" ]
}

@test "Pass when parsing a blank config file with a comment" {
   run bin/kubeconform -summary fixtures/comment.yaml
   [ "$status" -eq 0 ]
   [ "$output" = "Summary: 0 resource found in 1 file - Valid: 0, Invalid: 0, Errors: 0 Skipped: 0" ]
}

@test "Fail when parsing a config with additional properties and strict set" {
  run bin/kubeconform -strict -kubernetes-version 1.16.0 fixtures/extra_property.yaml
  [ "$status" -eq 1 ]
}

@test "Fail when parsing a config with CRD" {
  run bin/kubeconform fixtures/test_crd.yaml
  [ "$status" -eq 1 ]
}

@test "Pass when parsing a config with Custom Resource and ignoring missing schemas" {
  run bin/kubeconform -ignore-missing-schemas fixtures/test_crd.yaml
  [ "$status" -eq 0 ]
}

@test "Pass when parsing a Custom Resource and using a local schema registry with appropriate CRD" {
  run bin/kubeconform -schema-location './fixtures/registry/{{ .ResourceKind }}{{ .KindSuffix }}.json' fixtures/test_crd.yaml
  [ "$status" -eq 0 ]
}

@test "Pass when parsing a config with additional properties" {
  run bin/kubeconform -summary fixtures/extra_property.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "Summary: 1 resource found in 1 file - Valid: 1, Invalid: 0, Errors: 0 Skipped: 0" ]
}

@test "Pass when using a valid, preset -schema-location" {
  run bin/kubeconform -schema-location https://kubernetesjsonschema.dev fixtures/valid.yaml
  [ "$status" -eq 0 ]
}

@test "Pass when using a valid HTTP -schema-location" {
  run bin/kubeconform -schema-location 'https://kubernetesjsonschema.dev/{{ .NormalizedVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json' fixtures/valid.yaml
  [ "$status" -eq 0 ]
}

@test "Fail when using an invalid HTTP -schema-location" {
  run bin/kubeconform -schema-location 'http://foo' fixtures/valid.yaml
  [ "$status" -eq 1 ]
}

@test "Fail when using an invalid non-HTTP -schema-location" {
  run bin/kubeconform -schema-location 'foo' fixtures/valid.yaml
  [ "$status" -eq 1 ]
}

@test "Pass with a valid input when validating against openshift manifests" {
  run bin/kubeconform -kubernetes-version 3.8.0  -schema-location 'https://raw.githubusercontent.com/garethr/openshift-json-schema/master/{{ .NormalizedVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}.json'  -summary fixtures/valid.yaml
  [ "$status" -eq 0 ]
}

@test "Fail with an invalid input when validating against openshift manifests" {
  run bin/kubeconform -kubernetes-version 3.8.0  -schema-location 'https://raw.githubusercontent.com/garethr/openshift-json-schema/master/{{ .NormalizedVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}.json'  -summary fixtures/invalid.yaml
  [ "$status" -eq 1 ]
}

@test "Pass when parsing a valid Kubernetes config YAML file on stdin" {
  run bash -c "cat fixtures/valid.yaml | bin/kubeconform -summary"
  [ "$status" -eq 0 ]
  [ "$output" = "Summary: 1 resource found parsing stdin - Valid: 1, Invalid: 0, Errors: 0 Skipped: 0" ]
}

@test "Pass when parsing a valid Kubernetes config YAML file explicitly on stdin" {
  run bash -c "cat fixtures/valid.yaml | bin/kubeconform -summary"
  [ "$status" -eq 0 ]
  [ "$output" = "Summary: 1 resource found parsing stdin - Valid: 1, Invalid: 0, Errors: 0 Skipped: 0" ]
}

@test "Fail when parsing an invalid Kubernetes config file on stdin" {
   run bash -c "cat fixtures/invalid.yaml | bin/kubeconform -"
   [ "$status" -eq 1 ]
 }
