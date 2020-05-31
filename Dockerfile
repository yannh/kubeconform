FROM golang:1.14 AS builder

RUN mkdir -p github.com/yannh/kubeconform
COPY . github.com/yannh/kubeconform/
WORKDIR github.com/yannh/kubeconform
RUN make build-static

FROM scratch AS kubeconform
MAINTAINER Yann HAMON <yann@mandragor.org>
COPY --from=builder /go/github.com/yannh/kubeconform/bin/kubeconform /
ENTRYPOINT ["/kubeconform"]