<img width="50%" alt="Kubeconform-GitHub-Hero" src="https://user-images.githubusercontent.com/19731161/142411871-f695e40c-bfa8-43ca-97c0-94c256749732.png">
<hr>

[![Build status](https://github.com/yannh/kubeconform/workflows/build/badge.svg?branch=master)](https://github.com/yannh/kubeconform/actions?query=branch%3Amaster)
[![Homebrew](https://img.shields.io/badge/dynamic/json.svg?url=https://formulae.brew.sh/api/formula/kubeconform.json&query=$.versions.stable&label=homebrew)](https://formulae.brew.sh/formula/kubeconform)
[![Go Report card](https://goreportcard.com/badge/github.com/yannh/kubeconform)](https://goreportcard.com/report/github.com/yannh/kubeconform)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/yannh/kubeconform/pkg/validator)](https://pkg.go.dev/github.com/yannh/kubeconform/pkg/validator)

`Kubeconform` is a Kubernetes manifest validation tool. Incorporate it into your CI, or use it locally to validate your Kubernetes configuration!

It is inspired by, contains code from and is designed to stay close to
[Kubeval](https://github.com/instrumenta/kubeval), but with the following improvements:
 * **high performance**: will validate & download manifests over multiple routines, caching
   downloaded files in memory
 * configurable list of **remote, or local schemas locations**, enabling validating Kubernetes
   custom resources (CRDs) and offline validation capabilities
 * uses by default a [self-updating fork](https://github.com/yannh/kubernetes-json-schema) of the schemas registry maintained
   by the kubernetes-json-schema project - which guarantees
   up-to-date **schemas for all recent versions of Kubernetes**.
   
<details><summary><h4>Speed comparison with Kubeval</h4></summary><p>
Running on a pretty large kubeconfigs setup, on a laptop with 4 cores:
   
```bash
$ time kubeconform -ignore-missing-schemas -n 8 -summary  preview staging production
Summary: 50714 resources found in 35139 files - Valid: 27334, Invalid: 0, Errors: 0 Skipped: 23380
real	0m6,710s
user	0m38,701s
sys	0m1,161s
$ time kubeval -d preview,staging,production --ignore-missing-schemas --quiet
[... Skipping output]
real	0m35,336s
user	0m0,717s
sys	0m1,069s
```
</p></details>

## Table of contents

* [A small overview of Kubernetes manifest validation](#a-small-overview-of-kubernetes-manifest-validation)
  * [Limits of Kubeconform validation](#Limits-of-Kubeconform-validation)
* [Installation](#Installation)
* [Usage](#Usage)
  * [Usage examples](#Usage-examples)
  * [Proxy support](#Proxy-support)
* [Overriding schemas location](#Overriding-schemas-location)
  * [CustomResourceDefinition (CRD) Support](#CustomResourceDefinition-(CRD)-Support)
  * [OpenShift schema Support](#OpenShift-schema-Support)
* [Integrating Kubeconform in the CI](#Integrating-Kubeconform-in-the-CI)
  * [Github Workflow](#Github-Workflow)
  * [Gitlab-CI](#Gitlab-CI)
* [Helm charts](#helm-charts)
  * [Helm plugin](#helm-plugin)
  * [Helm `pre-commit` hook](#helm-pre-commit-hook)
* [Using kubeconform as a Go Module](#Using-kubeconform-as-a-Go-Module)
* [Credits](#Credits)

## A small overview of Kubernetes manifest validation

Kubernetes's API is described using the [OpenAPI (formerly swagger) specification](https://www.openapis.org),
in a [file](https://github.com/kubernetes/kubernetes/blob/master/api/openapi-spec/swagger.json) checked into
the main Kubernetes repository.

Because of the state of the tooling to perform validation against OpenAPI schemas, projects usually convert
the OpenAPI schemas to [JSON schemas](https://json-schema.org/) first. Kubeval relies on
[instrumenta/OpenApi2JsonSchema](https://github.com/instrumenta/openapi2jsonschema) to convert Kubernetes' Swagger file
and break it down into multiple JSON schemas, stored in github at
[instrumenta/kubernetes-json-schema](https://github.com/instrumenta/kubernetes-json-schema) and published on
[kubernetesjsonschema.dev](https://kubernetesjsonschema.dev/).

`Kubeconform` relies on [a fork of kubernetes-json-schema](https://github.com/yannh/kubernetes-json-schema/)
that is more meticulously kept up-to-date, and contains schemas for all recent versions of Kubernetes.

### Limits of Kubeconform validation

`Kubeconform`, similar to `kubeval`, only validates manifests using the official Kubernetes OpenAPI specifications. The Kubernetes controllers still perform additional server-side validations that are not part of the OpenAPI specifications. Those server-side validations are not covered by `Kubeconform` (examples: [#65](https://github.com/yannh/kubeconform/issues/65), [#122](https://github.com/yannh/kubeconform/issues/122), [#142](https://github.com/yannh/kubeconform/issues/142)). You can use a 3rd-party tool or the `kubectl --dry-run=server` command to fill the missing (validation) gap.

## Installation

If you are a [Homebrew](https://brew.sh/) user, you can install by running:

```bash
$ brew install kubeconform
```

You can also download the latest version from the [release page](https://github.com/yannh/kubeconform/releases).

Another way of installation is via Golang's package manager:

```bash
# With a specific version tag
$ go install github.com/yannh/kubeconform/cmd/kubeconform@v0.4.13

# Latest version
$ go install github.com/yannh/kubeconform/cmd/kubeconform@latest
```

## Usage

```
$ kubeconform -h
Usage: ./bin/kubeconform [OPTION]... [FILE OR FOLDER]...
  -cache string
        cache schemas downloaded via HTTP to this folder
  -debug
        print debug information
  -exit-on-error
        immediately stop execution when the first error is encountered
  -h    show help information
  -ignore-filename-pattern value
        regular expression specifying paths to ignore (can be specified multiple times)
  -ignore-missing-schemas
        skip files with missing schemas instead of failing
  -insecure-skip-tls-verify
        disable verification of the server\'s SSL certificate. This will make your HTTPS connections insecure
  -kubernetes-version string
        version of Kubernetes to validate against, e.g.: 1.18.0 (default "master")
  -n int
        number of goroutines to run concurrently (default 4)
  -output string
        output format - json, junit, tap, text (default "text")
  -reject string
        comma-separated list of kinds or GVKs to reject
  -schema-location value
        override schemas location search path (can be specified multiple times)
  -skip string
        comma-separated list of kinds or GVKs to ignore
  -strict
        disallow additional properties not in schema or duplicated keys
  -summary
        print a summary at the end (ignored for junit output)
  -v    show version information
  -verbose
        print results for all resources (ignored for tap and junit output)
```

### Usage examples

* Validating a single, valid file
```bash
$ kubeconform fixtures/valid.yaml
$ echo $?
0
```

* Validating a single invalid file, setting output to json, and printing a summary
```bash
$ kubeconform -summary -output json fixtures/invalid.yaml
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

* Passing manifests via Stdin
```bash
cat fixtures/valid.yaml  | ./bin/kubeconform -summary
Summary: 1 resource found parsing stdin - Valid: 1, Invalid: 0, Errors: 0 Skipped: 0
```

* Validating a file, ignoring its resource using both Kind, and GVK (Group, Version, Kind) notations
```
# This will ignore ReplicationController for all apiVersions
$ kubeconform -summary -skip ReplicationController fixtures/valid.yaml
Summary: 1 resource found in 1 file - Valid: 0, Invalid: 0, Errors: 0, Skipped: 1

# This will ignore ReplicationController only for apiVersion v1
$ kubeconform -summary -skip v1/ReplicationController fixtures/valid.yaml
Summary: 1 resource found in 1 file - Valid: 0, Invalid: 0, Errors: 0, Skipped: 1
```

* Validating a folder, increasing the number of parallel workers
```
$ kubeconform -summary -n 16 fixtures
fixtures/crd_schema.yaml - CustomResourceDefinition trainingjobs.sagemaker.aws.amazon.com failed validation: could not find schema for CustomResourceDefinition
fixtures/invalid.yaml - ReplicationController bob is invalid: Invalid type. Expected: [integer,null], given: string
[...]
Summary: 65 resources found in 34 files - Valid: 55, Invalid: 2, Errors: 8 Skipped: 0
```

### Proxy support

`Kubeconform` will respect the **HTTPS_PROXY** variable when downloading schema files.

```bash
$ HTTPS_PROXY=proxy.local bin/kubeconform fixtures/valid.yaml
```

## Overriding schemas location

When the `-schema-location` parameter is not used, or set to `default`, kubeconform will default to downloading
schemas from https://github.com/yannh/kubernetes-json-schema. Kubeconform however supports passing one, or multiple,
schemas locations - HTTP(s) URLs, or local filesystem paths, in which case it will lookup for schema definitions
in each of them, in order, stopping as soon as a matching file is found.

 * If the `-schema-location` value does not end with `.json`, Kubeconform will assume filenames / a file
 structure identical to that of [kubernetesjsonschema.dev](https://kubernetesjsonschema.dev/) or [yannh/kubernetes-json-schema](https://github.com/yannh/kubernetes-json-schema).
 * if the `-schema-location` value ends with `.json` - Kubeconform assumes the value is a **Go templated
 string** that indicates how to search for JSON schemas.
* the `-schema-location` value of `default` is an alias for `https://raw.githubusercontent.com/yannh/kubernetes-json-schema/master/{{.NormalizedKubernetesVersion}}-standalone{{.StrictSuffix}}/{{.ResourceKind}}{{.KindSuffix}}.json`.

**The following command lines are equivalent:**
```bash
$ kubeconform fixtures/valid.yaml
$ kubeconform -schema-location default fixtures/valid.yaml
$ kubeconform -schema-location 'https://raw.githubusercontent.com/yannh/kubernetes-json-schema/master/{{.NormalizedKubernetesVersion}}-standalone{{.StrictSuffix}}/{{.ResourceKind}}{{.KindSuffix}}.json' fixtures/valid.yaml
```
Here are the variables you can use in -schema-location:
 * *NormalizedKubernetesVersion* - Kubernetes Version, prefixed by v
 * *StrictSuffix* - "-strict" or "" depending on whether validation is running in strict mode or not
 * *ResourceKind* - Kind of the Kubernetes Resource
 * *ResourceAPIVersion* - Version of API used for the resource - "v1" in "apiVersion: monitoring.coreos.com/v1"
 * *Group* - the group name as stated in this resource's definition - "monitoring.coreos.com" in "apiVersion: monitoring.coreos.com/v1"
 * *KindSuffix* - suffix computed from apiVersion - for compatibility with `Kubeval` schema registries

### CustomResourceDefinition (CRD) Support

Because Custom Resources (CR) are not native Kubernetes objects, they are not included in the default schema.  
If your CRs are present in [Datree's CRDs-catalog](https://github.com/datreeio/CRDs-catalog), you can specify this project as an additional registry to lookup:
  
```bash
# Look in the CRDs-catalog for the desired schema/s
$ kubeconform -schema-location default -schema-location 'https://raw.githubusercontent.com/datreeio/CRDs-catalog/main/{{.Group}}/{{.ResourceKind}}_{{.ResourceAPIVersion}}.json' [MANIFEST]
```

If your CRs are not present in the CRDs-catalog, you will need to manually pull the CRDs manifests from your cluster and convert the `OpenAPI.spec` to JSON schema format.

<details><summary>Converting an OpenAPI file to a JSON Schema</summary>
<p>

`Kubeconform` uses JSON schemas to validate Kubernetes resources. For Custom Resource, the CustomResourceDefinition
first needs to be converted to JSON Schema. A script is provided to convert these CustomResourceDefinitions
to JSON schema. Here is an example how to use it:

```bash
$ python ./scripts/openapi2jsonschema.py https://raw.githubusercontent.com/aws/amazon-sagemaker-operator-for-k8s/master/config/crd/bases/sagemaker.aws.amazon.com_trainingjobs.yaml
JSON schema written to trainingjob_v1.json
```

By default, the file name output format is `{kind}_{version}`. The `FILENAME_FORMAT` environment variable can be used to change the output file name (Available variables: `kind`, `group`, `version`):

```
$ export FILENAME_FORMAT='{kind}-{group}-{version}'
$ ./scripts/openapi2jsonschema.py https://raw.githubusercontent.com/aws/amazon-sagemaker-operator-for-k8s/master/config/crd/bases/sagemaker.aws.amazon.com_trainingjobs.yaml
JSON schema written to trainingjob-sagemaker-v1.json
```

After converting your CRDs to JSON schema files, you can use `kubeconform` to validate your CRs against them:

```
# If the resource Kind is not found in default, also lookup in the schemas/ folder for a matching file
$ kubeconform -schema-location default -schema-location 'schemas/{{ .ResourceKind }}{{ .KindSuffix }}.json' fixtures/custom-resource.yaml
```

ℹ️ Datree's [CRD Extractor](https://github.com/datreeio/CRDs-catalog#crd-extractor) is a utility that can be used instead of this manual process.

</p>
</details>

### OpenShift schema Support

You can validate Openshift manifests using a custom schema location. Set the OpenShift version (v3.10.0-4.1.0) to validate
against using `-kubernetes-version`.

```
kubeconform -kubernetes-version 3.8.0  -schema-location 'https://raw.githubusercontent.com/garethr/openshift-json-schema/master/{{ .NormalizedKubernetesVersion }}-standalone{{ .StrictSuffix }}/{{ .ResourceKind }}.json'  -summary fixtures/valid.yaml
Summary: 1 resource found in 1 file - Valid: 1, Invalid: 0, Errors: 0 Skipped: 0
```

## Integrating Kubeconform in the CI

`Kubeconform` publishes Docker Images to Github's new Container Registry (ghcr.io). These images
can be used directly in a Github Action, once logged in using a [_Github Token_](https://github.blog/changelog/2021-03-24-packages-container-registry-now-supports-github_token/).

### Github Workflow

Example:
```yaml
name: kubeconform
on: push
jobs:
  kubeconform:
    runs-on: ubuntu-latest
    steps:
      - name: login to Github Packages
        run: echo "${{ github.token }}" | docker login https://ghcr.io -u ${GITHUB_ACTOR} --password-stdin
      - uses: actions/checkout@v2
      - uses: docker://ghcr.io/yannh/kubeconform:master
        with:
          entrypoint: '/kubeconform'
          args: "-summary -output json kubeconfigs/"
```

_Note on pricing_: Kubeconform relies on Github Container Registry which is currently in Beta. During that period,
[bandwidth is free](https://docs.github.com/en/packages/guides/about-github-container-registry). After that period,
bandwidth costs might be applicable. Since bandwidth from Github Packages within Github Actions is free, I expect
Github Container Registry to also be usable for free within Github Actions in the future. If that were not to be the
case, I might publish the Docker image to a different platform.

### Gitlab-CI

The Kubeconform Docker image can be used in Gitlab-CI. Here is an example of a Gitlab-CI job:

```yaml
lint-kubeconform:
  stage: validate
  image:
    name: ghcr.io/yannh/kubeconform:latest-alpine
    entrypoint: [""]
  script:
  - kubeconform
```

See [issue 106](https://github.com/yannh/kubeconform/issues/106) for more details.

## Helm charts

`kubeconform` supports automation for [Helm charts](https://helm.sh) in the form
of [Helm plugin](https://helm.sh/docs/topics/plugins/) and [`pre-commit`
hook](https://pre-commit.com/).

### Helm plugin

The `kubeconform` [Helm plugin](https://helm.sh/docs/topics/plugins/) can be
installed using this command:

```shell
helm plugin install https://github.com/yannh/kubeconform
```

Once installed, the plugin can be used from any Helm chart directory:

```shell
# Enter the chart directory
cd charts/mychart
# Run kubeconform plugin
helm kubeconform .
```

The plugin uses `helm template` internally and passes its output to the
`kubeconform`. There is several `helm template` command line options supported
by the plugin that can be specified:

```shell
helm kubeconform --namespace myns .
```

There is also several `kubeconform` command line options supported by the plugin
that can be specified:

```shell
# Kubeconform options
helm kubeconform --verbose --summary .
```

It's also possible to create `.kubeconform` file in the Helm chart directory
that can contain default `kubeconform` settings:

```yaml
# Command line options that can be set multiple times can be defined as an array
schema-location:
  - default
  - https://raw.githubusercontent.com/datreeio/CRDs-catalog/main/{{.Group}}/{{.ResourceKind}}_{{.ResourceAPIVersion}}.json
# Command line options that can be specified without a value must have boolean
# value in the config file
summary: true
verbose: true
```

The full list of options for the
[plugin](https://github.com/yannh/kubeconform/blob/master/scripts/helm/plugin_wrapper.py)
is as follows:

```text
$ ./plugin_wrapper.py --help
usage: plugin_wrapper.py [-h] [--cache] [--cache-dir DIR] [--config FILE] [--values-dir DIR]
                         [--values-pattern PATTERN] [--debug] [--skip-refresh] [--verify]
                         [-f FILE] [-n NAME] [-r NAME] [--ignore-missing-schemas]
                         [--insecure-skip-tls-verify] [--kubernetes-version VERSION]
                         [--goroutines NUMBER] [--output {json,junit,tap,text}]
                         [--reject LIST] [--schema-location LOCATION] [--skip LIST]
                         [--strict] [--summary] [--verbose]
                         CHART

Wrapper to run kubeconform for a Helm chart.

options:
  -h, --help            show this help message and exit
  --cache               whether to use kubeconform cache
  --cache-dir DIR       path to the cache directory (default: ~/.cache/kubeconform)
  --config FILE         config file name (default: .kubeconform)
  --values-dir DIR      directory with optional values files for the tests (default: ci)
  --values-pattern PATTERN
                        pattern to select the values files (default: *-values.yaml)
  --debug               debug output

helm build:
  Options passed to the 'helm build' command

  --skip-refresh        do not refresh the local repository cache
  --verify              verify the packages against signatures

helm template:
  Options passed to the 'helm template' command

  -f FILE, --values FILE
                        values YAML file or URL (can specified multiple)
  -n NAME, --namespace NAME
                        namespace
  -r NAME, --release NAME
                        release name
  CHART                 chart path (e.g. '.')

kubeconform:
  Options passsed to the 'kubeconform' command

  --ignore-missing-schemas
                        skip files with missing schemas instead of failing
  --insecure-skip-tls-verify
                        disable verification of the server's SSL certificate
  --kubernetes-version VERSION
                        version of Kubernetes to validate against, e.g. 1.18.0 (default:
                        master)
  --goroutines NUMBER   number of goroutines to run concurrently (default: 4)
  --output {json,junit,tap,text}
                        output format (default: text)
  --reject LIST         comma-separated list of kinds or GVKs to reject
  --schema-location LOCATION
                        override schemas location search path (can specified multiple)
  --skip LIST           comma-separated list of kinds or GVKs to ignore
  --strict              disallow additional properties not in schema or duplicated keys
  --summary             print a summary at the end (ignored for junit output)
  --verbose             print results for all resources (ignored for tap and junit output)
```

## Helm `pre-commit` hook

The `kubeconform` [`pre-commit` hook](https://pre-commit.com) can be added into the
`.pre-commit-config.yaml` file like this:

```yaml
repos:
  - repo: https://github.com/yannh/kubeconform
    rev: v0.5.0
    hooks:
      - id: kubeconform-helm
```

The hook uses `helm template` internally and passes its output to the
`kubeconform`. There is several `helm template` command line options supported
by the hook that can be specified:

```yaml
  - repo: https://github.com/yannh/kubeconform
    rev: v0.5.0
    hooks:
      - id: kubeconform-helm
        args:
          - --namespace=myns
          - --release=myrelease
```

There is also several `kubeconform` command line options supported by the hook
that can be specified:

```yaml
  - repo: https://github.com/yannh/kubeconform
    rev: v0.5.0
    hooks:
      - id: kubeconform-helm
        args:
          - --kubernetes-version=1.24.0
          - --verbose
          - --summary
```

The full list of options for the
[hook](https://github.com/yannh/kubeconform/blob/master/scripts/helm/pre-commit.py)
is as follows:

```text
$ ./pre-commit.py --help
usage: pre-commit.py [-h] [--charts-path PATH] [--include-charts LIST]
                     [--exclude-charts LIST] [--cache] [--cache-dir DIR] [--config FILE]
                     [--values-dir DIR] [--values-pattern PATTERN] [--debug]
                     [--skip-refresh] [--verify] [-f FILE] [-n NAME] [-r NAME]
                     [--ignore-missing-schemas] [--insecure-skip-tls-verify]
                     [--kubernetes-version VERSION] [--goroutines NUMBER]
                     [--output {json,junit,tap,text}] [--reject LIST]
                     [--schema-location LOCATION] [--skip LIST] [--strict] [--summary]
                     [--verbose]
                     FILES [FILES ...]

Wrapper to run kubeconform for a Helm chart.

positional arguments:
  FILES                 files that have changed

options:
  -h, --help            show this help message and exit
  --charts-path PATH    path to the directory with charts (default: charts)
  --include-charts LIST
                        comma-separated list of chart names to include in the testing
  --exclude-charts LIST
                        comma-separated list of chart names to exclude from the testing
  --cache               whether to use kubeconform cache
  --cache-dir DIR       path to the cache directory (default: ~/.cache/kubeconform)
  --config FILE         config file name (default: .kubeconform)
  --values-dir DIR      directory with optional values files for the tests (default: ci)
  --values-pattern PATTERN
                        pattern to select the values files (default: *-values.yaml)
  --debug               debug output

helm build:
  Options passed to the 'helm build' command

  --skip-refresh        do not refresh the local repository cache
  --verify              verify the packages against signatures

helm template:
  Options passed to the 'helm template' command

  -f FILE, --values FILE
                        values YAML file or URL (can specified multiple)
  -n NAME, --namespace NAME
                        namespace
  -r NAME, --release NAME
                        release name

kubeconform:
  Options passsed to the 'kubeconform' command

  --ignore-missing-schemas
                        skip files with missing schemas instead of failing
  --insecure-skip-tls-verify
                        disable verification of the server's SSL certificate
  --kubernetes-version VERSION
                        version of Kubernetes to validate against, e.g. 1.18.0 (default:
                        master)
  --goroutines NUMBER   number of goroutines to run concurrently (default: 4)
  --output {json,junit,tap,text}
                        output format (default: text)
  --reject LIST         comma-separated list of kinds or GVKs to reject
  --schema-location LOCATION
                        override schemas location search path (can specified multiple)
  --skip LIST           comma-separated list of kinds or GVKs to ignore
  --strict              disallow additional properties not in schema or duplicated keys
  --summary             print a summary at the end (ignored for junit output)
  --verbose             print results for all resources (ignored for tap and junit output)
```

## Using kubeconform as a Go Module

**Warning**: This is a work-in-progress, the interface is not yet considered stable. Feedback is encouraged.

`Kubeconform` contains a package that can be used as a library.
An example of usage can be found in [examples/main.go](examples/main.go)

Additional documentation on [pkg.go.dev](https://pkg.go.dev/github.com/yannh/kubeconform/pkg/validator)

## Credits

 * @garethr for the [Kubeval](https://github.com/instrumenta/kubeval) and
 [kubernetes-json-schema](https://github.com/instrumenta/kubernetes-json-schema) projects ❤️
