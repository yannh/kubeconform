# Kubeconform

[![Build status](https://github.com/yannh/kubeconform/workflows/build/badge.svg?branch=master)](https://github.com/yannh/kubeconform/actions?query=branch%3Amaster)
[![Go Report card](https://goreportcard.com/badge/github.com/yannh/kubeconform)](https://goreportcard.com/report/github.com/yannh/kubeconform)

Kubeconform is a Kubernetes manifests validation tool. Build it into your CI to validate your Kubernetes
configuration using the schemas from the registry maintained by the
[kubernetes-json-schema](https://github.com/instrumenta/kubernetes-json-schema) project!

It is inspired by, contains code from and is designed to stay close to
[Kubeval](https://github.com/instrumenta/kubeval), but with the following improvements:
 * **high performance**: will validate & download manifests over multiple routines, caching
   downloaded files in memory
 * configurable list of schemas locations, enabling validating Kubernetes custom resources (CRDs)

### A small overview of Kubernetes manifest validation

Kubernetes's API is described using the [OpenAPI (formerly swagger) specification](https://www.openapis.org),
in a [file](https://github.com/kubernetes/kubernetes/blob/master/api/openapi-spec/swagger.json) checked into
the main Kubernetes repository.

Because of the state of the tooling to perform validation against OpenAPI schemas, projects usually convert
the OpenAPI schemas to [JSON schemas](https://json-schema.org/) first. Kubeval relies on 
[instrumenta/OpenApi2JsonSchema](https://github.com/instrumenta/openapi2jsonschema) to convert Kubernetes' Swagger file
and break it down into multiple JSON schemas, stored in github at
[instrumenta/kubernetes-json-schema](https://github.com/instrumenta/kubernetes-json-schema) and published on
[kubernetesjsonschema.dev](https://kubernetesjsonschema.dev/).

Kubeconform relies on the same JSON schemas from kubernetesjsonschema.dev, and will download required
schemas at runtime as required.

### Usage

```
$ ./bin/kubeconform -h
Usage of ./bin/kubeconform:
  -ignore-missing-schemas
        skip files with missing schemas instead of failing
  -kubernetes-sversion string
        version of Kubernetes to validate against (default "1.18.0")
  -n int
        number of routines to run in parallel (default 4)
  -output string
        output format - text, json (default "text")
  -schema-location value
        override schemas location search path (can be specified multiple times)
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
$ ./bin/kubeconform fixtures/valid.yaml
$ echo $?
0
```

* Validating a single invalid file, setting output to json, and printing a summary
```
$ ./bin/kubeconform -summary -output json fixtures/invalid.yaml
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
$ ./bin/kubeconform -summary -n 16 fixtures
fixtures/crd_schema.yaml - CustomResourceDefinition trainingjobs.sagemaker.aws.amazon.com failed validation: could not find schema for CustomResourceDefinition
fixtures/invalid.yaml - ReplicationController bob is invalid: Invalid type. Expected: [integer,null], given: string
[...]
Summary: 65 resources found in 34 files - Valid: 55, Invalid: 2, Errors: 8 Skipped: 0
```

### Overriding schemas location - CRD support

When the `-schema-location` parameter is not used, kubeconform will default to downloading schemas from
`https://kubernetesjsonschema.dev`. Kubeconform however supports passing one, or multiple, schemas
locations - HTTP URLs, or local filesystem paths, in which case it will lookup for schema definitions
in each of them, in order, stopping as soon as a matching file is found.

 * If the -schema-location value does not end with '.json', Kubeconform will assume filenames / a file
 structure identical to that of kubernetesjsonschema.dev
 * if the -schema-location value ends with '.json' - Kubeconform assumes the value is a Go templated
 string that indicates how to search for JSON schemas.

All 3 following command lines are equivalent:
```
$ ./bin/kubeconform fixtures/valid.yaml
$ ./bin/kubeconform -schema-location https://kubernetesjsonschema.dev fixtures/valid.yaml
$ ./bin/kubeconform -schema-location 'https://kubernetesjsonschema.dev/{{ .NormalizedVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json' fixtures/valid.yaml
```

To support validating CRDs, we need to convert OpenAPI files to JSON schema, storing the JSON schemas
in a local folder - for example schemas. Then we specify this folder as an additional registry to lookup:

```
# If the resource Kind is not found in kubernetesjsonschema.dev, also lookup in the schemas/ folder for a matching file
$ ./bin/kubeconform -registry kubernetesjsonschema.dev -schema-location 'schemas/{{ .ResourceKind }}{{ .KindSuffix }}.json' fixtures/custom-resource.yaml
```

### Generating a JSON schema from an OpenAPI file

Kubeconform uses JSON schemas to validate Kubernetes resources. For Custom Resource, the CustomResourceDefinition
first needs to be converted to JSON Schema. A script is provided to convert these CustomResourceDefinitions 
to JSON schema. Here is an example how to use it:

```
$ ./cmd/openapi2jsonschema/main.py https://raw.githubusercontent.com/aws/amazon-sagemaker-operator-for-k8s/master/config/crd/bases/sagemaker.aws.amazon.com_trainingjobs.yaml > fixtures/registry/trainingjob-sagemaker-v1.json
```

### Credits

 * @garethr for the [Kubeval](https://github.com/instrumenta/kubeval) and
 [kubernetes-json-schema](https://github.com/instrumenta/kubernetes-json-schema) projects ❤️
