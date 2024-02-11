FROM openapitools/openapi-generator-cli as openapi-generator
COPY yqaas.yaml /home/yqaas.yaml

RUN docker-entrypoint.sh generate \
    -i /home/yqaas.yaml \
    -g go-server \
    -o /home/gen \
    --additional-properties=packageName=api,sourceFolder=api,outputAsLibrary=true

FROM golang:1.22 AS build
RUN go install golang.org/x/tools/cmd/goimports@latest

WORKDIR /go/src
COPY --from=openapi-generator /home/gen ./gen
COPY impl ./impl
COPY main.go .
COPY go.mod .
COPY go.sum .

ENV CGO_ENABLED=0
RUN goimports -w gen/api/*.go
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o yqaas .

FROM scratch AS runtime
LABEL org.opencontainers.image.authors="https://github.com/thorgull"
LABEL org.opencontainers.image.url="https://github.com/thorgull/yqaas"
LABEL org.opencontainers.image.source="https://github.com/thorgull/yqaas"
LABEL org.opencontainers.image.licenses="AGPL-3.0-or-later"
LABEL org.opencontainers.image.title="YQ As A Service"
WORKDIR /
COPY --from=ghcr.io/jqlang/jq /jq ./
COPY --from=build /go/src/yqaas ./
COPY yqaas.yaml ./
EXPOSE 8080/tcp
ENV PATH=/
ENTRYPOINT ["./yqaas"]
