#!/usr/bin/make -f

.PHONY: test build-python-bats docker-test docker-acceptance build-image

test: build-python-bats docker-test docker-acceptance

build-python-bats:
	docker build -t python-bats -f Dockerfile.bats .

docker-test: build-python-bats
	docker run --entrypoint "/usr/local/bin/pytest" -t python-bats openapi2jsonschema.py

docker-acceptance: build-python-bats
	docker run --entrypoint "/usr/bin/bats" -t python-bats /code/acceptance.bats

build-image: docker-test docker-acceptance
	docker build -t openapi2jsonschema -f Dockerfile .
