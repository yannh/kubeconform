# Build context must be the repo root so that scripts/fixtures/ is reachable.
# See the Makefile in this folder.
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY openapi2jsonschema-go/go.mod openapi2jsonschema-go/go.sum ./
COPY openapi2jsonschema-go/vendor ./vendor
COPY openapi2jsonschema-go/openapi2jsonschema.go openapi2jsonschema-go/openapi2jsonschema_test.go ./
RUN go build -mod=vendor -o /openapi2jsonschema .

FROM alpine:3.20
RUN apk --no-cache add bats diffutils
COPY --from=build /openapi2jsonschema /code/openapi2jsonschema-go/openapi2jsonschema
COPY scripts/fixtures /code/fixtures
COPY openapi2jsonschema-go/acceptance.bats /code/openapi2jsonschema-go/
WORKDIR /code/openapi2jsonschema-go
