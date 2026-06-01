#!/usr/bin/make -f

test: unit-test docker-acceptance

build:
	go build -o openapi2jsonschema .

unit-test:
	go test ./...

build-go-bats:
	docker build -t go-bats -f Dockerfile.bats ..

docker-acceptance: build-go-bats
	docker run --entrypoint "/usr/bin/bats" -t go-bats /code/go/acceptance.bats

clean:
	rm -f openapi2jsonschema prometheus_v1.json prometheus-monitoring-v1.json
