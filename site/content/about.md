---
title: "About"
date: 2021-07-02T00:00:00Z
draft: false
tags: ["Kubeconform", "About"]
---

Kubeconform is a Kubernetes manifests validation tool. Build it into your CI to validate your Kubernetes
configuration!

It is inspired by, contains code from and is designed to stay close to
[Kubeval](https://github.com/instrumenta/kubeval), but with the following improvements:
* **high performance**: will validate & download manifests over multiple routines, caching
  downloaded files in memory
* configurable list of **remote, or local schemas locations**, enabling validating Kubernetes
  custom resources (CRDs) and offline validation capabilities
* uses by default a [self-updating fork](https://github.com/yannh/kubernetes-json-schema) of the schemas registry maintained
  by the [kubernetes-json-schema](https://github.com/instrumenta/kubernetes-json-schema) project - which guarantees
  up-to-date **schemas for all recent versions of Kubernetes**.