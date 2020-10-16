#!/usr/bin/make -f

.PHONY: test-build test build build-static docker-test docker-build-static build-bats docker-acceptance docker-image

test-build: test build

test:
	go test ./...

build:
	go build -o bin/ ./...

docker-image:
	docker build -t kubeconform .

build-static:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bin/ ./...

docker-test:
	docker run -t -v $$PWD:/go/src/github.com/yannh/kubeconform -w /go/src/github.com/yannh/kubeconform golang:1.14 make test

docker-build-static:
	docker run -t -v $$PWD:/go/src/github.com/yannh/kubeconform -w /go/src/github.com/yannh/kubeconform golang:1.14 make build-static

build-bats:
	docker build -t bats -f Dockerfile.bats .

docker-acceptance: build-bats
	docker run -t bats acceptance.bats

release:
	docker run -e GITHUB_TOKEN -t -v $$PWD:/go/src/github.com/yannh/kubeconform -w /go/src/github.com/yannh/kubeconform goreleaser/goreleaser:v0.138 goreleaser release --rm-dist