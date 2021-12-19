---
title: "Overview"
date: 2021-07-02T00:00:00Z
draft: false
tags: ["Kubeconform", "Overview"]
weight: 1
---

Kubeconform is a Kubernetes manifests validation tool, and checks whether your Kubernetes manifests
are valid, according to Kubernetes resources definitions.

It is inspired by, contains code from and is designed to stay close to
[Kubeval](https://github.com/instrumenta/kubeval), but with the following improvements:
* **high performance**: will validate & download manifests over multiple routines, caching
  downloaded files in memory
* configurable list of **remote, or local schemas locations**, enabling validating Kubernetes
  custom resources (CRDs) and offline validation capabilities
* uses by default a [self-updating fork](https://github.com/yannh/kubernetes-json-schema) of the schemas registry maintained
  by the [kubernetes-json-schema](https://github.com/instrumenta/kubernetes-json-schema) project - which guarantees
  up-to-date **schemas for all recent versions of Kubernetes**.
* improved logging: support for more formats (Tap, Junit, JSON).

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

Kubeconform relies on [a fork of kubernetes-json-schema](https://github.com/yannh/kubernetes-json-schema/)
that is more aggressively kept up-to-date, and contains schemas for all recent versions of Kubernetes.

### Limits of Kubeconform validation

Kubeconform, similarly to kubeval, only validates manifests using the OpenAPI specifications. In some
cases, the Kubernetes controllers might perform additional validation - so that manifests passing kubeval
validation would still error when being deployed. See for example these bugs against kubeval:
[#253](https://github.com/instrumenta/kubeval/issues/253)
[#256](https://github.com/instrumenta/kubeval/issues/256)
[#257](https://github.com/instrumenta/kubeval/issues/257)
[#259](https://github.com/instrumenta/kubeval/issues/259). The validation logic mentioned in these
bug reports is not part of Kubernetes' OpenAPI spec, and therefore kubeconform/kubeval will not detect the
configuration errors.