FROM bats/bats:v1.1.0
RUN apk --no-cache add ca-certificates
COPY bin/kubeconform /bin/
COPY acceptance.bats /acceptance.bats
COPY fixtures /fixtures
