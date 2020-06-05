# Kubeconform

[![Build status](https://github.com/yannh/kubeconform/workflows/build/badge.svg?branch=master)](https://github.com/yannh/kubeconform/actions?query=branch%3Amaster)
[![Go Report card](https://goreportcard.com/badge/github.com/yannh/kubeconform)](https://goreportcard.com/report/github.com/yannh/kubeconform)

Kubeconform is a Kubernetes manifests validation tool. Build it into your CI to validate your Kubernetes
configuration using the schemas from the registry maintained by the
[kubernetes-json-schema](https://github.com/instrumenta/kubernetes-json-schema) project!

It is inspired by and similar to [Kubeval](https://github.com/instrumenta/kubeval), but with the
following improvements:
 * **high performance**: will validate & download manifests over multiple routines
 * support for **Kubernetes CRDs**

### Usage

```
$ ./bin/kubeconform -h
Usage of ./bin/kubeconform:
  -dir value
        directory to validate (can be specified multiple times)
  -ignore-missing-schemas
        skip files with missing schemas instead of failing
  -k8sversion string
        version of Kubernetes to test against (default "1.18.0")
  -n int
        number of routines to run in parallel (default 4)
  -output string
        output format - text, json (default "text")
  -schema value
        file containing an additional Schema (can be specified multiple times)
  -skip string
        comma-separated list of kinds to ignore
  -strict
        disallow additional properties not in schema
  -summary
        print a summary at the end
  -verbose
        print results for all resources
```

### Usage examples

* Validating a single, valid file
```
$ ./bin/kubeconform -file fixtures/valid.yaml
$ echo $?
0
```

* Validating a single invalid file, setting output to json, and printing a summary
```
$ ./bin/kubeconform -file fixtures/invalid.yaml -summary -output json
{
  "resources": [
    {
      "filename": "fixtures/invalid.yaml",
      "kind": "ReplicationController",
      "version": "v1",
      "status": "INVALID",
      "msg": "Additional property templates is not allowed - Invalid type. Expected: [integer,null], given: string"
    }
  ],
  "summary": {
    "valid": 0,
    "invalid": 1,
    "errors": 0,
    "skipped": 0
  }
}
$ echo $?
1
```

* Validating a folder, increasing the number of parallel workers
```
$ ./bin/kubeconform -dir fixtures -summary -n 16
fixtures/multi_invalid.yaml - Service is invalid: Invalid type. Expected: integer, given: string
fixtures/invalid.yaml - ReplicationController is invalid: Invalid type. Expected: [integer,null], given: string
Run summary - Valid: 40, Invalid: 2, Errors: 0 Skipped: 6
```

* Validating a custom resources, using a local schema

```
$ bin/kubeconform -file fixtures/test_crd.yaml -schema fixtures/crd_schema.yaml -verbose
fixtures/test_crd.yaml - TrainingJob is valid
```

### Credits

 * @garethr for the [Kubeval](https://github.com/instrumenta/kubeval) and
 [kubernetes-json-schema](https://github.com/instrumenta/kubernetes-json-schema) projects
