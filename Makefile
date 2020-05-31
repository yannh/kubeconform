#!/usr/bin/make -f

build:
	go build -o bin/kubeconform

build-static:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bin/kubeconform
