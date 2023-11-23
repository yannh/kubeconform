#!/usr/bin/make -f

# This is really early days

test: build-python-bats docker-test docker-acceptance

build-python-bats:
	docker build -t python-bats -f Dockerfile.bats .

docker-test: build-python-bats
	docker run --entrypoint "/usr/local/bin/pytest" -t python-bats openapi2jsonschema.py

docker-acceptance: build-python-bats
	docker run --entrypoint "/usr/bin/bats" -t python-bats /code/acceptance.bats

get-crds:
# get emissary mapping v2 CRD
	curl -Lnk  https://app.getambassador.io/yaml/emissary/3.8.0/emissary-crds.yaml | yq -o json |\
	jq -s '.[] | select(.metadata.name=="mappings.getambassador.io") | del(.spec.versions[] | select( .name != "v2"))' |\
	yq -p json -o yaml > fixtures/emissary-mapping-crd.yaml