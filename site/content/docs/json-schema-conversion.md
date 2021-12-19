---
title: "OpenAPI to JSON Schema conversion"
date: 2021-07-02T00:00:00Z
draft: false
tags: ["Kubeconform", "Usage"]
weight: 5
---

Kubeconform uses JSON schemas to validate Kubernetes resources. For custom resources, the CustomResourceDefinition
first needs to be converted to JSON Schema. A script is provided to convert these CustomResourceDefinitions
to JSON schema. Here is an example how to use it:

{{< prism >}}#!/bin/bash
$ ./scripts/openapi2jsonschema.py https://raw.githubusercontent.com/aws/amazon-sagemaker-operator-for-k8s/master/config/crd/bases/sagemaker.aws.amazon.com_trainingjobs.yaml
JSON schema written to trainingjob_v1.json
{{< /prism >}}

The `FILENAME_FORMAT` environment variable can be used to change the output file name (Available variables: `kind`, `group`, `version`) (Default: `{kind}_{version}`).

{{< prism >}}$ export FILENAME_FORMAT='{kind}-{group}-{version}'
$ ./scripts/openapi2jsonschema.py https://raw.githubusercontent.com/aws/amazon-sagemaker-operator-for-k8s/master/config/crd/bases/sagemaker.aws.amazon.com_trainingjobs.yaml
JSON schema written to trainingjob-sagemaker-v1.json
{{< /prism >}}

Some CRD schemas do not have explicit validation for fields implicitly validated by the Kubernetes API like `apiVersion`, `kind`, and `metadata`, thus additional properties are allowed at the root of the JSON schema by default, if this is not desired the `DENY_ROOT_ADDITIONAL_PROPERTIES` environment variable can be set to any non-empty value.