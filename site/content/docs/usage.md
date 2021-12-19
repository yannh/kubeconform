---
title: "Usage"
date: 2021-07-02T00:00:00Z
draft: false
tags: ["Kubeconform", "Usage"]
weight: 3
---

{{< prism >}}$ ./bin/kubeconform -h
Usage: ./bin/kubeconform [OPTION]... [FILE OR FOLDER]...
  -cache string
        cache schemas downloaded via HTTP to this folder
  -cpu-prof string
        debug - log CPU profiling to file
  -exit-on-error
        immediately stop execution when the first error is encountered
  -h    show help information
  -ignore-filename-pattern value
        regular expression specifying paths to ignore (can be specified multiple times)
  -ignore-missing-schemas
        skip files with missing schemas instead of failing
  -insecure-skip-tls-verify
        disable verification of the server's SSL certificate. This will make your HTTPS connections insecure
  -kubernetes-version string
        version of Kubernetes to validate against, e.g.: 1.18.0 (default "master")
  -n int
        number of goroutines to run concurrently (default 4)
  -output string
        output format - json, junit, tap, text (default "text")
  -reject string
        comma-separated list of kinds to reject
  -schema-location value
        override schemas location search path (can be specified multiple times)
  -skip string
        comma-separated list of kinds to ignore
  -strict
        disallow additional properties not in schema
  -summary
        print a summary at the end (ignored for junit output)
  -v	show version information
  -verbose
        print results for all resources (ignored for tap and junit output)
{{< /prism >}}

### Validating a single, valid file

{{< prism >}}$ ./bin/kubeconform fixtures/valid.yaml
$ echo $?
0
{{< /prism >}}

### Validating a single invalid file, setting output to json, and printing a summary
{{< prism >}}$ ./bin/kubeconform -summary -output json fixtures/invalid.yaml
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
{{< /prism >}}

### Passing manifests via Stdin
{{< prism >}}cat fixtures/valid.yaml  | ./bin/kubeconform -summary
Summary: 1 resource found parsing stdin - Valid: 1, Invalid: 0, Errors: 0 Skipped: 0
{{< /prism >}}

### Validating a folder, increasing the number of parallel workers
{{< prism >}}$ ./bin/kubeconform -summary -n 16 fixtures
fixtures/crd_schema.yaml - CustomResourceDefinition trainingjobs.sagemaker.aws.amazon.com failed validation: could not find schema for CustomResourceDefinition
fixtures/invalid.yaml - ReplicationController bob is invalid: Invalid type. Expected: [integer,null], given: string
[...]
Summary: 65 resources found in 34 files - Valid: 55, Invalid: 2, Errors: 8 Skipped: 0
{{< /prism >}}