FROM bats/bats:v1.2.1
RUN apk --no-cache add ca-certificates parallel
COPY dist/kubeconform_linux_amd64/kubeconform /code/bin/
COPY acceptance.bats acceptance-nonetwork.bats /code/
COPY fixtures /code/fixtures
