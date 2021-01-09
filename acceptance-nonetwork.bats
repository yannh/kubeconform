#!/usr/bin/env bats

@test "Fail when parsing a valid Kubernetes config YAML file without network access" {
  run bin/kubeconform fixtures/valid.yaml
  [ "$status" -eq 1 ]
}

@test "Pass when parsing a valid config YAML file without network access, with cache" {
  run bin/kubeconform -cache fixtures/cache/ fixtures/valid.yaml
  [ "$status" -eq 0 ]
}

@test "Pass when parsing a Custom Resource and using a local schema registry with appropriate CRD" {
  run bin/kubeconform -schema-location './fixtures/registry/{{ .ResourceKind }}{{ .KindSuffix }}.json' fixtures/test_crd.yaml
  [ "$status" -eq 0 ]
}

@test "Pass when parsing a Custom Resource and specifying several local registries, the last one having the appropriate CRD" {
  run bin/kubeconform -schema-location 'fixtures/{{ .ResourceKind }}.json' -schema-location './fixtures/registry/{{ .ResourceKind }}{{ .KindSuffix }}.json' fixtures/test_crd.yaml
  [ "$status" -eq 0 ]
}
