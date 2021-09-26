#!/usr/bin/make -f

# This is really early days

build-python-bats:
	docker build -t python-bats -f Dockerfile.bats .

docker-acceptance: build-python-bats
	docker run --entrypoint "/usr/bin/bats" -t python-bats /code/acceptance.bats
