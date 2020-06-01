# Kubeconform

[![Build status](https://github.com/yannh/kubeconform/workflows/build/badge.svg?branch=master)](https://github.com/yannh/kubeconform/actions?query=branch%3Amaster)

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
  -file value
        file to validate (can be specified multiple times)
  -k8sversion string
        version of Kubernetes to test against (default "1.18.0")
  -output string
        output format - text, json (default "text")
  -printsummary
        print a summary at the end
  -quiet
        quiet output - only print invalid files, and errors
  -schema value
        file containing an additional Schema (can be specified multiple times)
  -skipKinds string
        comma-separated list of kinds to ignore
  -strict
        disallow additional properties not in schema
  -workers int
        number of routines to run in parallel (default 4)
```

### Credits

 * @garethr for the [Kubeval](https://github.com/instrumenta/kubeval) and
 [kubernetes-json-schema](https://github.com/instrumenta/kubernetes-json-schema) projects
