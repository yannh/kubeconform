FROM python:3.9.7-alpine3.14
RUN apk --no-cache add bats
COPY acceptance.bats openapi2jsonschema.py requirements.txt /code/
COPY fixtures /code/fixtures
WORKDIR /code
RUN pip install -r requirements.txt
