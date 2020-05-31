#!/usr/bin/make -f

test-build: test build

test:
	go test ./...

build:
	go build -o bin/kubeconform

build-static:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bin/kubeconform

docker-test:
	docker run -t -v $$PWD:/go/src/github.com/yannh/kubeconform -w /go/src/github.com/yannh/kubeconform golang:1.14 make test

docker-build-static:
	docker run -t -v $$PWD:/go/src/github.com/yannh/kubeconform -w /go/src/github.com/yannh/kubeconform golang:1.14 make build-static

