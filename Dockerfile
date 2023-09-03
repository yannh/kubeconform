FROM alpine:3.18.3 as certs
RUN apk add ca-certificates

FROM scratch AS kubeconform
LABEL org.opencontainers.image.authors="yann@mandragor.org" \
      org.opencontainers.image.source="https://github.com/yannh/kubeconform/" \
      org.opencontainers.image.description="A Kubernetes manifests validation tool" \
      org.opencontainers.image.documentation="https://github.com/yannh/kubeconform/" \
      org.opencontainers.image.licenses="Apache License 2.0" \
      org.opencontainers.image.title="kubeconform" \
      org.opencontainers.image.url="https://github.com/yannh/kubeconform/"
MAINTAINER Yann HAMON <yann@mandragor.org>
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY kubeconform /
ENTRYPOINT ["/kubeconform"]
