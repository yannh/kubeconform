#!/usr/bin/make -f

GO_VERSION ?= 1.24.3

.PHONY: test local-test local-build local-build-static docker-test docker-build docker-build-static build-bats docker-acceptance update-deps clean

test: local-test docker-acceptance

local-test:
	go test -race ./... -count=1

local-build:
	go build -o bin/openapi2jsonschema .

local-build-static:
	CGO_ENABLED=0 GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH) GO111MODULE=on go build -trimpath -tags=netgo -ldflags "-extldflags=\"-static\""  -a -o bin/openapi2jsonschema .

# Mount the repo root so scripts/fixtures/ is reachable from the container.
docker-test:
	docker run -t -v $$PWD/..:/src -w /src/openapi2jsonschema-go golang:$(GO_VERSION) make local-test

docker-build:
	docker run -t -v $$PWD/..:/src -w /src/openapi2jsonschema-go golang:$(GO_VERSION) make local-build

docker-build-static:
	docker run -t -v $$PWD/..:/src -w /src/openapi2jsonschema-go golang:$(GO_VERSION) make local-build-static

# Build context is the repo root so the image can pull in scripts/fixtures/.
build-bats:
	docker build -t openapi2jsonschema-go-bats -f Dockerfile.bats ..

docker-acceptance: build-bats
	docker run --entrypoint "/usr/bin/bats" -t openapi2jsonschema-go-bats /code/openapi2jsonschema-go/acceptance.bats

update-deps:
	go get -u ./...
	go mod tidy

clean:
	rm -rf bin
	rm -f openapi2jsonschema prometheus_v1.json prometheus-monitoring-v1.json
