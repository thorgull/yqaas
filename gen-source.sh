#!/usr/bin/env bash

set -e

rm -rf gen

docker run --rm -v $(pwd):/local openapitools/openapi-generator-cli generate \
  -i /local/yqaas.yaml \
  -g go-server \
  -o /local/gen \
  --additional-properties=packageName=api,sourceFolder=api,outputAsLibrary=true

goimports -w gen/api/*.go