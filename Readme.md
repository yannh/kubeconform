# Kubeconform

A Kubernetes manifests validation tool, inspired by & similar to [Kubeval](https://github.com/instrumenta/kubeval)

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
  -skipKinds string
        comma-separated list of kinds to ignore
```
