# Build context must be the parent `scripts/` directory so that
# `fixtures/` is reachable. See the Makefile in this folder.
FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go/go.mod go/go.sum ./
RUN go mod download
COPY go/openapi2jsonschema.go go/openapi2jsonschema_test.go ./
RUN go build -o /openapi2jsonschema .

FROM alpine:3.20
RUN apk --no-cache add bats diffutils
COPY --from=build /openapi2jsonschema /code/go/openapi2jsonschema
COPY fixtures /code/fixtures
COPY go/acceptance.bats /code/go/
WORKDIR /code/go
