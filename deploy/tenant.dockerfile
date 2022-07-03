FROM golang:1.18-alpine as base

LABEL org.opencontainers.image.authors="devpie"

ENV CGO_ENABLED=0

WORKDIR /core

RUN mkdir log
COPY ../go.* ./
COPY ../cmd/tenant ./cmd/tenant
COPY ../internal/tenant ./internal/tenant
COPY ../pkg ./pkg

RUN go mod download && go build ./cmd/tenant

CMD ["./tenant"]
