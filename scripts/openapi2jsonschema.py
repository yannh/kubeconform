#!/usr/bin/env python

# Derived from https://github.com/instrumenta/openapi2jsonschema
import yaml
import json
import sys
import os
import urllib.request

def iteritems(d):
    if hasattr(dict, "iteritems"):
        return d.iteritems()
    else:
        return iter(d.items())


def additional_properties(data):
    "This recreates the behaviour of kubectl at https://github.com/kubernetes/kubernetes/blob/225b9119d6a8f03fcbe3cc3d590c261965d928d0/pkg/kubectl/validation/schema.go#L312"
    new = {}
    try:
        for k, v in iteritems(data):
            new_v = v
            if isinstance(v, dict):
                if "properties" in v:
                    if "additionalProperties" not in v:
                        v["additionalProperties"] = False
                new_v = additional_properties(v)
            else:
                new_v = v
            new[k] = new_v
        return new
    except AttributeError:
        return data


def replace_int_or_string(data):
    new = {}
    try:
        for k, v in iteritems(data):
            new_v = v
            if isinstance(v, dict):
                if "format" in v and v["format"] == "int-or-string":
                    new_v = {"oneOf": [{"type": "string"}, {"type": "integer"}]}
                else:
                    new_v = replace_int_or_string(v)
            elif isinstance(v, list):
                new_v = list()
                for x in v:
                    new_v.append(replace_int_or_string(x))
            else:
                new_v = v
            new[k] = new_v
        return new
    except AttributeError:
        return data


def allow_null_optional_fields(data, parent=None, grand_parent=None, key=None):
    new = {}
    try:
        for k, v in iteritems(data):
            new_v = v
            if isinstance(v, dict):
                new_v = allow_null_optional_fields(v, data, parent, k)
            elif isinstance(v, list):
                new_v = list()
                for x in v:
                    new_v.append(allow_null_optional_fields(x, v, parent, k))
            elif isinstance(v, str):
                is_non_null_type = k == "type" and v != "null"
                has_required_fields = grand_parent and "required" in grand_parent
                if is_non_null_type and not has_required_field:
                    new_v = [v, "null"]
            new[k] = new_v
        return new
    except AttributeError:
        return data


def append_no_duplicates(obj, key, value):
    """
    Given a dictionary, lookup the given key, if it doesn't exist create a new array.
    Then check if the given value already exists in the array, if it doesn't add it.
    """
    if key not in obj:
        obj[key] = []
    if value not in obj[key]:
        obj[key].append(value)


def write_schema_file(schema, filename):
    schemaJSON = ""

    schema = additional_properties(schema)
    schema = replace_int_or_string(schema)
    schemaJSON = json.dumps(schema, indent=2)

    # Dealing with user input here..
    filename = os.path.basename(filename)
    f = open(filename, "w")
    f.write(schemaJSON)
    f.close()
    print("JSON schema written to {filename}".format(filename=filename))


if len(sys.argv) == 0:
    print("missing file")
    exit(1)

for crdFile in sys.argv[1:]:
    if crdFile.startswith("http"):
      f = urllib.request.urlopen(crdFile)
    else:
      f = open(crdFile)
    with f:
        defs = []
        for y in yaml.load_all(f, Loader=yaml.SafeLoader):
            if y is None:
                continue
            if "items" in y:
                defs.extend(y["items"])
            if "kind" not in y:
                continue
            if y["kind"] != "CustomResourceDefinition":
                continue
            else:
                defs.append(y)

        for y in defs:
            filename_format = os.getenv("FILENAME_FORMAT", "{kind}_{version}")
            filename = ""
            if "spec" in y and "validation" in y["spec"] and "openAPIV3Schema" in y["spec"]["validation"]:
                filename = filename_format.format(
                    kind=y["spec"]["names"]["kind"],
                    group=y["spec"]["group"].split(".")[0],
                    version=y["spec"]["version"],
                ).lower() + ".json"

                schema = y["spec"]["validation"]["openAPIV3Schema"]
                write_schema_file(schema, filename)
            elif "spec" in y and "versions" in y["spec"]:
                for version in y["spec"]["versions"]:
                    if "schema" in version and "openAPIV3Schema" in version["schema"]:
                        filename = filename_format.format(
                            kind=y["spec"]["names"]["kind"],
                            group=y["spec"]["group"].split(".")[0],
                            version=version["name"],
                        ).lower() + ".json"

                        schema = version["schema"]["openAPIV3Schema"]
                        write_schema_file(schema, filename)

exit(0)
