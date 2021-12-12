#!/usr/bin/env bats

setup() {
  rm -f prometheus_v1.json
  rm -f prometheus-monitoring-v1.json
}

@test "Should generate expected prometheus resource" {
  run ./openapi2jsonschema.py fixtures/prometheus-operator-0prometheusCustomResourceDefinition.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "JSON schema written to prometheus_v1.json" ]
  run diff prometheus_v1.json ./fixtures/prometheus_v1-expected.json
  [ "$status" -eq 0 ]
}

@test "Should generate expected prometheus resource from an HTTP resource" {
  run ./openapi2jsonschema.py https://raw.githubusercontent.com/yannh/kubeconform/aebc298047c386116eeeda9b1ada83671a58aedd/scripts/fixtures/prometheus-operator-0prometheusCustomResourceDefinition.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "JSON schema written to prometheus_v1.json" ]
  run diff prometheus_v1.json ./fixtures/prometheus_v1-expected.json
  [ "$status" -eq 0 ]
}

@test "Should output filename in {kind}-{group}-{version} format" {
  FILENAME_FORMAT='{kind}-{group}-{version}' run ./openapi2jsonschema.py fixtures/prometheus-operator-0prometheusCustomResourceDefinition.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "JSON schema written to prometheus-monitoring-v1.json" ]
  run diff prometheus-monitoring-v1.json ./fixtures/prometheus_v1-expected.json
  [ "$status" -eq 0 ]
}

@test "Should set 'additionalProperties: false' at the root" {
  DENY_ROOT_ADDITIONAL_PROPERTIES='true' run ./openapi2jsonschema.py fixtures/prometheus-operator-0prometheusCustomResourceDefinition.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "JSON schema written to prometheus_v1.json" ]
  run diff prometheus_v1.json ./fixtures/prometheus_v1-denyRootAdditionalProperties.json
  [ "$status" -eq 0 ]
}
