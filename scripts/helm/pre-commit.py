#!/usr/bin/env python3

import os
import sys

from contextlib import contextmanager

import plugin_wrapper as pw


@contextmanager
def cd(newdir):
    prevdir = os.getcwd()

    os.chdir(os.path.expanduser(newdir))

    try:
        yield
    finally:
        os.chdir(prevdir)


def main():
    # Parse args
    args = pw.parse_args(
        add_chart=False,
        add_files=True,
        add_path=True,
        add_incl_excl=True,
    )

    # We gonna change directory into the chart directory so we add it as local
    # path for helm dependency build and helm template
    args["helm_build"].append(".")
    args["helm_tmpl"].append(".")

    # Ger logger
    log = pw.get_logger(
        args["wrapper"].debug,
    )

    # Here we store paths fo the changed charts
    charts = {}

    # Calculate length of the path to the directory with charts
    path_items = args["wrapper"].charts_path.split(os.sep)

    # Take only paths pointing to files in the chart
    path_items_len = len(path_items) + 1

    # Includes and excludes
    if args["wrapper"].include_charts is not None:
        include_charts = list(map(str.strip, args["wrapper"].include_charts.split(",")))
    else:
        include_charts = []

    if args["wrapper"].exclude_charts is not None:
        exclude_charts = list(map(str.strip, args["wrapper"].exclude_charts.split(",")))
    else:
        exclude_charts = []

    for f in args["wrapper"].FILES:
        if f.startswith("%s%s" % (args["wrapper"].charts_path, os.sep)):
            items = f.split(os.sep)
            name = items[path_items_len - 1]

            # Skip chart if it's not included or is excluded
            if (
                include_charts and name not in include_charts
            ) or name in exclude_charts:
                continue

            if len(items) > path_items_len:
                path = os.sep.join(items[0:path_items_len])

                if path not in charts:
                    charts[name] = path

    # Change directory to the chart and run tests
    for name, path in charts.items():
        print("Testing chart '%s'" % name)

        with cd(path):
            # Parse config file
            config_args = pw.parse_config(
                args["wrapper"].config,
            )

            # Merge the args from config file and from command line
            if config_args:
                args["kubeconform"] = config_args + args["kubeconform"]

            # Get list of values files
            values_files = pw.get_values_files(
                args["wrapper"].values_dir,
                args["wrapper"].values_pattern,
            )

            # Run tests
            try:
                if values_files:
                    for values_file in values_files:
                        log.debug("Testing with an extra values file %s" % values_file)

                        pw.run_test(args, values_file)
                else:
                    log.debug("Testing without any extra values files")

                    pw.run_test(args)
            except Exception as e:
                log.error("Testing failed: %s" % e)

                sys.exit(1)


if __name__ == "__main__":
    main()
