#!/usr/bin/env python3

import argparse
import logging
import os
import subprocess
import sys
import yaml

from glob import glob


def parse_args(add_chart=True, add_files=False, add_path=False, add_incl_excl=False):
    # Command line options for helm, kubeconform and the plugin itself
    args = {
        "helm_tmpl": [],
        "helm_build": [],
        "kubeconform": [],
        "wrapper": [],
    }

    # Define parser
    parser = argparse.ArgumentParser(
        description="Wrapper to run kubeconform for a Helm chart."
    )

    if add_path:
        parser.add_argument(
            "--charts-path",
            metavar="PATH",
            help="path to the directory with charts (default: charts)",
            default="charts",
        )

    if add_incl_excl:
        parser.add_argument(
            "--include-charts",
            metavar="LIST",
            help="comma-separated list of chart names to include in the testing",
        )
        parser.add_argument(
            "--exclude-charts",
            metavar="LIST",
            help="comma-separated list of chart names to exclude from the testing",
        )

    parser.add_argument(
        "--cache", help="whether to use kubeconform cache", action="store_true"
    )
    parser.add_argument(
        "--cache-dir",
        metavar="DIR",
        help="path to the cache directory (default: ~/.cache/kubeconform)",
        default="~/.cache/kubeconform",
    )
    parser.add_argument(
        "--config",
        metavar="FILE",
        help="config file name (default: .kubeconform)",
        default=".kubeconform",
    )
    parser.add_argument(
        "--values-dir",
        metavar="DIR",
        help="directory with optional values files for the tests (default: ci)",
        default="ci",
    )
    parser.add_argument(
        "--values-pattern",
        metavar="PATTERN",
        help="pattern to select the values files (default: *-values.yaml)",
        default="*-values.yaml",
    )

    parser.add_argument(
        "--debug",
        help="debug output",
        action="store_true",
    )

    group_helm_build = parser.add_argument_group(
        "helm build", "Options passed to the 'helm build' command"
    )

    group_helm_build.add_argument(
        "--skip-refresh",
        help="do not refresh the local repository cache",
        action="store_true",
    )
    group_helm_build.add_argument(
        "--verify", help="verify the packages against signatures", action="store_true"
    )

    group_helm_tmpl = parser.add_argument_group(
        "helm template", "Options passed to the 'helm template' command"
    )

    group_helm_tmpl.add_argument(
        "-f",
        "--values",
        metavar="FILE",
        help="values YAML file or URL (can specified multiple)",
        action="append",
    )
    group_helm_tmpl.add_argument(
        "-n",
        "--namespace",
        metavar="NAME",
        help="namespace",
    )
    group_helm_tmpl.add_argument(
        "-r",
        "--release",
        metavar="NAME",
        help="release name",
    )

    if add_chart:
        group_helm_tmpl.add_argument(
            "CHART",
            help="chart path (e.g. '.')",
        )

    group_kubeconform = parser.add_argument_group(
        "kubeconform", "Options passsed to the 'kubeconform' command"
    )

    group_kubeconform.add_argument(
        "--ignore-missing-schemas",
        help="skip files with missing schemas instead of failing",
        action="store_true",
    )
    group_kubeconform.add_argument(
        "--insecure-skip-tls-verify",
        help="disable verification of the server's SSL certificate",
        action="store_true",
    )
    group_kubeconform.add_argument(
        "--kubernetes-version",
        metavar="VERSION",
        help="version of Kubernetes to validate against, e.g. 1.18.0 (default: master)",
    )
    group_kubeconform.add_argument(
        "--goroutines",
        metavar="NUMBER",
        help="number of goroutines to run concurrently (default: 4)",
    )
    group_kubeconform.add_argument(
        "--output",
        help="output format (default: text)",
        choices=["json", "junit", "tap", "text"],
    )
    group_kubeconform.add_argument(
        "--reject",
        metavar="LIST",
        help="comma-separated list of kinds or GVKs to reject",
    )
    group_kubeconform.add_argument(
        "--schema-location",
        metavar="LOCATION",
        help="override schemas location search path (can specified multiple)",
        action="append",
    )
    group_kubeconform.add_argument(
        "--skip",
        metavar="LIST",
        help="comma-separated list of kinds or GVKs to ignore",
    )
    group_kubeconform.add_argument(
        "--strict",
        help="disallow additional properties not in schema or duplicated keys",
        action="store_true",
    )
    group_kubeconform.add_argument(
        "--summary",
        help="print a summary at the end (ignored for junit output)",
        action="store_true",
    )
    group_kubeconform.add_argument(
        "--verbose",
        help="print results for all resources (ignored for tap and junit output)",
        action="store_true",
    )

    if add_files:
        parser.add_argument(
            "FILES",
            help="files that have changed",
            nargs="+",
        )

    # Parse the args
    a = parser.parse_args()

    # ### Populate the helm build options
    if a.skip_refresh:
        args["helm_build"] = ["--skip-refresh"]

    if a.verify:
        args["helm_build"] = ["--verify"]

    # This must stay the last item from 'helm_build'!
    if add_chart:
        args["helm_build"] += [a.CHART]

    # ### Populate the helm template options
    if a.values:
        for v in a.values:
            args["helm_tmpl"] += ["--values", v]

    if a.namespace is not None:
        args["helm_tmpl"] += ["--namespace", a.namespace]

    if a.release is not None:
        args["helm_tmpl"] += [a.release]

    # This must stay the last item from 'helm_tmpl'!
    if add_chart:
        args["helm_tmpl"] += [a.CHART]

    # ### Polulate the kubeconform options
    if a.cache:
        args["kubeconform"] += ["-cache", os.path.expanduser(a.cache_dir)]

    if a.ignore_missing_schemas is True:
        args["kubeconform"] += ["-ignore-missing-schemas"]

    if a.insecure_skip_tls_verify is True:
        args["kubeconform"] += ["-insecure-skip-tls-verify"]

    if a.kubernetes_version is not None:
        args["kubeconform"] += ["-kubernetes-version", a.kubernetes_version]

    if a.goroutines is not None:
        args["kubeconform"] += ["-n", a.goroutines]

    if a.output is not None:
        args["kubeconform"] += ["-output", a.output]

    if a.reject is True:
        args["kubeconform"] += ["-reject"]

    if a.schema_location:
        for v in a.schema_location:
            args["kubeconform"] += ["-schema-location", v]

    if a.skip is True:
        args["kubeconform"] += ["-skip"]

    if a.strict is True:
        args["kubeconform"] += ["-strict"]

    if a.summary is True:
        args["kubeconform"] += ["-summary"]

    if a.verbose is True:
        args["kubeconform"] += ["-verbose"]

    # ### All args are wrapper options
    args["wrapper"] = a

    return args


def get_logger(debug):
    if debug:
        level = logging.DEBUG
    else:
        level = logging.ERROR

    format = "%(levelname)s: %(message)s"

    logging.basicConfig(level=level, format=format)

    return logging.getLogger(__name__)


def parse_config(filename):
    args = []

    # Check if file exists
    if not os.path.isfile(filename):
        return args

    # Read and parse the file
    try:
        with open(filename, "r") as stream:
            try:
                data = yaml.load(stream, Loader=yaml.Loader)
            except yaml.YAMLError as e:
                raise Exception("cannot parse YAML file '%s': %s" % (filename, e))
    except IOError as e:
        raise Exception("cannot open file '%s': %s" % (filename, e))

    # Produce extra args out of the config file
    if isinstance(data, dict):
        for key, val in data.items():
            if isinstance(val, list):
                for v in val:
                    args.append("-%s=%s" % (key, v))
            elif isinstance(v, dict):
                # No deep dicts allowed in the config
                continue
            else:
                args.append("-%s=%s" % (key, val))

    return args


def get_values_files(values_dir, values_pattern, chart_dir=None):
    values_files = []

    if os.path.isdir(values_dir):
        # Get values files from a specific path
        values_files = glob(os.path.join(values_dir, values_pattern))
    elif chart_dir is not None and os.path.isdir(os.path.join(chart_dir, values_dir)):
        # Get values files from the chart directory
        values_files = glob(os.path.join(chart_dir, values_dir, values_pattern))

    return values_files


def run_helm_dependecy_build(args):
    # Check if it's local chart
    if not os.path.isfile(os.path.join(args[-1], "Chart.yaml")):
        return

    charts_dir = os.path.join(args[-1], "charts")

    # Check if the dependency charts are already there
    if os.path.isdir(charts_dir):
        with open(os.path.join(args[-1], "Chart.yaml"), "r") as f:
            try:
                data = yaml.safe_load(f)
            except yaml.YAMLError as e:
                raise Exception("failed to parse Chart.yaml: %s" % e)

            if "dependencies" in data:
                for d in data["dependencies"]:
                    if not (
                        "name" in d
                        and (
                            (
                                "version" in d
                                and os.path.isfile(
                                    os.path.join(
                                        charts_dir,
                                        "%s-%s.tgz" % (d["name"], d["version"]),
                                    )
                                )
                            )
                            or os.path.islink(os.path.join(charts_dir, d["name"]))
                        )
                    ):
                        # Dependency missing, let's get it
                        break
                else:
                    # All dependencies seem to be there so don't run anything
                    return

    # Check if there is Chart.lock
    if os.path.isfile(os.path.join(args[-1], "Chart.yaml")):
        action = "update"
    else:
        action = "build"

    # Run process
    result = subprocess.run(
        [
            os.getenv("HELM_BIN", "helm"),
            "dependency",
            action,
        ]
        + args,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
    )

    # Check for errors
    if result.returncode != 0:
        raise Exception(
            "failed to run helm dependency build: rc=%d %s"
            % (result.returncode, result.stderr)
        )


def run_helm_template(args):
    # Run process
    result = subprocess.run(
        [
            os.getenv("HELM_BIN", "helm"),
            "template",
        ]
        + args,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
    )

    # Check for errors
    if result.returncode != 0:
        raise Exception(
            "failed to run helm template: rc=%d %s" % (result.returncode, result.stderr)
        )

    return result


def run_kubeconform(args, input):
    bin_file = "kubeconform"

    # Try to use `HELM_PLUGIN_DIR` env var to determine location of kubeconform
    helm_plugin_dir = os.getenv("HELM_PLUGIN_DIR", "")
    helm_plugin_bin = os.path.join(helm_plugin_dir, "bin", bin_file)

    if os.path.isfile(helm_plugin_bin):
        bin_file = helm_plugin_bin
    else:
        helm_error = False

        # Try to use `helm env` to determine location of kubeconform
        try:
            result = subprocess.run(
                [
                    "helm",
                    "env",
                ],
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True,
            )
        except Exception:
            helm_error = True

        if not helm_error:
            for line in result.stdout.split("\n"):
                if line.startswith("HELM_PLUGINS") and "=" in line:
                    _, plugins_path = line.split("=")

                    helm_plugin_bin = os.path.join(
                        plugins_path.strip('"'), bin_file, "bin", bin_file
                    )

                    if os.path.isfile(helm_plugin_bin):
                        bin_file = helm_plugin_bin

    # Create the cache dir
    if "-cache" in args:
        c_idx = args.index("-cache")
        c_dir = args[c_idx + 1]

        if not os.path.isdir(c_dir):
            try:
                os.mkdir(c_dir, 0o755)
            except OSError as e:
                raise Exception("failed to create cache directory: %s" % e)

    # Run process
    result = subprocess.run(
        [
            bin_file,
        ]
        + args,
        input=input,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
    )

    # Check for errors
    if result.returncode != 0:
        raise Exception(
            "failed to run kubeconform: rc=%d\n%s%s"
            % (result.returncode, result.stderr, result.stdout)
        )

    return result


def run_test(args, values_file=None):
    values_args = []

    # Add extra values parameter if any file is specified
    if values_file:
        values_args = [
            "--values",
            values_file,
        ]

    # Build Helm dependencies
    try:
        run_helm_dependecy_build(
            args["helm_build"],
        )
    except Exception as e:
        raise Exception("dependency build failed: %s" % e)

    # Get templated output
    try:
        helm_result = run_helm_template(
            args["helm_tmpl"] + values_args,
        )
    except Exception as e:
        raise Exception("templating failed: %s" % e)

    # Get kubeconform output
    try:
        kubeconform_result = run_kubeconform(
            args["kubeconform"],
            helm_result.stdout,
        )
    except Exception as e:
        raise Exception("kubeconform failed: %s" % e)

    # Print results
    if kubeconform_result.stdout:
        print(kubeconform_result.stdout.rstrip())


def main():
    # Parse args
    args = parse_args()

    # Ger logger
    log = get_logger(args["wrapper"].debug)

    # Parse config file
    config_args = parse_config(
        args["wrapper"].config,
    )

    # Merge the args from config file and from command line
    if config_args:
        args["kubeconform"] = config_args + args["kubeconform"]

    # Get list of values files
    values_files = get_values_files(
        args["wrapper"].values_dir,
        args["wrapper"].values_pattern,
        args["helm_tmpl"][-1],
    )

    # Run tests
    try:
        if values_files:
            for values_file in values_files:
                log.debug("Testing with CI values file %s" % values_file)

                run_test(args, values_file)
        else:
            log.debug("Testing without CI values files")

            run_test(args)
    except Exception as e:
        log.error("Testing failed: %s" % e)

        sys.exit(1)


if __name__ == "__main__":
    main()
