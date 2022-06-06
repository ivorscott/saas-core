FROM golang:1.18-alpine as base

LABEL org.opencontainers.image.authors="devpie"

ENV CGO_ENABLED=0

WORKDIR /core

COPY ../go.* ./
RUN go mod download

COPY ../cmd/admin ./cmd/admin
COPY ../internal/admin ./internal/admin
COPY ../pkg ./pkg
RUN go build ./cmd/admin

CMD ["./admin"]
