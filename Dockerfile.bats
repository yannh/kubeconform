FROM bats/bats:1.11.0
RUN apk --no-cache add ca-certificates parallel libxml2-utils
COPY bin/kubeconform /code/bin/
COPY acceptance.bats acceptance-nonetwork.bats /code/
COPY fixtures /code/fixtures
