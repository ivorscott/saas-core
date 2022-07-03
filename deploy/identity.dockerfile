FROM golang:1.18-alpine as base

LABEL org.opencontainers.image.authors="devpie"

ENV CGO_ENABLED=0

WORKDIR /core

RUN mkdir log
COPY ../go.* ./
COPY ../cmd/identity ./cmd/identity
COPY ../internal/identity ./internal/identity
COPY ../pkg ./pkg

RUN go mod download && go build ./cmd/identity

CMD ["./identity"]
