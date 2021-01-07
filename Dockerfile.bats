FROM bats/bats:v1.2.1
RUN apk --no-cache add ca-certificates parallel python
COPY bin/kubeconform /code/bin/
COPY acceptance.bats /code/acceptance.bats
COPY fixtures /code/fixtures
