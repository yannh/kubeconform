---
title: "GitHub Action"
date: 2021-07-02T00:00:00Z
draft: false
tags: ["Kubeconform", "Usage"]
weight: 6
---

Kubeconform is publishes Docker Images to GitHub's new Container Registry, ghcr.io. These images
can be used directly in a GitHub Action, once logged in using a [_GitHub Token_](https://github.blog/changelog/2021-03-24-packages-container-registry-now-supports-github_token/).

{{< prism >}}name: kubeconform
on: push
jobs:
  kubeconform:
    runs-on: ubuntu-latest
    steps:
      - name: login to GitHub Packages
        run: echo "${{ github.token }}" | docker login https://ghcr.io -u ${GITHUB_ACTOR} --password-stdin
      - uses: actions/checkout@v2
      - uses: docker://ghcr.io/yannh/kubeconform:master
        with:
          entrypoint: '/kubeconform'
          args: "-summary -output json kubeconfigs/"
{{< /prism >}}

_Note on pricing_: Kubeconform relies on GitHub Container Registry which is currently in Beta. During that period,
[bandwidth is free](https://docs.github.com/en/packages/guides/about-github-container-registry). After that period,
bandwidth costs might be applicable. Since bandwidth from GitHub Packages within GitHub Actions is free, I expect
GitHub Container Registry to also be usable for free within GitHub Actions in the future. If that were not to be the
case, I might publish the Docker image to a different platform.