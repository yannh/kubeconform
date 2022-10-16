FROM python:3.9.15-alpine3.16
LABEL org.opencontainers.image.authors="Various contributors" \
      org.opencontainers.image.source="https://github.com/yannh/kubeconform/tree/master/scripts" \
      org.opencontainers.image.description="Convert OpenAPI CRD to JSON Schema" \
      org.opencontainers.image.licenses="Apache License 2.0" \
      org.opencontainers.image.title="openapi2jsonschema" \
      org.opencontainers.image.url="https://github.com/yannh/kubeconform/tree/master/scripts/"
MAINTAINER Yann HAMON <yann@mandragor.org>
ENV FILENAME_FORMAT='{kind}-{group}-{version}'

COPY openapi2jsonschema.py /
RUN pip install pyyaml

ENTRYPOINT ["/openapi2jsonschema"]
