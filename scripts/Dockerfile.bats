FROM python:3.9.7-alpine3.14
RUN apk --no-cache add bats
COPY requirements.txt /code/
RUN pip install -r /code/requirements.txt
COPY fixtures /code/fixtures
COPY acceptance.bats openapi2jsonschema.py /code/
WORKDIR /code
