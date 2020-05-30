# Kubeconform

A Kubernetes manifests validation tool, inspired by & similar to [Kubeval](https://github.com/instrumenta/kubeval)

Notable features:
 * high performance: will validate & download manifests over multiple routines
 * support for Kubernetes CRDs

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
  -schema value
        file containing an additional Schema (can be specified multiple times)
  -skipKinds string
        comma-separated list of kinds to ignore
  -workers int
        number of routines to run in parallel (default 4)
```
