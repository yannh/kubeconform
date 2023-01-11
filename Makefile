#!/usr/bin/make -f

RELEASE_VERSION ?= latest

.PHONY: local-test local-build local-build-static docker-test docker-build docker-build-static build-bats docker-acceptance release update-deps build-single-target

local-test:
	go test -race ./...

local-build:
	go build -o bin/ ./...

local-build-static:
	CGO_ENABLED=0 GOFLAGS=-mod=vendor GOOS=linux GOARCH=amd64 GO111MODULE=on go build -trimpath -tags=netgo -ldflags "-extldflags=\"-static\""  -a -o bin/ ./...

# These only used for development. Release artifacts and docker images are produced by goreleaser.
docker-test:
	docker run -t -v $$PWD:/go/src/github.com/yannh/kubeconform -w /go/src/github.com/yannh/kubeconform golang:1.17 make local-test

docker-build:
	docker run -t -v $$PWD:/go/src/github.com/yannh/kubeconform -w /go/src/github.com/yannh/kubeconform golang:1.17 make local-build

docker-build-static:
	docker run -t -v $$PWD:/go/src/github.com/yannh/kubeconform -w /go/src/github.com/yannh/kubeconform golang:1.17 make local-build-static

build-bats:
	docker build -t bats -f Dockerfile.bats .

docker-acceptance: build-bats
	docker run -t bats -p acceptance.bats
	docker run --network none -t bats -p acceptance-nonetwork.bats

goreleaser-build-static:
	docker run -t -e GOOS=linux -e GOARCH=amd64 -v $$PWD:/go/src/github.com/yannh/kubeconform -w /go/src/github.com/yannh/kubeconform goreleaser/goreleaser:v1.11.5 build --single-target --skip-post-hooks --rm-dist --snapshot
	cp dist/kubeconform_linux_amd64_v1/kubeconform bin/

release:
	docker run -e GITHUB_TOKEN -e GIT_OWNER -t -v /var/run/docker.sock:/var/run/docker.sock -v $$PWD:/go/src/github.com/yannh/kubeconform -w /go/src/github.com/yannh/kubeconform goreleaser/goreleaser:v1.11.5 release --rm-dist

update-deps:
	go get -u ./...
	go mod tidy

update-junit-xsd:
	curl https://raw.githubusercontent.com/junit-team/junit5/main/platform-tests/src/test/resources/jenkins-junit.xsd > fixtures/junit.xsd
