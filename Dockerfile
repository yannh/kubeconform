FROM alpine:3.14 as certs
RUN apk add ca-certificates

FROM scratch AS kubeconform
MAINTAINER Yann HAMON <yann@mandragor.org>
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY kubeconform /
ENTRYPOINT ["/kubeconform"]
