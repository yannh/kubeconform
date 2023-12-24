#!/usr/bin/env bats

resetCacheFolder() {
  rm -rf cache
  mkdir -p cache
}

@test "Pass when displaying help with -h" {
  run bin/kubeconform -h
  [ "$status" -eq 0 ]
  [ "${lines[0]}" == 'Usage: bin/kubeconform [OPTION]... [FILE OR FOLDER]...' ]
}

@test "Fail and display help when using an incorrect flag" {
  run bin/kubeconform -xyz
  [ "$status" -eq 1 ]
  [ "${lines[0]}" == 'flag provided but not defined: -xyz' ]
}

@test "Pass when parsing a valid Kubernetes config YAML file" {
  run bin/kubeconform -summary fixtures/valid.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "Summary: 1 resource found in 1 file - Valid: 1, Invalid: 0, Errors: 0, Skipped: 0" ]
}

@test "Pass when parsing a folder containing valid YAML files" {
  run bin/kubeconform -summary fixtures/folder
  [ "$status" -eq 0 ]
  [ "$output" = "Summary: 7 resources found in 2 files - Valid: 7, Invalid: 0, Errors: 0, Skipped: 0" ]
}

@test "Pass when parsing a valid Kubernetes config file with int_to_string vars" {
  run bin/kubeconform -verbose fixtures/int_or_string.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "fixtures/int_or_string.yaml - Service heapster is valid" ]
}

@test "Pass when parsing a valid Kubernetes config JSON file" {
  run bin/kubeconform -kubernetes-version 1.20.0 -summary fixtures/valid.json
  [ "$status" -eq 0 ]
  [ "$output" = "Summary: 1 resource found in 1 file - Valid: 1, Invalid: 0, Errors: 0, Skipped: 0" ]
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
  [ "$output" = "Summary: 1 resource found in 1 file - Valid: 1, Invalid: 0, Errors: 0, Skipped: 0" ]
}

@test "Pass when parsing a multi-document config file" {
  run bin/kubeconform -summary fixtures/multi_valid.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "Summary: 6 resources found in 1 file - Valid: 6, Invalid: 0, Errors: 0, Skipped: 0" ]
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
   [ "$output" = "Summary: 0 resource found in 1 file - Valid: 0, Invalid: 0, Errors: 0, Skipped: 0" ]
}

@test "Pass when parsing a blank config file with a comment" {
   run bin/kubeconform -summary fixtures/comment.yaml
   [ "$status" -eq 0 ]
   [ "$output" = "Summary: 0 resource found in 1 file - Valid: 0, Invalid: 0, Errors: 0, Skipped: 0" ]
}

@test "Fail when parsing a config that is missing a Kind" {
   run bin/kubeconform -summary fixtures/missing_kind.yaml
   [ "$status" -eq 1 ]
   [[ "$output" == *"missing 'kind' key"* ]]
}

@test "Fail when parsing a config that is missing an apiVersion" {
   run bin/kubeconform -summary fixtures/missing_apiversion.yaml
   [ "$status" -eq 1 ]
   [[ "$output" == *"missing 'apiVersion' key"* ]]
}

@test "Fail when parsing a config that is missing a Kind value" {
   run bin/kubeconform -summary fixtures/missing_kind_value.yaml
   [ "$status" -eq 1 ]
   [[ "$output" == *"missing 'kind' key"* ]]
}

@test "Fail when parsing a config with CRD" {
  run bin/kubeconform fixtures/test_crd.yaml
  [ "$status" -eq 1 ]
}

@test "Pass when parsing a config with Custom Resource and ignoring missing schemas" {
  run bin/kubeconform -ignore-missing-schemas fixtures/test_crd.yaml
  [ "$status" -eq 0 ]
}

@test "Pass when parsing a config with additional properties" {
  run bin/kubeconform -summary fixtures/extra_property.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "Summary: 1 resource found in 1 file - Valid: 1, Invalid: 0, Errors: 0, Skipped: 0" ]
}

@test "Fail when parsing a config with additional properties and strict set" {
  run bin/kubeconform -strict -kubernetes-version 1.20.0 fixtures/extra_property.yaml
  [ "$status" -eq 1 ]
}

@test "Fail when parsing a config with duplicate properties and strict set" {
  run bin/kubeconform -strict -kubernetes-version 1.20.0 fixtures/duplicate_property.yaml
  [ "$status" -eq 1 ]
}

@test "Pass when parsing a config with duplicate properties and strict NOT set" {
  run bin/kubeconform -kubernetes-version 1.20.0 fixtures/duplicate_property.yaml
  [ "$status" -eq 0 ]
}

@test "Pass when using a valid, preset -schema-location" {
  run bin/kubeconform -schema-location default fixtures/valid.yaml
  [ "$status" -eq 0 ]
}

@test "Pass when using a valid HTTP -schema-location" {
  run bin/kubeconform -schema-location 'https://raw.githubusercontent.com/yannh/kubernetes-json-schema/master/{{ .NormalizedKubernetesVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json' fixtures/valid.yaml
  [ "$status" -eq 0 ]
}

@test "Pass when using schemas with HTTP references" {
  run bin/kubeconform -summary -schema-location 'https://raw.githubusercontent.com/yannh/kubernetes-json-schema/master/{{ .NormalizedKubernetesVersion }}{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json' fixtures/valid.yaml
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

@test "Fail early when passing a non valid -schema-location template" {
  run bin/kubeconform -schema-location 'foo {{ .Foo }}' fixtures/valid.yaml
  [[ "$output" == "failed initialising"* ]]
  [[ `echo "$output" | wc -l` -eq 1 ]]
  [ "$status" -eq 1 ]
}

@test "Fail early when passing a non valid -kubernetes-version" {
  run bin/kubeconform -kubernetes-version 1.25 fixtures/valid.yaml
  [ "${lines[0]}" == 'invalid value "1.25" for flag -kubernetes-version: 1.25 is not a valid version. Valid values are "master" (default) or full version x.y.z (e.g. "1.27.2")' ]
  [[ "${lines[1]}" == "Usage:"* ]]
  [ "$status" -eq 1 ]
}

@test "Pass with a valid input when validating against openshift manifests" {
  run bin/kubeconform -kubernetes-version 3.8.0  -schema-location 'https://raw.githubusercontent.com/garethr/openshift-json-schema/master/{{ .NormalizedKubernetesVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}.json'  -summary fixtures/valid.yaml
  [ "$status" -eq 0 ]
}

@test "Fail with an invalid input when validating against openshift manifests" {
  run bin/kubeconform -kubernetes-version 3.8.0  -schema-location 'https://raw.githubusercontent.com/garethr/openshift-json-schema/master/{{ .NormalizedKubernetesVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}.json'  -summary fixtures/invalid.yaml
  [ "$status" -eq 1 ]
}

@test "Pass when parsing a valid Kubernetes config YAML file on stdin" {
  run bash -c "cat fixtures/valid.yaml | bin/kubeconform -summary"
  [ "$status" -eq 0 ]
  [ "$output" = "Summary: 1 resource found parsing stdin - Valid: 1, Invalid: 0, Errors: 0, Skipped: 0" ]
}

@test "Pass when parsing a valid Kubernetes config YAML file explicitly on stdin" {
  run bash -c "cat fixtures/valid.yaml | bin/kubeconform -summary -"
  [ "$status" -eq 0 ]
  [ "$output" = "Summary: 1 resource found parsing stdin - Valid: 1, Invalid: 0, Errors: 0, Skipped: 0" ]
}

@test "Fail when parsing an invalid Kubernetes config file on stdin" {
  run bash -c "cat fixtures/invalid.yaml | bin/kubeconform -"
  [ "$status" -eq 1 ]
}

@test "Fail when not passing data to stdin, when implicitly configured to read from stdin" {
  run bash -c "bin/kubeconform -summary"
  [ "$status" -eq 1 ]
}

@test "Fail when not passing data to stdin, when explicitly configured to read from stdin" {
  run bash -c "bin/kubeconform -summary -"
  [ "$status" -eq 1 ]
}

@test "Skip when parsing a resource from a kind to skip" {
  run bin/kubeconform -verbose -skip ReplicationController fixtures/valid.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "fixtures/valid.yaml - bob ReplicationController skipped" ]
}

@test "Skip when parsing a resource with a GVK to skip" {
  run bin/kubeconform -verbose -skip v1/ReplicationController fixtures/valid.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "fixtures/valid.yaml - bob ReplicationController skipped" ]
}

@test "Do not skip when parsing a resource with a GVK to skip, where the Kind matches but not the version" {
  run bin/kubeconform -verbose -skip v2/ReplicationController fixtures/valid.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "fixtures/valid.yaml - ReplicationController bob is valid" ]
}

@test "Fail when parsing a resource from a kind to reject" {
  run bin/kubeconform -verbose -reject ReplicationController fixtures/valid.yaml
  [ "$status" -eq 1 ]
  [ "$output" = "fixtures/valid.yaml - ReplicationController bob failed validation: prohibited resource kind ReplicationController" ]
}

@test "Ignores file that match the --ignore-filename-pattern given" {
  run bin/kubeconform -summary --ignore-filename-pattern 'crd' --ignore-filename-pattern '.*invalid.*' fixtures/multi_invalid.yaml fixtures/list_invalid.yaml fixtures/quantity.yaml fixtures/crd_schema.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "Summary: 1 resource found in 1 file - Valid: 1, Invalid: 0, Errors: 0, Skipped: 0" ]
}

@test "Pass when parsing a valid Kubernetes config YAML file and store cache" {
  resetCacheFolder
  run bin/kubeconform -cache cache -summary fixtures/valid.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "Summary: 1 resource found in 1 file - Valid: 1, Invalid: 0, Errors: 0, Skipped: 0" ]
  [ "`ls cache/ | wc -l`" -eq 1 ]
}

@test "Fail when no schema found, ensure 404 is not cached on disk" {
  resetCacheFolder
  run bin/kubeconform -cache cache -schema-location 'https://raw.githubusercontent.com/yannh/kubernetes-json-schema/master/doesnotexist.json' fixtures/valid.yaml
  [ "$status" -eq 1 ]
  [ "$output" == 'fixtures/valid.yaml - ReplicationController bob failed validation: could not find schema for ReplicationController' ]
  [ "`ls cache/ | wc -l`" -eq 0 ]
}

@test "Fail when cache folder does not exist" {
  run bin/kubeconform -cache cache_does_not_exist -summary fixtures/valid.yaml
  [ "$status" -eq 1 ]
  [ "$output" = "failed opening cache folder cache_does_not_exist: stat cache_does_not_exist: no such file or directory" ]
}

@test "Produces correct TAP output" {
  run bin/kubeconform -output tap fixtures/valid.yaml
  [ "$status" -eq 0 ]
  [ "${lines[0]}" == 'TAP version 13' ]
  [ "${lines[1]}" == 'ok 1 - fixtures/valid.yaml (v1/ReplicationController//bob)' ]
  [ "${lines[2]}" == '1..1' ]
}

@test "Pass when parsing a file containing a List" {
  run bin/kubeconform -summary fixtures/list_valid.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "Summary: 6 resources found in 1 file - Valid: 6, Invalid: 0, Errors: 0, Skipped: 0" ]
}

@test "Pass when parsing a List resource from stdin" {
  run bash -c "cat fixtures/list_valid.yaml | bin/kubeconform -summary"
  [ "$status" -eq 0 ]
  [ "$output" = 'Summary: 6 resources found parsing stdin - Valid: 6, Invalid: 0, Errors: 0, Skipped: 0' ]
}

@test "Fail when parsing a List that contains an invalid resource" {
  run bin/kubeconform -summary fixtures/list_invalid.yaml
  [ "$status" -eq 1 ]
  [ "${lines[0]}" == 'fixtures/list_invalid.yaml - ReplicationController bob is invalid: problem validating schema. Check JSON formatting: jsonschema: '\''/spec/replicas'\'' does not validate with https://raw.githubusercontent.com/yannh/kubernetes-json-schema/master/master-standalone/replicationcontroller-v1.json#/properties/spec/properties/replicas/type: expected integer or null, but got string' ]
  [ "${lines[1]}" == 'Summary: 2 resources found in 1 file - Valid: 1, Invalid: 1, Errors: 0, Skipped: 0' ]
}

@test "Fail when parsing a List that contains an invalid resource from stdin" {
  run bash -c "cat fixtures/list_invalid.yaml | bin/kubeconform -summary -"
  [ "$status" -eq 1 ]
  [ "${lines[0]}" == 'stdin - ReplicationController bob is invalid: problem validating schema. Check JSON formatting: jsonschema: '\''/spec/replicas'\'' does not validate with https://raw.githubusercontent.com/yannh/kubernetes-json-schema/master/master-standalone/replicationcontroller-v1.json#/properties/spec/properties/replicas/type: expected integer or null, but got string' ]
  [ "${lines[1]}" == 'Summary: 2 resources found parsing stdin - Valid: 1, Invalid: 1, Errors: 0, Skipped: 0' ]
}

@test "Pass on valid, empty list" {
  run bin/kubeconform -summary fixtures/list_empty_valid.yaml
  [ "$status" -eq 0 ]
  [ "$output" = 'Summary: 0 resource found in 1 file - Valid: 0, Invalid: 0, Errors: 0, Skipped: 0' ]
}

@test "Pass on multi-yaml containing one resource, one list" {
  run bin/kubeconform -summary fixtures/multi_with_list.yaml
  [ "$status" -eq 0 ]
  [ "$output" = 'Summary: 2 resources found in 1 file - Valid: 2, Invalid: 0, Errors: 0, Skipped: 0' ]
}

@test "Fail when using HTTPS_PROXY with a failing proxy" {
  # This only tests that the HTTPS_PROXY variable is picked up and that it tries to use it
  run bash -c "HTTPS_PROXY=127.0.0.1:1234 bin/kubeconform fixtures/valid.yaml"
  [ "$status" -eq 1 ]
  [[ "$output" == *"proxyconnect tcp: dial tcp 127.0.0.1:1234: connect: connection refused"* ]]
}

@test "Pass when parsing a very large file" {
  run bin/kubeconform -summary fixtures/valid_large.yaml
  [ "$status" -eq 0 ]
  [ "$output" = 'Summary: 100000 resources found in 1 file - Valid: 100000, Invalid: 0, Errors: 0, Skipped: 0' ]
}

@test "Pass when parsing a very long stream from stdin" {
  run bash -c "cat fixtures/valid_large.yaml | bin/kubeconform -summary"
  [ "$status" -eq 0 ]
  [ "$output" = 'Summary: 100000 resources found parsing stdin - Valid: 100000, Invalid: 0, Errors: 0, Skipped: 0' ]
}

@test "JUnit output can be validated against the Junit schema definition" {
  run bash -c "bin/kubeconform -output junit -summary fixtures/valid.yaml > output.xml"
  [ "$status" -eq 0 ]
  run xmllint --noout --schema fixtures/junit.xsd output.xml
  [ "$status" -eq 0 ]
}

@test "passes when trying to use a CRD that does not have the JSONSchema set" {
  run bash -c "bin/kubeconform -schema-location default -schema-location 'https://raw.githubusercontent.com/datreeio/CRDs-catalog/main/{{.Group}}/{{.ResourceKind}}_{{.ResourceAPIVersion}}.json' fixtures/httpproxy.yaml"
  [ "$status" -eq 0 ]
}
