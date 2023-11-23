FROM bats/bats:v1.2.1
RUN apk --no-cache add ca-certificates parallel libxml2-utils
COPY dist/kubeconform_linux_amd64_v1/kubeconform /code/bin/
COPY acceptance.bats acceptance-nonetwork.bats /code/
COPY fixtures /code/fixtures
COPY scripts/fixtures/mapping_v2-expected.json /code/fixtures/registry/mapping_v2.json
