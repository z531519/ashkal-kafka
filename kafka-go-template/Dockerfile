
## Build
FROM golang:alpine AS build

RUN apk update && apk add gcc g++ libc-dev librdkafka-dev pkgconf

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY pkg pkg
COPY cmd cmd
COPY config.yaml config.yaml

RUN go build -tags musl -o demo ./cmd/demo

## Deploy
FROM alpine:latest

RUN apk update && apk add gcc g++ libc-dev librdkafka-dev pkgconf
WORKDIR /app
COPY --from=build /app/demo /app/demo
COPY --from=build /app/config.yaml /app/config.yaml

ENTRYPOINT ["/app/demo"]