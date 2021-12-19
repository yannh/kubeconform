---
title: "Custom Resources support"
date: 2021-07-02T00:00:00Z
draft: false
tags: ["Kubeconform", "Usage"]
weight: 4
---

When the `-schema-location` parameter is not used, or set to "default", kubeconform will default to downloading
schemas from `https://github.com/yannh/kubernetes-json-schema`. Kubeconform however supports passing one, or multiple,
schemas locations - HTTP(s) URLs, or local filesystem paths, in which case it will lookup for schema definitions
in each of them, in order, stopping as soon as a matching file is found.

* If the -schema-location value does not end with '.json', Kubeconform will assume filenames / a file
  structure identical to that of kubernetesjsonschema.dev or github.com/yannh/kubernetes-json-schema.
* if the -schema-location value ends with '.json' - Kubeconform assumes the value is a Go templated
  string that indicates how to search for JSON schemas.
* the -schema-location value of "default" is an alias for https://raw.githubusercontent.com/yannh/kubernetes-json-schema/master/{{ .NormalizedKubernetesVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json.
  Both following command lines are equivalent:

{{< prism >}}$ ./bin/kubeconform fixtures/valid.yaml
$ ./bin/kubeconform -schema-location default fixtures/valid.yaml
$ ./bin/kubeconform -schema-location 'https://raw.githubusercontent.com/yannh/kubernetes-json-schema/master/{{ .NormalizedKubernetesVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}{{ .KindSuffix }}.json' fixtures/valid.yaml
{{< /prism >}}

To support validating CRDs, we need to convert OpenAPI files to JSON schema, storing the JSON schemas
in a local folder - for example schemas. Then we specify this folder as an additional registry to lookup:

{{< prism >}}# If the resource Kind is not found in kubernetesjsonschema.dev, also lookup in the schemas/ folder for a matching file
$ ./bin/kubeconform -schema-location default -schema-location 'schemas/{{ .ResourceKind }}{{ .KindSuffix }}.json' fixtures/custom-resource.yaml
{{< /prism >}}

You can validate Openshift manifests using a custom schema location. Set the OpenShift version to validate
against using -kubernetes-version.

{{< prism >}}$ ./bin/kubeconform -kubernetes-version 3.8.0  -schema-location 'https://raw.githubusercontent.com/garethr/openshift-json-schema/master/{{ .NormalizedKubernetesVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}.json'  -summary fixtures/valid.yaml
Summary: 1 resource found in 1 file - Valid: 1, Invalid: 0, Errors: 0 Skipped: 0
{{< /prism >}}

Here are the variables you can use in -schema-location:
* *NormalizedKubernetesVersion* - Kubernetes Version, prefixed by v
* *StrictSuffix* - "-strict" or "" depending on whether validation is running in strict mode or not
* *ResourceKind* - Kind of the Kubernetes Resource
* *ResourceAPIVersion* - Version of API used for the resource - "v1" in "apiVersion: monitoring.coreos.com/v1"
* *KindSuffix* - suffix computed from apiVersion - for compatibility with Kubeval schema registries
