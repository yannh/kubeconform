#!/usr/bin/env bats

@test "Should generate expected prometheus resource" {
  run ./openapi2jsonschema.py fixtures/prometheus-operator-0prometheusCustomResourceDefinition.yaml
  [ "$status" -eq 0 ]
  [ "$output" = "JSON schema written to prometheus_v1.json" ]
  run diff prometheus_v1.json ./fixtures/prometheus_v1-expected.json
  [ "$status" -eq 0 ]
}
